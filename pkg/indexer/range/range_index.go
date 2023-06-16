package rangeindex

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/google/uuid"
	"github.com/spaghettifunk/norman/pkg/indexer"
)

type RangeIndex[T indexer.ValidTypes] struct {
	index map[T]*roaring.Bitmap
}

func New[T indexer.ValidTypes]() *RangeIndex[T] {
	return &RangeIndex[T]{
		index: make(map[T]*roaring.Bitmap, 1_000),
	}
}

func (i *RangeIndex[T]) Build(id uuid.UUID, value T) bool {
	return true
}
