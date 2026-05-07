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
		defer close(jobs)
		for _, url := range urls {
			select {
			case <-ctx.Done():
				return // context cancelled - stop sending jobs
			case jobs <- url:
				// job sent successfully - continue to next URL
			}
		}
	}()

	// Collect all results.
	// We know we'll receive exactly len(urls) results or stop if context is cancelled.
	gathered := make([]Result, 0, len(urls))
	for {
		select {
		case <-ctx.Done():
			// context cancelled - stop waiting and return what we've collected
			wg.Wait() // wait for workers to finish cleanup
			return gathered
		case r, ok := <-results:
			if !ok {
				// results channel closed - shouldn't happen before we get all results
				break
			}
			gathered = append(gathered, r)
			if len(gathered) == len(urls) {
				// received all results - we're done
				wg.Wait() // wait for workers to finish
				return gathered
			}
		}
	}
}
