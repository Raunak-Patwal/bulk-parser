package main

import (
	"fmt"
	"sync"
	"time"

	"bulk-parser/internal/limiter"
	"bulk-parser/internal/orchestrator"
)

func main() {
	maxRequests := 5

	requestLimiter := limiter.NewRequestLimiter(maxRequests)

	var wg sync.WaitGroup

	orch := orchestrator.NewOrchestrator(5, 20, process)

	for i := 1; i <= 10; i++ {
		wg.Add(1)

		go func(requestID int) {
			defer wg.Done()

			requestLimiter.Acquire()

			fmt.Printf("Request %d STARTED\n", requestID)

			time.Sleep(3 * time.Second)

			fmt.Printf("Request %d FINISHED\n", requestID)

			requestLimiter.Release()
		}(i)
	}

	wg.Wait()

	fmt.Println("All requests finished")
}

func process(request string) string {
	print("hi")
}
