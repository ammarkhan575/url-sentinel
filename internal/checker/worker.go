package checker

import (
	"context"
	"sync"
	"time"
)

// Pool holds the configuration for a worker pool run.
type Pool struct {
	Workers int
	Timeout time.Duration
}

// Run is the heart of the system. It:
// 1. Spawns N worker goroutines
// 2. Sends all URLs into a jobs channel
// 3. Collects all results from a results channel
// 4. Returns when every URL is checked OR context is cancelled
func (p *Pool) Run(ctx context.Context, urls []string) []Result {
	// Buffered channels — size = number of URLs so producers never block waiting for workers
	jobs := make(chan string, len(urls))
	results := make(chan Result, len(urls))

	// --- Launch workers FIRST (before sending jobs) ---
	// sync.WaitGroup tracks how many goroutines are still running.
	// Think of it as a counter: Add(1) increments, Done() decrements, Wait() blocks until 0.
	var wg sync.WaitGroup
	for i := 0; i < p.Workers; i++ {
		wg.Add(1) // increment BEFORE launching goroutine — avoids a race with Wait()
		go func() {
			defer wg.Done() // decrement when this goroutine exits — defer guarantees it runs

			// range over jobs blocks until a job arrives OR jobs is closed.
			// When jobs is closed + empty, the loop exits and the goroutine finishes
			for url := range jobs {
				// select: wait on multiple channel operations simultaneously.
				// The first one ready wins. Both cases are checked on every iteration.
				select {
				case <-ctx.Done():
					// Context cancelled (timeout or user signal) - stop working
					return
				default:
					results <- Check(ctx, url, p.Timeout)
				}
			}
		}()
	}
	// send jobs in a separate goroutine so we can simultaneously collect results
	go func() {
		// ensure jobs channel is closed once we've enqueued all URLs or context cancelled
		defer close(jobs)
		for _, url := range urls {
			select {
			case <-ctx.Done():
				return // context cancelled - stop sending jobs
			case jobs <- url:
				// job sent successfully - continue to next URL
			}
		}
		close(jobs); // CRITICAL: closing jobs causes worker for-range loops to exit
	}()

	// close results when all workers have finished
	go func() {
		wg.Wait()
		close(results)
	}()

	// collect results until results channel is closed
	gathered := make([]Result, 0, len(urls))
	for r := range results {
		gathered = append(gathered, r)
	}

	return gathered
}
