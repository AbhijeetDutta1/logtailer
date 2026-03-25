package tailer

import (
	"fmt"
	"io"
	"os"

	"github.com/nxadm/tail"
)

// Line represents a single line from a log file.
type Line struct {
	Text  string
	Error error
}

// TailFile starts tailing a file and returns a channel of lines.
func TailFile(filePath string) (<-chan Line, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", filePath)
	}

	config := tail.Config{
		Follow:    true,
		ReOpen:    true,
		MustExist: true,
		Poll:      false,
		Location:  &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd},
	}

	t, err := tail.TailFile(filePath, config)
	if err != nil {
		return nil, err
	}

	lineChan := make(chan Line)

	go func() {
		defer close(lineChan)
		for line := range t.Lines {
			if line.Err != nil {
				lineChan <- Line{Error: line.Err}
				return
			}
			lineChan <- Line{Text: line.Text}
		}
	}()

	return lineChan, nil
}
