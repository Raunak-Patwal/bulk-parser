package parser

import "bulk-parser/internal/model"

type GoParser struct{}

func (p GoParser) Parse(job model.Job) model.Result {
	return model.Result{
		JobID:  job.ID,
		Output: "Parsed Go data: " + job.Data,
	}
}
