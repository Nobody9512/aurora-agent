package utils

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

// ActiveCmd stores the currently active command (for CTRL+C handling)
var ActiveCmd *exec.Cmd

// RunCommandWithPTY runs a command with PTY to preserve colors
func RunCommandWithPTY(cmd *exec.Cmd) {
	ptmx, err := pty.Start(cmd)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer ptmx.Close()

	// Store the active process
	ActiveCmd = cmd

	// Redirect PTY output to terminal
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				break
			}
			os.Stdout.Write(buf[:n])
		}
	}()

	// Wait for command to complete
	cmd.Wait()

	// Clear activeCmd after process completes
	ActiveCmd = nil
}
