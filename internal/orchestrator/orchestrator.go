package orchestrator

import (
	"bulk-parser/internal/model"
	"log"
	"sync"
	"time"
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

func NewOrchestrator(workercount int, jobBuffer int, processorProvider ProcessorProvider) *Orchestrator {
	return &Orchestrator{
		jobs:              make(chan task, jobBuffer),
		workerCount:       workercount,
		processorProvider: processorProvider,
	}
}

func (o *Orchestrator) Start() {
	for i := 0; i < o.workerCount; i++ {
		processor := o.processorProvider.Create()
		workerID := i + 1

		o.wg.Add(1)

		go func() {
			defer o.wg.Done()

			log.Printf("[WORKER %d] Started and waiting for jobs", workerID)

			for task := range o.jobs {
				log.Printf(
					"[WORKER %d] PICKED JobID=%d | Queue remaining=%d",
					workerID,
					task.job.ID,
					task.job.Data,
				)

				start := time.Now()
				result := processor.Process(task.job)

				log.Printf(
					"[WORKER %d] COMPLETED JobID=%d | Duration=%v | Error=%v",
					workerID,
					task.job.ID,
					task.job.Data,
					time.Since(start),
					result.Error,
				)

				task.result <- result
			}

			log.Printf("[WORKER %d] Stopped", workerID)
		}()
	}
}

func (o *Orchestrator) Process(jobs []model.Job) []model.Result {
	results := make(chan model.Result, len(jobs))

	for _, job := range jobs {
		log.Printf(
			"[QUEUE] Submitting JobID=%d | Language=%s",
			job.ID,
			job.Language,
		)

		o.jobs <- task{
			job:    job,
			result: results,
		}

		log.Printf(
			"[QUEUE] Submitted JobID=%d | Buffered jobs=%d/%d",
			job.ID,
			len(o.jobs),
			cap(o.jobs),
		)
	}

	output := make([]model.Result, 0, len(jobs))

	for i := 0; i < len(jobs); i++ {
		output = append(output, <-results)
	}

	return output
}

func (o *Orchestrator) Shutdown() {
	close(o.jobs)
	o.wg.Wait()
}
