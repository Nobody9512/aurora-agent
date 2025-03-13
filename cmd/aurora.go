package cmd

import (
	"fmt"
	"strings"
)

// ProcessAuroraCommand handles Aurora-specific commands
func ProcessAuroraCommand(input string) bool {
	// Check if input contains "aurora"
	if strings.Contains(strings.ToLower(input), "aurora") {
		// TODO: Aurora is processing...
		fmt.Println("Aurora is processing...")
		return true
	}
	return false
}
