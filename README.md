
# Go Website Uptime Monitor

A lightweight command-line website monitoring tool written in Go.

The application checks multiple websites concurrently,
measures HTTP response times, detects unavailable websites,
and generates JSON monitoring reports.

## Features

- Monitor multiple websites
- HTTP status code detection
- Response time measurement
- Concurrent website monitoring
- Configurable request timeout
- Configurable worker count
- Invalid URL detection
- JSON report generation
- Persistent monitoring history
- Availability statistics
- Automated tests

## Technologies

- Go
- HTTP
- Goroutines
- Channels
- Context
- JSON
- Go testing package

## Project Structure

```text
cmd/
  uptime-monitor/
    main.go

internal/
  monitor/
    checker.go
    checker_test.go

  reporter/
    reporter.go
    reporter_test.go

reports/
```

## Requirements

- Go 1.23 or later

## Installation

Clone the repository:

```bash
git clone YOUR_REPOSITORY_URL
cd go-website-monitor
```

Run the application:

```bash
go run ./cmd/uptime-monitor https://google.com
```

## Usage

Check multiple websites:

```bash
go run ./cmd/uptime-monitor \
  https://google.com \
  https://github.com \
  https://stackoverflow.com
```

Configure timeout:

```bash
go run ./cmd/uptime-monitor \
  -timeout=5s \
  https://google.com
```

Configure concurrent workers:

```bash
go run ./cmd/uptime-monitor \
  -workers=10 \
  https://google.com \
  https://github.com
```

Customize report location:

```bash
go run ./cmd/uptime-monitor \
  -output=reports/custom.json \
  https://google.com
```

Customize monitoring history:

```bash
go run ./cmd/uptime-monitor \
  -history=reports/custom-history.jsonl \
  https://google.com
```

View available options:

```bash
go run ./cmd/uptime-monitor -h
```

## Testing

Run all tests:

```bash
go test ./...
```

Run tests with verbose output:

```bash
go test -v ./...
```

Run the race detector:

```bash
go test -race ./...
```

## Build

Build the executable:

```bash
go build -o uptime-monitor ./cmd/uptime-monitor
```

Run the executable:

```bash
./uptime-monitor https://google.com
```

## Monitoring Logic

Websites returning HTTP status codes between
200 and 399 are considered healthy.

Other HTTP responses are considered unhealthy.

Response times are measured in milliseconds.

The default timeout is 10 seconds.

The default worker count is 5.

## Reports

The application generates:

- `reports/latest.json`
- `reports/history.jsonl`

The JSON report contains the latest monitoring summary.

The JSONL file stores previous monitoring runs.

## Limitations

- Monitoring is initiated manually.
- Historical reports are stored locally.
- Availability is based on individual HTTP checks.
- A single monitoring run does not represent long-term uptime.
- The tool does not yet support scheduled monitoring.

## License

MIT License.
