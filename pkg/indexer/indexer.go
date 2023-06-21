package indexer

import (
	"reflect"

	"golang.org/x/exp/constraints"
)

type IndexType string

const (
	TextInvertedIndex IndexType = "TEXT_INVERTED_INDEX"
	BitmapIndex       IndexType = "BITMAP_INDEX"
	RangeIndex        IndexType = "RANGE_INDEX"
	SortedIndex       IndexType = "SORTED_INDEX"
)

type Indexer interface {
	GetColumnName() string
	GetIndexType() IndexType
	AddValue(id string, value interface{}) bool
	Search(value interface{}) []uint32
	Deserialize(data []byte) error
}

type ValidType interface {
	constraints.Float | constraints.Integer | string
}

type IndexMetadata[T ValidType] struct {
	CastType   reflect.Kind `json:"castType"`
	IndexType  IndexType    `json:"type"`
	ColumnName string       `json:"column"`
}
