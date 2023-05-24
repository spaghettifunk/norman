package segment

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"sync"

	"github.com/goccy/go-json"
	"github.com/spaghettifunk/norman/internal/common/types"
)

var (
	// enc defines the encoding that we persist record sizes and index entries in
	enc = binary.BigEndian
)

const (
	// lenWidth defines the number of bytes used to store the record’s length
	lenWidth = 8
)

type Column struct {
	Name      string          `json:"-"`
	FieldType types.FieldType `json:"-"`
	DataType  types.DataType  `json:"-"`
	Values    []interface{}   `json:"-"`
	// Store the data as a columnar storage
	file *os.File
	mu   sync.Mutex
	buf  *bufio.Writer
	size uint64
}

func NewColumn(dir string, name string, ft types.FieldType, dt types.DataType) (*Column, error) {
	fn := fmt.Sprintf("%s%s.segment", dir, name)

	var err error
	columnFile, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(columnFile.Name())
	if err != nil {
		return nil, err
	}
	size := uint64(fi.Size())

	return &Column{
		Name:      name,
		FieldType: ft,
		DataType:  dt,
		Values:    make([]interface{}, MaxNumberOfEntriesPerSegment),
		buf:       bufio.NewWriter(columnFile),
		size:      size,
	}, nil
}

func (col *Column) InsertData(val interface{}) (n uint64, pos uint64, err error) {
	// transform the incoming object into bytes
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(val); err != nil {
		return 0, 0, err
	}

	// write to the file and return the last position
	col.mu.Lock()
	defer col.mu.Unlock()
	pos = col.size
	if err := binary.Write(col.buf, enc, uint64(len(b.Bytes()))); err != nil {
		return 0, 0, err
	}
	w, err := col.buf.Write(b.Bytes())
	if err != nil {
		return 0, 0, err
	}
	w += lenWidth
	col.size += uint64(w)

	// append to list
	col.Values = append(col.Values, val)

	return uint64(w), pos, nil
}

func (col *Column) Read(pos uint64) ([]byte, error) {
	col.mu.Lock()
	defer col.mu.Unlock()
	if err := col.buf.Flush(); err != nil {
		return nil, err
	}
	size := make([]byte, lenWidth)
	if _, err := col.file.ReadAt(size, int64(pos)); err != nil {
		return nil, err
	}
	b := make([]byte, enc.Uint64(size))
	if _, err := col.file.ReadAt(b, int64(pos+lenWidth)); err != nil {
		return nil, err
	}
	return b, nil
}

func (col *Column) ReadAt(p []byte, off int64) (int, error) {
	col.mu.Lock()
	defer col.mu.Unlock()
	if err := col.buf.Flush(); err != nil {
		return 0, err
	}
	return col.file.ReadAt(p, off)
}

// Flush persist the segment on disk
func (col *Column) Flush() error {
	col.mu.Lock()
	defer col.mu.Unlock()
	err := col.buf.Flush()
	if err != nil {
		return err
	}
	return col.file.Close()
}
