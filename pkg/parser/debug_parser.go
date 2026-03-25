package parser

import (
	"fmt"
)

func main() {
	line := "2026-03-24 19:50:20.828902-0700  localhost sharingd[459]: (CoreUtils) [com.apple.sharing:SDNearbyAgentCore] NearbyInfo received activity level: 0x3 after decryption"
	entry, ok := ParseSyslogLine(line)
	if !ok {
		fmt.Println("FAILED to parse")
	} else {
		fmt.Printf("SUCCESS: %+v\n", entry)
	}
}
