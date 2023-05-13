package ingestion

import (
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/spaghettifunk/norman/internal/common/model"
	"github.com/spaghettifunk/norman/pkg/realtime/kafka"
	"github.com/spaghettifunk/norman/pkg/realtime/kinesis"
)

type Ingestion interface {
	Initialize() error
	GetEvents() error
	Shutdown(failure bool) error
}

type Job struct {
	Client        Ingestion
	Configuration *model.IngestionJobConfiguration
	Schema        *model.Schema
	JobStatus     model.JobStatus
	wg            *sync.WaitGroup
}

// NewJob creates a new Job based on the configuration in input
func NewJob(cfg *model.IngestionJobConfiguration) (*Job, error) {
	ingestor := &Job{
		Configuration: cfg,
		JobStatus:     model.JobCreated,
	}
	var err error
	switch cfg.Type {
	case model.StreamKafka:
		ingestor.Client, err = kafka.NewIngestor(cfg.IngestionConfiguration.Realtime.KafkaConfiguration)
	case model.StreamKinesis:
		ingestor.Client = kinesis.New()
	default:
		return nil, fmt.Errorf("invalid Realtime Ingestion type: %s", cfg.Type)
	}
	if err != nil {
		return nil, err
	}
	return ingestor, nil
}

// Initialize initiate the client and prepare for data reading process
func (j *Job) Initialize() error {
	j.JobStatus = model.JobInitialized
	return j.Client.Initialize()
}

// Shutdown terminate the process of data incoming
func (j *Job) Shutdown(failure bool) error {
	if failure {
		j.JobStatus = model.JobPartiallyDone
	} else {
		j.JobStatus = model.JobDone
	}
	return j.Client.Shutdown(failure)
}

func (j *Job) Execute() error {
	defer j.wg.Done()
	j.JobStatus = model.JobInProgress
	return j.Client.GetEvents()
}

func (j *Job) OnFailure(e error) {
	j.JobStatus = model.JobFailed
	log.Error().Err(e).Msg("error while consuming messages")
}
