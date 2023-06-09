package textinvertedindex

import (
	"encoding/json"
	"reflect"

	"github.com/RoaringBitmap/roaring"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/fasthash/fnv1a"
	"github.com/spaghettifunk/norman/internal/storage/indexer"
	"github.com/spaghettifunk/norman/pkg/containers/mapset"
)

type TextInvertedIndex[T indexer.ValidType] struct {
	Metadata  indexer.IndexMetadata[T]   `json:"metadata"`
	Index     map[string]*roaring.Bitmap `json:"index"`
	StopWords mapset.Set[string]         `json:"-"`
}

// New creates a new Text InvertedIndex object
// English is the only language supported
func New[T indexer.ValidType](columnName string) *TextInvertedIndex[T] {
	stopWords := mapset.New[string]()
	for _, sw := range stopWordsEN {
		stopWords.Put(sw)
	}
	var t T
	return &TextInvertedIndex[T]{
		Metadata: indexer.IndexMetadata[T]{
			CastType:   reflect.TypeOf(t).Kind(),
			IndexType:  indexer.TextInvertedIndex,
			ColumnName: columnName,
		},
		Index:     make(map[string]*roaring.Bitmap, 1_000),
		StopWords: stopWords,
	}
}

// AddValue adds the current value for the given id to the index
func (i *TextInvertedIndex[T]) AddValue(id string, value interface{}) bool {
	// make sure we operate only with strings
	document, ok := value.(string)
	if !ok {
		log.Error().Msg("value cannot be casted to string")
		return false
	}

	tokens := i.analyze(document)
	visited := make(map[string]bool, len(tokens))

	for _, word := range tokens {
		if _, ok := visited[word]; ok {
			continue
		}

		rb, ok := i.Index[word]
		if !ok {
			rb = roaring.NewBitmap()
			i.Index[word] = rb
		}

		rb.Add(fnv1a.HashString32(id))
		visited[word] = true
	}
	return true
}

// Search queries the index for the given text.
func (i *TextInvertedIndex[T]) Search(value interface{}) []uint32 {
	// make sure we operate only with strings
	text, ok := value.(string)
	if !ok {
		log.Error().Msgf("value %x is not accepted", value)
		return nil
	}
	var r *roaring.Bitmap
	for _, token := range i.analyze(text) {
		if ids, ok := i.Index[token]; ok {
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
		} else {
			// Token doesn't exist.
			return nil
		}
	}
	return r.ToArray()
}

func (i *TextInvertedIndex[T]) GetColumnName() string {
	return i.Metadata.ColumnName
}

func (i *TextInvertedIndex[T]) GetIndexType() indexer.IndexType {
	return i.Metadata.IndexType
}

type TextInvertedIndexJSON[T indexer.ValidType] struct {
	Metadata indexer.IndexMetadata[T] `json:"metadata"`
	Indexes  []TextMapIdsJSON         `json:"index"`
}

type TextMapIdsJSON struct {
	Key string   `json:"key"`
	Ids []uint32 `json:"ids"`
}

func (i *TextInvertedIndex[T]) MarshalJSON() ([]byte, error) {
	textIndexes := TextInvertedIndexJSON[T]{
		Metadata: i.Metadata,
	}
	for val, ids := range i.Index {
		textIndexes.Indexes = append(textIndexes.Indexes, TextMapIdsJSON{val, ids.ToArray()})
	}
	return json.Marshal(textIndexes)
}

func (i *TextInvertedIndex[T]) Deserialize(data []byte) error {
	textIndexes := []TextMapIdsJSON{}
	if err := json.Unmarshal(data, &textIndexes); err != nil {
		return err
	}
	for _, bi := range textIndexes {
		rb := roaring.NewBitmap()
		rb.AddMany(bi.Ids)
		i.Index[bi.Key] = rb
	}
	return nil
}
