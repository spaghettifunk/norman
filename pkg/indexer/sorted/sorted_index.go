package sortedindex

import (
	"github.com/google/uuid"
	"github.com/spaghettifunk/norman/pkg/indexer"
)

type SortedIndex[T indexer.ValidTypes] struct {
}

func New[T indexer.ValidTypes]() *SortedIndex[T] {
	return &SortedIndex[T]{}
}

func (i *SortedIndex[T]) Build(id uuid.UUID, value T) bool {
	return true
}
