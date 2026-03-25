package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mac-log-tailer/pkg/parser"
	"mac-log-tailer/pkg/tailer"
)

const (
	ColorRed    = "\033[31m"
	ColorReset  = "\033[0m"
	MaxLogSize  = 1 * 1024 * 1024 // 1MB
	MaxDuration = 30 * time.Second
)

func main() {
	fmt.Println("--- Streaming macOS Unified Logging System ---")
	fmt.Printf("Batch rotation: %d bytes or %v\n", MaxLogSize, MaxDuration)

	// Ensure logs directory exists
	if err := os.MkdirAll("logs", 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating logs directory: %v\n", err)
		os.Exit(1)
	}

	manager := tailer.NewLogManager("logs", MaxLogSize, MaxDuration)
	manager.StartTimer()

	// Channel to handle interrupt signals for clean shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	lineChan, err := tailer.StreamLogs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	for {
		select {
		case line, ok := <-lineChan:
			if !ok {
				fmt.Println("Streaming stopped.")
				return
			}
			if line.Error != nil {
				fmt.Fprintf(os.Stderr, "Streaming error: %v\n", line.Error)
				return
			}

			// Parse the line
			entry, parsed := parser.ParseSyslogLine(line.Text)
			if !parsed {
				// If not parsed, just print as is
				fmt.Println(line.Text)
				continue
			}

			// Add to manager for batching/rotation
			if err := manager.AddEntry(entry); err != nil {
				fmt.Fprintf(os.Stderr, "Manager error: %v\n", err)
			}

			// Console display with color-coding
			displayLine := line.Text
			if entry.IsError() {
				displayLine = ColorRed + line.Text + ColorReset
			}
			fmt.Println(displayLine)

		case <-sigChan:
			fmt.Println("\nStopping stream and flushing logs...")
			return
		}
	}
}
