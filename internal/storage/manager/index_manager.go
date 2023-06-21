package manager

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
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
	internal  *internalData
	file      *os.File
	mu        sync.Mutex
}

type internalData struct {
	SegmentID      string                     `json:"segmentId"`
	PartitionStart string                     `json:"partitionStart"`
	PartitionEnd   string                     `json:"partitionEnd"`
	Indexes        map[string]indexer.Indexer `json:"indexes"`
	Metadata       map[string]interface{}     `json:"metadata"`
}

func NewIndexManager(dir string) *IndexManager {
	// register gob interfaces
	gob.Register(map[string]indexer.Indexer{})

	return &IndexManager{
		directory: dir,
		internal: &internalData{
			Indexes: make(map[string]indexer.Indexer, 10),
		},
	}
}

// CreateIndex is needed to be designed this way to avoid many complications in type casting and generics
func CreateIndex[T indexer.ValidType](m *IndexManager, columnName string, indexType indexer.IndexType) error {
	if m.internal.Indexes[columnName] != nil {
		return fmt.Errorf("index already existing for column %s", columnName)
	}

	switch indexType {
	case indexer.TextInvertedIndex:
		m.internal.Indexes[columnName] = textinvertedindex.New[T](columnName)
	case indexer.BitmapIndex:
		m.internal.Indexes[columnName] = bitmapindex.New[T](columnName)
	case indexer.RangeIndex:
		m.internal.Indexes[columnName] = rangeindex.New[T](columnName)
	case indexer.SortedIndex:
		m.internal.Indexes[columnName] = sortedindex.New[T](columnName)
	default:
		return fmt.Errorf("wrong index type %s", indexType)
	}
	return nil
}

func (m *IndexManager) Add(columnName string, id string, value interface{}) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	idx, ok := m.internal.Indexes[columnName]
	if !ok {
		log.Error().Msgf("no index for column %s", columnName)
		return false
	}

	return idx.AddValue(id, value)
}

func (m *IndexManager) QueryIndex(columnName string, value interface{}) []uint32 {
	idx, ok := m.internal.Indexes[columnName]
	if !ok {
		log.Error().Msgf("no index for column %s", columnName)
		return nil
	}
	return idx.Search(value)
}

func (m *IndexManager) PersistToDisk(segmentID, partitionStart, partitionEnd string) error {
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

	// store extra infos before marshaling
	m.internal.SegmentID = segmentID
	m.internal.PartitionStart = partitionStart
	m.internal.PartitionEnd = partitionEnd

	// marshal and compress with Brotli to save space
	buffer, err := json.Marshal(&m.internal)
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

	im := &IndexManager{
		directory: dir,
		mu:        sync.Mutex{},
	}
	if err = json.Unmarshal(buffer, im); err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}

	return im, err
}

type internalDataJSON struct {
	SegmentID      string                 `json:"segmentId"`
	PartitionStart string                 `json:"partitionStart"`
	PartitionEnd   string                 `json:"partitionEnd"`
	Metadata       map[string]interface{} `json:"metadata"`
	Indexes        map[string]struct {
		Metadata indexMetadataJSON `json:"metadata"`
		Index    json.RawMessage   `json:"index"`
	} `json:"indexes"`
}

type indexMetadataJSON struct {
	CastType   uint              `json:"castType"`
	IndexType  indexer.IndexType `json:"type"`
	ColumnName string            `json:"column"`
}

func (m *IndexManager) UnmarshalJSON(data []byte) error {
	var err error
	id := internalDataJSON{}
	if err = json.Unmarshal(data, &id); err != nil {
		return err
	}

	m.internal = &internalData{
		SegmentID:      id.SegmentID,
		PartitionStart: id.PartitionStart,
		PartitionEnd:   id.PartitionEnd,
		Indexes:        make(map[string]indexer.Indexer, len(id.Indexes)),
	}
	for dim, idx := range id.Indexes {
		if err := creatNewIndexUnmarshal(m, dim, idx.Metadata.IndexType, idx.Metadata.CastType); err != nil {
			return err
		}
		if err := m.internal.Indexes[dim].Deserialize(idx.Index); err != nil {
			return err
		}
	}
	return nil
}

func creatNewIndexUnmarshal(m *IndexManager, dim string, it indexer.IndexType, ct uint) error {
	var err error
	switch ct {
	case uint(reflect.Int):
		err = CreateIndex[int](m, dim, it)
	case uint(reflect.Int8):
		err = CreateIndex[int8](m, dim, it)
	case uint(reflect.Int16):
		err = CreateIndex[int16](m, dim, it)
	case uint(reflect.Int32):
		err = CreateIndex[int32](m, dim, it)
	case uint(reflect.Int64):
		err = CreateIndex[int64](m, dim, it)
	case uint(reflect.Uint):
		err = CreateIndex[uint](m, dim, it)
	case uint(reflect.Uint8):
		err = CreateIndex[uint8](m, dim, it)
	case uint(reflect.Uint16):
		err = CreateIndex[uint16](m, dim, it)
	case uint(reflect.Uint32):
		err = CreateIndex[uint32](m, dim, it)
	case uint(reflect.Uint64):
		err = CreateIndex[uint64](m, dim, it)
	case uint(reflect.Float32):
		err = CreateIndex[float32](m, dim, it)
	case uint(reflect.Float64):
		err = CreateIndex[float64](m, dim, it)
	case uint(reflect.String):
		err = CreateIndex[string](m, dim, it)
	default:
		return fmt.Errorf("wrong type cast (uint %x) for index", ct)
	}
	if err != nil {
		return err
	}
	return nil
}
