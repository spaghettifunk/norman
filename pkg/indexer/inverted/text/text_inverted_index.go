package textinvertedindex

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/fasthash/fnv1a"
	"github.com/spaghettifunk/norman/pkg/containers/mapset"
	"github.com/spaghettifunk/norman/pkg/indexer"
)

type TextInvertedIndex[T indexer.ValidType] struct {
	columnName string
	index      map[string]*roaring.Bitmap
	stopWords  mapset.Set[string]
}

// New creates a new Text InvertedIndex object
// English is the only language supported
func New[T indexer.ValidType](columnName string) *TextInvertedIndex[T] {
	stopWords := mapset.New[string]()
	for _, sw := range stopWordsEN {
		stopWords.Put(sw)
	}

	return &TextInvertedIndex[T]{
		columnName: columnName,
		index:      make(map[string]*roaring.Bitmap, 1_000),
		stopWords:  stopWords,
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

		rb, ok := i.index[word]
		if !ok {
			rb = roaring.NewBitmap()
			i.index[word] = rb
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
		if ids, ok := i.index[token]; ok {
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
	return i.columnName
}
