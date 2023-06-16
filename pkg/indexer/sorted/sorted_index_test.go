package sortedindex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSortedIndex(t *testing.T) {
	ri := New("dimension-a")
	assert.NotNil(t, ri)
}
