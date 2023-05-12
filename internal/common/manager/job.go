package manager

import (
	"fmt"

	"github.com/spaghettifunk/norman/internal/common/model"
	realtime_ingestion "github.com/spaghettifunk/norman/internal/storage/ingestion/realtime"
)

type IngestionJobManager struct {
	// Aqua gRPC client
}

func NewIngestionJobManager() *IngestionJobManager {
	return &IngestionJobManager{}
}

func (ijm *IngestionJobManager) Initialize() error {
	return nil
}

func (ijm *IngestionJobManager) Start() error {
	return nil
}

// TODO: add a Final State Machine for handling the job
func (ijm *IngestionJobManager) CreateJob(config []byte) error {
	j, err := model.NewIngestionJob(config)
	if err != nil {
		return err
	}

	switch j.Type {
	case model.Offline:
		// submit to offline queue
		break
	case model.Realtime:
		// submit to realtime queue
		rt, err := realtime_ingestion.New(j.IngestionConfiguration.StreamIngestionConfiguration)
		if err != nil {
			return err
		}
		return rt.ReadEvents()
	default:
		return fmt.Errorf("incorrect ingestion job type: %s", j.Type)
	}

	return nil
}

func (ijm *IngestionJobManager) Shutdown() error {
	return nil
}
