package cmd

import (
	"fmt"
	"strings"

	"aurora-agent/config"
)

// processConfigCommand handles all configuration related commands
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
			// if api key is changed, reload agent
			AgentMgr = NewAgentManager()
		case "model":
			config.CurrentConfig.OpenAI.Model = value
			fmt.Printf("\033[32mOpenAI.Model = %s\033[0m\n", value)
			// Update the model in the active agent
			if agent, ok := AgentMgr.activeAgent.(*OpenAIAgent); ok {
				agent.SetModel(value)
			}
		default:
			fmt.Printf("\033[31mError: '%s' key not found in OpenAI section\033[0m\n", key)
		}

	case "interface":
		switch strings.ToLower(key) {
		case "theme":
			config.CurrentConfig.Interface.Theme = value
			fmt.Printf("\033[32mInterface.Theme = %s\033[0m\n", value)
		case "systemprompt":
			config.CurrentConfig.Interface.SystemPrompt = value
			fmt.Printf("\033[32mInterface.SystemPrompt = %s\033[0m\n", value)
			// if system prompt is changed, reload agent
			AgentMgr = NewAgentManager()
		default:
			fmt.Printf("\033[31mError: '%s' key not found in Interface section\033[0m\n", key)
		}

	default:
		fmt.Printf("\033[31mError: '%s' section not found. Available sections: General, OpenAI, Interface\033[0m\n", section)
	}

	fmt.Println("\033[33mNote: Remember to save changes using 'config save'\033[0m")
}
