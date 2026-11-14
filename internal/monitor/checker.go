// This is the core of our application.

// It will:
// - Validate website URLs.
// - Send HTTP requests.
// - Measure response times.
// - Detect unavailable websites.
// - Run multiple checks concurrently.
// - Calculate monitoring statistics.

package monitor

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Result represents the result of checking one website.
type Result struct {
	URL            string    `json:"url"`
	StatusCode     int       `json:"status_code"`
	ResponseTimeMS float64   `json:"response_time_ms"`
	IsUp           bool      `json:"is_up"`
	CheckedAt      time.Time `json:"checked_at"`
	Error          string    `json:"error,omitempty"`
}

// Summary contains statistics for a monitoring run.
type Summary struct {
	Total                 int      `json:"total"`
	Online                int      `json:"online"`
	Offline               int      `json:"offline"`
	AvailabilityPercent   float64  `json:"availability_percent"`
	AverageResponseTimeMS float64  `json:"average_response_time_ms"`
	Results               []Result `json:"results"`
}

// CheckURL checks the availability of a single website.
func CheckURL(
	ctx context.Context,
	client *http.Client,
	rawURL string,
) Result {

	result := Result{
		URL:       rawURL,
		CheckedAt: time.Now().UTC(),
	}

	// Validate URL.
	parsedURL, err := url.ParseRequestURI(rawURL)

	if err != nil ||
		parsedURL.Host == "" ||
		(parsedURL.Scheme != "http" &&
			parsedURL.Scheme != "https") {

		result.Error = "Invalid URL"
		return result
	}

	// Create HTTP request.
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		rawURL,
		nil,
	)

	if err != nil {
		result.Error = err.Error()
		return result
	}

	// Measure response time.
	start := time.Now()

	response, err := client.Do(request)

	result.ResponseTimeMS = float64(
		time.Since(start).Microseconds(),
	) / 1000

	if err != nil {
		result.Error = err.Error()
		return result
	}

	defer response.Body.Close()

	result.StatusCode = response.StatusCode

	// Treat HTTP 200-399 as healthy.
	result.IsUp = response.StatusCode >= 200 &&
		response.StatusCode < 400

	return result
}

// CheckURLs checks multiple websites concurrently.
func CheckURLs(
	ctx context.Context,
	urls []string,
	workers int,
	timeout time.Duration,
) []Result {

	if workers < 1 {
		workers = 1
	}

	if workers > len(urls) {
		workers = len(urls)
	}

	results := make([]Result, len(urls))

	if len(urls) == 0 {
		return results
	}

	client := &http.Client{
		Timeout: timeout,
	}

	type job struct {
		index int
		url   string
	}

	jobs := make(chan job)

	var wg sync.WaitGroup

	// Start worker goroutines.
	for i := 0; i < workers; i++ {

		wg.Add(1)

		go func() {

			defer wg.Done()

			for currentJob := range jobs {

				results[currentJob.index] = CheckURL(
					ctx,
					client,
					currentJob.url,
				)
			}
		}()
	}

	// Send jobs to workers.
	for index, website := range urls {

		jobs <- job{
			index: index,
			url:   website,
		}
	}

	close(jobs)

	// Wait for all workers.
	wg.Wait()

	return results
}

// BuildSummary calculates overall monitoring statistics.
func BuildSummary(results []Result) Summary {

	summary := Summary{
		Total:   len(results),
		Results: results,
	}

	var totalResponseTime float64

	for _, result := range results {

		if result.IsUp {
			summary.Online++
		} else {
			summary.Offline++
		}

		totalResponseTime += result.ResponseTimeMS
	}

	if summary.Total > 0 {

		summary.AvailabilityPercent =
			float64(summary.Online) /
				float64(summary.Total) * 100

		summary.AverageResponseTimeMS =
			totalResponseTime / float64(summary.Total)
	}

	return summary
}

// ValidateTimeout validates the configured timeout.
func ValidateTimeout(timeout time.Duration) error {

	if timeout <= 0 {
		return fmt.Errorf(
			"timeout must be greater than zero",
		)
	}

	return nil
}
