package workerpool

import (
	"fmt"

	"sync"

	"github.com/rs/zerolog/log"
)

type Pool interface {
	// Start gets the workerpool ready to process jobs, and should only be called once
	Start()
	// Stop stops the workerpool, tears down any required resources,
	// and should only be called once
	Stop()
	// AddWork adds a task for the worker pool to process. It is only valid after
	// Start() has been called and before Stop() has been called.
	AddWork(Task)
}

type Task interface {
	// Execute performs the work
	Execute() error
	// OnFailure handles any error returned from Execute()
	OnFailure(error)
}

type WorkerPool struct {
	numWorkers int
	tasks      chan Task
	// ensure the pool can only be started once
	start sync.Once
	// ensure the pool can only be stopped once
	stop sync.Once
	// close to signal the workers to stop working
	quit chan struct{}
}

var ErrNoWorkers = fmt.Errorf("attempting to create worker pool with less than 1 worker")
var ErrNegativeChannelSize = fmt.Errorf("attempting to create worker pool with a negative channel size")

func NewWorkerPool(numWorkers int, channelSize int) (Pool, error) {
	if numWorkers <= 0 {
		return nil, ErrNoWorkers
	}
	if channelSize < 0 {
		return nil, ErrNegativeChannelSize
	}

	tasks := make(chan Task, channelSize)

	return &WorkerPool{
		numWorkers: numWorkers,
		tasks:      tasks,
		start:      sync.Once{},
		stop:       sync.Once{},
		quit:       make(chan struct{}),
	}, nil
}

func (p *WorkerPool) Start() {
	p.start.Do(func() {
		log.Info().Msgf("starting worker pool")
		p.startWorkers()
	})
}

func (p *WorkerPool) startWorkers() {
	for i := 0; i < p.numWorkers; i++ {
		go func(workerNum int) {
			log.Info().Msgf("starting worker %d", workerNum)
			for {
				select {
				case <-p.quit:
					log.Info().Msgf("stopping worker %d with quit channel\n", workerNum)
					return
				case task, ok := <-p.tasks:
					if !ok {
						log.Info().Msgf("stopping worker %d with closed tasks channel\n", workerNum)
						return
					}
					if err := task.Execute(); err != nil {
						task.OnFailure(err)
					}
				}
			}
		}(i)
	}
}

func (p *WorkerPool) Stop() {
	p.stop.Do(func() {
		log.Info().Msgf("stopping simple worker pool")
		close(p.quit)
	})
}

// AddWork adds work to the SimplePool. If the channel buffer is full (or 0) and
// all workers are occupied, this will hang until work is consumed or Stop() is called.
func (p *WorkerPool) AddWork(t Task) {
	select {
	case p.tasks <- t:
	case <-p.quit:
	}
}

// AddWorkNonBlocking adds work to the SimplePool and returns immediately
func (p *WorkerPool) AddWorkNonBlocking(t Task) {
	go p.AddWork(t)
}
