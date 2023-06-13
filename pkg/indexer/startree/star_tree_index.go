package startreeindex

import "github.com/rs/zerolog/log"

type StarTreeNode struct {
	Level           int
	Dimension       string
	DimensionValue  float64
	AggregatedValue float64
	Children        []*StarTreeNode
}

// NewStarTreeIndex initialize the root node
func NewStarTreeIndex() *StarTreeNode {
	return &StarTreeNode{
		Level:           0,
		Dimension:       "root",
		AggregatedValue: 0.0,
		Children:        []*StarTreeNode{},
	}
}

func (n *StarTreeNode) ProcessEvent(node *StarTreeNode, event map[string]interface{}, dimensions []string, level int) {
	if level == len(dimensions)-1 {
		// Leaf level, update the aggregated value
		if val, ok := event[dimensions[level]].(float64); ok {
			// TODO: change this with the type of aggreagtion we want to support
			node.AggregatedValue += val
		}
		return
	}

	currentDimension := dimensions[level+1]
	dimensionValue, ok := event[currentDimension].(float64)
	if !ok {
		log.Error().Msgf("couldn't retried value for dimension %s in event", currentDimension)
		return
	}

	// Find the child node corresponding to the dimension value
	var child *StarTreeNode
	for _, c := range node.Children {
		if c.Dimension == currentDimension && c.DimensionValue == dimensionValue {
			child = c
			break
		}
	}

	// If the child node doesn't exist, create it
	if child == nil {
		child = &StarTreeNode{
			Level:           level + 1,
			Dimension:       currentDimension,
			DimensionValue:  dimensionValue,
			AggregatedValue: 0.0,
			Children:        []*StarTreeNode{},
		}
		node.Children = append(node.Children, child)
	}

	// Recursively process the child node
	n.ProcessEvent(child, event, dimensions, level+1)
}
