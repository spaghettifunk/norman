package manager

import (
	"github.com/spaghettifunk/norman/internal/common/model"
	"github.com/spaghettifunk/norman/internal/common/segment"
)

type SegmentManager struct {
	Segment *segment.Segment
	Schema  *model.Schema
	// Aqua gRPC client
}

func NewSegmentManager() *SegmentManager {
	return &SegmentManager{}
}

func (sm *SegmentManager) CreateNewSegment() error {
	s, err := segment.NewSegment()
	if err != nil {
		return err
	}
	sm.Segment = s
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

func (sm *SegmentManager) FlushSegment() error {
	return nil
}

func (sm *SegmentManager) CompressSegment() error {
	return nil
}

func (sm *SegmentManager) ResetSegment() error {
	return nil
}
