package manager

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/spaghettifunk/norman/pkg/indexer"

	bitmapindex "github.com/spaghettifunk/norman/pkg/indexer/bitmap"
	textinvertedindex "github.com/spaghettifunk/norman/pkg/indexer/inverted/text"
	rangeindex "github.com/spaghettifunk/norman/pkg/indexer/range"
	sortedindex "github.com/spaghettifunk/norman/pkg/indexer/sorted"
)

type IndexType string

const (
	indexFileName     string    = "indexes.norman"
	TextInvertedIndex IndexType = "TEXT_INVERTED_INDEX"
	BitmapIndex       IndexType = "BITMAP_INDEX"
	RangeIndex        IndexType = "RANGE_INDEX"
	SortedIndex       IndexType = "SORTED_INDEX"
)

type IndexManager struct {
	directory string
	Indexes   map[string]indexer.Indexer
	file      *os.File
	mu        sync.Mutex
}

func NewIndexManager(dir string) *IndexManager {
	return &IndexManager{
		directory: dir,
		Indexes:   make(map[string]indexer.Indexer, 10),
	}
}

func (m *IndexManager) AddIndex(columnName string, indexType IndexType) error {
	if m.Indexes[columnName] != nil {
		return fmt.Errorf("index already existing for column %s", columnName)
	}

	switch indexType {
	case TextInvertedIndex:
		m.Indexes[columnName] = textinvertedindex.New(columnName)
	case BitmapIndex:
		m.Indexes[columnName] = bitmapindex.New(columnName)
	case RangeIndex:
		m.Indexes[columnName] = rangeindex.New(columnName)
	case SortedIndex:
		m.Indexes[columnName] = sortedindex.New(columnName)
	default:
		return fmt.Errorf("wrong index type %s", indexType)
	}
	return nil
}

func (m *IndexManager) BuildIndex(columnName string, id string, value interface{}) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.Indexes[columnName].Build(id, value)
}

func (m *IndexManager) QueryIndex(columnName string, value interface{}) []uint32 {
	return m.Indexes[columnName].Search(value)
}

func (m *IndexManager) PersistToDisk() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(m.directory, os.ModePerm); err != nil {
		return err
	}

	// Create the file if it doesn't exist
	filePath := filepath.Join(m.directory, indexFileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if m.file, err = os.Create(filePath); err != nil {
			return err
		}
		defer func() {
			err = errors.Join(err, m.file.Close())
		}()
	}

	return binary.Write(m.file, binary.BigEndian, m.Indexes)
}

func ReadIndexFile(dir string) (*IndexManager, error) {
	filePath := filepath.Join(dir, indexFileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	var im *IndexManager
	err = binary.Read(file, binary.BigEndian, im)
	return im, err
}
