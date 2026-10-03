package worker

import (
	"bulk-parser/internal/factory"
	"bulk-parser/internal/model"
)

type Worker struct {
	factory *factory.ParserFactory
}

func NewWorker(factory *factory.ParserFactory) *Worker {
	return &Worker{
		factory: factory,
	}
}

func (w *Worker) Process(job model.Job) model.Result {
	parser, err := w.factory.GetParser(job.Language)
	if err != nil {
		return model.Result{
			JobID: job.ID,
			Error: err,
		}
	}

	return parser.Parse(job)
}
