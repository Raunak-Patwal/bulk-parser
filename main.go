package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bulk-parser/internal/factory"
	"bulk-parser/internal/handler"
	"bulk-parser/internal/limiter"
	"bulk-parser/internal/orchestrator"
	"bulk-parser/internal/worker"
)

func main() {

	// 1. Create parser dependencies
	parserFactory := factory.Newparsefactory()

	// 2. Create processor provider
	processorProvider := worker.NewProvider(parserFactory)

	// 3. Create orchestrator
	orch := orchestrator.NewOrchestrator(
		5,  // worker count
		20, // job buffer size
		processorProvider,
	)

	// 4. Start worker pool
	orch.Start()

	// 5. Create request limiter
	requestLimiter := limiter.NewRequestLimiter(5)

	// 6. Create HTTP handler
	h := handler.NewHandler(
		requestLimiter,
		orch,
	)

	// 7. Register routes
	mux := http.NewServeMux()
	mux.HandleFunc("/process", h.Process)

	// 8. Create HTTP server
	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	// 9. Start HTTP server
	go func() {
		log.Println("server running on :8081")

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// 10. Wait for shutdown signal
	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-stop

	// 11. Gracefully shutdown HTTP server
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	// 12. Shutdown worker pool
	orch.Shutdown()

	log.Println("server stopped")
}
