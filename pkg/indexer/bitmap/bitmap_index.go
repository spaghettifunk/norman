package bitmapindex

import (
	"encoding/json"

	"github.com/kelindar/bitmap"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/fasthash/fnv1a"
	"github.com/spaghettifunk/norman/pkg/indexer"
)

type BitmapIndex[T indexer.ValidType] struct {
	ColumnName string
	Index      map[T]*bitmap.Bitmap
}

func New[T indexer.ValidType](columnName string) *BitmapIndex[T] {
	return &BitmapIndex[T]{
		ColumnName: columnName,
		Index:      make(map[T]*bitmap.Bitmap, 1_000),
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
	return i.ColumnName
}

func (i *BitmapIndex[T]) GetIndexType() indexer.IndexType {
	return indexer.BitmapIndex
}

type BitmapIndexJSON[T indexer.ValidType] struct {
	IndexType  indexer.IndexType     `json:"type"`
	ColumnName string                `json:"column"`
	Indexes    []BitmapMapIdsJSON[T] `json:"index"`
}

type BitmapMapIdsJSON[T indexer.ValidType] struct {
	Key T      `json:"key"`
	Ids []byte `json:"ids"`
}

func (i *BitmapIndex[T]) MarshalJSON() ([]byte, error) {
	bitmapIndexes := BitmapIndexJSON[T]{IndexType: i.GetIndexType(), ColumnName: i.GetColumnName()}
	for val, ids := range i.Index {
		bitmapIndexes.Indexes = append(bitmapIndexes.Indexes, BitmapMapIdsJSON[T]{val, ids.ToBytes()})
	}
	return json.Marshal(bitmapIndexes)
}

func (i *BitmapIndex[T]) UnmarshalJSON(data []byte) error {
	bitmapIndexes := BitmapIndexJSON[T]{}
	if err := json.Unmarshal(data, &bitmapIndexes); err != nil {
		return err
	}
	i.ColumnName = bitmapIndexes.ColumnName
	for _, bi := range bitmapIndexes.Indexes {
		bm := bitmap.FromBytes(bi.Ids)
		i.Index[bi.Key] = &bm
	}
	return nil
}
