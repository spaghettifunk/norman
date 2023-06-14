package rangeindex

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/google/uuid"
)

type ValidTypes interface {
	~float32 | ~float64 | uint | uint16 | uint32 | uint64 | int | int16 | int32 | int64
}

type RangeIndex[T ValidTypes] struct {
	index map[T]*roaring.Bitmap
}

func NewRangeIndex[T ValidTypes]() *RangeIndex[T] {
	return &RangeIndex[T]{
		index: make(map[T]*roaring.Bitmap, 1_000),
	}
}

func (i *RangeIndex[T]) Build(id uuid.UUID, document T) bool {
	return true
}
