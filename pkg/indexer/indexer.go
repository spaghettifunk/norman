package indexer

import (
	"github.com/google/uuid"
	"golang.org/x/exp/constraints"
)

type ValidTypes interface {
	constraints.Float | constraints.Integer
}

type Indexer[T ValidTypes] interface {
	Build(id uuid.UUID, value T) bool
}
