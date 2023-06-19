package sortedindex

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/spaghettifunk/norman/pkg/indexer"
)

type SortedIndex[T indexer.ValidType] struct {
	ColumnName string `json:"columnName"`
}

func New[T indexer.ValidType](columnName string) *SortedIndex[T] {
	return &SortedIndex[T]{
		ColumnName: columnName,
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
	return i.ColumnName
}

func (i *SortedIndex[T]) GetIndexType() indexer.IndexType {
	return indexer.SortedIndex
}

type SortedIndexJSON[T indexer.ValidType] struct {
	IndexType  indexer.IndexType `json:"type"`
	ColumnName string            `json:"column"`
}

func (i *SortedIndex[T]) MarshalJSON() ([]byte, error) {
	si := &SortedIndexJSON[T]{IndexType: indexer.SortedIndex, ColumnName: i.ColumnName}
	return json.Marshal(si)
}

func (i *SortedIndex[T]) UnmarshalJSON(data []byte) error {
	si := &SortedIndexJSON[T]{}
	if err := json.Unmarshal(data, si); err != nil {
		return err
	}
	i.ColumnName = si.ColumnName
	return nil
}
