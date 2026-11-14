// This is the main application.

// It connects the monitoring engine and reporting module.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Yohannes022/go-website-monitor/internal/monitor"
	"github.com/Yohannes022/go-website-monitor/internal/reporter"
)

func main() {

	// Define command-line flags.
	timeout := flag.Duration(
		"timeout",
		10*time.Second,
		"Maximum time allowed for each website check",
	)

	workers := flag.Int(
		"workers",
		5,
		"Number of concurrent workers",
	)

	output := flag.String(
		"output",
		"reports/latest.json",
		"Path to the JSON report",
	)

	history := flag.String(
		"history",
		"reports/history.jsonl",
		"Path to monitoring history",
	)

	flag.Parse()

	// Get website URLs.
	urls := flag.Args()

	if len(urls) == 0 {

		fmt.Println("GO WEBSITE UPTIME MONITOR")
		fmt.Println("========================================")
		fmt.Println("Usage:")
		fmt.Println(
			"uptime-monitor [flags] <URL1> <URL2> ...",
		)
		fmt.Println("\nExample:")
		fmt.Println(
			"uptime-monitor https://google.com https://github.com",
		)

		flag.PrintDefaults()
		os.Exit(1)
	}

	// Validate timeout.
	if err := monitor.ValidateTimeout(*timeout); err != nil {
		fmt.Fprintf(os.Stderr, "\nError: %v\n", err)
		os.Exit(1)
	}

	if *workers < 1 {
		fmt.Fprintln(
			os.Stderr,
			"\nError: workers must be at least 1",
		)
		os.Exit(1)
	}

	// Start monitoring.
	fmt.Println("\nGO WEBSITE UPTIME MONITOR")
	fmt.Println("========================================")

	fmt.Printf(
		"\nChecking %d websites...\n\n",
		len(urls),
	)

	results := monitor.CheckURLs(
		context.Background(),
		urls,
		*workers,
		*timeout,
	)

	// Build monitoring summary.
	summary := monitor.BuildSummary(results)

	// Display individual results.
	for _, result := range summary.Results {

		if result.IsUp {

			fmt.Printf(
				"[UP]   %s\n",
				result.URL,
			)

		} else {

			fmt.Printf(
				"[DOWN] %s\n",
				result.URL,
			)
		}

		if result.StatusCode > 0 {

			fmt.Printf(
				"       Status: %d\n",
				result.StatusCode,
			)

		} else {

			fmt.Printf(
				"       Status: unavailable\n",
			)
		}

		fmt.Printf(
			"       Response: %.2f ms\n",
			result.ResponseTimeMS,
		)

		if result.Error != "" {

			fmt.Printf(
				"       Error: %s\n",
				result.Error,
			)
		}

		fmt.Println()
	}

	// Display summary.
	fmt.Println("========================================")
	fmt.Println("MONITORING SUMMARY")
	fmt.Println()

	fmt.Printf("Total websites: %d\n", summary.Total)
	fmt.Printf("Online: %d\n", summary.Online)
	fmt.Printf("Offline: %d\n", summary.Offline)

	fmt.Printf(
		"Availability: %.2f%%\n",
		summary.AvailabilityPercent,
	)

	fmt.Printf(
		"Average response time: %.2f ms\n",
		summary.AverageResponseTimeMS,
	)

	// Save JSON report.
	if err := reporter.SaveJSON(
		*output,
		summary,
	); err != nil {

		fmt.Fprintf(
			os.Stderr,
			"\nError saving report: %v\n",
			err,
		)

		os.Exit(1)
	}

	fmt.Printf(
		"\nReport saved: %s\n",
		*output,
	)

	// Save monitoring history.
	if err := reporter.AppendHistory(
		*history,
		summary,
	); err != nil {

		fmt.Fprintf(
			os.Stderr,
			"\nError saving history: %v\n",
			err,
		)

		os.Exit(1)
	}

	fmt.Printf(
		"History updated: %s\n",
		*history,
	)

	fmt.Println("========================================")
}
