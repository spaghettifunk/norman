package bitmapindex

import (
	"github.com/google/uuid"
	"github.com/kelindar/bitmap"
	"github.com/spaghettifunk/norman/pkg/indexer"
)

type BitmapIndex[T indexer.ValidTypes] struct {
	index map[T]*bitmap.Bitmap
}

func New[T indexer.ValidTypes]() *BitmapIndex[T] {
	return &BitmapIndex[T]{
		index: make(map[T]*bitmap.Bitmap, 1_000),
	}
}

func (i *BitmapIndex[T]) Build(id uuid.UUID, value T) bool {
	// append the ID to the string list
	if _, ok := i.index[value]; !ok {
		i.index[value] = &bitmap.Bitmap{}
	}
	i.index[value].Set(id.ID())

	return true
}

func (i *BitmapIndex[T]) Search(value T) []uint32 {
	return nil
}

func (i *BitmapIndex[T]) SearchRange(min, max T) []uint32 {
	return nil
}
