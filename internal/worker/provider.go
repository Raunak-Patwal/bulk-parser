package worker

import (
	"bulk-parser/internal/factory"
	"bulk-parser/internal/orchestrator"
)

type Provider struct {
	factory *factory.ParserFactory
}

func NewProvider(factory *factory.ParserFactory) *Provider {
	return &Provider{
		factory: factory,
	}
}

func (p *Provider) Create() orchestrator.Processor {
	return &Worker{
		factory: p.factory,
	}
}
