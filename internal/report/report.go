package report

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ammarkhan575/url-sentinel/internal/checker"
)

type Summary struct {
	Total     int              `json:"total"`
	Up        int              `json:"up"`
	Down      int              `json:"down"`
	Timeouts  int              `json:"timeouts"`
	TotalTime time.Duration    `json:"total_time_ms"`
	Results   []checker.Result `json:"results"`
}

func Build(results []checker.Result, elapsed time.Duration) Summary {
	s := Summary{Results: results, TotalTime: elapsed}
	for _, r := range results {
		switch r.Status {
		case checker.StatusUp:
			s.Up++
		case checker.StatusDown:
			s.Down++
		case checker.StatusTimeout:
			s.Timeouts++
		}
	}
	return s
}

func Print(s Summary) {
	for _, r := range s.Results {
		icon := "✓"
		if r.Status != checker.StatusUp {
			icon = "✗"
		}
		fmt.Printf("%s %-50s %s %dms\n", icon, r.URL, r.Status, r.Latency.Milliseconds())
	}
	fmt.Printf("\nDone: %d up, %d down, %d timeout. Time: %s\n",
		s.Up, s.Down, s.Timeouts, s.TotalTime.Round(time.Millisecond))
}

func SaveJson(s Summary, path string) error {
	data, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
