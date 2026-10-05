package orchestrator

import (
	"bulk-parser/internal/model"
	"sync"
)

type task struct {
	job    model.Job
	result chan<- model.Result
}

type Orchestrator struct {
	jobs              chan task
	workerCount       int
	wg                sync.WaitGroup
	processorProvider ProcessorProvider
}

func NewOrchestrator(workercount int, bufferSize int, processorProvider ProcessorProvider) *Orchestrator {
	return &Orchestrator{
		jobs:              make(chan task, bufferSize),
		workerCount:       workercount,
		processorProvider: processorProvider,
	}
}

func (o *Orchestrator) Start() {
	for i := 0; i < o.workerCount; i++ {
		processor := o.processorProvider.Create()

		o.wg.Add(1)

		go func() {
			defer o.wg.Done()

			for task := range o.jobs {
				result := processor.Process(task.job)
				task.result <- result
			}
		}()
	}

}

func (o *Orchestrator) Process(jobs []model.Job) []model.Result {
	results := make(chan model.Result, len(jobs))

	for _, job := range jobs {
		o.jobs <- task{
			job:    job,
			result: results,
		}
	}

	output := make([]model.Result, 0, len(jobs))

	for i := 0; i < len(jobs); i++ {
		output = append(output, <-results)
	}

	return output
}
