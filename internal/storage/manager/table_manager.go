package manager

import (
	"fmt"
	"os"

	"github.com/spaghettifunk/norman/internal/common/entities"
	"github.com/spaghettifunk/norman/internal/storage/segment"
)

type TableManager struct {
	Table         *entities.Table
	activeSegment *segment.Segment
	baseDir       string
	segments      []*segment.Segment
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
	return &TableManager{
		Table:   table,
		baseDir: baseDir,
	}, nil
}

func (t *TableManager) CreateNewSegment() error {
	s, err := segment.NewSegment(t.baseDir, t.Table.EventSchema)
	if err != nil {
		return err
	}
	t.activeSegment = s
	return nil
}

func (t *TableManager) InsertData(data []byte) error {
	return t.activeSegment.InsertData(data)
}

func (t *TableManager) GetSegmentSize() uint32 {
	return t.activeSegment.GetCounter()
}

// FlushSegment first persist on disk the current segment
// secondly, it compresses the segment to save space and lastly
// it reset the memory object so that it can start over
func (t *TableManager) FlushSegment() error {
	if err := t.activeSegment.Flush(); err != nil {
		return err
	}
	// store the active segments in the list of segments
	t.segments = append(t.segments, t.activeSegment)
	return nil
}