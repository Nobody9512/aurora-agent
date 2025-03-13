package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/chzyer/readline"
)

// GetShellCommands retrieves available commands for tab completion
func GetShellCommands() []readline.PrefixCompleterInterface {
	out, err := exec.Command("bash", "-c", "compgen -c").Output()
	if err != nil {
		fmt.Println("Error: Problem getting commands.")
		return nil
	}

	commands := strings.Split(string(out), "\n")
	var completions []readline.PrefixCompleterInterface

	for _, cmd := range commands {
		if cmd != "" {
			completions = append(completions, readline.PcItem(cmd))
		}
	}

	extraCommands := []string{"exit", "quit", "clear"}
	for _, cmd := range extraCommands {
		completions = append(completions, readline.PcItem(cmd))
	}

	return completions
}

// GetDefaultShell determines the user's default shell
func GetDefaultShell() string {
	userShell := os.Getenv("SHELL")
	if userShell == "" {
		userShell = "/bin/bash" // Use bash if not detected
		fmt.Println("SHELL not detected, using bash")
	}
	return userShell
}
