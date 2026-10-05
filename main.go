package main

import (
	"log"
	"net/http"

	"bulk-parser/internal/factory"
	"bulk-parser/internal/handler"
	"bulk-parser/internal/limiter"
	"bulk-parser/internal/orchestrator"
	"bulk-parser/internal/worker"
)

func main() {

	// Parser dependencies
	parserFactory := factory.Newparsefactory()

	// Worker dependency
	processorProvider := worker.NewProvider(parserFactory)

	orch := orchestrator.NewOrchestrator(
		5,
		20,
		processorProvider,
	)

	orch.Start()

	requestLimiter := limiter.NewRequestLimiter(5)

	h := handler.NewHandler(
		requestLimiter,
		orch,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/process", h.Process)

	log.Println("server running on :8081")

	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		log.Fatal(err)
	}
}
