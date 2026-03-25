package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"mac-log-tailer/pkg/tailer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <file_path>\n", os.Args[0])
		os.Exit(1)
	}

	filePath := os.Args[1]
	fmt.Printf("--- Tailing %s ---\n", filePath)

	// Channel to handle interrupt signals for clean shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	lineChan, err := tailer.TailFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	for {
		select {
		case line, ok := <-lineChan:
			if !ok {
				fmt.Println("Tailing stopped.")
				return
			}
			if line.Error != nil {
				fmt.Fprintf(os.Stderr, "Tailing error: %v\n", line.Error)
				return
			}
			fmt.Println(line.Text)
		case <-sigChan:
			fmt.Println("\nStopping tailer...")
			return
		}
	}
}
