package manager

type SegmentManager struct {
}

func NewSegmentManager() *SegmentManager {
	return &SegmentManager{}
}

func (sm *SegmentManager) Initialize() error {
	return nil
}

func (sm *SegmentManager) Start() error {
	return nil
}

func (sm *SegmentManager) Shutdown() error {
	return nil
}
