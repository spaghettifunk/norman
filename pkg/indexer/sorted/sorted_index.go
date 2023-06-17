package sortedindex

import (
	"github.com/rs/zerolog/log"
	"github.com/spaghettifunk/norman/pkg/indexer"
)

type SortedIndex[T indexer.ValidType] struct {
	columnName string
}

func New[T indexer.ValidType](columnName string) *SortedIndex[T] {
	return &SortedIndex[T]{
		columnName: columnName,
	}
}

func (i *SortedIndex[T]) AddValue(id string, value interface{}) bool {
	_, ok := value.(T)
	if !ok {
		log.Error().Msg("value cannot be casted to ValidType")
		return false
	}
	return true
}

func (i *SortedIndex[T]) Search(value interface{}) []uint32 {
	_, ok := value.(T)
	if !ok {
		log.Error().Msg("value cannot be casted to ValidType")
		return nil
	}
	return nil
}

func (i *SortedIndex[T]) GetColumnName() string {
	return i.columnName
}
