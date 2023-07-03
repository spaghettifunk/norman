package sortedindex

import (
	"encoding/json"
	"reflect"

	"github.com/google/btree"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/fasthash/fnv1a"
	"github.com/spaghettifunk/norman/internal/storage/indexer"
)

const (
	BTreeDegree int = 32
)

type SortedIndex[T indexer.ValidType] struct {
	Metadata indexer.IndexMetadata[T]    `json:"metadata"`
	Index    *btree.BTreeG[indexItem[T]] `json:"index"`
}

type indexItem[T indexer.ValidType] struct {
	IDs   []uint32 `json:"ids"`
	Value T        `json:"value"`
}

func Less[T indexer.ValidType]() btree.LessFunc[indexItem[T]] {
	return func(a, b indexItem[T]) bool {		
		return a.Value < b.Value
	}
}

func New[T indexer.ValidType](columnName string) *SortedIndex[T] {
	lessFn := Less[T]()

	var t T
	return &SortedIndex[T]{
		Metadata: indexer.IndexMetadata[T]{
			CastType:   reflect.TypeOf(t).Kind(),
			IndexType:  indexer.SortedIndex,
			ColumnName: columnName,
		},
		Index: btree.NewG[indexItem[T]](BTreeDegree, lessFn),
	}
}

func (i *SortedIndex[T]) AddValue(id string, value interface{}) bool {
	val, ok := value.(T)
	if !ok {
		log.Error().Msg("value cannot be casted to ValidType")
		return false
	}

	// check if item is present in the tree
	// update the IDs in case or create new array
	it := indexItem[T]{Value: val}
	key, found := i.Index.Get(it)
	if found {
		key.IDs = append(key.IDs, fnv1a.HashString32(id))
	} else {
		it.IDs = append(it.IDs, fnv1a.HashString32(id))
	}

	if _, success := i.Index.ReplaceOrInsert(it); !success {
		log.Warn().Msgf("value already added to the index")
	}
	return true
}

func (i *SortedIndex[T]) Search(value interface{}) []uint32 {
	val, ok := value.(T)
	if !ok {
		log.Error().Msg("value cannot be casted to ValidType")
		return nil
	}

	it := indexItem[T]{Value: val}
	key, found := i.Index.Get(it)
	if !found {
		return nil
	}

	return key.IDs
}

func (i *SortedIndex[T]) GetColumnName() string {
	return i.Metadata.ColumnName
}

func (i *SortedIndex[T]) GetIndexType() indexer.IndexType {
	return i.Metadata.IndexType
}

type SortedIndexJSON[T indexer.ValidType] struct {
	Metadata indexer.IndexMetadata[T] `json:"metadata"`
	Indexes  []SortedMapIdsJSON[T]    `json:"index"`
}

type SortedMapIdsJSON[T indexer.ValidType] struct {
	Value T        `json:"value"`
	Ids   []uint32 `json:"ids"`
}

func (i *SortedIndex[T]) MarshalJSON() ([]byte, error) {
	si := &SortedIndexJSON[T]{
		Metadata: i.Metadata,
	}
	i.Index.Ascend(func(item indexItem[T]) bool {
		si.Indexes = append(si.Indexes, SortedMapIdsJSON[T]{
			Value: item.Value,
			Ids:   item.IDs,
		})
		return true
	})
	return json.Marshal(si)
}

func (i *SortedIndex[T]) Deserialize(data []byte) error {
	sortedIndexes := []SortedMapIdsJSON[T]{}
	if err := json.Unmarshal(data, &sortedIndexes); err != nil {
		return err
	}
	for _, item := range sortedIndexes {
		i.Index.ReplaceOrInsert(indexItem[T]{
			IDs:   item.Ids,
			Value: item.Value,
		})
	}
	return nil
}
