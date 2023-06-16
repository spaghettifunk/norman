package rangeindex

import (
	"github.com/RoaringBitmap/roaring"
)

type RangeIndex struct {
	columnName string
	index      map[interface{}]*roaring.Bitmap
}

func New(columnName string) *RangeIndex {
	return &RangeIndex{
		columnName: columnName,
		index:      make(map[interface{}]*roaring.Bitmap, 1_000),
	}
}

func (i *RangeIndex) Build(id string, value interface{}) bool {
	return true
}

func (i *RangeIndex) Search(value interface{}) []uint32 {
	return nil
}

func (i *RangeIndex) GetColumnName() string {
	return i.columnName
}
