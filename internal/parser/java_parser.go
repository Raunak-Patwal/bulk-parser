package parser

import "bulk-parser/internal/model"

type JavaParser struct{}

func (p JavaParser) Parse(job model.Job) model.Result {
	//WE will implement java language parser logic here
	return model.Result{
		JobID:  job.ID,
		Output: "Java code parsed",
	}
}
