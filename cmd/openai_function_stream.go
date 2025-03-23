package cmd

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

// streamResponseWithFunctions creates a completion stream and processes the response
func (a *OpenAIAgent) streamResponseWithFunctions(ctx context.Context) (string, bool, string, string, error) {
	functions := a.getAvailableFunctions()

	stream, err := a.client.CreateChatCompletionStream(
		ctx,
		openai.ChatCompletionRequest{
			Model:     a.model,
			Messages:  a.messages,
			Stream:    true,
			Functions: functions,
		},
	)
	if err != nil {
		return "", false, "", "", fmt.Errorf("OpenAI API stream error: %v", err)
	}
	defer stream.Close()

	return a.processStream(stream)
}
