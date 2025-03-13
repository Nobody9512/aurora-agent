package cmd

import (
	"fmt"
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

		response, err := AgentMgr.Query(input)
		if err != nil {
			fmt.Printf("Error querying AI agent: %v\n", err)
		} else {
			fmt.Printf("AI Agent (%s): %s\n", AgentMgr.GetActiveAgentName(), response)
		}

		return true
	}
	return false
}

// SetAIAgent sets the active AI agent
func SetAIAgent(agentType string) error {
	return AgentMgr.SetActiveAgent(AgentType(agentType))
}
