package manager

import (
	"github.com/spaghettifunk/norman/pkg/indexer"
	bitmapindex "github.com/spaghettifunk/norman/pkg/indexer/bitmap"
	textinvertedindex "github.com/spaghettifunk/norman/pkg/indexer/inverted/text"
	rangeindex "github.com/spaghettifunk/norman/pkg/indexer/range"
	sortedindex "github.com/spaghettifunk/norman/pkg/indexer/sorted"
	startreeindex "github.com/spaghettifunk/norman/pkg/indexer/startree"
)

type IndexManager[T indexer.ValidTypes] struct {
	TextInvertedIndexes map[string]*textinvertedindex.TextInvertedIndex
	BitmapIndexes       map[string]*bitmapindex.BitmapIndex[T]
	RangeIndexes        map[string]*rangeindex.RangeIndex[T]
	SortedIndexes       map[string]*sortedindex.SortedIndex[T]
	StarTreeIndexes     map[string]*startreeindex.StarTreeNode[T]
}

func NewIndexerManager[T indexer.ValidTypes]() *IndexManager[T] {
	return nil
}
