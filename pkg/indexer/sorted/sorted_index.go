package sortedindex

import (
	"encoding/json"
	"reflect"

	"github.com/rs/zerolog/log"
	"github.com/spaghettifunk/norman/pkg/indexer"
)

type SortedIndex[T indexer.ValidType] struct {
	Metadata indexer.IndexMetadata[T] `json:"metadata"`
}

func New[T indexer.ValidType](columnName string) *SortedIndex[T] {
	var t T
	return &SortedIndex[T]{
		Metadata: indexer.IndexMetadata[T]{
			CastType:   reflect.TypeOf(t).Kind(),
			IndexType:  indexer.SortedIndex,
			ColumnName: columnName,
		},
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
	return i.Metadata.ColumnName
}

func (i *SortedIndex[T]) GetIndexType() indexer.IndexType {
	return i.Metadata.IndexType
}

type SortedIndexJSON[T indexer.ValidType] struct {
	Metadata indexer.IndexMetadata[T] `json:"metadata"`
}

func (i *SortedIndex[T]) MarshalJSON() ([]byte, error) {
	si := &SortedIndexJSON[T]{
		Metadata: i.Metadata,
	}
	return json.Marshal(si)
}

// TODO: implement this once there is an actual index system
func (i *SortedIndex[T]) Deserialize(data []byte) error {
	return nil
}
