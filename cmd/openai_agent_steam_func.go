package cmd

// This file is being refactored into several smaller files for better code organization:
// - openai_function_executor.go - For function execution code
// - openai_stream_processor.go - For stream processing
// - openai_function_stream.go - For the main function streaming functionality
// - utils/ansi.go - For ANSI code processing

// The implementation remains in this file temporarily to ensure backward compatibility
// while the refactoring is being completed.

// In the future, this file will be removed and replaced with the imports of the
// smaller, more focused files.

import (
	"context"
	"fmt"
	"time"

	"github.com/sashabaranov/go-openai"
)

// StreamQueryWithFunctionCalls sends a prompt to OpenAI, handles function calls, and streams the response
func (a *OpenAIAgent) StreamQueryWithFunctionCalls(prompt string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// Add user message to history
	a.messages = append(a.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	})

	// Validate model before making the API call
	if a.model == "" {
		// If model is empty, use a default model
		a.model = openai.GPT4o
	}

	// Main processing loop - allows multiple commands in sequence
	for {
		fullResponse, isFunctionCall, functionName, functionCall, err := a.streamResponseWithFunctions(ctx)
		if err != nil {
			return err
		}

		// Handle function call if present
		if isFunctionCall {
			if functionName == "execute_command" {
				if err := a.handleExecuteCommand(functionName, functionCall); err != nil {
					return err
				}
				// Continue the loop to get more function calls
				continue
			} else if functionName == "pwd" {
				if err := a.handlePwd(functionName, functionCall); err != nil {
					return err
				}
				// Continue the loop to get more function calls
				continue
			}
		}

		// Add assistant response to history if no function call
		if !isFunctionCall {
			a.messages = append(a.messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: fullResponse,
			})
		}

		// Add a newline at the end of the response for better readability
		fmt.Print("\n")

		// If we reach here with no function call, exit the loop
		if !isFunctionCall {
			break
		}
	}

	return nil
}
