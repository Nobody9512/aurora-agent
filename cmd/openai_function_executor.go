package cmd

import (
	"aurora-agent/utils"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/sashabaranov/go-openai"
)

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
		{
			Name:        "pwd",
			Description: "Print current working directory",
			Parameters: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
				"required":   []string{},
			},
		},
	}
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

// handlePwd gets the current working directory and adds the result to the message history
func (a *OpenAIAgent) handlePwd(functionName string, functionCall string) error {
	// Execute the command
	cmd := exec.Command("pwd")
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

	// Add a newline after command output for better readability
	fmt.Print("\n")

	return nil
}
