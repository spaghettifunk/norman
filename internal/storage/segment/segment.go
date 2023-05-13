package segment

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"os"
	"sync"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

var (
	// enc defines the encoding that we persist record sizes and index entries in
	enc = binary.BigEndian
)

const (
	// lenWidth defines the number of bytes used to store the record’s length
	lenWidth = 8
)

// TODO: add mutex here
type Segment struct {
	ID uuid.UUID `json:"-"`
	// below are related to the file where the records are saved
	*os.File
	mu   sync.Mutex
	buf  *bufio.Writer
	size uint64
}

func NewSegment(f *os.File) (*Segment, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	size := uint64(fi.Size())

	return &Segment{
		ID:   id,
		File: f,
		size: size,
		buf:  bufio.NewWriter(f),
	}, nil
}

func (s *Segment) InsertRow(values map[*Column]interface{}) (n uint64, pos uint64, err error) {
	rec := NewEmptyRecord()

	// iterate through the values and run a validation
	for k, v := range values {
		// TODO: check error here once we have actual errors to check
		_ = rec.AddValue(v, k)
	}

	// transform the incoming object into bytes
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(rec); err != nil {
		return 0, 0, err
	}

	// write to the file and return the last position
	s.mu.Lock()
	defer s.mu.Unlock()
	pos = s.size
	if err := binary.Write(s.buf, enc, uint64(len(b.Bytes()))); err != nil {
		return 0, 0, err
	}
	w, err := s.buf.Write(b.Bytes())
	if err != nil {
		return 0, 0, err
	}
	w += lenWidth
	s.size += uint64(w)

	return uint64(w), pos, nil
}

func (s *Segment) Read(pos uint64) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.buf.Flush(); err != nil {
		return nil, err
	}
	size := make([]byte, lenWidth)
	if _, err := s.File.ReadAt(size, int64(pos)); err != nil {
		return nil, err
	}
	b := make([]byte, enc.Uint64(size))
	if _, err := s.File.ReadAt(b, int64(pos+lenWidth)); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Segment) ReadAt(p []byte, off int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.buf.Flush(); err != nil {
		return 0, err
	}
	return s.File.ReadAt(p, off)
}

// Flush persist the segment on disk
func (s *Segment) Flush() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := s.buf.Flush()
	if err != nil {
		return err
	}
	return s.File.Close()
}
