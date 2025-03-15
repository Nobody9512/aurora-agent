package cmd

import (
	"io"
)

// AgentType represents the type of AI agent
type AgentType string

const (
	// OpenAI represents the OpenAI agent type
	OpenAI AgentType = "openai"
	// Claude represents the Claude agent type
	Claude AgentType = "claude"
)

// AIAgent interface for different AI providers
type AIAgent interface {
	// Query sends a prompt to the AI and returns the response
	Query(prompt string) (string, error)
	// StreamQuery sends a prompt to the AI and streams the response to the writer
	StreamQuery(prompt string, writer io.Writer) error
	// Name returns the name of the agent
	Name() string
}