package rangeindex

import (
	"encoding/json"
	"reflect"

	"github.com/RoaringBitmap/roaring"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/fasthash/fnv1a"
	"github.com/spaghettifunk/norman/internal/storage/indexer"
)

type RangeIndex[T indexer.ValidType] struct {
	Metadata indexer.IndexMetadata[T] `json:"metadata"`
	Index    map[T]*roaring.Bitmap    `json:"index"`
}

func New[T indexer.ValidType](columnName string) *RangeIndex[T] {
	var t T
	return &RangeIndex[T]{
		Metadata: indexer.IndexMetadata[T]{
			CastType:   reflect.TypeOf(t).Kind(),
			IndexType:  indexer.RangeIndex,
			ColumnName: columnName,
		},
		Index: make(map[T]*roaring.Bitmap, 1_000),
	}
}

func (i *RangeIndex[T]) AddValue(id string, value interface{}) bool {
	val, ok := value.(T)
	if !ok {
		log.Error().Msg("value cannot be casted to ValidType")
		return false
	}
	rb, ok := i.Index[val]
	if !ok {
		rb = roaring.NewBitmap()
		i.Index[val] = rb
	}

	rb.Add(fnv1a.HashString32(id))
	return true
}

func (i *RangeIndex[T]) Search(value interface{}) []uint32 {
	val, ok := value.(T)
	if !ok {
		log.Error().Msg("value cannot be casted to ValidType")
		return nil
	}

	var r *roaring.Bitmap
	if ids, ok := i.Index[val]; ok {
		if r == nil {
			r = roaring.NewBitmap()
			iterator := ids.Iterator()
			for iterator.HasNext() {
				r.Add(iterator.Next())
			}
		} else {
			// run intersection
			r = roaring.ParAnd(0, r, ids)
		}
		return r.ToArray()
	}
	// value doesn't exist.
	return nil
}

func (i *RangeIndex[T]) GetColumnName() string {
	return i.Metadata.ColumnName
}

func (i *RangeIndex[T]) GetIndexType() indexer.IndexType {
	return i.Metadata.IndexType
}

type RangeIndexJSON[T indexer.ValidType] struct {
	Metadata indexer.IndexMetadata[T] `json:"metadata"`
	Indexes  []RangeMapIdsJSON[T]     `json:"index"`
}

type RangeMapIdsJSON[T indexer.ValidType] struct {
	Key T        `json:"key"`
	Ids []uint32 `json:"ids"`
}

func (i *RangeIndex[T]) MarshalJSON() ([]byte, error) {
	rangeIndexes := RangeIndexJSON[T]{
		Metadata: i.Metadata,
	}
	for val, ids := range i.Index {
		rangeIndexes.Indexes = append(rangeIndexes.Indexes, RangeMapIdsJSON[T]{val, ids.ToArray()})
	}
	return json.Marshal(rangeIndexes)
}

func (i *RangeIndex[T]) Deserialize(data []byte) error {
	rangeIndexes := []RangeMapIdsJSON[T]{}
	if err := json.Unmarshal(data, &rangeIndexes); err != nil {
		return err
	}
	for _, ri := range rangeIndexes {
		rb := roaring.NewBitmap()
		rb.AddMany(ri.Ids)
		i.Index[ri.Key] = rb
	}
	return nil
}
