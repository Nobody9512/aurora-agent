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

// FunctionCallResult represents the result of a function call
type FunctionCallResult struct {
	Name    string
	Output  string
	Success bool
}

// AIAgent interface for different AI providers
type AIAgent interface {
	// Query sends a prompt to the AI and returns the response
	Query(prompt string) (string, error)
	// StreamQuery sends a prompt to the AI and streams the response to the writer
	StreamQuery(prompt string, writer io.Writer) error
	// StreamQueryWithFunctionCalls sends a prompt to the AI, handles function calls, and streams the response
	StreamQueryWithFunctionCalls(prompt string) error
	// Name returns the name of the agent
	Name() string
}
