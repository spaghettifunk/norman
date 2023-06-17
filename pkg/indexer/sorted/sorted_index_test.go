package sortedindex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSortedIndex(t *testing.T) {
	ri := New[int]("dimension-a")
	assert.NotNil(t, ri)
}
