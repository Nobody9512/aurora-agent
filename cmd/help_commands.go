package cmd

import (
	"fmt"
)

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
