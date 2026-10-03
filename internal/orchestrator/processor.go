package orchestrator

import "bulk-parser/internal/model"

type Processor interface {
	Process(job model.Job) model.Result
}
