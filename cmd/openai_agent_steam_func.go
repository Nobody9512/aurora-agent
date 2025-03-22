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
func (a *OpenAIAgent) StreamQueryWithFunctionCalls(prompt string, writer io.Writer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Add user message to history
	a.messages = append(a.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	})

	// Define the function for executing shell commands
	functions := []openai.FunctionDefinition{
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

	// Validate model before making the API call
	if a.model == "" {
		// If model is empty, use a default model
		a.model = openai.GPT4o
	}

	// Create a streaming request with function calling
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
		return fmt.Errorf("OpenAI API stream error: %v", err)
	}
	defer stream.Close()

	// Variable to collect the full response
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
			return fmt.Errorf("stream error: %v", err)
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

			// If buffer contains ANSI code
			if config.AnsiPattern.MatchString(ansiBuffer) {
				// Buffer has ANSI code, process it
				processedBuffer := utils.ProcessANSICodes(ansiBuffer)
				fmt.Print(processedBuffer)
				ansiBuffer = ""
			} else if config.AnsiStartPattern.MatchString(ansiBuffer) && len(ansiBuffer) > 100 {
				// If buffer contains the start of an ANSI code, but not the end
				// and buffer length is more than 100, process it
				// This can happen when ANSI code is in incorrect format
				processedBuffer := utils.ProcessANSICodes(ansiBuffer)
				fmt.Print(processedBuffer)
				ansiBuffer = ""
			} else if len(ansiBuffer) > 80 && !config.AnsiStartPattern.MatchString(ansiBuffer) {
				// If buffer length is more than 80 and no ANSI code start is found,
				// process it
				processedBuffer := utils.ProcessANSICodes(ansiBuffer)
				fmt.Print(processedBuffer)
				ansiBuffer = ""
			}
		}
	}

	// Process remaining buffer
	if ansiBuffer != "" {
		processedBuffer := utils.ProcessANSICodes(ansiBuffer)
		fmt.Print(processedBuffer)
	}

	// Handle function call if present
	if isFunctionCall && functionName == "execute_command" {
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

		// Get the final response from the AI with the function result using streaming
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Create a streaming request for the final response
		stream, err := a.client.CreateChatCompletionStream(
			ctx,
			openai.ChatCompletionRequest{
				Model:    a.model,
				Messages: a.messages,
				Stream:   true,
			},
		)
		if err != nil {
			return fmt.Errorf("OpenAI API stream error: %v", err)
		}
		defer stream.Close()

		// Variable to collect the full response
		finalResponse := ""

		// Ansi buffer for the final response
		ansiBuffer := ""

		// Print a newline before the final response
		fmt.Print("\n")

		// Stream the final response
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("stream error: %v", err)
			}

			// Get the content delta
			content := response.Choices[0].Delta.Content
			if content != "" {
				// Collect the full response
				finalResponse += content

				// add to ansi buffer
				ansiBuffer += content

				// If buffer contains ANSI code
				if config.AnsiPattern.MatchString(ansiBuffer) {
					// Buffer has ANSI code, process it
					processedBuffer := utils.ProcessANSICodes(ansiBuffer)
					fmt.Print(processedBuffer)
					ansiBuffer = ""
				} else if config.AnsiStartPattern.MatchString(ansiBuffer) && len(ansiBuffer) > 100 {
					// If buffer contains the start of an ANSI code, but not the end
					// and buffer length is more than 100, process it
					// This can happen when ANSI code is in incorrect format
					processedBuffer := utils.ProcessANSICodes(ansiBuffer)
					fmt.Print(processedBuffer)
					ansiBuffer = ""
				} else if len(ansiBuffer) > 80 && !config.AnsiStartPattern.MatchString(ansiBuffer) {
					// If buffer length is more than 80 and no ANSI code start is found,
					// process it
					processedBuffer := utils.ProcessANSICodes(ansiBuffer)
					fmt.Print(processedBuffer)
					ansiBuffer = ""
				}
			}
		}

		// Process remaining buffer
		if ansiBuffer != "" {
			processedBuffer := utils.ProcessANSICodes(ansiBuffer)
			fmt.Print(processedBuffer)
		}

		// Print a newline after the final response
		fmt.Print("\n")

		// Add the final response to message history
		a.messages = append(a.messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: finalResponse,
		})

		return nil
	}

	// Add assistant response to history if no function call
	if !isFunctionCall {
		a.messages = append(a.messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: fullResponse,
		})
	}

	return nil
}
