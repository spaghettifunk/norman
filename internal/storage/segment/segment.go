package segment

import (
	"errors"
	"fmt"
	"sync"

	"github.com/apache/arrow/go/v12/arrow"

	"github.com/apache/arrow/go/v12/parquet"
	"github.com/apache/arrow/go/v12/parquet/compress"
	"github.com/apache/arrow/go/v12/parquet/pqarrow"

	"github.com/google/uuid"
)

type Segment struct {
	ID     uuid.UUID `json:"-"`
	mu     sync.Mutex
	pFile  *LocalParquet
	writer *pqarrow.FileWriter
	schema *arrow.Schema
}

func NewSegment(dir string, partition int, schema *arrow.Schema) (*Segment, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	// create a parquet file
	pFile, err := NewLocalParquet(dir, fmt.Sprintf("%s_%d.parquet", id.String(), partition))
	if err != nil {
		return nil, err
	}

	// create the apache arrow writer
	props := parquet.NewWriterProperties(
		parquet.WithCompression(compress.Codecs.Snappy),
		parquet.WithDictionaryDefault(false),
		parquet.WithDataPageVersion(parquet.DataPageV1),
		parquet.WithVersion(parquet.V1_0),
	)
	w, err := pqarrow.NewFileWriter(schema, pFile, props, pqarrow.DefaultWriterProps())
	if err != nil {
		panic(err)
	}

	return &Segment{
		ID:     id,
		mu:     sync.Mutex{},
		pFile:  pFile,
		writer: w,
		schema: schema,
	}, nil
}

// Flush persist the segment on disk
func (s *Segment) Flush(record arrow.Record) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// closable
	var err error
	defer func() {
		err = errors.Join(err, s.pFile.Close())
	}()

	defer func() {
		err = errors.Join(err, s.writer.Close())
	}()

	defer record.Release()

	// write arrow file
	err = errors.Join(s.writer.WriteBuffered(record))
	return err
}
