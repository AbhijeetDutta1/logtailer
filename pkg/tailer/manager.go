package tailer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"logtailer/pkg/parser"
)

// LogManager handles batching log entries and rotating files based on size and time.
type LogManager struct {
	logsDir      string
	bucketName   string
	maxSize      int64
	maxDuration  time.Duration
	buffer       []*parser.LogEntry
	bufferSize   int64
	mu           sync.Mutex
	lastRotation time.Time
}

// NewLogManager creates a new LogManager.
func NewLogManager(logsDir string, bucketName string, maxSize int64, maxDuration time.Duration) *LogManager {
	return &LogManager{
		logsDir:      logsDir,
		bucketName:   bucketName,
		maxSize:      maxSize,
		maxDuration:  maxDuration,
		buffer:       make([]*parser.LogEntry, 0),
		lastRotation: time.Now(),
	}
}

// AddEntry adds a log entry to the buffer and checks if rotation is needed.
func (m *LogManager) AddEntry(entry *parser.LogEntry) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Estimate size of this entry in JSON
	entryBytes, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}

	m.buffer = append(m.buffer, entry)
	m.bufferSize += int64(len(entryBytes))

	if m.shouldRotate() {
		return m.rotate()
	}

	return nil
}

// shouldRotate checks if either size or time limit has been reached.
func (m *LogManager) shouldRotate() bool {
	if m.bufferSize >= m.maxSize {
		return true
	}
	if time.Since(m.lastRotation) >= m.maxDuration && len(m.buffer) > 0 {
		return true
	}
	return false
}

// rotate writes the current buffer to a JSON file and resets.
func (m *LogManager) rotate() error {
	if len(m.buffer) == 0 {
		m.lastRotation = time.Now()
		return nil
	}

	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(m.logsDir, fmt.Sprintf("logs_%s.json", timestamp))

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(m.buffer); err != nil {
		return err
	}
	file.Close() // Ensure file is closed before uploading

	// Upload to GCS if bucket is configured
	if m.bucketName != "" {
		go func(fname string) {
			if err := m.uploadToGCS(fname); err != nil {
				fmt.Fprintf(os.Stderr, "[LogManager] GCS Upload Error: %v\n", err)
			}
		}(filename)
	}

	// Reset buffer
	m.buffer = make([]*parser.LogEntry, 0)
	m.bufferSize = 0
	m.lastRotation = time.Now()

	fmt.Printf("[LogManager] Rotated logs to %s\n", filename)
	return nil
}

// StartTimer starts a background goroutine to trigger rotation based on time.
func (m *LogManager) StartTimer() {
	go func() {
		ticker := time.NewTicker(m.maxDuration / 2)
		for range ticker.C {
			m.mu.Lock()
			if m.shouldRotate() {
				m.rotate()
			}
			m.mu.Unlock()
		}
	}()
}

// uploadToGCS uploads a file to the configured Google Cloud Storage bucket.
func (m *LogManager) uploadToGCS(filename string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	objectName := filepath.Base(filename)
	wc := client.Bucket(m.bucketName).Object(objectName).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	fmt.Printf("[LogManager] Successfully uploaded %s to gs://%s/\n", objectName, m.bucketName)
	return nil
}
