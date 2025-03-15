package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"aurora-agent/config"
)

// Global agent manager instance
var AgentMgr *AgentManager

func init() {
	// Initialize the agent manager
	AgentMgr = NewAgentManager()
}

// ProcessAuroraCommand handles Aurora-specific commands
func ProcessAuroraCommand(input string) bool {
	// Check if input contains "aurora" or is not a shell command
	if isAuroraCommand(input) || !isShellCommand(input) {
		// Use streaming response
		fmt.Print("\n") // Add a newline before the response for better readability

		// Print a colored prompt to indicate AI response
		fmt.Print("\033[36mAurora: \033[0m") // Cyan color for Aurora name

		// Use function calls for natural language processing
		err := AgentMgr.StreamQueryWithFunctionCalls(input, os.Stdout)
		if err != nil {
			fmt.Printf("\n\033[31mError querying AI agent: %v\033[0m\n", err) // Red error message
		} else {
			fmt.Println() // Add a newline after the streamed response
		}

		return true
	}
	return false
}

func isShellCommand(input string) bool {
	input = strings.ToLower(strings.TrimSpace(input)) // Trim and convert to lowercase
	words := strings.Fields(input)                    // Split into words

	if len(words) == 0 {
		return false
	}

	if slices.Contains(config.ShellCommands, words[0]) {
		return true
	}

	return false
}

func isAuroraCommand(input string) bool {
	return strings.Contains(strings.ToLower(input), "aurora")
}

// SetAIAgent sets the active AI agent
func SetAIAgent(agentType string) error {
	return AgentMgr.SetActiveAgent(AgentType(agentType))
}