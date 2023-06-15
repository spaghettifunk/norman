package textinvertedindex

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewInvertedIndex(t *testing.T) {
	ii := New()
	assert.NotNil(t, ii)
}

func TestBuildInvertedIndex(t *testing.T) {
	documents := make(map[uuid.UUID]string)
	documents[uuid.New()] = "This is the first document."
	documents[uuid.New()] = "This document is the second document."
	documents[uuid.New()] = "And this is the third one."
	documents[uuid.New()] = "Is this the first document?"

	ii := New()
	assert.NotNil(t, ii)

	for id, doc := range documents {
		res := ii.Build(id, doc)
		assert.Equal(t, true, res)
	}
}

func TestRetrieveIndex(t *testing.T) {
	documents := make(map[uuid.UUID]string)
	documents[uuid.New()] = "This is the first document."
	documents[uuid.New()] = "This document is the second document."
	documents[uuid.New()] = "And this is the third one."
	documents[uuid.New()] = "Is this the first document?"

	ii := New()
	assert.NotNil(t, ii)

	for id, doc := range documents {
		res := ii.Build(id, doc)
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
