package manager

import (
	"runtime"
	"sync"

	"github.com/spaghettifunk/norman/internal/common/model"
	"github.com/spaghettifunk/norman/internal/storage/ingestion"
	"github.com/spaghettifunk/norman/pkg/workerpool"
)

var (
	MaxNumberOfWorkers int = runtime.NumCPU()
)

type IngestionJobManager struct {
	WorkerPool workerpool.Pool
	wg         *sync.WaitGroup
	// Aqua gRPC client
}

func NewIngestionJobManager() (*IngestionJobManager, error) {
	wp, err := workerpool.NewWorkerPool(MaxNumberOfWorkers, 0)
	if err != nil {
		return nil, err
	}
	return &IngestionJobManager{
		WorkerPool: wp,
		wg:         &sync.WaitGroup{},
	}, nil
}

func (ijm *IngestionJobManager) Initialize() error {
	ijm.WorkerPool.Start()
	return nil
}

// TODO: add a Final State Machine for handling the job
func (ijm *IngestionJobManager) Execute(config []byte) error {
	// parse config and transform into an IngestionJob
	j, err := model.NewIngestionJob(config)
	if err != nil {
		return err
	}
	// create the actual job to be exectued
	job, err := ingestion.NewJob(j)
	if err != nil {
		return err
	}

	// Set initialization of the job
	job.Initialize()

	// add new task to the workerpool and wait until completion
	ijm.wg.Add(1)
	go func() {
		ijm.WorkerPool.AddWork(job)
		// Notify Aqua that a new IngestionJob is in progress
	}()
	ijm.wg.Wait()

	return nil
}

func (ijm *IngestionJobManager) Shutdown() error {
	ijm.WorkerPool.Stop()
	return nil
}
