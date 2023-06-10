package offline_blob_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog/log"

	offline_blob "github.com/spaghettifunk/norman/internal/storage/ingestion/offline/blob"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
)

var (
	s3TestBucket = "test"
)

var MockSession = func() *session.Session {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	return session.Must(session.NewSession(&aws.Config{
		DisableSSL:  aws.Bool(true),
		Endpoint:    aws.String(server.URL),
		Credentials: credentials.NewStaticCredentials("AKID", "SECRET", "SESSION"),
		Region:      aws.String("mock-region"),
	}))
}()

// CreateTestS3Bucket returns a bucket and defer a drop
func CreateTestS3Bucket(t *testing.T, bucket *offline_blob.S3Bucket, sess *session.Session) func() {
	s := offline_blob.NewS3Client(bucket, sess)
	if _, err := s.CreateS3Bucket(&offline_blob.S3Bucket{Bucket: bucket.Bucket}); err != nil {
		log.Error().Msg(err.Error())
		t.FailNow()
	}
	return func() {
		if ok, err := s.DeleteS3Bucket(bucket); !ok || err != nil {
			log.Error().Msg(err.Error())
			t.FailNow()
		}
	}
}

func TestNewS3Client(t *testing.T) {
	t.Skip("unable to mock the session")

	bucket := &offline_blob.S3Bucket{Bucket: s3TestBucket}
	sess := MockSession

	s := offline_blob.NewS3Client(bucket, sess)
	assert.NotNil(t, s)
}

func TestGetObject(t *testing.T) {
	t.Skip("unable to mock the session")

	bucket := &offline_blob.S3Bucket{Bucket: s3TestBucket}
	sess := MockSession

	drop := CreateTestS3Bucket(t, bucket, sess)
	defer drop()

	s := offline_blob.NewS3Client(bucket, sess)

	_, err := s.Service.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s3TestBucket),
		Body:   bytes.NewReader([]byte("What is the meaning of life? 42.")),
		Key:    aws.String("foo/bar.txt"),
	})
	if err != nil {
		t.Failed()
	}

	f, err := s.GetObject("foo/bar.txt")
	if err != nil {
		t.Failed()
	}

	// convert body to string
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(*f)
	if err != nil {
		t.Failed()
	}

	assert.Equal(t, "What is the meaning of life? 42.", buf.String())
}

func TestGetObjectFails(t *testing.T) {
	t.Skip("unable to mock the session")

	bucket := &offline_blob.S3Bucket{Bucket: s3TestBucket}
	sess := MockSession

	drop := CreateTestS3Bucket(t, bucket, sess)
	defer drop()

	s := offline_blob.NewS3Client(bucket, sess)

	f, err := s.GetObject("foo/bar2.txt")
	if err == nil {
		t.Fail()
	}

	assert.Equal(t, (*io.ReadCloser)(nil), f)
}

func TestExistsObject(t *testing.T) {
	t.Skip("unable to mock the session")

	bucket := &offline_blob.S3Bucket{Bucket: s3TestBucket}
	sess := MockSession

	drop := CreateTestS3Bucket(t, bucket, sess)
	defer drop()

	s := offline_blob.NewS3Client(bucket, sess)

	_, err := s.Service.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s3TestBucket),
		Body:   bytes.NewReader([]byte("What is the meaning of life? 42.")),
		Key:    aws.String("foo/bar.txt"),
	})
	if err != nil {
		t.Failed()
	}

	if s.ExistsObject("foo/bar.txt") == false {
		t.Failed()
	}
}

func TestExistsObjectFails(t *testing.T) {
	t.Skip("unable to mock the session")

	bucket := &offline_blob.S3Bucket{Bucket: s3TestBucket}
	sess := MockSession

	drop := CreateTestS3Bucket(t, bucket, sess)
	defer drop()

	s := offline_blob.NewS3Client(bucket, sess)

	// Key should not exists
	if s.ExistsObject("foo/bar2.txt") {
		t.Failed()
	}
}

func TestDeleteBucket(t *testing.T) {
	t.Skip("unable to mock the session")

	bucket := &offline_blob.S3Bucket{Bucket: "test1"}
	sess := MockSession

	s := offline_blob.NewS3Client(bucket, sess)
	_, err := s.CreateS3Bucket(bucket)
	if err != nil {
		t.Failed()
	}

	_, err = s.Service.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s3TestBucket),
		Body:   bytes.NewReader([]byte("lorem ipsum dolor")),
		Key:    aws.String("foo/bar.txt"),
	})
	if err != nil {
		t.Failed()
	}

	_, err = s.DeleteS3Bucket(bucket)
	if err != nil {
		t.Failed()
	}
}
