package textinvertedindex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInvertedIndex(t *testing.T) {
	ii := New[string]("dimension-a")
	assert.NotNil(t, ii)
}

func TestBuildInvertedIndex(t *testing.T) {
	documents := make(map[string]string)
	documents["event-1"] = "This is the first document."
	documents["event-2"] = "This document is the second document."
	documents["event-3"] = "And this is the third one."
	documents["event-4"] = "Is this the first document?"

	ii := New[string]("dimension-a")
	assert.NotNil(t, ii)

	for id, doc := range documents {
		res := ii.AddValue(id, doc)
		assert.Equal(t, true, res)
	}
}

func TestRetrieveIndex(t *testing.T) {
	documents := make(map[string]string)
	documents["event-1"] = "This is the first document."
	documents["event-2"] = "This document is the second document."
	documents["event-3"] = "And this is the third one."
	documents["event-4"] = "Is this the first document?"

	ii := New[string]("dimension-a")
	assert.NotNil(t, ii)

	for id, doc := range documents {
		res := ii.AddValue(id, doc)
		assert.Equal(t, true, res)
	}

	idx := ii.Search("document")
	assert.Equal(t, 3, len(idx))

	idx = ii.Search("first")
	assert.Equal(t, 2, len(idx))

	idx = ii.Search("first document")
	assert.Equal(t, 2, len(idx))

	idx = ii.Search("second document")
	assert.Equal(t, 1, len(idx))

	idx = ii.Search("apple")
	assert.Nil(t, idx)
}
