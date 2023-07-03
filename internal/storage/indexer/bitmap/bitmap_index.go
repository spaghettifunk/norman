package bitmapindex

import (
	"encoding/json"
	"reflect"

	"github.com/kelindar/bitmap"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/fasthash/fnv1a"
	"github.com/spaghettifunk/norman/internal/storage/indexer"
)

type BitmapIndex[T indexer.ValidType] struct {
	Metadata indexer.IndexMetadata[T] `json:"metadata"`
	Index    map[T]*bitmap.Bitmap     `json:"index"`
}

func New[T indexer.ValidType](columnName string) *BitmapIndex[T] {
	var t T
	return &BitmapIndex[T]{
		Metadata: indexer.IndexMetadata[T]{
			CastType:   reflect.TypeOf(t).Kind(),
			IndexType:  indexer.BitmapIndex,
			ColumnName: columnName,
		},
		Index: make(map[T]*bitmap.Bitmap, 1_000),
	}
}

func (i *BitmapIndex[T]) AddValue(id string, value interface{}) bool {
	val, ok := value.(T)
	if !ok {
		log.Error().Msg("value cannot be casted to ValidType")
		return false
	}

	// append the ID to the string list
	if _, ok := i.Index[val]; !ok {
		i.Index[val] = &bitmap.Bitmap{}
	}
	i.Index[val].Set(fnv1a.HashString32(id))

	return true
}

func (i *BitmapIndex[T]) Search(value interface{}) []uint32 {
	return nil
}

func (i *BitmapIndex[T]) SearchRange(min, max interface{}) []uint32 {
	return nil
}

func (i *BitmapIndex[T]) GetColumnName() string {
	return i.Metadata.ColumnName
}

func (i *BitmapIndex[T]) GetIndexType() indexer.IndexType {
	return i.Metadata.IndexType
}

type BitmapIndexJSON[T indexer.ValidType] struct {
	Metadata indexer.IndexMetadata[T] `json:"metadata"`
	Indexes  []BitmapMapIdsJSON[T]    `json:"index"`
}

type BitmapMapIdsJSON[T indexer.ValidType] struct {
	Key T      `json:"key"`
	Ids []byte `json:"ids"`
}

func (i *BitmapIndex[T]) MarshalJSON() ([]byte, error) {
	bitmapIndexes := BitmapIndexJSON[T]{
		Metadata: i.Metadata,
	}
	for val, ids := range i.Index {
		bitmapIndexes.Indexes = append(bitmapIndexes.Indexes, BitmapMapIdsJSON[T]{val, ids.ToBytes()})
	}
	return json.Marshal(bitmapIndexes)
}

func (i *BitmapIndex[T]) Deserialize(data []byte) error {
	bitmapIndexes := []BitmapMapIdsJSON[T]{}
	if err := json.Unmarshal(data, &bitmapIndexes); err != nil {
		return err
	}
	for _, bi := range bitmapIndexes {
		bm := bitmap.FromBytes(bi.Ids)
		i.Index[bi.Key] = &bm
	}
	return nil
}
