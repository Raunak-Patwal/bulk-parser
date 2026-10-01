package orchestrator

import (
	"bulk-parser/internal/model"
	"bulk-parser/internal/worker"
	"sync"
)

type Orchestrator struct {
	jobs        chan model.Job
	results     chan model.Result
	workerCount int
	wg          sync.WaitGroup
	workerFunc  wokerFunction
}

type wokerFunction func()

func newOrchestrator(workercount int, bufferSize int, workerFunc wokerFunction) *Orchestrator {
	return &Orchestrator{
		jobs:        make(chan model.Job, bufferSize),
		results:     make(chan model.Result, bufferSize),
		workerfunc:  func() {},
		workerCount: workercount,
	}
}

func (o *Orchestrator) Start() {
	for i := 0; i < o.workerCount; i++ {
		w := worker.NewWorker()

		o.wg.Add(1)

		go func() {
			defer o.wg.Done()

			w.Start(o.jobs, o.results)
		}()
	}

	go func() {
		o.wg.Wait()
		close(o.results)
	}()
}

func (o *Orchestrator) Process(request model.BulkRequest) []model.Result {

	for i, source := range request.Sources {

		job := model.Job{
			ID:       i + 1,
			Language: source.Language,
			Data:     source.Data,
		}

		o.jobs <- job
	}

	close(o.jobs)

	var results []model.Result

	for result := range o.results {
		results = append(results, result)
	}

	return results
}
