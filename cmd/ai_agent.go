package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

// AgentType represents the type of AI agent
type AgentType string

const (
	OpenAI AgentType = "openai"
	Claude AgentType = "claude"

	// System prompt -> TODO: move to a config file
	// You are a helpful assistant that provides SHORT and CONCISE answers.
	SystemPrompt = `
Your name is Aurora.
You are a helpful assistant that provides SHORT and CONCISE answers.
You are currently in a terminal environment.
You can use ANSI escape codes to color text:
- Red: \033[31m
- Green: \033[32m
- Yellow: \033[33m
- Blue: \033[34m
- Magenta: \033[35m
- Cyan: \033[36m
- Reset: \033[0m
- Bold: \033[1m
- Underline: \033[4m

Use appropriate colors to highlight important information, warnings, and errors.
Example usage: \033[31mThis is red text\033[0m

You can use the following commands:
clear - to clear the terminal
	`
)

// Ansi pattern
var ansiPattern = regexp.MustCompile(`\\033\[\d+(;\d+)*m`)

// Ansi code start
var ansiStartPattern = regexp.MustCompile(`\\033\[`)

// AIAgent interface for different AI providers
type AIAgent interface {
	Query(prompt string) (string, error)
	StreamQuery(prompt string, writer io.Writer) error
	Name() string
}

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
				Content: SystemPrompt,
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
			if ansiPattern.MatchString(ansiBuffer) {
				// Buffer has ANSI code, process it
				processedBuffer := processANSICodes(ansiBuffer)
				fmt.Print(processedBuffer)
				ansiBuffer = ""
			} else if ansiStartPattern.MatchString(ansiBuffer) && len(ansiBuffer) > 30 {
				// If buffer contains the start of an ANSI code, but not the end
				// and buffer length is more than 30, process it
				// This can happen when ANSI code is in incorrect format
				processedBuffer := processANSICodes(ansiBuffer)
				fmt.Print(processedBuffer)
				ansiBuffer = ""
			} else if len(ansiBuffer) > 20 && !ansiStartPattern.MatchString(ansiBuffer) {
				// If buffer length is more than 20 and no ANSI code start is found,
				// process it
				processedBuffer := processANSICodes(ansiBuffer)
				fmt.Print(processedBuffer)
				ansiBuffer = ""
			}
		}
	}

	// Process remaining buffer
	if ansiBuffer != "" {
		processedBuffer := processANSICodes(ansiBuffer)
		fmt.Print(processedBuffer)
	}

	// Add assistant response to history
	a.messages = append(a.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: fullResponse,
	})

	return nil
}

// processANSICodes barcha turdagi ANSI kodlarini qayta ishlaydi
func processANSICodes(text string) string {
	// Oddiy escape qilingan kodlarni almashtiramiz (\\033)
	result := strings.ReplaceAll(text, "\\033", "\033")

	// Qo'sh escape qilingan kodlarni almashtiramiz (\\\\033)
	result = strings.ReplaceAll(result, "\\\\033", "\033")

	// Unicode escape qilingan kodlarni almashtiramiz (\u001b)
	result = strings.ReplaceAll(result, "\\u001b", "\033")

	// Agar hali ham escape qilingan kodlar qolgan bo'lsa
	result = strings.ReplaceAll(result, "\\e", "\033")

	// Agar hali ham escape qilingan kodlar qolgan bo'lsa
	result = strings.ReplaceAll(result, "\\x1b", "\033")

	// Agar hali ham escape qilingan kodlar qolgan bo'lsa
	result = strings.ReplaceAll(result, "\\x1B", "\033")

	// Agar hali ham escape qilingan kodlar qolgan bo'lsa
	result = strings.ReplaceAll(result, "\\u001B", "\033")

	// Agar hali ham escape qilingan kodlar qolgan bo'lsa
	result = strings.ReplaceAll(result, "\\27", "\033")

	// Agar hali ham escape qilingan kodlar qolgan bo'lsa
	result = strings.ReplaceAll(result, "\\33", "\033")

	return result
}

// Name returns the name of the agent
func (a *OpenAIAgent) Name() string {
	return string(OpenAI)
}

// SetModel sets the OpenAI model to use
func (a *OpenAIAgent) SetModel(model string) {
	a.model = model
}

// AgentManager manages different AI agents
type AgentManager struct {
	activeAgent AIAgent
	agents      map[AgentType]AIAgent
}

// NewAgentManager creates a new agent manager
func NewAgentManager() *AgentManager {
	// Create a default OpenAI agent
	// In a real implementation, you would get the API key from environment or config
	openAIAgent := NewOpenAIAgent("")

	agents := make(map[AgentType]AIAgent)
	agents[OpenAI] = openAIAgent

	return &AgentManager{
		activeAgent: openAIAgent,
		agents:      agents,
	}
}

// SetActiveAgent sets the active AI agent
func (m *AgentManager) SetActiveAgent(agentType AgentType) error {
	agent, exists := m.agents[agentType]
	if !exists {
		return fmt.Errorf("agent type %s not found", agentType)
	}

	m.activeAgent = agent
	return nil
}

// AddAgent adds a new AI agent
func (m *AgentManager) AddAgent(agentType AgentType, agent AIAgent) {
	m.agents[agentType] = agent
}

// Query sends a prompt to the active AI agent
func (m *AgentManager) Query(prompt string) (string, error) {
	if m.activeAgent == nil {
		return "", fmt.Errorf("no active agent set")
	}

	return m.activeAgent.Query(prompt)
}

// GetActiveAgentName returns the name of the active agent
func (m *AgentManager) GetActiveAgentName() string {
	if m.activeAgent == nil {
		return "none"
	}

	return m.activeAgent.Name()
}

// StreamQuery sends a prompt to the active AI agent and streams the response
func (m *AgentManager) StreamQuery(prompt string, writer io.Writer) error {
	if m.activeAgent == nil {
		return fmt.Errorf("no active agent set")
	}

	return m.activeAgent.StreamQuery(prompt, writer)
}
