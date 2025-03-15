package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// AppConfig - asosiy konfiguratsiya strukturasi
type AppConfig struct {
	General   GeneralConfig   `yaml:"general"`
	OpenAI    OpenAIConfig    `yaml:"openai"`
	Interface InterfaceConfig `yaml:"interface"`
}

// GeneralConfig - umumiy sozlamalar
type GeneralConfig struct {
	DefaultShell    string   `yaml:"default_shell"`
	HistorySize     int      `yaml:"history_size"`
	ShellCommands   []string `yaml:"shell_commands"`
	IgnoredCommands []string `yaml:"ignored_commands"`
}

// OpenAIConfig - OpenAI sozlamalari
type OpenAIConfig struct {
	APIKey      string  `yaml:"api_key"`
	Model       string  `yaml:"model"`
	Temperature float64 `yaml:"temperature"`
}

// InterfaceConfig - interfeys sozlamalari
type InterfaceConfig struct {
	Theme       string `yaml:"theme"`
	PromptStyle string `yaml:"prompt_style"`
}

// DefaultConfig - standart konfiguratsiya
var DefaultConfig = AppConfig{
	General: GeneralConfig{
		DefaultShell:    "",
		HistorySize:     1000,
		ShellCommands:   []string{},
		IgnoredCommands: []string{},
	},
	OpenAI: OpenAIConfig{
		APIKey:      "",
		Model:       "gpt-3.5-turbo",
		Temperature: 0.7,
	},
	Interface: InterfaceConfig{
		Theme:       "default",
		PromptStyle: "default",
	},
}

// CurrentConfig - current configuration
var CurrentConfig AppConfig

// GetConfigPath - get configuration file path
func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error: User home directory not found:", err)
		return ""
	}

	return filepath.Join(homeDir, ".config", "aurora", "config.yaml")
}

// LoadConfig - load configuration file or create a new one
func LoadConfig() error {
	configPath := GetConfigPath()

	// Check and create configuration directory
	configDir := filepath.Dir(configPath)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create configuration directory: %w", err)
		}
	}

	// Read configuration file
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// File does not exist, save default configuration
			CurrentConfig = DefaultConfig
			return SaveConfig()
		}
		return fmt.Errorf("failed to read configuration file: %w", err)
	}

	// Read YAML format
	if err := yaml.Unmarshal(data, &CurrentConfig); err != nil {
		return fmt.Errorf("failed to parse configuration file: %w", err)
	}

	return nil
}

// SaveConfig - save current configuration to file
func SaveConfig() error {
	configPath := GetConfigPath()

	// Convert to YAML format
	data, err := yaml.Marshal(&CurrentConfig)
	if err != nil {
		return fmt.Errorf("failed to convert configuration to YAML: %w", err)
	}

	// Save to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to save configuration file: %w", err)
	}

	return nil
}

// DefaultShellCommands - default terminal commands list
var DefaultShellCommands = []string{
	// User information and system data
	"whoami", "id", "uname", "hostname", "uptime", "w", "who", "groups",
	"whois", "finger", "last", "lastlog", "lastb", "lastcomm", "lastcomm",
	"clear", "help", "exit", "quit", "cls", "clr",

	// User and permissions
	"su", "sudo", "passwd", "chown", "chmod", "chgrp", "umask",

	// Folders and files
	"ls", "ll", "la", "pwd", "cd", "mkdir", "rmdir", "touch", "rm", "mv", "cp",
	"find", "locate", "updatedb", "tree",

	// File and text manipulation
	"cat", "tac", "less", "more", "head", "tail", "grep", "awk", "sed", "cut",
	"sort", "uniq", "tee", "wc", "diff", "cmp",

	// Processes and system monitoring
	"ps", "top", "htop", "nice", "renice", "kill", "pkill", "jobs", "fg", "bg",
	"nohup", "time", "strace", "lsof",

	// Disk and file system
	"df", "du", "mount", "umount", "fsck", "mkfs", "blkid", "fdisk", "parted",
	"lsblk", "e2fsck", "sync", "tune2fs",

	// Network and internet
	"ping", "traceroute", "netstat", "ss", "ifconfig", "ip", "iwconfig",
	"curl", "wget", "scp", "rsync", "nc", "telnet", "ftp", "sftp",

	// Archiving and compression
	"tar", "zip", "unzip", "gzip", "gunzip", "bzip2", "bunzip2", "xz", "7z",
	"rar", "unrar",

	// Software and package management
	"apt", "apt-get", "yum", "dnf", "pacman", "zypper", "brew", "snap", "flatpak",
	"dpkg", "rpm", "pip", "gem", "npm", "cargo", "go", "python", "ruby", "perl",

	// Docker and containers
	"docker", "docker-compose", "podman", "kubectl", "minikube",

	// Git and version control
	"git", "git clone", "git pull", "git push", "git commit", "git status",
	"git log", "git branch", "git checkout", "git merge", "git rebase",

	// Programming languages and compilers
	"python", "python3", "node", "nodejs", "npm", "npx", "java", "javac",
	"ruby", "perl", "php", "php-cli", "go", "rustc", "gcc", "g++",

	// Shell and scripts
	"bash", "sh", "zsh", "fish", "dash", "tcsh", "csh", "ksh",
}

// GetShellCommands - terminal commands list
func GetShellCommands() []string {
	// If there are additional commands in the configuration, add them
	allCommands := make([]string, 0, len(DefaultShellCommands)+len(CurrentConfig.General.ShellCommands))

	// Add default commands
	allCommands = append(allCommands, DefaultShellCommands...)

	// Add user-added commands
	allCommands = append(allCommands, CurrentConfig.General.ShellCommands...)

	// Remove ignored commands
	if len(CurrentConfig.General.IgnoredCommands) > 0 {
		filteredCommands := make([]string, 0, len(allCommands))
		for _, cmd := range allCommands {
			ignored := false
			for _, ignoredCmd := range CurrentConfig.General.IgnoredCommands {
				if cmd == ignoredCmd {
					ignored = true
					break
				}
			}
			if !ignored {
				filteredCommands = append(filteredCommands, cmd)
			}
		}
		return filteredCommands
	}

	return allCommands
}

// AddShellCommand - add new terminal command
func AddShellCommand(command string) error {
	// Command must not be empty
	if command = strings.TrimSpace(command); command == "" {
		return fmt.Errorf("command cannot be empty")
	}

	// Check if command already exists
	commands := GetShellCommands()
	for _, cmd := range commands {
		if cmd == command {
			return fmt.Errorf("'%s' command already exists", command)
		}
	}

	// Remove command from ignored list (if it exists)
	for i, cmd := range CurrentConfig.General.IgnoredCommands {
		if cmd == command {
			CurrentConfig.General.IgnoredCommands = append(
				CurrentConfig.General.IgnoredCommands[:i],
				CurrentConfig.General.IgnoredCommands[i+1:]...,
			)
			break
		}
	}

	// Add command
	CurrentConfig.General.ShellCommands = append(CurrentConfig.General.ShellCommands, command)
	return nil
}

// RemoveShellCommand - remove terminal command or ignore it
func RemoveShellCommand(command string) error {
	// Command must not be empty
	if command = strings.TrimSpace(command); command == "" {
		return fmt.Errorf("command cannot be empty")
	}

	// Check if command exists
	found := false
	for _, cmd := range DefaultShellCommands {
		if cmd == command {
			found = true
			break
		}
	}

	if !found {
		// Search in user-added commands
		for i, cmd := range CurrentConfig.General.ShellCommands {
			if cmd == command {
				// Remove command
				CurrentConfig.General.ShellCommands = append(
					CurrentConfig.General.ShellCommands[:i],
					CurrentConfig.General.ShellCommands[i+1:]...,
				)
				return nil
			}
		}
		return fmt.Errorf("'%s' command not found", command)
	}

	// If it's a default command, add it to the ignored list
	for _, cmd := range CurrentConfig.General.IgnoredCommands {
		if cmd == command {
			return fmt.Errorf("'%s' command already ignored", command)
		}
	}

	CurrentConfig.General.IgnoredCommands = append(CurrentConfig.General.IgnoredCommands, command)
	return nil
}

// ResetShellCommands - reset terminal commands list to default
func ResetShellCommands() {
	CurrentConfig.General.ShellCommands = []string{}
	CurrentConfig.General.IgnoredCommands = []string{}
}

// DefaultSystemPrompt - standart tizim prompti
const DefaultSystemPrompt = `
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

You can execute terminal commands when asked. For example, if someone asks about the version of a program installed, you can run the appropriate command to check and provide the answer.

{{USER_INPUT}}
`

// GetSystemPrompt - user prompt style-based system prompt
func GetSystemPrompt() string {
	// If user prompt style exists, add it
	if CurrentConfig.Interface.PromptStyle != "default" && CurrentConfig.Interface.PromptStyle != "" {
		// Add user prompt style to system prompt
		// Do not change {{USER_INPUT}}
		customPrompt := strings.Replace(DefaultSystemPrompt, "{{USER_INPUT}}",
			CurrentConfig.Interface.PromptStyle+"\n\n{{USER_INPUT}}", 1)
		return customPrompt
	}

	return strings.Replace(DefaultSystemPrompt, "{{USER_INPUT}}", "", 1)
}

var (
	AnsiPattern      = regexp.MustCompile(`\\033\[\d+(;\d+)*m`)
	AnsiStartPattern = regexp.MustCompile(`\\033\[`)
)
