package reporter

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Yohannes022/go-website-monitor/internal/monitor"
)

func TestSaveJSON(t *testing.T) {

	directory := t.TempDir()

	outputPath := filepath.Join(
		directory,
		"reports",
		"report.json",
	)

	summary := monitor.Summary{
		Total:               2,
		Online:              2,
		Offline:             0,
		AvailabilityPercent: 100,
	}

	err := SaveJSON(outputPath, summary)

	if err != nil {
		t.Fatalf("Failed to save JSON: %v", err)
	}

	data, err := os.ReadFile(outputPath)

	if err != nil {
		t.Fatalf("Failed to read JSON: %v", err)
	}

	var result monitor.Summary

	err = json.Unmarshal(data, &result)

	if err != nil {
		t.Fatalf("Invalid JSON: %v", err)
	}

	if result.Total != 2 {
		t.Errorf(
			"Expected 2 total, got %d",
			result.Total,
		)
	}
}

func TestAppendHistory(t *testing.T) {

	directory := t.TempDir()

	historyPath := filepath.Join(
		directory,
		"history.jsonl",
	)

	summary := monitor.Summary{
		Total:               1,
		Online:              1,
		AvailabilityPercent: 100,
	}

	// First monitoring run.
	err := AppendHistory(historyPath, summary)

	if err != nil {
		t.Fatalf("Failed to append history: %v", err)
	}

	// Second monitoring run.
	err = AppendHistory(historyPath, summary)

	if err != nil {
		t.Fatalf("Failed to append second history: %v", err)
	}

	data, err := os.ReadFile(historyPath)

	if err != nil {
		t.Fatalf("Failed to read history: %v", err)
	}

	lines := 0

	for _, character := range string(data) {
		if character == '\n' {
			lines++
		}
	}

	if lines != 2 {
		t.Errorf(
			"Expected 2 history records, got %d",
			lines,
		)
	}
}
