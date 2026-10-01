package parser

import "bulk-parser/internal/model"

type Parser interface {
	Parse(job model.Job) model.Result
}
