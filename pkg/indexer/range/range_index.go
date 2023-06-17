package rangeindex

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/rs/zerolog/log"
	"github.com/spaghettifunk/norman/pkg/indexer"
)

type RangeIndex[T indexer.ValidType] struct {
	columnName string
	index      map[T]*roaring.Bitmap
}

func New[T indexer.ValidType](columnName string) *RangeIndex[T] {
	return &RangeIndex[T]{
		columnName: columnName,
		index:      make(map[T]*roaring.Bitmap, 1_000),
	}
}

func (i *RangeIndex[T]) AddValue(id string, value interface{}) bool {
	_, ok := value.(T)
	if !ok {
		log.Error().Msg("value cannot be casted to ValidType")
		return false
	}
	return true
}

func (i *RangeIndex[T]) Search(value interface{}) []uint32 {
	_, ok := value.(T)
	if !ok {
		log.Error().Msg("value cannot be casted to ValidType")
		return nil
	}
	return nil
}

func (i *RangeIndex[T]) GetColumnName() string {
	return i.columnName
}
