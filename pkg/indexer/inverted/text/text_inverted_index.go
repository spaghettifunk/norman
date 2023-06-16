package textinvertedindex

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/segmentio/fasthash/fnv1a"
	"github.com/spaghettifunk/norman/pkg/containers/mapset"
)

type TextInvertedIndex struct {
	columnName string
	index      map[string]*roaring.Bitmap
	stopWords  mapset.Set[string]
}

// New creates a new Text InvertedIndex object
// English is the only language supported
func New(columnName string) *TextInvertedIndex {
	stopWords := mapset.New[string]()
	for _, sw := range stopWordsEN {
		stopWords.Put(sw)
	}

	return &TextInvertedIndex{
		index:     make(map[string]*roaring.Bitmap, 1_000),
		stopWords: stopWords,
	}
}

// Build builds the inverted index
// id is a UUID as string
func (i *TextInvertedIndex) Build(id string, document interface{}) bool {
	tokens := i.analyze(document.(string))
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
func (i *TextInvertedIndex) Search(value interface{}) []uint32 {
	text := value.(string)
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

func (i *TextInvertedIndex) GetColumnName() string {
	return i.columnName
}
