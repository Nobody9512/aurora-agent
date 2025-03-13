package cmd

import (
	"fmt"
	"os/exec"
)

// CheckSudoPassword verifies if the provided sudo password is correct
func CheckSudoPassword(password string) bool {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("echo %s | sudo -S -v", password))
	err := cmd.Run()
	return err == nil
}

// CreateSudoCommand creates a command with sudo privileges
func CreateSudoCommand(shell string, password string, args []string) *exec.Cmd {
	return exec.Command(shell, "-i", "-c", fmt.Sprintf("echo %s | sudo -S -p '' %s", password, args))
}
