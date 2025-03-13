package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

// Global agent manager instance
var AgentMgr *AgentManager

func init() {
	// Initialize the agent manager
	AgentMgr = NewAgentManager()
}

var shellCommands = []string{
	"ls", "cd", "pwd", "echo", "rm",
	"mkdir", "rmdir", "mv", "cp", "touch",
	"cat", "head", "tail", "grep", "find",
	"chmod", "chown", "tar", "zip", "unzip",
	"git", "docker", "kubectl", "python", "go",
	"htop", "ping", "kill", "jobs", "fg", "bg",
	"whoami", "uname", "df", "du", "ps", "top",
	"clear", "help", "exit", "quit", "cls", "clr",
}

// ProcessAuroraCommand handles Aurora-specific commands
func ProcessAuroraCommand(input string) bool {
	// Check if input contains "aurora"
	if isAuroraCommand(input) || !isShellCommand(input) {
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

func isShellCommand(input string) bool {
	input = strings.ToLower(strings.TrimSpace(input)) // Trim and convert to lowercase
	words := strings.Fields(input)                    // Split into words

	if len(words) == 0 {
		return false
	}

	if slices.Contains(shellCommands, words[0]) {
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
