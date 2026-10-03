package orchestrator

type ProcessorProvider interface {
	Create() Processor
}
