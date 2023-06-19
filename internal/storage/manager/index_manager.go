package manager

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/spaghettifunk/norman/pkg/indexer"

	bitmapindex "github.com/spaghettifunk/norman/pkg/indexer/bitmap"
	textinvertedindex "github.com/spaghettifunk/norman/pkg/indexer/inverted/text"
	rangeindex "github.com/spaghettifunk/norman/pkg/indexer/range"
	sortedindex "github.com/spaghettifunk/norman/pkg/indexer/sorted"
)

const (
	indexFileName string = "indexes.norman"
)

type IndexManager struct {
	directory string
	Indexes   map[string]indexer.Indexer `json:"indexes"`
	file      *os.File
	mu        sync.Mutex
}

func NewIndexManager(dir string) *IndexManager {
	// register gob interfaces
	gob.Register(map[string]indexer.Indexer{})

	return &IndexManager{
		directory: dir,
		Indexes:   make(map[string]indexer.Indexer, 10),
	}
}

// CreateIndex is needed to be designed this way to avoid many complications in type casting and generics
func CreateIndex[T indexer.ValidType](m *IndexManager, columnName string, indexType indexer.IndexType) error {
	if m.Indexes[columnName] != nil {
		return fmt.Errorf("index already existing for column %s", columnName)
	}

	switch indexType {
	case indexer.TextInvertedIndex:
		m.Indexes[columnName] = textinvertedindex.New[T](columnName)
	case indexer.BitmapIndex:
		m.Indexes[columnName] = bitmapindex.New[T](columnName)
	case indexer.RangeIndex:
		m.Indexes[columnName] = rangeindex.New[T](columnName)
	case indexer.SortedIndex:
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

	// marshal and compress with Brotli to save space
	buffer, err := json.Marshal(&m)
	if err != nil {
		return err
	}
	// buffer, err := utils.CompressBrotli(buf)
	// if err != nil {
	// 	return err
	// }

	// TODO: handle when file already exists. Potential solution is to create a tmp file
	// delete the current index file and then change the name of the tmp file
	if _, err = m.file.Write(buffer); err != nil {
		return err
	}
	return nil
}

func ReadIndexFile(dir string) (*IndexManager, error) {
	filePath := filepath.Join(dir, indexFileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, err
	}

	buffer, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// buffer, err := utils.DecompressBrotli(buf)
	// if err != nil {
	// 	return nil, err
	// }

	im := &IndexManager{}
	if err = json.Unmarshal(buffer, &im); err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}

	return im, err
}

type IndexManagerJSON struct {
}

// func (i *IndexManager) UnmarshalJSON(data []byte) error {
// 	bitmapIndexes := []aliasJSON[T]{}
// 	if err := json.Unmarshal(data, &bitmapIndexes); err != nil {
// 		return err
// 	}
// 	for _, bi := range bitmapIndexes {
// 		bm := bitmap.FromBytes(bi.Ids)
// 		i.Index[bi.Key] = &bm
// 	}
// 	return nil
// }
