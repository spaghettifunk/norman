package segment

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/spaghettifunk/norman/internal/common/schema"
)

type SegmentManager struct {
	segment *Segment
	schema  *schema.Schema
	// Aqua gRPC client
}

func NewSegmentManager(schema *schema.Schema) *SegmentManager {
	return &SegmentManager{
		schema: schema,
	}
}

func (sm *SegmentManager) CreateNewSegment() error {
	// get directory where to store the segment from Aqua
	dir := fmt.Sprintf("/tmp/%s_%s", time.Now().Format("2023-05-01T10:00:00"), sm.schema.Name)

	var err error
	segmentFile, err := os.OpenFile(
		path.Join(dir, ".segment"),
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0644,
	)
	if err != nil {
		return err
	}
	s, err := NewSegment(segmentFile)
	if err != nil {
		return err
	}
	sm.segment = s
	return nil
}

func (sm *SegmentManager) InsertRowInSegment(values []byte) error {
	err := sm.validateEvent(values)
	if err != nil {
		return nil
	}
	// add to segment here
	// ....

	return nil
}

func (sm *SegmentManager) validateEvent(values []byte) error {
	return nil
}

func (sm *SegmentManager) GetSegmentLength() int {
	return 0
}

// FlushSegment first persist on disk the current segment
// secondly, it compresses the segment to save space and lastly
// it reset the memory object so that it can start over
func (sm *SegmentManager) FlushSegment() error {
	if err := sm.segment.Flush(); err != nil {
		return err
	}
	return sm.compressSegment()
}

func (sm *SegmentManager) compressSegment() error {
	return nil
}
