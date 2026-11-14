// This module handles report generation and monitoring history.

// We'll produce two output formats:
// 	- quality_report.json equivalent: a JSON summary for the latest run.
// 	- history.jsonl: an append-only monitoring history file.

// JSONL means JSON Lines: every line contains one independent JSON object.

package reporter

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Yohannes022/go-website-monitor/internal/monitor"
)

// SaveJSON saves a monitoring summary as a JSON file.
func SaveJSON(
	outputPath string,
	summary monitor.Summary,
) error {

	// Create parent directories if necessary.
	directory := filepath.Dir(outputPath)

	if err := os.MkdirAll(directory, 0755); err != nil {
		return err
	}

	// Convert summary into formatted JSON.
	data, err := json.MarshalIndent(
		summary,
		"",
		"    ",
	)

	if err != nil {
		return err
	}

	// Write JSON report.
	return os.WriteFile(
		outputPath,
		data,
		0644,
	)
}

// AppendHistory appends one monitoring run to a JSONL file.
func AppendHistory(
	outputPath string,
	summary monitor.Summary,
) error {

	directory := filepath.Dir(outputPath)

	if err := os.MkdirAll(directory, 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(
		outputPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)

	return encoder.Encode(summary)
}
