package parser

import (
	"regexp"
	"strings"
)

// LogEntry represents a structured log entry from the macOS Unified Logging System.
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Hostname  string `json:"hostname"`
	Process   string `json:"process"`
	PID       string `json:"pid"`
	Type      string `json:"type"`
	Message   string `json:"message"`
}

// SyslogRegex matches the standard macOS 'log stream --style syslog' format.
// Format: 2026-03-24 19:14:10.871021-0700  localhost kernel[0]: [Error] some message
var syslogRegex = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}\.\d{6}-\d{4})\s+(\S+)\s+(.+?)\[(\d+)\]:\s+(?:\[(.*?)\]\s+)?(.*)$`)

// ParseSyslogLine parses a single line of syslog output into a LogEntry.
func ParseSyslogLine(line string) (*LogEntry, bool) {
	matches := syslogRegex.FindStringSubmatch(line)
	if matches == nil {
		return nil, false
	}

	entry := &LogEntry{
		Timestamp: matches[1],
		Hostname:  matches[2],
		Process:   matches[3],
		PID:       matches[4],
		Type:      matches[5],
		Message:   matches[6],
	}

	// Default type to "Default" if not found
	if entry.Type == "" {
		entry.Type = "Default"
	}

	return entry, true
}

// IsError checks if the log entry is an error level.
func (e *LogEntry) IsError() bool {
	t := strings.ToLower(e.Type)
	return t == "error" || t == "fault" || strings.Contains(strings.ToLower(e.Message), "error")
}
