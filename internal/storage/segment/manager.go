package segment

import "github.com/spaghettifunk/norman/internal/common/schema"

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
	s, err := NewSegment()
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
	return len(sm.segment.Rows)
}

// FlushSegment first persist on disk the current segment
// secondly, it compresses the segment to save space and lastly
// it reset the memory object so that it can start over
func (sm *SegmentManager) FlushSegment() error {
	// fmt.Sprintf("/tmp/norman/%s/%s", k.Topic, time.Now().String())
	if err := sm.segment.Flush(""); err != nil {
		return err
	}

	if err := sm.compressSegment(); err != nil {
		return err
	}

	return sm.segment.Reset()
}

func (sm *SegmentManager) compressSegment() error {
	return nil
}
