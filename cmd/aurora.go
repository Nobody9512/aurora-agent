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
		// Extract the query part (everything after "aurora")
		query := extractQuery(input)

		// Use the extracted query if available, otherwise use the full input
		promptToSend := input
		if query != "" {
			promptToSend = query
		}

		response, err := AgentMgr.Query(promptToSend)
		if err != nil {
			fmt.Printf("Error querying AI agent: %v\n", err)
		} else {
			fmt.Println(response)
		}

		return true
	}
	return false
}

// extractQuery extracts the query part from the input
// For example, "aurora what is the weather" -> "what is the weather"
func extractQuery(input string) string {
	// Convert to lowercase for case-insensitive matching
	lowerInput := strings.ToLower(input)

	// Find the position of "aurora"
	auroraPos := strings.Index(lowerInput, "aurora")
	if auroraPos == -1 {
		return ""
	}

	// Extract everything after "aurora" + its length
	afterAurora := input[auroraPos+len("aurora"):]

	// Trim spaces
	return strings.TrimSpace(afterAurora)
}

// SetAIAgent sets the active AI agent
func SetAIAgent(agentType string) error {
	return AgentMgr.SetActiveAgent(AgentType(agentType))
}
