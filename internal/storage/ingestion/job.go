package ingestion

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/spaghettifunk/norman/internal/common/entities"
	cingestion "github.com/spaghettifunk/norman/internal/common/ingestion"
	"github.com/spaghettifunk/norman/internal/storage/manager"

	"github.com/spaghettifunk/norman/pkg/realtime/kafka"
	"github.com/spaghettifunk/norman/pkg/realtime/kinesis"
)

type Ingestion interface {
	Initialize() error
	GetEvents() error
	Shutdown(failure bool) error
}

type Job struct {
	ID              uuid.UUID
	IngestionClient Ingestion
	Configuration   *cingestion.IngestionJobConfiguration
	JobStatus       cingestion.JobStatus
	wg              *sync.WaitGroup
}

// NewJob creates a new Job based on the configuration in input
func NewJob(cfg *cingestion.IngestionJobConfiguration, table *entities.Table) (*Job, error) {
	// UUID of the job
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	// create new TableManager
	tm, err := manager.NewTableManager(table)
	if err != nil {
		return nil, err
	}

	// get the correct client based on ingestion type
	var client Ingestion
	switch cfg.Type {
	case cingestion.StreamKafka:
		if client, err = kafka.NewIngestor(cfg.IngestionConfiguration.Realtime.KafkaConfiguration, tm); err != nil {
			return nil, err
		}
	case cingestion.StreamKinesis:
		if client, err = kinesis.New(); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid Realtime Ingestion type: %s", cfg.Type)
	}
	// create new job
	return &Job{
		ID:              id,
		IngestionClient: client,
		Configuration:   cfg,
		JobStatus:       cingestion.JobCreated,
	}, nil
}

// Initialize initiate the client and prepare for data reading process
func (j *Job) Initialize() error {
	j.JobStatus = cingestion.JobInitialized
	return j.IngestionClient.Initialize()
}

// Shutdown terminate the process of data incoming
func (j *Job) Shutdown(failure bool) error {
	if failure {
		j.JobStatus = cingestion.JobPartiallyDone
	} else {
		j.JobStatus = cingestion.JobDone
	}
	return j.IngestionClient.Shutdown(failure)
}

func (j *Job) Execute() error {
	defer j.wg.Done()
	j.JobStatus = cingestion.JobInProgress
	return j.IngestionClient.GetEvents()
}

func (j *Job) OnFailure(e error) {
	j.JobStatus = cingestion.JobFailed
	log.Error().Err(e).Msg("error while consuming messages")
}
