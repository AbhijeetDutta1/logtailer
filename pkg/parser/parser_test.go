package parser

import (
	"testing"
)

func TestParseSyslogLine(t *testing.T) {
	line := "2026-03-24 19:14:10.871021-0700  localhost kernel[0]: [Error] some error message"
	entry, ok := ParseSyslogLine(line)
	if !ok {
		t.Fatal("Failed to parse valid syslog line")
	}

	if entry.Timestamp != "2026-03-24 19:14:10.871021-0700" {
		t.Errorf("Expected timestamp %s, got %s", "2026-03-24 19:14:10.871021-0700", entry.Timestamp)
	}
	if entry.Hostname != "localhost" {
		t.Errorf("Expected hostname localhost, got %s", entry.Hostname)
	}
	if entry.Process != "kernel" {
		t.Errorf("Expected process kernel, got %s", entry.Process)
	}
	if entry.PID != "0" {
		t.Errorf("Expected PID 0, got %s", entry.PID)
	}
	if entry.Type != "Error" {
		t.Errorf("Expected type Error, got %s", entry.Type)
	}
	if entry.Message != "some error message" {
		t.Errorf("Expected message 'some error message', got %s", entry.Message)
	}
}

func TestParseSyslogLineWithoutType(t *testing.T) {
	line := "2026-03-24 19:14:10.871021-0700  localhost process[123]: just a message"
	entry, ok := ParseSyslogLine(line)
	if !ok {
		t.Fatal("Failed to parse valid syslog line without type")
	}

	if entry.Type != "Default" {
		t.Errorf("Expected type Default, got %s", entry.Type)
	}
	if entry.Message != "just a message" {
		t.Errorf("Expected message 'just a message', got %s", entry.Message)
	}
}

func TestParseLiveSample(t *testing.T) {
	line := "2026-03-24 19:50:20.828902-0700  localhost sharingd[459]: (CoreUtils) [com.apple.sharing:SDNearbyAgentCore] NearbyInfo received activity level: 0x3 after decryption"
	entry, ok := ParseSyslogLine(line)
	if !ok {
		t.Fatal("Failed to parse live sample syslog line")
	}

	if entry.Process != "sharingd" {
		t.Errorf("Expected process sharingd, got %s", entry.Process)
	}
	if entry.Message != "(CoreUtils) [com.apple.sharing:SDNearbyAgentCore] NearbyInfo received activity level: 0x3 after decryption" {
		t.Errorf("Expected message to include CoreUtils, got %s", entry.Message)
	}
}
