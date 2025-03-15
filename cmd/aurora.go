package cmd

import (
	"fmt"
	"os"
	"strings"

	"aurora-agent/config"
)

// Global agent manager instance
var AgentMgr *AgentManager

func init() {
	// Initialize the agent manager
	AgentMgr = NewAgentManager()
}

// ProcessAuroraCommand handles Aurora-specific commands
func ProcessAuroraCommand(input string) bool {
	// Check configuration commands
	if processConfigCommand(input) {
		return true
	}

	// Check if input contains "aurora" or is not a shell command
	if isAuroraCommand(input) || !isShellCommand(input) {
		// Use streaming response
		fmt.Print("\n") // Add a newline before the response for better readability

		// Print a colored prompt to indicate AI response
		fmt.Print("\033[36mAurora: \033[0m") // Cyan color for Aurora name

		// Use function calls for natural language processing
		err := AgentMgr.StreamQueryWithFunctionCalls(input, os.Stdout)
		if err != nil {
			fmt.Printf("\n\033[31mError querying AI agent: %v\033[0m\n", err) // Red error message
		} else {
			fmt.Println() // Add a newline after the streamed response
		}

		return true
	}
	return false
}

func processConfigCommand(input string) bool {
	input = strings.TrimSpace(input)
	words := strings.Fields(input)

	if len(words) == 0 {
		return false
	}

	// Help command
	if words[0] == "help" {
		showHelp()
		return true
	}

	// Check configuration commands
	if words[0] == "config" {
		if len(words) == 1 {
			// Show configuration information
			showConfig()
			return true
		}

		switch words[1] {
		case "show":
			// Show configuration information
			showConfig()
			return true

		case "set":
			// Change configuration
			if len(words) < 4 {
				fmt.Println("\033[31mError: Wrong format. Use: config set <section> <key> <value>\033[0m")
				return true
			}
			setConfigValue(words[2], words[3], strings.Join(words[4:], " "))
			return true

		case "save":
			// Save configuration
			if err := config.SaveConfig(); err != nil {
				fmt.Printf("\033[31mError: Failed to save configuration: %v\033[0m\n", err)
			} else {
				fmt.Printf("\033[32mConfiguration saved successfully: %s\033[0m\n", config.GetConfigPath())
			}
			return true

		case "reload":
			// Reload configuration
			if err := config.LoadConfig(); err != nil {
				fmt.Printf("\033[31mError: Failed to load configuration: %v\033[0m\n", err)
			} else {
				fmt.Println("\033[32mConfiguration loaded successfully\033[0m")
			}
			return true

		case "commands":
			// Shell commands
			if len(words) < 3 {
				fmt.Println("\033[31mError: Wrong format. Use: config commands [list|add|remove|reset]\033[0m")
				return true
			}

			switch words[2] {
			case "list":
				// Show commands list
				showShellCommands()
				return true

			case "add":
				// Add command
				if len(words) < 4 {
					fmt.Println("\033[31mError: Wrong format. Use: config commands add <command>\033[0m")
					return true
				}
				addShellCommand(words[3])
				return true

			case "remove":
				// Remove command
				if len(words) < 4 {
					fmt.Println("\033[31mError: Wrong format. Use: config commands remove <command>\033[0m")
					return true
				}
				removeShellCommand(words[3])
				return true

			case "reset":
				// Reset commands list to default
				resetShellCommands()
				return true

			default:
				fmt.Println("\033[31mUnknown command. Available commands: list, add, remove, reset\033[0m")
				return true
			}

		default:
			fmt.Println("\033[31mUnknown configuration command. Available commands: show, set, save, reload, commands\033[0m")
			return true
		}
	}

	return false
}

// showConfig - show config information
func showConfig() {
	fmt.Println("\n\033[1mCurrent configuration:\033[0m")
	fmt.Println("\n\033[1m[General]\033[0m")
	fmt.Printf("  DefaultShell: %s\n", config.CurrentConfig.General.DefaultShell)
	fmt.Printf("  HistorySize: %d\n", config.CurrentConfig.General.HistorySize)
	fmt.Printf("  ShellCommands: %d additional commands\n", len(config.CurrentConfig.General.ShellCommands))
	fmt.Printf("  IgnoredCommands: %d ignored commands\n", len(config.CurrentConfig.General.IgnoredCommands))

	fmt.Println("\033[1m[OpenAI]\033[0m")
	apiKey := config.CurrentConfig.OpenAI.APIKey
	if apiKey != "" {
		apiKey = "********" // Hide API key
	}
	fmt.Printf("  APIKey: %s\n", apiKey)
	fmt.Printf("  Model: %s\n", config.CurrentConfig.OpenAI.Model)
	fmt.Printf("  Temperature: %.1f\n", config.CurrentConfig.OpenAI.Temperature)

	fmt.Println("\033[1m[Interface]\033[0m")
	fmt.Printf("  Theme: %s\n", config.CurrentConfig.Interface.Theme)

	// Prompt style information
	promptStyle := config.CurrentConfig.Interface.PromptStyle
	if promptStyle == "default" || promptStyle == "" {
		fmt.Printf("  PromptStyle: default\n")
	} else {
		// If prompt style is long, shorten it
		if len(promptStyle) > 50 {
			fmt.Printf("  PromptStyle: \"%s...\"\n", promptStyle[:47])
		} else {
			fmt.Printf("  PromptStyle: \"%s\"\n", promptStyle)
		}
	}

	fmt.Printf("\nConfiguration file: \033[32m%s\033[0m\n", config.GetConfigPath())
	fmt.Println("\nTo see the commands list, use `\033[32mconfig commands list\033[0m`")
	fmt.Println()
}

// setConfigValue - change configuration value
func setConfigValue(section, key, value string) {
	switch strings.ToLower(section) {
	case "general":
		switch strings.ToLower(key) {
		case "defaultshell":
			config.CurrentConfig.General.DefaultShell = value
			fmt.Printf("\033[32mGeneral.DefaultShell = %s\033[0m\n", value)
		case "historysize":
			var size int
			if _, err := fmt.Sscanf(value, "%d", &size); err != nil {
				fmt.Println("\033[31mError: HistorySize must be an integer\033[0m")
				return
			}
			config.CurrentConfig.General.HistorySize = size
			fmt.Printf("\033[32mGeneral.HistorySize = %d\033[0m\n", size)
		default:
			fmt.Printf("\033[31mError: '%s' key not found in General section\033[0m\n", key)
		}

	case "openai":
		switch strings.ToLower(key) {
		case "apikey":
			config.CurrentConfig.OpenAI.APIKey = value
			fmt.Println("\033[32mOpenAI.APIKey updated\033[0m")
			// API kalit o'zgartirilganda agentni qayta ishga tushirish
			AgentMgr = NewAgentManager()
		case "model":
			config.CurrentConfig.OpenAI.Model = value
			fmt.Printf("\033[32mOpenAI.Model = %s\033[0m\n", value)
		case "temperature":
			var temp float64
			if _, err := fmt.Sscanf(value, "%f", &temp); err != nil {
				fmt.Println("\033[31mError: Temperature must be a float\033[0m")
				return
			}
			config.CurrentConfig.OpenAI.Temperature = temp
			fmt.Printf("\033[32mOpenAI.Temperature = %.1f\033[0m\n", temp)
		default:
			fmt.Printf("\033[31mError: '%s' key not found in OpenAI section\033[0m\n", key)
		}

	case "interface":
		switch strings.ToLower(key) {
		case "theme":
			config.CurrentConfig.Interface.Theme = value
			fmt.Printf("\033[32mInterface.Theme = %s\033[0m\n", value)
		case "promptstyle":
			config.CurrentConfig.Interface.PromptStyle = value
			fmt.Printf("\033[32mInterface.PromptStyle = %s\033[0m\n", value)
			// Prompt style o'zgartirilganda agentni qayta ishga tushirish
			AgentMgr = NewAgentManager()
		default:
			fmt.Printf("\033[31mError: '%s' key not found in Interface section\033[0m\n", key)
		}

	default:
		fmt.Printf("\033[31mError: '%s' section not found. Available sections: General, OpenAI, Interface\033[0m\n", section)
	}

	fmt.Println("\033[33mNote: Remember to save changes using 'config save'\033[0m")
}

// showHelp - show help information
func showHelp() {
	fmt.Println("\n\033[1mAurora Agent help information\033[0m")
	fmt.Println("\n\033[1mMain commands:\033[0m")
	fmt.Println("  \033[32mhelp\033[0m                - Show help information")
	fmt.Println("  \033[32mexit, quit\033[0m          - Exit the program")
	fmt.Println("  \033[32mclear\033[0m               - Clear the screen")

	fmt.Println("\033[1mConfiguration commands:\033[0m")
	fmt.Println("  \033[32mconfig\033[0m              - Show current configuration")
	fmt.Println("  \033[32mconfig show\033[0m         - Show current configuration")
	fmt.Println("  \033[32mconfig set <section> <key> <value>\033[0m - Change configuration value")
	fmt.Println("  \033[32mconfig save\033[0m         - Save configuration")
	fmt.Println("  \033[32mconfig reload\033[0m       - Reload configuration")

	fmt.Println("\033[1mWorking with shell commands:\033[0m")
	fmt.Println("  \033[32mconfig commands list\033[0m    - Show all commands")
	fmt.Println("  \033[32mconfig commands add <command>\033[0m - Add new command")
	fmt.Println("  \033[32mconfig commands remove <command>\033[0m - Remove or ignore command")
	fmt.Println("  \033[32mconfig commands reset\033[0m   - Reset commands list to default")

	fmt.Println("\033[1mExample:\033[0m")
	fmt.Println("  \033[32mconfig set openai apikey sk-your-api-key\033[0m")
	fmt.Println("  \033[32mconfig set general defaultshell /bin/zsh\033[0m")
	fmt.Println("  \033[32mconfig commands add mycommand\033[0m")
	fmt.Println("  \033[32mconfig commands remove ls\033[0m")
	fmt.Println("  \033[32mconfig save\033[0m")
	fmt.Println()
}

func isShellCommand(input string) bool {
	input = strings.ToLower(strings.TrimSpace(input)) // Trim and convert to lowercase
	words := strings.Fields(input)                    // Split into words

	if len(words) == 0 {
		return false
	}

	commands := config.GetShellCommands()
	for _, cmd := range commands {
		if cmd == words[0] {
			return true
		}
	}

	return false
}

func isAuroraCommand(input string) bool {
	return strings.Contains(strings.ToLower(input), "aurora")
}

// SetAIAgent sets the active AI agent
func SetAIAgent(agentType string) error {
	return AgentMgr.SetActiveAgent(AgentType(agentType))
}

// showShellCommands - show commands list
func showShellCommands() {
	commands := config.GetShellCommands()

	fmt.Println("\033[1mCommands list:\033[0m")

	// Show commands in 5 columns
	cols := 5
	for i, cmd := range commands {
		if i > 0 && i%cols == 0 {
			fmt.Println()
		}
		fmt.Printf("%-15s", cmd)
	}
	fmt.Println()

	// Additional commands
	if len(config.CurrentConfig.General.ShellCommands) > 0 {
		fmt.Println("\n\033[1mUser added commands:\033[0m")
		for _, cmd := range config.CurrentConfig.General.ShellCommands {
			fmt.Printf("  %s\n", cmd)
		}
	}

	// Ignored commands
	if len(config.CurrentConfig.General.IgnoredCommands) > 0 {
		fmt.Println("\n\033[1mIgnored commands:\033[0m")
		for _, cmd := range config.CurrentConfig.General.IgnoredCommands {
			fmt.Printf("  %s\n", cmd)
		}
	}

	fmt.Printf("\nTotal commands: %d\n", len(commands))
}

// addShellCommand - add new command
func addShellCommand(command string) {
	if err := config.AddShellCommand(command); err != nil {
		fmt.Printf("\033[31mError: %v\033[0m\n", err)
		return
	}

	fmt.Printf("\033[32m'%s' command added successfully\033[0m\n", command)
	fmt.Println("\033[33mNote: Remember to save changes using 'config save'\033[0m")
}

// removeShellCommand - remove command
func removeShellCommand(command string) {
	if err := config.RemoveShellCommand(command); err != nil {
		fmt.Printf("\033[31mError: %v\033[0m\n", err)
		return
	}

	fmt.Printf("\033[32m'%s' command removed successfully\033[0m\n", command)
	fmt.Println("\033[33mNote: Remember to save changes using 'config save'\033[0m")
}

// resetShellCommands - reset commands list to default
func resetShellCommands() {
	config.ResetShellCommands()
	fmt.Println("\033[32mCommands list reset to default\033[0m")
	fmt.Println("\033[33mNote: Remember to save changes using 'config save'\033[0m")
}
