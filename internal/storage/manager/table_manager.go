package manager

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/goccy/go-json"
	"github.com/spaghettifunk/norman/internal/common/entities"
)

const (
	eventIDName         string = "_normanID"
	partitionTimeFormat string = "2006-01-02T15:04:05"
)

type TableManager struct {
	Table          *entities.Table
	SegmentManager *SegmentManager

	baseDir           string
	datetimeFieldName string
	wg                sync.WaitGroup
	granularity       *entities.GranularitySpec
}

func NewTableManager(table *entities.Table) (*TableManager, error) {
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

	return &TableManager{
		Table:             table,
		SegmentManager:    sm,
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

	// TODO: advertise that a new segment is created
	// ...

	return nil
}

func (t *TableManager) InsertData(data []byte) error {
	event := make(map[string]interface{}, len(t.Table.EventSchema.Fields()))
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	// add the data to the current segment
	if err := t.SegmentManager.AppendData(event, t.datetimeFieldName, t.granularity.UnitSpec); err != nil {
		return err
	}

	// TODO: index the segment
	// ...

	return nil
}

// FlushSegment first persist on disk the current segment
// secondly, it compresses the segment to save space and lastly
// it reset the memory object so that it can start over
func (t *TableManager) FlushSegment() error {
	if err := t.SegmentManager.Flush(); err != nil {
		return err
	}

	// TODO: publish segment to consul
	// ...

	return nil
}
