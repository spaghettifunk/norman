package manager

import (
	startreeindex "github.com/spaghettifunk/norman/internal/storage/indexer/startree"
)

type StarTreeIndexManager[T startreeindex.ValidTypes] struct {
	directory string
	index     *startreeindex.StarTreeNode[T]
}

func NewStarTreeIndexManager[T startreeindex.ValidTypes](dir string) *StarTreeIndexManager[T] {
	return &StarTreeIndexManager[T]{
		directory: dir,
		index:     startreeindex.New[T](),
	}
}

func (s *StarTreeIndexManager[T]) ProcessEvent(event map[string]interface{}, dimensions []string) error {
	return s.index.ProcessEvent(s.index, event, dimensions, 0)
}
