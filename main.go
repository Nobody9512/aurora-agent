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
		Prompt:          "> ",
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

		// Process Aurora commands
		if cmd.ProcessAuroraCommand(input) {
			continue
		}

		if input == "exit" || input == "quit" {
			fmt.Println("Exiting program.")
			break
		}

		args := strings.Fields(input)

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
