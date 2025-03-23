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
		{
			Name:        "read_file",
			Description: "Read the contents of a file, either the entire file or a specific range of lines",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"file_path": map[string]interface{}{
						"type":        "string",
						"description": "The path to the file to read",
					},
					"start_line": map[string]interface{}{
						"type":        "integer",
						"description": "The starting line number to read (optional, 1-based indexing)",
					},
					"end_line": map[string]interface{}{
						"type":        "integer",
						"description": "The ending line number to read (optional, 1-based indexing)",
					},
					"read_entire": map[string]interface{}{
						"type":        "boolean",
						"description": "Whether to read the entire file regardless of size (optional, defaults to false)",
					},
				},
				"required": []string{"file_path"},
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

// handleReadFile reads the contents of a file and adds the result to the message history
func (a *OpenAIAgent) handleReadFile(functionName string, functionCall string) error {
	// Parse the function call arguments
	var args struct {
		FilePath   string `json:"file_path"`
		StartLine  int    `json:"start_line"`
		EndLine    int    `json:"end_line"`
		ReadEntire bool   `json:"read_entire"`
	}
	if err := json.Unmarshal([]byte(functionCall), &args); err != nil {
		return fmt.Errorf("error parsing function call arguments: %v", err)
	}

	// Print what file is being read
	fmt.Printf("\n\033[33mReading file: %s\033[0m\n", args.FilePath)

	var outputStr string
	var err error

	// Check if file exists
	checkCmd := exec.Command("bash", "-c", fmt.Sprintf("test -f '%s' && echo 'exists' || echo 'not found'", args.FilePath))
	checkOutput, checkErr := checkCmd.CombinedOutput()
	if checkErr != nil || string(checkOutput) == "not found\n" {
		outputStr = fmt.Sprintf("Error: File not found - %s", args.FilePath)
		err = fmt.Errorf("file not found")
		// Print the error message
		fmt.Print(outputStr)
	} else {
		// Determine if we need to read specific lines or the entire file
		if args.StartLine > 0 && args.EndLine > 0 && !args.ReadEntire {
			// Read specific range of lines
			cmd := exec.Command("bash", "-c", fmt.Sprintf("sed -n '%d,%dp' '%s'", args.StartLine, args.EndLine, args.FilePath))
			output, readErr := cmd.CombinedOutput()
			outputStr = string(output)

			// Get file info
			infoCmd := exec.Command("bash", "-c", fmt.Sprintf("wc -l < '%s'", args.FilePath))
			infoOutput, _ := infoCmd.CombinedOutput()
			totalLines := string(infoOutput)

			// Only create the description for messages
			outputInfoStr := fmt.Sprintf("File: %s (reading lines %d-%d of %s)",
				args.FilePath, args.StartLine, args.EndLine, totalLines)

			// Print only the file info, not the content
			fmt.Println(outputInfoStr)

			err = readErr
		} else {
			// First check file size
			sizeCmd := exec.Command("bash", "-c", fmt.Sprintf("wc -c < '%s'", args.FilePath))
			sizeOutput, _ := sizeCmd.CombinedOutput()
			sizeStr := string(sizeOutput)

			// Convert size to int
			var size int
			fmt.Sscanf(sizeStr, "%d", &size)

			// Get line count
			lineCmd := exec.Command("bash", "-c", fmt.Sprintf("wc -l < '%s'", args.FilePath))
			lineOutput, _ := lineCmd.CombinedOutput()
			lineStr := string(lineOutput)

			// If file is large (> 1MB) and not forced to read entire, read first 100 lines
			const MAX_SIZE = 1048576 // 1MB
			if size > MAX_SIZE && !args.ReadEntire {
				headCmd := exec.Command("bash", "-c", fmt.Sprintf("head -n 100 '%s'", args.FilePath))
				headOutput, headErr := headCmd.CombinedOutput()
				outputStr = string(headOutput)
				err = headErr

				// Only print file info
				fmt.Printf("File: %s (reading first 100 lines, file is large: %s bytes, %s lines total)\n",
					args.FilePath, sizeStr, lineStr)
			} else {
				// Read entire file
				cmd := exec.Command("bash", "-c", fmt.Sprintf("cat '%s'", args.FilePath))
				output, readErr := cmd.CombinedOutput()
				outputStr = string(output)
				err = readErr

				// Only print file info
				fmt.Printf("File: %s (reading entire file, %s bytes, %s lines)\n",
					args.FilePath, sizeStr, lineStr)
			}
		}
	}

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

	// Add a newline after output for better readability
	fmt.Print("\n")

	return nil
}
