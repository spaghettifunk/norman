package rangeindex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRangeIndex(t *testing.T) {
	ri := New[int]()
	assert.NotNil(t, ri)
}
