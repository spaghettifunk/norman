package offline_blob_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/rs/zerolog/log"

	paws "github.com/spaghettifunk/norman/pkg/aws"
	offline_blob "github.com/spaghettifunk/norman/pkg/offline/blob"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
)

var (
	s3TestEndpoint = "localhost:4572" // TODO: this is the endpoint of localstack
	s3TestBucket   = "test"
	s3TestRegion   = "eu-west-1"
)

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
	bucket := &offline_blob.S3Bucket{Bucket: s3TestBucket}
	sess := paws.NewAWSSession(s3TestRegion, s3TestEndpoint, true)

	s := offline_blob.NewS3Client(bucket, sess)
	assert.NotNil(t, s)
}

func TestGetObject(t *testing.T) {
	bucket := &offline_blob.S3Bucket{Bucket: s3TestBucket}
	sess := paws.NewAWSSession(s3TestRegion, s3TestEndpoint, true)

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
	bucket := &offline_blob.S3Bucket{Bucket: s3TestBucket}
	sess := paws.NewAWSSession(s3TestRegion, s3TestEndpoint, true)

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
	bucket := &offline_blob.S3Bucket{Bucket: s3TestBucket}
	sess := paws.NewAWSSession(s3TestRegion, s3TestEndpoint, true)

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
	bucket := &offline_blob.S3Bucket{Bucket: s3TestBucket}
	sess := paws.NewAWSSession(s3TestRegion, s3TestEndpoint, true)

	drop := CreateTestS3Bucket(t, bucket, sess)
	defer drop()

	s := offline_blob.NewS3Client(bucket, sess)

	// Key should not exists
	if s.ExistsObject("foo/bar2.txt") {
		t.Failed()
	}
}

func TestDeleteBucket(t *testing.T) {
	bucket := &offline_blob.S3Bucket{Bucket: "test1"}
	sess := paws.NewAWSSession(s3TestRegion, s3TestEndpoint, true)

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
