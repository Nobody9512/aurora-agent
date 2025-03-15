package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/chzyer/readline"
	"golang.org/x/term"

	"aurora-agent/cmd"
	"aurora-agent/utils"
)

var sudoPassword string
var sudoEnabled bool

// Global signal channel
var sigs chan os.Signal

func init() {
	// Create a single signal channel (to avoid multiple calls)
	sigs = make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	// Catch CTRL+C signal (for all processes)
	go func() {
		for range sigs {
			if utils.ActiveCmd != nil {
				fmt.Println("\n[!] Process terminated")
				utils.ActiveCmd.Process.Signal(syscall.SIGINT) // Only kill the active process
			}
		}
	}()
}

func main() {
	// Determine user's default shell
	userShell := cmd.GetDefaultShell()

	// Check for --sudo flag
	if len(os.Args) > 1 && os.Args[1] == "--sudo" {
		fmt.Print("[sudo] Enter password: ")

		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			fmt.Println("Error: Could not read password.")
			os.Exit(1)
		}
		sudoPassword = strings.TrimSpace(string(bytePassword))
		sudoEnabled = true

		// Verify password
		if !cmd.CheckSudoPassword(sudoPassword) {
			fmt.Println("Error: Incorrect password!")
			os.Exit(1)
		}

		fmt.Println("Sudo mode activated!")
	}

	// Readline settings (with Tab completion)
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          getPrompt(),
		HistoryFile:     "/tmp/shell-history.tmp",
		AutoComplete:    readline.NewPrefixCompleter(cmd.GetShellCommands()...),
		InterruptPrompt: "^C",
	})
	if err != nil {
		fmt.Println("Error: Could not read terminal.")
		os.Exit(1)
	}
	defer rl.Close()

	for {
		input, err := rl.Readline()
		if err != nil {
			fmt.Println("\nExiting program.")
			break
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// Handle AI agent commands
		if handleAIAgentCommands(input) {
			continue
		}

		// Process Aurora commands
		if cmd.ProcessAuroraCommand(input) {
			continue
		}

		if input == "exit" || input == "quit" {
			fmt.Println("Exiting program.")
			break
		}

		args := strings.Fields(input)

		// Handle cd command specially
		if args[0] == "cd" {
			// Default to home directory if no argument is provided
			if len(args) == 1 {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					fmt.Println("Error:", err)
					continue
				}
				err = os.Chdir(homeDir)
				if err != nil {
					fmt.Println("Error:", err)
				} else {
					rl.SetPrompt(getPrompt())
				}
				continue
			}
			
			// Handle cd with path argument
			err := os.Chdir(args[1])
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				rl.SetPrompt(getPrompt())
			}
			continue
		}

		// Check for sudo
		if args[0] == "sudo" {
			if !sudoEnabled {
				fmt.Println("Sudo mode not enabled. Start the program with --sudo or enter password.")
				continue
			}

			args = args[1:]
			command := exec.Command(userShell, "-i", "-c", fmt.Sprintf("echo %s | sudo -S -p '' %s", sudoPassword, strings.Join(args, " ")))
			utils.RunCommandWithPTY(command)
			continue
		}

		// Run in shell environment (to preserve colors)
		command := exec.Command(userShell, "-i", "-c", input)
		utils.RunCommandWithPTY(command)
	}
}

// getPrompt returns a prompt string with the current working directory
func getPrompt() string {
	pwd, err := os.Getwd()
	if err != nil {
		return "> "
	}
	return pwd + " > "
}

// handleAIAgentCommands handles commands related to AI agents
func handleAIAgentCommands(input string) bool {
	// Check for agent switching command
	if strings.HasPrefix(input, "use agent") {
		parts := strings.Fields(input)
		if len(parts) < 3 {
			fmt.Println("Usage: use agent <agent_type>")
			fmt.Println("Available agents: openai, claude")
			return true
		}

		agentType := parts[2]
		err := cmd.SetAIAgent(agentType)
		if err != nil {
			fmt.Printf("Error setting agent: %v\n", err)
		} else {
			fmt.Printf("Switched to %s agent\n", agentType)
		}
		return true
	}

	// Check for setting OpenAI API key
	if strings.HasPrefix(input, "set openai key") {
		parts := strings.Fields(input)
		if len(parts) < 4 {
			fmt.Println("Usage: set openai key <your_api_key>")
			return true
		}

		apiKey := parts[3]
		// Set the API key in environment variable
		os.Setenv("OPENAI_API_KEY", apiKey)

		// Reinitialize the agent manager to use the new key
		cmd.AgentMgr = cmd.NewAgentManager()

		fmt.Println("OpenAI API key set successfully")
		return true
	}

	// Check for agent status command
	if input == "agent status" {
		fmt.Printf("Current AI agent: %s\n", cmd.AgentMgr.GetActiveAgentName())
		return true
	}

	return false
}