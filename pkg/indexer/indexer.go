package indexer

import "golang.org/x/exp/constraints"

type Indexer interface {
	GetColumnName() string
	AddValue(id string, value interface{}) bool
	Search(value interface{}) []uint32
}

type ValidType interface {
	constraints.Float | constraints.Integer | string
}
