package indexer

import (
	"reflect"
)

type IndexType string

const (
	TextInvertedIndex IndexType = "TEXT_INVERTED_INDEX"
	BitmapIndex       IndexType = "BITMAP_INDEX"
	RangeIndex        IndexType = "RANGE_INDEX"
	SortedIndex       IndexType = "SORTED_INDEX"
	GeospatialIndex   IndexType = "GEOSPATIAL_INDEX"
)

type Indexer interface {
	GetColumnName() string
	GetIndexType() IndexType
	AddValue(id string, value interface{}) bool
	Search(value interface{}) []uint32
	// SearchRange(from, to interface{}) []uint32
	Deserialize(data []byte) error
}

type ValidType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 | ~string
}

type IndexMetadata[T ValidType] struct {
	CastType   reflect.Kind `json:"castType"`
	IndexType  IndexType    `json:"type"`
	ColumnName string       `json:"column"`
}
