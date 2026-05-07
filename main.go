package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ammarkhan575/url-sentinel/internal/checker"
	"github.com/ammarkhan575/url-sentinel/internal/report"
)

func main() {
	fmt.Println("Welcome to URLSentinel v0.1")

	// flags
	file := flag.String("file", "", "path to file containing URLs (one per line)")
	workers := flag.Int("workers", 10, "number of concurrent workers")
	timeout := flag.Duration("timeout", 5*time.Second, "HTTP request timeout duration")
	output := flag.String("output", "", "path to save JSON report (optional)")

	flag.Parse()

	if *file == "" {
		fmt.Println("Error: -file flag is required")
		flag.Usage()
		os.Exit(1)
	}

	urls, err := readLines(*file)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Checking %d URLs with %d workers, timeout %s...\n", len(urls), *workers, *timeout)

	// Top-level context - cancel everything if main returns of if we add signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // ensure all goroutines are cleaned up when main exits

	pool := &checker.Pool{Workers: *workers, Timeout: *timeout}

	start := time.Now()
	results := pool.Run(ctx, urls)
	elapsed := time.Since(start)

	summary := report.Build(results, elapsed)
	report.Print(summary)

	if *output != "" {
		if err := report.SaveJson(summary, *output); err != nil {
			fmt.Printf("Error saving report: %v\n", err)
		} else {
			fmt.Printf("Report saved to %s\n", *output)
		}
	}
}

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var lines []string
	sc := bufio.NewScanner(f)

	for sc.Scan() {
		if line := sc.Text(); line != "" {
			lines = append(lines, line)
		}
	}
	return lines, sc.Err()
}
