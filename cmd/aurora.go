package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

// Global agent manager instance
var AgentMgr *AgentManager

func init() {
	// Initialize the agent manager
	AgentMgr = NewAgentManager()
}

var shellCommands = []string{
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

// ProcessAuroraCommand handles Aurora-specific commands
func ProcessAuroraCommand(input string) bool {
	// Check if input contains "aurora"
	if isAuroraCommand(input) || !isShellCommand(input) {
		// Use streaming response
		err := AgentMgr.StreamQuery(input, os.Stdout)
		if err != nil {
			fmt.Printf("\nError querying AI agent: %v\n", err)
		} else {
			fmt.Println() // Add a newline after the streamed response
		}

		return true
	}
	return false
}

func isShellCommand(input string) bool {
	input = strings.ToLower(strings.TrimSpace(input)) // Trim and convert to lowercase
	words := strings.Fields(input)                    // Split into words

	if len(words) == 0 {
		return false
	}

	if slices.Contains(shellCommands, words[0]) {
		return true
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
