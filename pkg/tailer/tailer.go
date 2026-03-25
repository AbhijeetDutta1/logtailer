package tailer

import (
	"bufio"
	"fmt"
	"os/exec"
)

// Line represents a single line from the log stream.
type Line struct {
	Text  string
	Error error
}

// StreamLogs starts streaming from 'log stream --style syslog' and returns a channel of lines.
func StreamLogs() (<-chan Line, error) {
	cmd := exec.Command("log", "stream", "--style", "syslog")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start log stream: %v", err)
	}

	lineChan := make(chan Line)

	go func() {
		defer close(lineChan)
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			lineChan <- Line{Text: scanner.Text()}
		}
		if err := scanner.Err(); err != nil {
			lineChan <- Line{Error: err}
		}
		cmd.Wait()
	}()

	return lineChan, nil
}
