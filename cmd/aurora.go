// Package cmd contains command functionality for the Aurora Agent
//
// This file is the main entry point for Aurora command processing.
// All core functionality is implemented in separate files:
//
// - config_commands.go: Configuration command processors
// - display_config.go: Display configuration functions
// - help_commands.go: Help system commands
// - shell_command_utils.go: Shell command utilities
package cmd

import (
	"fmt"
)

// Global agent manager instance
var AgentMgr *AgentManager

func init() {
	// Initialize the agent manager
	AgentMgr = NewAgentManager()
}

// ProcessAuroraCommand handles Aurora-specific commands
func ProcessAuroraCommand(input string) bool {
	// Check configuration commands
	if processConfigCommand(input) {
		return true
	}

	// Check if input contains "aurora" or is not a shell command
	if isAuroraCommand(input) || !isShellCommand(input) {
		// Use streaming response
		fmt.Print("\n") // Add a newline before the response for better readability

		// Print a colored prompt to indicate AI response
		fmt.Print("\033[36mAurora: \033[0m") // Cyan color for Aurora name

		// Use function calls for natural language processing
		err := AgentMgr.StreamQueryWithFunctionCalls(input)
		if err != nil {
			fmt.Printf("\n\033[31mError querying AI agent: %v\033[0m\n", err) // Red error message
		}

		// No need to add a newline here as StreamQueryWithFunctionCalls now adds one

		return true
	}
	return false
}

// SetAIAgent sets the active AI agent
func SetAIAgent(agentType string) error {
	return AgentMgr.SetActiveAgent(AgentType(agentType))
}
