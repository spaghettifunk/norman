package manager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/spaghettifunk/norman/pkg/indexer"
	"github.com/vmihailenco/msgpack/v5"

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
	Indexes   map[string]indexer.Indexer `msgpack:"indexes,inline"`
	file      *os.File
	mu        sync.Mutex
}

func NewIndexManager(dir string) *IndexManager {
	return &IndexManager{
		directory: dir,
		Indexes:   make(map[string]indexer.Indexer, 10),
	}
}

// CreateIndex is needed to be designed this way to avoid many complications in type casting and generics
func CreateIndex[T indexer.ValidType](m *IndexManager, columnName string, indexType IndexType) error {
	if m.Indexes[columnName] != nil {
		return fmt.Errorf("index already existing for column %s", columnName)
	}

	switch indexType {
	case TextInvertedIndex:
		m.Indexes[columnName] = textinvertedindex.New[T](columnName)
	case BitmapIndex:
		m.Indexes[columnName] = bitmapindex.New[T](columnName)
	case RangeIndex:
		m.Indexes[columnName] = rangeindex.New[T](columnName)
	case SortedIndex:
		m.Indexes[columnName] = sortedindex.New[T](columnName)
	default:
		return fmt.Errorf("wrong index type %s", indexType)
	}
	return nil
}

func (m *IndexManager) Add(columnName string, id string, value interface{}) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	idx, ok := m.Indexes[columnName]
	if !ok {
		log.Error().Msgf("no index for column %s", columnName)
		return false
	}

	return idx.AddValue(id, value)
}

func (m *IndexManager) QueryIndex(columnName string, value interface{}) []uint32 {
	idx, ok := m.Indexes[columnName]
	if !ok {
		log.Error().Msgf("no index for column %s", columnName)
		return nil
	}
	return idx.Search(value)
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

	// marshall the object using messagepack
	b, err := msgpack.Marshal(m.Indexes)
	if err != nil {
		return err
	}

	offset, err := m.file.Write(b)
	if offset <= 0 || err != nil {
		return fmt.Errorf("error in writing indexing file")
	}
	return nil
}

func ReadIndexFile(dir string) (*IndexManager, error) {
	filePath := filepath.Join(dir, indexFileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, err
	}

	buff, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	im := IndexManager{}
	if err = msgpack.Unmarshal(buff, &im.Indexes); err != nil {
		return nil, err
	}
	return &im, err
}
