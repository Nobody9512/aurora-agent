package utils

import (
	"aurora-agent/config"
	"fmt"
	"strings"
)

// ProcessANSICodes processes ANSI codes in a string
func ProcessANSICodes(text string) string {
	result := strings.ReplaceAll(text, "\\033", "\033")
	result = strings.ReplaceAll(result, "\\\\033", "\033")
	result = strings.ReplaceAll(result, "\\u001b", "\033")
	result = strings.ReplaceAll(result, "\\e", "\033")
	result = strings.ReplaceAll(result, "\\x1b", "\033")
	result = strings.ReplaceAll(result, "\\x1B", "\033")
	result = strings.ReplaceAll(result, "\\u001B", "\033")
	result = strings.ReplaceAll(result, "\\27", "\033")
	result = strings.ReplaceAll(result, "\\33", "\033")

	return result
}

// processAnsiBuffer processes ANSI codes in the buffer and prints them
func ProcessAnsiBuffer(ansiBuffer string) string {
	if config.AnsiPattern.MatchString(ansiBuffer) {
		// Buffer has ANSI code, process it
		processedBuffer := ProcessANSICodes(ansiBuffer)
		fmt.Print(processedBuffer)
		return ""
	} else if config.AnsiStartPattern.MatchString(ansiBuffer) && len(ansiBuffer) > 100 {
		// If buffer contains the start of an ANSI code, but not the end
		// and buffer length is more than 100, process it
		// This can happen when ANSI code is in incorrect format
		processedBuffer := ProcessANSICodes(ansiBuffer)
		fmt.Print(processedBuffer)
		return ""
	} else if len(ansiBuffer) > 80 && !config.AnsiStartPattern.MatchString(ansiBuffer) {
		// If buffer length is more than 80 and no ANSI code start is found,
		// process it
		processedBuffer := ProcessANSICodes(ansiBuffer)
		fmt.Print(processedBuffer)
		return ""
	}

	return ansiBuffer
}
