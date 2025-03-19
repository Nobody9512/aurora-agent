package cmd

import (
	"fmt"

	"aurora-agent/config"
)

// showConfig - show config information
func showConfig() {
	fmt.Println("\n\033[1mCurrent configuration:\033[0m")
	fmt.Println("\n\033[1m[General]\033[0m")
	fmt.Printf("  DefaultShell: %s\n", config.CurrentConfig.General.DefaultShell)
	fmt.Printf("  HistorySize: %d\n", config.CurrentConfig.General.HistorySize)
	fmt.Printf("  ShellCommands: %d additional commands\n", len(config.CurrentConfig.General.ShellCommands))
	fmt.Printf("  IgnoredCommands: %d ignored commands\n", len(config.CurrentConfig.General.IgnoredCommands))

	fmt.Println("\033[1m[OpenAI]\033[0m")
	apiKey := config.CurrentConfig.OpenAI.APIKey
	if apiKey != "" {
		apiKey = "********" // Hide API key
	}
	fmt.Printf("  APIKey: %s\n", apiKey)
	fmt.Printf("  Model: %s\n", config.CurrentConfig.OpenAI.Model)

	fmt.Println("\033[1m[Interface]\033[0m")
	fmt.Printf("  Theme: %s\n", config.CurrentConfig.Interface.Theme)

	// System prompt information
	systemPrompt := config.CurrentConfig.Interface.SystemPrompt
	if systemPrompt == "default" || systemPrompt == "" {
		fmt.Printf("  SystemPrompt: default\n")
	} else {
		// If system prompt is long, shorten it
		if len(systemPrompt) > 50 {
			fmt.Printf("  SystemPrompt: \"%s...\"\n", systemPrompt[:47])
		} else {
			fmt.Printf("  SystemPrompt: \"%s\"\n", systemPrompt)
		}
	}

	fmt.Printf("\nConfiguration file: \033[32m%s\033[0m\n", config.GetConfigPath())
	fmt.Println("\nTo see the commands list, use `\033[32mconfig commands list\033[0m`")
	fmt.Println()
}
