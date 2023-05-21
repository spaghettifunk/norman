package ingestion

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	common_ingestion "github.com/spaghettifunk/norman/internal/common/ingestion"
	"github.com/spaghettifunk/norman/internal/common/schema"
	"github.com/spaghettifunk/norman/internal/storage/segment"
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
	Configuration   *common_ingestion.IngestionJobConfiguration
	JobStatus       common_ingestion.JobStatus
	wg              *sync.WaitGroup
}

// NewJob creates a new Job based on the configuration in input
func NewJob(cfg *common_ingestion.IngestionJobConfiguration, schema *schema.Schema) (*Job, error) {
	// new segment manager
	sm, err := segment.NewSegmentManager(schema)
	if err != nil {
		return nil, err
	}

	// UUID of the job
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	// get the correct client based on ingestion type
	var client Ingestion
	switch cfg.Type {
	case common_ingestion.StreamKafka:
		if client, err = kafka.NewIngestor(cfg.IngestionConfiguration.Realtime.KafkaConfiguration, sm); err != nil {
			return nil, err
		}
	case common_ingestion.StreamKinesis:
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
		JobStatus:       common_ingestion.JobCreated,
	}, nil
}

// Initialize initiate the client and prepare for data reading process
func (j *Job) Initialize() error {
	j.JobStatus = common_ingestion.JobInitialized
	return j.IngestionClient.Initialize()
}

// Shutdown terminate the process of data incoming
func (j *Job) Shutdown(failure bool) error {
	if failure {
		j.JobStatus = common_ingestion.JobPartiallyDone
	} else {
		j.JobStatus = common_ingestion.JobDone
	}
	return j.IngestionClient.Shutdown(failure)
}

func (j *Job) Execute() error {
	defer j.wg.Done()
	j.JobStatus = common_ingestion.JobInProgress
	return j.IngestionClient.GetEvents()
}

func (j *Job) OnFailure(e error) {
	j.JobStatus = common_ingestion.JobFailed
	log.Error().Err(e).Msg("error while consuming messages")
}
