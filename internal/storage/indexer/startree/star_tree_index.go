package startreeindex

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

type AggregationType string

type ValidTypes interface {
	constraints.Float | constraints.Integer
}

const (
	Sum     AggregationType = "SUM"
	Average AggregationType = "AVERAGE"
)

type StarTreeNode[T ValidTypes] struct {
	Level           int
	Dimension       string
	DimensionValue  T
	AggregatedValue T
	AggregationType AggregationType
	Children        []*StarTreeNode[T]
}

// New initialize the root node
func New[T ValidTypes]() *StarTreeNode[T] {
	return &StarTreeNode[T]{
		Level:           0,
		Dimension:       "root",
		DimensionValue:  0,
		AggregatedValue: 0,
		Children:        []*StarTreeNode[T]{},
	}
}

func (n *StarTreeNode[T]) ProcessEvent(node *StarTreeNode[T], event map[string]interface{}, dimensions []string, level int) error {
	if level == len(dimensions)-1 {
		// Leaf level, update the aggregated value
		if val, ok := event[dimensions[level]].(T); ok {
			// TODO: change this with the type of aggreagtion we want to support
			node.AggregatedValue += val
		}
		return nil
	}

	currentDimension := dimensions[level+1]
	dimensionValue, ok := event[currentDimension].(T)
	if !ok {
		return fmt.Errorf("couldn't retried value for dimension %s in event", currentDimension)
	}

	// Find the child node corresponding to the dimension value
	var child *StarTreeNode[T]
	for _, c := range node.Children {
		if c.Dimension == currentDimension && c.DimensionValue == dimensionValue {
			child = c
			break
		}
	}

	// If the child node doesn't exist, create it
	if child == nil {
		child = &StarTreeNode[T]{
			Level:           level + 1,
			Dimension:       currentDimension,
			DimensionValue:  dimensionValue,
			AggregatedValue: 0.0,
			Children:        []*StarTreeNode[T]{},
		}
		node.Children = append(node.Children, child)
	}

	// Recursively process the child node
	return n.ProcessEvent(child, event, dimensions, level+1)
}
