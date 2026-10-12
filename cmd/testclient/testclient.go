package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 50; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			body := fmt.Sprintf(
				`{"sources":[{"language":"go","data":"request-%d"}]}`,
				id,
			)

			resp, err := http.Post(
				"http://localhost:8081/process",
				"application/json",
				strings.NewReader(body),
			)

			if err != nil {
				fmt.Println("request", id, "error:", err)
				return
			}

			resp.Body.Close()
			fmt.Println("request", id, "completed")
		}(i)
	}

	wg.Wait()
	fmt.Println("All requests completed")
}
