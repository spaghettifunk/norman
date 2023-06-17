package bitmapindex

import (
	"github.com/kelindar/bitmap"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/fasthash/fnv1a"
	"github.com/spaghettifunk/norman/pkg/indexer"
)

type BitmapIndex[T indexer.ValidType] struct {
	columnName string
	index      map[T]*bitmap.Bitmap
}

func New[T indexer.ValidType](columnName string) *BitmapIndex[T] {
	return &BitmapIndex[T]{
		columnName: columnName,
		index:      make(map[T]*bitmap.Bitmap, 1_000),
	}
}

func (i *BitmapIndex[T]) AddValue(id string, value interface{}) bool {
	val, ok := value.(T)
	if !ok {
		log.Error().Msg("value cannot be casted to ValidType")
		return false
	}

	// append the ID to the string list
	if _, ok := i.index[val]; !ok {
		i.index[val] = &bitmap.Bitmap{}
	}
	i.index[val].Set(fnv1a.HashString32(id))

	return true
}

func (i *BitmapIndex[T]) Search(value interface{}) []uint32 {
	return nil
}

func (i *BitmapIndex[T]) SearchRange(min, max interface{}) []uint32 {
	return nil
}

func (i *BitmapIndex[T]) GetColumnName() string {
	return i.columnName
}
