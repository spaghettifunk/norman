package geospatial

import (
	"encoding/json"
	"reflect"

	"github.com/rs/zerolog/log"
	"github.com/spaghettifunk/norman/pkg/indexer"
)

type GeospatialIndex[T indexer.ValidType] struct {
	Metadata indexer.IndexMetadata[T] `json:"metadata"`
}

func New[T indexer.ValidType](columnName string) *GeospatialIndex[T] {
	return &GeospatialIndex[T]{
		Metadata: indexer.IndexMetadata[T]{
			CastType:   reflect.Float64,
			IndexType:  indexer.GeospatialIndex,
			ColumnName: columnName,
		},
	}
}

func (i *GeospatialIndex[T]) AddValue(id string, value interface{}) bool {
	_, ok := value.(float64)
	if !ok {
		log.Error().Msg("value cannot be casted to ValidType")
		return false
	}
	return true
}

func (i *GeospatialIndex[T]) Search(value interface{}) []uint32 {
	_, ok := value.(float64)
	if !ok {
		log.Error().Msg("value cannot be casted to ValidType")
		return nil
	}
	return nil
}

func (i *GeospatialIndex[T]) GetColumnName() string {
	return i.Metadata.ColumnName
}

func (i *GeospatialIndex[T]) GetIndexType() indexer.IndexType {
	return i.Metadata.IndexType
}

type SortedIndexJSON[T indexer.ValidType] struct {
	Metadata indexer.IndexMetadata[T] `json:"metadata"`
}

func (i *GeospatialIndex[T]) MarshalJSON() ([]byte, error) {
	si := &SortedIndexJSON[T]{
		Metadata: i.Metadata,
	}
	return json.Marshal(si)
}

// TODO: implement this once there is an actual index system
func (i *GeospatialIndex[T]) Deserialize(data []byte) error {
	return nil
}
