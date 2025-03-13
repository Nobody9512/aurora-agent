package cmd

import (
	"fmt"
	"os"
	"strings"
)

// Global agent manager instance
var AgentMgr *AgentManager

func init() {
	// Initialize the agent manager
	AgentMgr = NewAgentManager()
}

// ProcessAuroraCommand handles Aurora-specific commands
func ProcessAuroraCommand(input string) bool {
	// Check if input contains "aurora"
	if strings.Contains(strings.ToLower(input), "aurora") {
		// Use streaming response
		err := AgentMgr.StreamQuery(input, os.Stdout)
		if err != nil {
			fmt.Printf("\nError querying AI agent: %v\n", err)
		} else {
			fmt.Println() // Add a newline after the streamed response
		}

		return true
	}
	return false
}

// SetAIAgent sets the active AI agent
func SetAIAgent(agentType string) error {
	return AgentMgr.SetActiveAgent(AgentType(agentType))
}
