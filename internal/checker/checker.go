package checker

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Status string

const (
	StatusUp      Status = "up"
	StatusDown    Status = "down"
	StatusTimeout Status = "timeout"
)

// Result holds everything about one checked URL.
// This travels through channels — keep it a value type (no pointers).
// Small structs passed by value through channels are idiomatic and safe.
type Result struct {
	URL        string        `json:"url"`
	Status     Status        `json:"status"`
	StatusCode int           `json:"status_code,omitempty"`
	Latency    time.Duration `json:"latency_ms"`
	Error      string        `json:"error,omitempty"`
	CheckedAt  time.Time     `json:"checked_at"`
}

// Check performs a single HTTP GET and returns a Result.
// It respects context cancellation — if ctx is cancelled mid-flight, it returns immediately.
func Check(ctx context.Context, url string, timeout time.Duration) Result {
	start := time.Now()

	var start1 time.Time = time.Now()

	fmt.Println(start1, start)
	// Per-request context: the shorter of the global ctx and our timeout wins.
	// This is the idiomatic way to layer deadlines.
	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel() // always release resources - even if timeout fires first

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, url, nil)

	if err != nil {
		return Result{URL: url, Status: StatusDown, Error: err.Error(), CheckedAt: time.Now()}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	latency := time.Since(start)

	if err != nil {
		// context.DeadlineExceeded means our timeout fired
		// context.Canceled means the parent context was cancelled (e.g. user hit Ctrl+C)
		if ctx.Err() != nil {
			return Result{URL: url, Status: StatusTimeout, Latency: latency, CheckedAt: time.Now()}
		}

		return Result{URL: url, Status: StatusDown, Latency: latency, Error: err.Error(), CheckedAt: time.Now()}
	}
	defer resp.Body.Close() // defer cleanup — runs even if we return early below

	status := StatusUp
	if resp.StatusCode >= 400 {
		status = StatusDown
	}

	return Result{
		URL:        url,
		Status:     status,
		StatusCode: resp.StatusCode,
		Latency:    latency,
		CheckedAt:  time.Now(),
	}
}
