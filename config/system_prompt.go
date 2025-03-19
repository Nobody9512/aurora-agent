package config

import (
	"strings"
)

// GetSystemPrompt - user system prompt based system prompt
func GetSystemPrompt() string {
	// If user system prompt exists, add it
	if CurrentConfig.Interface.SystemPrompt != "default" && CurrentConfig.Interface.SystemPrompt != "" {
		// Add user system prompt to system prompt
		// Do not change {{USER_INPUT}}
		customPrompt := strings.Replace(DefaultSystemPrompt, "{{USER_INPUT}}",
			CurrentConfig.Interface.SystemPrompt+"\n\n{{USER_INPUT}}", 1)
		return customPrompt
	}

	return strings.Replace(DefaultSystemPrompt, "{{USER_INPUT}}", "", 1)
}
