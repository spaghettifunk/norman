package bitmapindex

import (
	"github.com/kelindar/bitmap"
	"github.com/segmentio/fasthash/fnv1a"
)

type BitmapIndex struct {
	columnName string
	index      map[interface{}]*bitmap.Bitmap
}

func New(columnName string) *BitmapIndex {
	return &BitmapIndex{
		columnName: columnName,
		index:      make(map[interface{}]*bitmap.Bitmap, 1_000),
	}
}

func (i *BitmapIndex) Build(id string, value interface{}) bool {
	// append the ID to the string list
	if _, ok := i.index[value]; !ok {
		i.index[value] = &bitmap.Bitmap{}
	}
	i.index[value].Set(fnv1a.HashString32(id))

	return true
}

func (i *BitmapIndex) Search(value interface{}) []uint32 {
	return nil
}

func (i *BitmapIndex) SearchRange(min, max interface{}) []uint32 {
	return nil
}

func (i *BitmapIndex) GetColumnName() string {
	return i.columnName
}
