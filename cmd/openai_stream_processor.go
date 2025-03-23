package cmd

import (
	"aurora-agent/utils"
	"fmt"
	"io"

	"github.com/sashabaranov/go-openai"
)

// processStream processes the completion stream and returns the extracted data
func (a *OpenAIAgent) processStream(stream *openai.ChatCompletionStream) (string, bool, string, string, error) {
	// Variables to collect the response
	fullResponse := ""
	functionCall := ""
	functionName := ""
	isFunctionCall := false

	// Ansi buffer
	ansiBuffer := ""

	// Stream the response
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", false, "", "", fmt.Errorf("stream error: %v", err)
		}

		// Check for function call
		if response.Choices[0].Delta.FunctionCall != nil {
			isFunctionCall = true
			if response.Choices[0].Delta.FunctionCall.Name != "" {
				functionName = response.Choices[0].Delta.FunctionCall.Name
			}
			if response.Choices[0].Delta.FunctionCall.Arguments != "" {
				functionCall += response.Choices[0].Delta.FunctionCall.Arguments
			}
			continue
		}

		// Get the content delta
		content := response.Choices[0].Delta.Content
		if content != "" {
			// Collect the full response
			fullResponse += content

			// add to ansi buffer
			ansiBuffer += content
			ansiBuffer = utils.ProcessAnsiBuffer(ansiBuffer)
		}
	}

	// Process remaining buffer
	if ansiBuffer != "" {
		processedBuffer := utils.ProcessANSICodes(ansiBuffer)
		fmt.Print(processedBuffer)
	}

	return fullResponse, isFunctionCall, functionName, functionCall, nil
}
