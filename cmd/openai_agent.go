package cmd

import (
	"aurora-agent/config"
	"aurora-agent/utils"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

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
	if apiKey == "" {
		// Try to get API key from environment variable
		apiKey = os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			log.Fatal("Warning: OPENAI_API_KEY not set. Using demo mode.")
			os.Exit(1)
		}
	}

	client := openai.NewClient(apiKey)

	return &OpenAIAgent{
		client: client,
		model:  openai.GPT4oLatest, // Default model
		messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: config.SystemPrompt,
			},
		},
	}
}

// Query sends a prompt to OpenAI and returns the response
func (a *OpenAIAgent) Query(prompt string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	a.messages = append(a.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	})

	resp, err := a.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    a.model,
			Messages: a.messages,
		},
	)

	if err != nil {
		return "", fmt.Errorf("OpenAI API error: %v", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}

// StreamQuery sends a prompt to OpenAI and streams the response to the writer
func (a *OpenAIAgent) StreamQuery(prompt string, writer io.Writer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Add user message to history
	a.messages = append(a.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	})

	// Create a streaming request
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
	fullResponse := ""

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

	// Add assistant response to history
	a.messages = append(a.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: fullResponse,
	})

	return nil
}

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
		fmt.Print(outputStr)

		// Get the final response from the AI with the function result
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		resp, err := a.client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model:    a.model,
				Messages: a.messages,
			},
		)

		if err != nil {
			return fmt.Errorf("OpenAI API error: %v", err)
		}

		if len(resp.Choices) == 0 {
			return fmt.Errorf("no response from OpenAI")
		}

		// Print the AI's response
		finalResponse := resp.Choices[0].Message.Content
		fmt.Print("\n" + finalResponse + "\n")

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

// Name returns the name of the agent
func (a *OpenAIAgent) Name() string {
	return string(OpenAI)
}

// SetModel sets the OpenAI model to use
func (a *OpenAIAgent) SetModel(model string) {
	a.model = model
}