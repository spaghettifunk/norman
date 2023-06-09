package manager

import (
	"runtime"
	"sync"

	"github.com/rs/zerolog/log"
	cingestion "github.com/spaghettifunk/norman/internal/common/ingestion"
	"github.com/spaghettifunk/norman/internal/storage/ingestion"
	"github.com/spaghettifunk/norman/pkg/consul"
	"github.com/spaghettifunk/norman/pkg/workerpool"
)

var (
	MaxNumberOfWorkers int = runtime.NumCPU()
)

type IngestionJobManager struct {
	WorkerPool workerpool.Pool
	wg         *sync.WaitGroup
	consul     *consul.Consul
}

func NewIngestionJobManager(c *consul.Consul) (*IngestionJobManager, error) {
	wp, err := workerpool.NewWorkerPool(MaxNumberOfWorkers, 0)
	if err != nil {
		return nil, err
	}
	return &IngestionJobManager{
		WorkerPool: wp,
		wg:         &sync.WaitGroup{},
		consul:     c,
	}, nil
}

func (ijm *IngestionJobManager) Initialize() error {
	ijm.WorkerPool.Start()
	return nil
}

// TODO: add a Final State Machine for handling the job
func (ijm *IngestionJobManager) Execute(config *cingestion.IngestionJobConfiguration) error {
	// fetch table information
	table, err := ijm.consul.GetTableConfiguration(config.SegmentConfiguration.TableName)
	if err != nil {
		return err
	}
	// get arrow schema
	evSchema, err := table.Schema.GetFullArrowSchema()
	if err != nil {
		return err
	}
	table.EventSchema = evSchema

	// create the actual job to be exectued
	job, err := ingestion.NewJob(config, table)
	if err != nil {
		return err
	}

	// Set initialization of the job
	if err := job.Initialize(); err != nil {
		return err
	}

	// add new task to the workerpool and wait until completion
	ijm.wg.Add(1)
	// TODO: looks really hacky!
	go func() {
		go func() {
			ijm.WorkerPool.AddWork(job)
			// Notify Consul that a new IngestionJob is in progress
			if err := ijm.consul.PutIngestionJobStatus(job.ID.String(), job.JobStatus); err != nil {
				log.Error().Err(err).Msgf("failed to update consul metadata for job with ID %s", job.ID.String())
			}
			log.Debug().Msgf("consul updated for job with id %s", job.ID.String())
		}()
		ijm.wg.Wait()
	}()

	return nil
}

func (ijm *IngestionJobManager) Shutdown() error {
	ijm.WorkerPool.Stop()
	return nil
}
