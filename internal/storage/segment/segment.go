package segment

import (
	"fmt"
	"path"
	"sync"
	"sync/atomic"

	"github.com/apache/arrow/go/arrow"
	"github.com/apache/arrow/go/arrow/array"
	"github.com/apache/arrow/go/arrow/memory"
	"github.com/google/uuid"
	"github.com/xitongsys/parquet-go/source"
	"github.com/xitongsys/parquet-go/writer"
)

type Segment struct {
	ID uuid.UUID `json:"-"`
	// count the inserted events
	counter uint32
	mu      sync.Mutex
	pFile   source.ParquetFile
	schema  *arrow.Schema
	builder *array.RecordBuilder
	writer  *writer.ArrowWriter
}

func NewSegment(dir string, schema *arrow.Schema) (*Segment, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	// create a parquet file
	fileName := path.Join(dir, fmt.Sprintf("%s.parquet", id.String()))
	pFile, err := newLocalParquetWrite(fileName)
	if err != nil {
		return nil, err
	}

	// create the apache arrow writer
	w, err := writer.NewArrowWriter(schema, pFile, 1)
	if err != nil {
		return nil, err
	}

	// create the new record builder for inserting data to arrow file
	mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
	b := array.NewRecordBuilder(mem, schema)

	return &Segment{
		ID:      id,
		schema:  schema,
		mu:      sync.Mutex{},
		pFile:   pFile,
		builder: b,
		writer:  w,
	}, nil
}

func (s *Segment) InsertData(data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// insert data to arrow here
	
	s.builder.Field(0)

	// increase counter
	atomic.AddUint32(&s.counter, 1)
	return nil
}

func (s *Segment) GetLength(colName string) int {
	return int(s.counter)
}

// Flush persist the segment on disk
func (s *Segment) Flush() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// write arrow file
	defer s.builder.Release()
	rec := s.builder.NewRecord()
	if err := s.writer.WriteArrow(rec); err != nil {
		return err
	}

	// store the parquet file
	if err := s.pFile.Close(); err != nil {
		return nil
	}

	// reset counter now that we flushed data to disk
	atomic.SwapUint32(&s.counter, 0)
	return nil
}
