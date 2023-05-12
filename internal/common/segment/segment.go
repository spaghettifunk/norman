package segment

import (
	"bytes"
	"fmt"
	"os"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

// TODO: add mutex here
type Segment struct {
	ID   uuid.UUID `json:"-"`
	Rows []*Row    `json:"-"`
}

func NewSegment() (*Segment, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	return &Segment{
		ID: id,
	}, nil
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
func (s *Segment) Flush(dir string, reset bool) error {
	fp := fmt.Sprintf("%s/%s.norman", dir, s.ID.String())

	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(s.Rows); err != nil {
		return err
	}

	if err := os.WriteFile(fp, b.Bytes(), 0644); err != nil {
		return err
	}

	if reset {
		return s.Reset()
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
