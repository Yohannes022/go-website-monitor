// This will create automated tests to make sure the monitoring engine
// and reporting system work correctly.

// Go has a built-in testing framework, so we don't need to install pytest or
// another testing library.

package monitor

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"math"
)

func TestCheckURLSuccess(t *testing.T) {

	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	defer server.Close()

	client := server.Client()

	result := CheckURL(
		context.Background(),
		client,
		server.URL,
	)

	if !result.IsUp {
		t.Error("Expected website to be online")
	}

	if result.StatusCode != http.StatusOK {
		t.Errorf(
			"Expected status 200, got %d",
			result.StatusCode,
		)
	}

	if result.ResponseTimeMS < 0 {
		t.Error("Response time cannot be negative")
	}

	if result.CheckedAt.IsZero() {
		t.Error("Expected a check timestamp")
	}
}

func TestCheckURLServerError(t *testing.T) {

	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)

	defer server.Close()

	result := CheckURL(
		context.Background(),
		server.Client(),
		server.URL,
	)

	if result.IsUp {
		t.Error("Expected HTTP 500 to be unhealthy")
	}

	if result.StatusCode != http.StatusInternalServerError {
		t.Errorf(
			"Expected status 500, got %d",
			result.StatusCode,
		)
	}
}

func TestCheckURLInvalidURL(t *testing.T) {

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	result := CheckURL(
		context.Background(),
		client,
		"not-a-valid-url",
	)

	if result.IsUp {
		t.Error("Invalid URL should not be online")
	}

	if result.Error == "" {
		t.Error("Expected an error message")
	}
}

func TestCheckURLs(t *testing.T) {

	server := httptest.NewServer(
		http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			if r.URL.Path == "/error" {
				w.WriteHeader(
					http.StatusInternalServerError,
				)
				return
			}

			w.WriteHeader(http.StatusOK)
		}),
	)

	defer server.Close()

	urls := []string{
		server.URL + "/success",
		server.URL + "/error",
		"invalid-url",
	}

	results := CheckURLs(
		context.Background(),
		urls,
		2,
		5*time.Second,
	)

	if len(results) != 3 {
		t.Fatalf(
			"Expected 3 results, got %d",
			len(results),
		)
	}

	if !results[0].IsUp {
		t.Error("First website should be online")
	}

	if results[1].IsUp {
		t.Error("Second website should be unhealthy")
	}

	if results[2].Error == "" {
		t.Error("Third result should contain an error")
	}
}

func TestBuildSummary(t *testing.T) {

	results := []Result{
		{
			IsUp:           true,
			ResponseTimeMS: 100,
		},
		{
			IsUp:           true,
			ResponseTimeMS: 200,
		},
		{
			IsUp:           false,
			ResponseTimeMS: 300,
		},
	}

	summary := BuildSummary(results)

	if summary.Total != 3 {
		t.Errorf(
			"Expected 3 total, got %d",
			summary.Total,
		)
	}

	if summary.Online != 2 {
		t.Errorf(
			"Expected 2 online, got %d",
			summary.Online,
		)
	}

	if summary.Offline != 1 {
		t.Errorf(
			"Expected 1 offline, got %d",
			summary.Offline,
		)
	}

	if math.Abs(summary.AvailabilityPercent-66.66666666666666) > 0.0001 {
		t.Errorf(
			"Unexpected availability: %f",
			summary.AvailabilityPercent,
		)
	}

	if summary.AverageResponseTimeMS != 200 {
		t.Errorf(
			"Expected average 200ms, got %f",
			summary.AverageResponseTimeMS,
		)
	}
}

func TestValidateTimeout(t *testing.T) {

	if err := ValidateTimeout(5 * time.Second); err != nil {
		t.Errorf("Expected valid timeout, got %v", err)
	}

	if err := ValidateTimeout(0); err == nil {
		t.Error("Expected zero timeout to be rejected")
	}

	if err := ValidateTimeout(-1); err == nil {
		t.Error("Expected negative timeout to be rejected")
	}
}
