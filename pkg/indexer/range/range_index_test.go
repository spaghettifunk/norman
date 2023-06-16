package rangeindex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRangeIndex(t *testing.T) {
	ri := New("dimension-a")
	assert.NotNil(t, ri)
}
