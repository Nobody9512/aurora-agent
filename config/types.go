package config

import (
	"regexp"

	"github.com/sashabaranov/go-openai"
)

// AppConfig - main configuration structure
type AppConfig struct {
	General   GeneralConfig   `yaml:"general"`
	OpenAI    OpenAIConfig    `yaml:"openai"`
	Interface InterfaceConfig `yaml:"interface"`
}

// GeneralConfig - general configuration
type GeneralConfig struct {
	DefaultShell    string   `yaml:"default_shell"`
	HistorySize     int      `yaml:"history_size"`
	ShellCommands   []string `yaml:"shell_commands"`
	IgnoredCommands []string `yaml:"ignored_commands"`
}

// OpenAIConfig - OpenAI configuration
type OpenAIConfig struct {
	APIKey string `yaml:"api_key"`
	Model  string `yaml:"model"`
}

// InterfaceConfig - interface configuration
type InterfaceConfig struct {
	Theme        string `yaml:"theme"`
	SystemPrompt string `yaml:"system_prompt"`
}

// DefaultConfig - standart configuration
var DefaultConfig = AppConfig{
	General: GeneralConfig{
		DefaultShell:    "",
		HistorySize:     1000,
		ShellCommands:   []string{},
		IgnoredCommands: []string{},
	},
	OpenAI: OpenAIConfig{
		APIKey: "",
		Model:  openai.GPT4o,
	},
	Interface: InterfaceConfig{
		Theme:        "default",
		SystemPrompt: "default",
	},
}

// CurrentConfig - current configuration
var CurrentConfig AppConfig

// Regular expressions for ANSI code processing
var (
	AnsiPattern      = regexp.MustCompile(`\\033\[\d+(;\d+)*m`)
	AnsiStartPattern = regexp.MustCompile(`\\033\[`)
)
