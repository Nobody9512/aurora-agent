# Aurora Agent

Aurora Agent is an interactive command-line shell application with integrated AI capabilities. It provides a seamless interface between traditional shell commands and AI-powered assistance.

## Features

- **Interactive Shell**: Command history, tab completion, and familiar shell experience
- **AI Integration**: Built-in support for AI agents (OpenAI, Claude)
- **Sudo Support**: Execute privileged commands with sudo authentication
- **ANSI Color Support**: Rich terminal output with color highlighting
- **Custom Command Processing**: Special handling for Aurora-specific commands

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

## Usage

### Basic Usage

Run Aurora Agent:

```bash
./aurora
```

### With Sudo Support

Run Aurora Agent with sudo privileges:

```bash
./aurora --sudo
```

You will be prompted to enter your sudo password.

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

## License

[MIT License](LICENSE)