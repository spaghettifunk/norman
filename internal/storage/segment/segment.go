package segment

import (
	"bytes"
	"fmt"
	"os"
	"sync"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

// TODO: add mutex here
type Segment struct {
	ID   uuid.UUID `json:"-"`
	Rows []*Row    `json:"-"`
}

// segPool is in charged of Pooling eventual requests in coming. This will help to reduce the alloc/s
// and efficiently improve the garbage collection operations
var segPool = sync.Pool{
	New: func() interface{} { return new(Segment) },
}

func NewSegment() (*Segment, error) {
	// get a new object from the pool and then dispose it
	sp := segPool.Get().(*Segment)
	defer segPool.Put(sp)

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	sp.ID = id
	return sp, nil
}

func (s *Segment) InsertRow(values map[*Column]interface{}) error {
	row := NewEmptyRow()
	for k, v := range values {
		// TODO: check error here once we have actual errors to check
		_ = row.AddValue(v, k)
	}
	s.Rows = append(s.Rows, row)
	return nil
}

// Flush persist the segment on disk on the given directory
// TODO: add mutex here
func (s *Segment) Flush(dir string) error {
	fp := fmt.Sprintf("%s/%s.norman", dir, s.ID.String())
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(s.Rows); err != nil {
		return err
	}
	if err := os.WriteFile(fp, b.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}

func (s *Segment) Reset() error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	s.ID = id
	s.Rows = make([]*Row, 0)
	return nil
}
