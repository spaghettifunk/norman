package manager

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/goccy/go-json"
	"github.com/spaghettifunk/norman/internal/common/entities"
	"github.com/spaghettifunk/norman/internal/common/types"
	"github.com/spaghettifunk/norman/internal/storage/indexer"
	"github.com/spaghettifunk/norman/pkg/eventmanager"
)

const (
	eventIDName         string = "_normanID"
	partitionTimeFormat string = "2006-01-02T15:04:05"
)

type TableManager struct {
	Table             *entities.Table
	SegmentManager    *SegmentManager
	IndexManager      *IndexManager
	baseDir           string
	datetimeFieldName string
	wg                sync.WaitGroup
	granularity       *entities.GranularitySpec
}

func NewTableManager(table *entities.Table, indexes map[indexer.IndexType][]string) (*TableManager, error) {
	// TODO: this should depend on a folder that comes from Configuration
	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// TODO: how to get the tenant name? --> default for now
	// Format: os_path + output/{tenantID}/
	baseDir := fmt.Sprintf("%s/output/default/%s", path, table.Name)
	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		return nil, err
	}
	// extract the datetime field
	dtField := table.GetDatetimeField()

	granularity, err := table.Schema.GetGranularity()
	if err != nil {
		return nil, err
	}

	interval := time.Duration(granularity.Size) * granularity.UnitSpec
	sm, err := NewSegmentManager(table.Name, baseDir, interval)
	if err != nil {
		return nil, err
	}

	// create a new index manager
	im := NewIndexManager(baseDir)
	for indexType, columns := range indexes {
		for _, column := range columns {
			df := table.Schema.GetDimensionField(column)
			if err := createColumnIndex(im, column, df.DataType, indexType); err != nil {
				return nil, err
			}
		}
	}

	return &TableManager{
		Table:             table,
		SegmentManager:    sm,
		IndexManager:      im,
		datetimeFieldName: dtField.Name,
		baseDir:           baseDir,
		wg:                sync.WaitGroup{},
		granularity:       granularity,
	}, nil
}

func (t *TableManager) CreateNewSegment() error {
	// create segment
	if err := t.SegmentManager.Create(t.Table.EventSchema); err != nil {
		return err
	}

	// fire segment initialized event
	eventmanager.GetEventManager().Notify(eventmanager.Event{
		Type: eventmanager.SegmentInitialized,
		Data: t.baseDir,
	})

	return nil
}

func (t *TableManager) InsertData(data []byte) error {
	event := make(map[string]interface{}, len(t.Table.EventSchema.Fields()))
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	// add the data to the current segment
	evtID, err := t.SegmentManager.AppendData(event, t.datetimeFieldName, t.granularity.UnitSpec)
	if err != nil {
		return err
	}

	// index the segment's value
	for columnName, val := range event {
		t.IndexManager.Add(columnName, evtID, val)
	}

	return nil
}

// FlushSegment first persist on disk the current segment
// secondly, it compresses the segment to save space and lastly
// it reset the memory object so that it can start over
func (t *TableManager) FlushSegment() error {
	if err := t.SegmentManager.Flush(); err != nil {
		return err
	}

	sID := t.SegmentManager.GetSegmentID()
	if err := t.IndexManager.PersistToDisk(sID, t.IndexManager.internal.PartitionStart,
		t.IndexManager.internal.PartitionEnd); err != nil {
		return err
	}

	// fire segment initialized event
	eventmanager.GetEventManager().Notify(eventmanager.Event{
		Type: eventmanager.SegmentCreated,
		Data: true,
	})

	return nil
}

func createColumnIndex(im *IndexManager, column, dataType string, indexType indexer.IndexType) error {
	switch dataType {
	case types.Integer:
		return CreateIndex[int](im, column, indexType)
	case types.Long:
		return CreateIndex[int64](im, column, indexType)
	case types.Float:
		return CreateIndex[float32](im, column, indexType)
	case types.Double:
		return CreateIndex[float64](im, column, indexType)
	case types.String:
		return CreateIndex[string](im, column, indexType)
	default:
		return fmt.Errorf("invalid type for creating an index")
	}
}
