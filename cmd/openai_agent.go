package cmd

import (
	"aurora-agent/config"
	"log"
	"os"

	"github.com/sashabaranov/go-openai"
)

// OpenAIAgent implements the AIAgent interface for OpenAI
type OpenAIAgent struct {
	client   *openai.Client
	model    string
	messages []openai.ChatCompletionMessage
}

// NewOpenAIAgent creates a new OpenAI agent
func NewOpenAIAgent(apiKey string) *OpenAIAgent {
	// If apiKey is empty, try to get it from config first, then from environment variable
	if apiKey == "" {
		// Try to get API key from config
		apiKey = config.CurrentConfig.OpenAI.APIKey

		// If still empty, try to get from environment variable
		if apiKey == "" {
			apiKey = os.Getenv("OPENAI_API_KEY")
			if apiKey == "" {
				log.Fatal("Error: OpenAI API key not found in config or environment variable (OPENAI_API_KEY).")
				os.Exit(1)
			}
		}
	}

	client := openai.NewClient(apiKey)

	// Get model from config, use default if empty or invalid
	model := config.CurrentConfig.OpenAI.Model
	if model == "" {
		// Use default model if not specified
		model = openai.GPT4o
	}

	return &OpenAIAgent{
		client: client,
		model:  model,
		messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: config.GetSystemPrompt(),
			},
		},
	}
}

// Name returns the name of the agent
func (a *OpenAIAgent) Name() string {
	return string(OpenAI)
}

// SetModel sets the OpenAI model to use
func (a *OpenAIAgent) SetModel(model string) {
	a.model = model
}

// StreamQueryWithFunctionCallsV2 is the new refactored version of StreamQueryWithFunctionCalls
// It will replace StreamQueryWithFunctionCalls once the refactoring is complete
func (a *OpenAIAgent) StreamQueryWithFunctionCallsV2(prompt string) error {
	// For now, just call the original function
	// This will be updated to use the new refactored functions
	return a.StreamQueryWithFunctionCalls(prompt)
}
