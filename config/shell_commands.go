package config

import (
	"fmt"
	"strings"
)

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
