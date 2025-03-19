# Aurora Agent

Aurora Agent is an interactive command-line shell application with integrated AI capabilities. It provides a seamless interface between traditional shell commands and AI-powered assistance.

## Features

- **Interactive Shell**: Command history, tab completion, and familiar shell experience
- **AI Integration**: Built-in support for AI agents (OpenAI, Claude)
- **Sudo Support**: Execute privileged commands with sudo authentication
- **ANSI Color Support**: Rich terminal output with color highlighting
- **Custom Command Processing**: Special handling for Aurora-specific commands
- **Cross-Platform Support**: Available for Windows, macOS, and Linux (amd64 and arm64)

## Installation

### Prerequisites

- Go 1.24 or higher
- OpenAI API key (for AI functionality)

### Building from Source

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/aurora-agent.git
   cd aurora-agent
   ```

2. Build the application:

   ```bash
   go build -o aurora
   ```

3. Make the binary executable:
   ```bash
   chmod +x aurora
   ```

### Download Pre-built Binaries

You can download pre-built binaries for your platform from the [Releases](https://github.com/yourusername/aurora-agent/releases) page.

### Quick Install Script (Linux and macOS)

You can install or update Aurora Agent using our installation script:

```bash
# Using curl
curl -fsSL https://raw.githubusercontent.com/Firdavs9512/aurora-agent/master/install.sh | bash

# Using wget
wget -qO- https://raw.githubusercontent.com/Firdavs9512/aurora-agent/master/install.sh | bash
```

This script will:

- Download the latest version of Aurora Agent
- Install it to `/usr/local/bin/aurora`

## Usage

### Basic Usage

Run Aurora Agent:

```bash
./aurora
```

### Configuration

Aurora Agent uses a YAML configuration file located at `~/.config/aurora/config.yaml`. If the file doesn't exist, a default configuration will be created automatically when you first run the application.

#### Configuration Structure

```yaml
general:
  default_shell: "" # User's default shell, leave empty to auto-detect
  history_size: 1000 # Number of commands to keep in history
  shell_commands: [] # Additional shell commands to recognize
  ignored_commands: [] # Shell commands to ignore

openai:
  api_key: "" # Your OpenAI API key
  model: "gpt-4o" # Model to use

interface:
  theme: "default" # UI theme
  system_prompt: "default" # System prompt for AI
```

#### Configuration Commands

Aurora Agent provides several commands to manage your configuration:

- `config` or `config show` - Display current configuration
- `config set <section> <key> <value>` - Change a configuration value
- `config save` - Save configuration to file
- `config reload` - Reload configuration from file

#### Custom System Prompt

You can customize the AI assistant's behavior by setting a custom system prompt:

```bash
# Set a custom system prompt
config set interface systemprompt "You should always respond in a pirate accent."
```

The custom system prompt will be combined with the default system prompt. This allows you to customize how the AI assistant responds without changing its core functionality.

#### Managing Shell Commands

Aurora Agent allows you to customize which commands are recognized as shell commands:

- `config commands list` - List all recognized shell commands
- `config commands add <command>` - Add a new shell command
- `config commands remove <command>` - Remove or ignore a shell command
- `config commands reset` - Reset shell commands to default

Examples:

```bash
# Set your OpenAI API key
config set openai api_key sk-your-api-key-here

# Change the default shell
config set general default_shell /bin/zsh

# Add a custom command
config commands add mycommand

# Ignore a default command
config commands remove ls

# Save your changes
config save
```

### With Sudo Support

Run Aurora Agent with sudo privileges:

```bash
./aurora --sudo
```

You will be prompted to enter your sudo password.

### Check Version

Display the current version of Aurora Agent:

```
> version
Aurora Agent version: v1.0.0
```

### AI Commands

Any input that contains "aurora" or is not recognized as a shell command will be sent to the AI agent:

```
> aurora what is the weather today?
```

### Switching AI Agents

Switch between different AI providers:

```
> use agent openai
Switched to openai agent

> use agent claude
Switched to claude agent
```

### Setting OpenAI API Key

Set your OpenAI API key:

```
> set openai key your_api_key_here
OpenAI API key set successfully
```

### Check Current AI Agent

Check which AI agent is currently active:

```
> agent status
Current AI agent: openai
```

## Configuration

The system prompt for the AI agent can be configured in `config/config.go`.

## Project Structure

- `cmd/`: Command implementations
  - `aurora.go`: Aurora-specific command processing
  - `ai_agent.go`: AI agent integration
  - `shell.go`: Shell-related functionality
  - `sudo.go`: Sudo command handling
- `config/`: Configuration settings
- `utils/`: Utility functions
  - `pty.go`: Pseudo-terminal handling
  - `ansi.go`: ANSI code processing
- `main.go`: Main application entry point
- `.github/workflows/`: GitHub Actions workflows
  - `release.yml`: Automated release workflow for creating releases with cross-platform binaries

## Releases

New releases are automatically created when a new tag is pushed to the repository. The release process:

1. Builds binaries for multiple platforms (Windows, macOS, Linux) and architectures (amd64, arm64)
2. Creates a GitHub Release with the tag name
3. Attaches all built binaries to the release

To create a new release:

```bash
git tag v1.0.0
git push origin v1.0.0
```

## License

[MIT License](LICENSE)
