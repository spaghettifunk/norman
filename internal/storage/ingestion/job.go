package ingestion

import (
	"fmt"
	"sync"

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
	Client        Ingestion
	Configuration *common_ingestion.IngestionJobConfiguration
	JobStatus     common_ingestion.JobStatus
	wg            *sync.WaitGroup
}

// NewJob creates a new Job based on the configuration in input
func NewJob(cfg *common_ingestion.IngestionJobConfiguration) (*Job, error) {
	// call Aqua to get the schema
	sc, err := fetchSchema(cfg.SegmentConfiguration.SchemaName)
	if err != nil {
		return nil, err
	}
	// new segment manager
	sm := segment.NewSegmentManager(sc)

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
		Client:        client,
		Configuration: cfg,
		JobStatus:     common_ingestion.JobCreated,
	}, nil
}

// fetchSchema returns the Schema object from Acqua based on its name
// if the schema is not found, return an error
func fetchSchema(schemaName string) (*schema.Schema, error) {
	return &schema.Schema{}, nil
}

// Initialize initiate the client and prepare for data reading process
func (j *Job) Initialize() error {
	j.JobStatus = common_ingestion.JobInitialized
	return j.Client.Initialize()
}

// Shutdown terminate the process of data incoming
func (j *Job) Shutdown(failure bool) error {
	if failure {
		j.JobStatus = common_ingestion.JobPartiallyDone
	} else {
		j.JobStatus = common_ingestion.JobDone
	}
	return j.Client.Shutdown(failure)
}

func (j *Job) Execute() error {
	defer j.wg.Done()
	j.JobStatus = common_ingestion.JobInProgress
	return j.Client.GetEvents()
}

func (j *Job) OnFailure(e error) {
	j.JobStatus = common_ingestion.JobFailed
	log.Error().Err(e).Msg("error while consuming messages")
}
