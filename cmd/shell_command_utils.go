package cmd

import (
	"fmt"
	"strings"

	"aurora-agent/config"
)

// Shell command utilities

// isShellCommand checks if input is a recognized shell command
func isShellCommand(input string) bool {
	input = strings.ToLower(strings.TrimSpace(input)) // Trim and convert to lowercase
	words := strings.Fields(input)                    // Split into words

	if len(words) == 0 {
		return false
	}

	commands := config.GetShellCommands()
	for _, cmd := range commands {
		if cmd == words[0] {
			return true
		}
	}

	return false
}

// isAuroraCommand checks if input contains "aurora"
func isAuroraCommand(input string) bool {
	return strings.Contains(strings.ToLower(input), "aurora")
}

// showShellCommands displays all available shell commands
func showShellCommands() {
	commands := config.GetShellCommands()

	fmt.Println("\033[1mCommands list:\033[0m")

	// Show commands in 5 columns
	cols := 5
	for i, cmd := range commands {
		if i > 0 && i%cols == 0 {
			fmt.Println()
		}
		fmt.Printf("%-15s", cmd)
	}
	fmt.Println()

	// Additional commands
	if len(config.CurrentConfig.General.ShellCommands) > 0 {
		fmt.Println("\n\033[1mUser added commands:\033[0m")
		for _, cmd := range config.CurrentConfig.General.ShellCommands {
			fmt.Printf("  %s\n", cmd)
		}
	}

	// Ignored commands
	if len(config.CurrentConfig.General.IgnoredCommands) > 0 {
		fmt.Println("\n\033[1mIgnored commands:\033[0m")
		for _, cmd := range config.CurrentConfig.General.IgnoredCommands {
			fmt.Printf("  %s\n", cmd)
		}
	}

	fmt.Printf("\nTotal commands: %d\n", len(commands))
}

// addShellCommand adds a new command to the shell command list
func addShellCommand(command string) {
	if err := config.AddShellCommand(command); err != nil {
		fmt.Printf("\033[31mError: %v\033[0m\n", err)
		return
	}

	fmt.Printf("\033[32m'%s' command added successfully\033[0m\n", command)
	fmt.Println("\033[33mNote: Remember to save changes using 'config save'\033[0m")
}

// removeShellCommand removes a command from the shell command list
func removeShellCommand(command string) {
	if err := config.RemoveShellCommand(command); err != nil {
		fmt.Printf("\033[31mError: %v\033[0m\n", err)
		return
	}

	fmt.Printf("\033[32m'%s' command removed successfully\033[0m\n", command)
	fmt.Println("\033[33mNote: Remember to save changes using 'config save'\033[0m")
}

// resetShellCommands resets shell commands list to default
func resetShellCommands() {
	config.ResetShellCommands()
	fmt.Println("\033[32mCommands list reset to default\033[0m")
	fmt.Println("\033[33mNote: Remember to save changes using 'config save'\033[0m")
}
