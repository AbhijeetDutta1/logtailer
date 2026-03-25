# logtailer

A professional macOS log streaming utility built in Go. `logtailer` hooks directly into the **macOS Unified Logging System**, providing real-time streaming, structured parsing, and automated rotation with optional **Google Cloud Storage (GCS)** integration.

## Description
`logtailer` is designed to provide a more structured and modern way to interact with macOS system logs. Instead of simple file tailing, it leverages the native `log stream` architecture to capture system-wide events, parses them into structured JSON, and manages log lifecycle through rotation and cloud backups.

## Features
- **Real-time Unified Logging**: Streams directly from the macOS Unified Logging System using `--style syslog`.
- **Structured Parsing**: Uses high-performance regex to extract timestamps, hostnames, processes, PIDs, and message types.
- **Color-coded CLI**: High-visibility console output with red highlighting for `[Error]` and `[Fault]` levels.
- **Automatic Rotation**: Rotates logs to JSON batches based on size (1MB) or time (30 seconds).
- **GCS Integration**: Background uploading of rotated log batches to a specified Google Cloud Storage bucket.
- **Zero-Key Authentication**: Supports Application Default Credentials (ADC) for secure GCS uploads without service account keys.

## Tech Stack
- **Language**: Go 1.25+ (optimized for 1.26.1)
- **Primary Dependencies**: 
  - `cloud.google.com/go/storage`: Official GCS client library.
- **Native Integration**: `os/exec` hooks for macOS `log` utility.

## Installation Instructions

### Prerequisites
- macOS
- Go 1.13+ (Go 1.25+ recommended for full GCS support)

### Build
From the project root, run:
```bash
go build -o mac-log-tailer ./cmd/tailer/main.go
```

## Usage Examples

### Basic Streaming
Start streaming logs to the console and rotating them to the `logs/` directory:
```bash
./mac-log-tailer
```

### With GCS Uploads
Upload rotated JSON batches to a GCS bucket automatically:
```bash
./mac-log-tailer --bucket your-bucket-name
```

## Project Structure
```text
logtailer/
├── cmd/
│   └── tailer/
│       └── main.go           # CLI entry point and flag parsing
├── pkg/
│   ├── parser/
│   │   ├── parser.go         # Syslog regex parsing logic
│   │   ├── parser_test.go    # Unit tests for the parser
│   │   └── debug_parser.go   # Debugging utility for parser development
│   └── tailer/
│       ├── tailer.go         # macOS log stream execution
│       └── manager.go        # Batching, rotation, and GCS upload manager
├── logs/                     # Local storage for rotated JSON files
├── USAGE.txt                 # Quick start guide
├── verify.sh                 # Verification script for local testing
└── README.md                 # Project documentation
```

## How It Works (Architecture Overview)
1.  **Streamer (`pkg/tailer/tailer.go`)**: Spawns a `log stream --style syslog` subprocess and pipes the output into a Go channel.
2.  **Parser (`pkg/parser/parser.go`)**: Consumes the raw text stream, applying structured regex to identifies fields like `Timestamp`, `Process`, and `Type`.
3.  **Manager (`pkg/tailer/manager.go`)**:
    - Aggregates parsed entries into a memory buffer.
    - Monitors buffer size and elapsed time.
    - **Rotation**: Flushes the buffer to a timestamped JSON file in the `logs/` directory.
    - **Upload**: If a bucket is configured, triggers a background goroutine to upload the new JSON file to GCS using the official client library.

## Configuration
- **Command Line Flags**:
  - `--bucket`: (Optional) The name of the GCS bucket for log uploads.
- **Authentication**:
  - Uses `GOOGLE_APPLICATION_CREDENTIALS` environment variable.
  - Supports `gcloud auth application-default login` for local developer authentication.

## Development Instructions

### Running Tests
To run unit tests for the parser:
```bash
go test ./pkg/parser/...
```

### Verification Script
Use the provided `verify.sh` for a quick end-to-end check of the binary's execution and log generation.

## Future Improvements
- [ ] Add log filtering by process or priority level via CLI flags.
- [ ] Implement compressed (GZip) uploads to GCS to save bandwidth.
- [ ] Add support for alternative cloud providers (AWS S3, Azure Blob Storage).
- [ ] Enhance the console UI with a more interactive dashboard using a TUI library.
