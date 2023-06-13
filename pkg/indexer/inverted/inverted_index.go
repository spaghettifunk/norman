package invertedindex

import (
	"encoding/binary"

	"github.com/RoaringBitmap/roaring"
	"github.com/google/uuid"
	"github.com/spaghettifunk/norman/pkg/containers/mapset"
)

type InvertedIndex struct {
	index     map[string]*roaring.Bitmap
	stopWords mapset.Set[string]
}

// NewInvertedIndex creates a new InvertedIndex object
// English is the only language supported
func NewInvertedIndex() *InvertedIndex {
	stopWords := mapset.New[string]()
	for _, sw := range stopWordsEN {
		stopWords.Put(sw)
	}

	return &InvertedIndex{
		index:     make(map[string]*roaring.Bitmap, 1_000),
		stopWords: stopWords,
	}
}

// Build builds the inverted index
// id is a UUID as string
func (i *InvertedIndex) Build(id uuid.UUID, document string) bool {
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

// intersection returns the set intersection between a and b.
// a and b have to be sorted in ascending order and contain no duplicates.
func (i *InvertedIndex) intersection(a []uint32, b []uint32) []uint32 {
	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}
	r := make([]uint32, 0, maxLen)
	var k, j int
	for k < len(a) && j < len(b) {
		if a[k] < b[j] {
			k++
		} else if a[k] > b[j] {
			j++
		} else {
			r = append(r, a[k])
			k++
			j++
		}
	}
	return r
}

// search queries the index for the given text.
func (i *InvertedIndex) Search(text string) []uint32 {
	var r []uint32
	for _, token := range i.analyze(text) {
		if ids, ok := i.index[token]; ok {
			if r == nil {
				iterator := ids.Iterator()
				for iterator.HasNext() {
					r = append(r, iterator.Next())
				}
			} else {
				r = i.intersection(r, ids.ToArray())
			}
		} else {
			// Token doesn't exist.
			return nil
		}
	}
	return r
}
