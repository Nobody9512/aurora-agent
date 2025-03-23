package cmd

import (
	"aurora-agent/config"
	"aurora-agent/utils"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
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
		if isFunctionCall && functionName == "execute_command" {
			if err := a.handleExecuteCommand(functionName, functionCall); err != nil {
				return err
			}
			// Continue the loop to get more function calls
			continue
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

// getAvailableFunctions returns the list of available functions for the AI to call
func (a *OpenAIAgent) getAvailableFunctions() []openai.FunctionDefinition {
	// Define the function for executing shell commands
	return []openai.FunctionDefinition{
		{
			Name:        "execute_command",
			Description: "Execute a shell command and return the output",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"command": map[string]interface{}{
						"type":        "string",
						"description": "The shell command to execute",
					},
				},
				"required": []string{"command"},
			},
		},
	}
}

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
			ansiBuffer = a.processAnsiBuffer(ansiBuffer)
		}
	}

	// Process remaining buffer
	if ansiBuffer != "" {
		processedBuffer := utils.ProcessANSICodes(ansiBuffer)
		fmt.Print(processedBuffer)
	}

	return fullResponse, isFunctionCall, functionName, functionCall, nil
}

// processAnsiBuffer processes ANSI codes in the buffer and prints them
func (a *OpenAIAgent) processAnsiBuffer(ansiBuffer string) string {
	if config.AnsiPattern.MatchString(ansiBuffer) {
		// Buffer has ANSI code, process it
		processedBuffer := utils.ProcessANSICodes(ansiBuffer)
		fmt.Print(processedBuffer)
		return ""
	} else if config.AnsiStartPattern.MatchString(ansiBuffer) && len(ansiBuffer) > 100 {
		// If buffer contains the start of an ANSI code, but not the end
		// and buffer length is more than 100, process it
		// This can happen when ANSI code is in incorrect format
		processedBuffer := utils.ProcessANSICodes(ansiBuffer)
		fmt.Print(processedBuffer)
		return ""
	} else if len(ansiBuffer) > 80 && !config.AnsiStartPattern.MatchString(ansiBuffer) {
		// If buffer length is more than 80 and no ANSI code start is found,
		// process it
		processedBuffer := utils.ProcessANSICodes(ansiBuffer)
		fmt.Print(processedBuffer)
		return ""
	}

	return ansiBuffer
}

// handleExecuteCommand executes a shell command and adds the result to the message history
func (a *OpenAIAgent) handleExecuteCommand(functionName string, functionCall string) error {
	// Parse the function call arguments
	var args struct {
		Command string `json:"command"`
	}
	if err := json.Unmarshal([]byte(functionCall), &args); err != nil {
		return fmt.Errorf("error parsing function call arguments: %v", err)
	}

	// Print the command being executed
	fmt.Printf("\n\033[33mRunning command: %s\033[0m\n", args.Command)

	// Execute the command
	cmd := exec.Command("bash", "-c", args.Command)
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	// Add function call to message history
	a.messages = append(a.messages, openai.ChatCompletionMessage{
		Role: openai.ChatMessageRoleAssistant,
		FunctionCall: &openai.FunctionCall{
			Name:      functionName,
			Arguments: functionCall,
		},
	})

	// Add function result to message history
	result := FunctionCallResult{
		Name:    functionName,
		Output:  outputStr,
		Success: err == nil,
	}
	resultJSON, _ := json.Marshal(result)
	a.messages = append(a.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleFunction,
		Name:    functionName,
		Content: string(resultJSON),
	})

	// Print the command output
	processedOutput := utils.ProcessANSICodes(outputStr)
	fmt.Print(processedOutput)

	// Add a newline after command output for better readability
	fmt.Print("\n")

	return nil
}
