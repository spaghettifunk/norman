package indexer

import "golang.org/x/exp/constraints"

type ValidTypes interface {
	constraints.Float | constraints.Integer
}
