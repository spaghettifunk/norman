package segment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
)

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
func (s *Segment) Flush(dir string) error {
	fp := fmt.Sprintf("%s/%s.norman", dir, s.ID.String())

	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(s.Rows); err != nil {
		return err
	}

	err := os.WriteFile(fp, b.Bytes(), 0644)
	return err
}
