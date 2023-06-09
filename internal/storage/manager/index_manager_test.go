package manager

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/spaghettifunk/norman/internal/storage/indexer"
	"github.com/stretchr/testify/assert"
)

func TestNewIndexManager(t *testing.T) {
	im := NewIndexManager("./output/default/transcript")
	assert.NotNil(t, im)
}

func TestCreateIndex(t *testing.T) {
	im := NewIndexManager("./output/default/transcript")
	assert.NotNil(t, im)

	tests := []struct {
		dimension string
		indexType indexer.IndexType
		expected  error
	}{
		{"dim-a", indexer.BitmapIndex, nil},
		{"dim-b", indexer.RangeIndex, nil},
		{"dim-c", indexer.SortedIndex, nil},
		{"dim-d", indexer.TextInvertedIndex, nil},
		{"dim-g", "unknown", fmt.Errorf("wrong index type unknown")},
		{"dim-d", indexer.TextInvertedIndex, fmt.Errorf("index already existing for column dim-d")},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%s,%s", tt.dimension, tt.indexType)
		t.Run(testname, func(t *testing.T) {
			err := CreateIndex[int](im, tt.dimension, tt.indexType)
			errors.Is(err, tt.expected)
		})
	}
}

func TestAddValue(t *testing.T) {
	im := NewIndexManager("./output/default/transcript")
	assert.NotNil(t, im)

	err := CreateIndex[string](im, "dim-a", indexer.TextInvertedIndex)
	assert.Nil(t, err)

	tests := []struct {
		column   string
		id       string
		value    interface{}
		expected bool
	}{
		{"dim-a", "id-1", "hellow world", true},
		{"dim-a", "id-2", "banana chocolate and strawberries?", true},
		{"dim-a", "id-3", "_***+++ hello??davide is it you?!111!!!!", true},
		{"dim-a", "id-4", 42, false},
		{"dim-a", "id-4", "42", true},
		{"dim-a", "id-4", "888888", true},
		{"dim-b", "id-4", "888888", false},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%s,%x", tt.id, tt.value)
		t.Run(testname, func(t *testing.T) {
			result := im.Add(tt.column, tt.id, tt.value)
			assert.Equal(t, tt.expected, result, fmt.Sprintf("failed test: %s", testname))
		})
	}
}

func TestSearchValue(t *testing.T) {
	im := NewIndexManager("./output/default/transcript")
	assert.NotNil(t, im)

	err := CreateIndex[string](im, "dim-a", indexer.TextInvertedIndex)
	assert.Nil(t, err)

	if im.Add("dim-a", "id-1", "hellow banana") == false {
		t.Errorf("failed adding event id %s", "id-1")
	}

	if im.Add("dim-a", "id-2", "banana chocolate and strawberries?") == false {
		t.Errorf("failed adding event id %s", "id-2")
	}

	tests := []struct {
		column string
		value  interface{}
		length int
	}{
		{"dim-a", "banana", 2},
		{"dim-a", "unknown", 0},
		{"dim-a", 42, 0},
		{"dim-b", "42", 0},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%s,%x", tt.column, tt.value)
		t.Run(testname, func(t *testing.T) {
			result := im.QueryIndex(tt.column, tt.value)
			assert.Equal(t, tt.length, len(result), fmt.Sprintf("failed test: %s", testname))
		})
	}
}

func TestPersistOnDisk(t *testing.T) {
	im := NewIndexManager("./output/default/transcript")
	assert.NotNil(t, im)

	err := CreateIndex[string](im, "dim-a", indexer.TextInvertedIndex)
	assert.Nil(t, err)

	if im.Add("dim-a", "id-1", "hellow banana") == false {
		t.Errorf("failed adding event id %s", "id-1")
	}

	if im.Add("dim-a", "id-2", "banana chocolate and strawberries?") == false {
		t.Errorf("failed adding event id %s", "id-2")
	}

	segmentID := uuid.New().String()
	partitionEnd := time.Now().Format(partitionTimeFormat)
	partitionStart := time.Now().Add(time.Duration(5) * time.Minute).Format(partitionTimeFormat)
	err = im.PersistToDisk(segmentID, partitionStart, partitionEnd)
	assert.Nil(t, err)

	deleteOutputFolder(t)
}

func TestReadIndexFile(t *testing.T) {
	// t.SkipNow()

	im := NewIndexManager("./output/default/transcript")
	assert.NotNil(t, im)

	err := CreateIndex[string](im, "dim-a", indexer.TextInvertedIndex)
	assert.Nil(t, err)
	err = CreateIndex[int](im, "dim-b", indexer.RangeIndex)
	assert.Nil(t, err)
	err = CreateIndex[int](im, "dim-c", indexer.SortedIndex)
	assert.Nil(t, err)
	err = CreateIndex[float32](im, "dim-d", indexer.BitmapIndex)
	assert.Nil(t, err)

	if im.Add("dim-a", "id-1", "hellow banana") == false {
		t.Errorf("failed adding event id %s", "id-1")
	}

	if im.Add("dim-a", "id-2", "banana chocolate and strawberries?") == false {
		t.Errorf("failed adding event id %s", "id-2")
	}

	if im.Add("dim-a", "id-3", "strawberries are the best, right?") == false {
		t.Errorf("failed adding event id %s", "id-3")
	}

	if im.Add("dim-b", "id-1", 42) == false {
		t.Errorf("failed adding event id %s", "id-1")
	}

	if im.Add("dim-b", "id-2", 125) == false {
		t.Errorf("failed adding event id %s", "id-2")
	}

	if im.Add("dim-b", "id-3", 25) == false {
		t.Errorf("failed adding event id %s", "id-3")
	}

	segmentID := uuid.New().String()
	partitionEnd := time.Now().Format(partitionTimeFormat)
	partitionStart := time.Now().Add(time.Duration(5) * time.Minute).Format(partitionTimeFormat)
	err = im.PersistToDisk(segmentID, partitionStart, partitionEnd)
	assert.Nil(t, err)

	mng, err := ReadIndexFile("./output/default/transcript")
	assert.NotNil(t, mng)
	assert.Nil(t, err)

	deleteOutputFolder(t)
}
