package utils

import "strings"

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
