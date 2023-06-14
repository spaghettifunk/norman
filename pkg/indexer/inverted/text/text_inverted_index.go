package textinvertedindex

import (
	"encoding/binary"

	"github.com/RoaringBitmap/roaring"
	"github.com/google/uuid"
	"github.com/spaghettifunk/norman/pkg/containers/mapset"
)

type TextInvertedIndex struct {
	index     map[string]*roaring.Bitmap
	stopWords mapset.Set[string]
}

// NewTextInvertedIndex creates a new InvertedIndex object
// English is the only language supported
func NewTextInvertedIndex() *TextInvertedIndex {
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
func (i *TextInvertedIndex) Build(id uuid.UUID, document string) bool {
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

		rb.Add(binary.BigEndian.Uint32(id[:]))
		visited[word] = true
	}
	return true
}

// Search queries the index for the given text.
func (i *TextInvertedIndex) Search(text string) []uint32 {
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
