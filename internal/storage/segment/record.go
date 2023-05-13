package segment

import (
	"github.com/spaghettifunk/norman/internal/common/types"
)

type Column struct {
	Name      string          `json:"-"`
	FieldType types.FieldType `json:"-"`
}

func NewColumn(name string, ft types.FieldType) *Column {
	return &Column{
		Name:      name,
		FieldType: ft,
	}
}

// Record is the data stored in our Segment - it can be seen as
// a row in a table
type Record struct {
	Values map[*Column]interface{} `json:"-"`
}

func NewEmptyRecord() *Record {
	return &Record{}
}

func (r *Record) AddValue(val interface{}, col *Column) error {
	r.Values[col] = val
	return nil
}
