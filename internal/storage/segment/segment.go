package segment

import (
	"fmt"
	"path"
	"sync"
	"sync/atomic"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/apache/arrow/go/v12/arrow/array"
	"github.com/apache/arrow/go/v12/arrow/memory"

	"github.com/apache/arrow/go/v12/parquet"
	"github.com/apache/arrow/go/v12/parquet/compress"
	"github.com/apache/arrow/go/v12/parquet/pqarrow"

	"github.com/google/uuid"
)

type Segment struct {
	ID         uuid.UUID `json:"-"`
	evtCounter uint32
	mu         sync.Mutex
	pFile      *LocalParquet
	schema     *arrow.Schema
	evtStruct  *arrow.StructType
	builder    *array.RecordBuilder
	writer     *pqarrow.FileWriter
}

func NewSegment(dir string, schema *arrow.Schema) (*Segment, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	// create a parquet file
	fileName := path.Join(dir, fmt.Sprintf("%s.parquet", id.String()))
	pFile, err := NewLocalParquet(fileName)
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

	// create the new record builder for inserting data to arrow file
	mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
	b := array.NewRecordBuilder(mem, schema)

	return &Segment{
		ID:        id,
		schema:    schema,
		evtStruct: arrow.StructOf(schema.Fields()...),
		mu:        sync.Mutex{},
		pFile:     pFile,
		builder:   b,
		writer:    w,
	}, nil
}

func (s *Segment) InsertData(data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.builder.UnmarshalJSON(data); err != nil {
		return err
	}

	// increment counter
	atomic.AddUint32(&s.evtCounter, 1)

	return nil
}

func (s *Segment) GetCounter() uint32 {
	return s.evtCounter
}

// Flush persist the segment on disk
func (s *Segment) Flush() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	rec := s.builder.NewRecord()

	// closable
	defer s.pFile.Close()
	defer s.builder.Release()
	defer s.writer.Close()
	defer rec.Release()

	// write arrow file
	if err := s.writer.WriteBuffered(rec); err != nil {
		return err
	}

	// reset counter
	atomic.CompareAndSwapUint32(&s.evtCounter, s.evtCounter, 0)
	return nil
}
