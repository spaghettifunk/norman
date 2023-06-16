package manager

import (
	"fmt"

	"github.com/spaghettifunk/norman/pkg/indexer"
	bitmapindex "github.com/spaghettifunk/norman/pkg/indexer/bitmap"
	textinvertedindex "github.com/spaghettifunk/norman/pkg/indexer/inverted/text"
	rangeindex "github.com/spaghettifunk/norman/pkg/indexer/range"
	sortedindex "github.com/spaghettifunk/norman/pkg/indexer/sorted"
	startreeindex "github.com/spaghettifunk/norman/pkg/indexer/startree"
)

type IndexType string

const (
	TextInvertedIndex IndexType = "TEXT_INVERTED_INDEX"
	BitmapIndex       IndexType = "BITMAP_INDEX"
	RangeIndex        IndexType = "RANGE_INDEX"
	SortedIndex       IndexType = "SORTED_INDEX"
	StarTreeIndex     IndexType = "STARTREE_INDEX"
)

type IndexManager[T indexer.ValidTypes] struct {
	TextInvertedIndexes map[string]*textinvertedindex.TextInvertedIndex
	BitmapIndexes       map[string]*bitmapindex.BitmapIndex[T]
	RangeIndexes        map[string]*rangeindex.RangeIndex[T]
	SortedIndexes       map[string]*sortedindex.SortedIndex[T]
	StarTreeIndexes     map[string]*startreeindex.StarTreeNode[T]
}

func NewIndexManager[T indexer.ValidTypes]() *IndexManager[T] {
	return &IndexManager[T]{}
}

func (m *IndexManager[T]) AddIndex(columnName string, indexType IndexType) error {
	if m.indexExists(columnName) {
		return fmt.Errorf("index already existing for column %s", columnName)
	}

	switch indexType {
	case TextInvertedIndex:
		m.TextInvertedIndexes[columnName] = textinvertedindex.New()
	case BitmapIndex:
		m.BitmapIndexes[columnName] = bitmapindex.New[T]()
	case RangeIndex:
		m.RangeIndexes[columnName] = rangeindex.New[T]()
	case SortedIndex:
		m.SortedIndexes[columnName] = sortedindex.New[T]()
	case StarTreeIndex:
		m.StarTreeIndexes[columnName] = startreeindex.New[T]()
	}
	return nil
}

func (m *IndexManager[T]) indexExists(columnName string) bool {
	return m.TextInvertedIndexes[columnName] != nil ||
		m.BitmapIndexes[columnName] != nil ||
		m.RangeIndexes[columnName] != nil ||
		m.SortedIndexes[columnName] != nil ||
		m.StarTreeIndexes[columnName] != nil
}
