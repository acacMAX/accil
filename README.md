# ACCIL

<div align="center">

```
  █████╗ ██████╗ ██████╗  ██████╗██╗  ██╗██╗     ███████╗
 ██╔══██╗██╔══██╗██╔══██╗██╔════╝██║  ██║██║     ██╔════╝
███████║██████╔╝██████╔╝██║     ███████║██║     █████╗
██╔══██║██╔══██╗██╔══██╗██║     ██╔══██║██║     ██╔══╝
 ██║  ██║██████╔╝██████╔╝╚██████╗██║  ██║███████╗███████╗
 ╚═╝  ╚═╝╚═════╝ ╚═════╝  ╚═════╝╚═╝  ╚═╝╚══════╝╚══════╝
```

**AI-Powered Autonomous Coding Assistant**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Version](https://img.shields.io/badge/release-1.4.0-blue.svg)](CHANGELOG.md)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Windows%20|%20macOS%20|%20Linux-lightgrey)](https://github.com/acacMAX/accil)
[![Website](https://img.shields.io/badge/Website-cli.acz.qzz.io-blue)](https://cli.acz.qzz.io/)

[English](README.md) | [中文](README_zh.md) | [Website](https://cli.acz.qzz.io/)

</div>

---
## ⚠️ Project Description
1.This project uses AI programming.
2.This project was created on a whim and may not be maintained long-term. If you need a reliable and powerful AI programming tool, I recommend using [Qoder](https://qoder.com/).
3.I don't want to write documentation. All previous versions were uploaded by AI, so there will be AI traces.

Finally, thank you for using this project. If you like it, feel free to give me a Star 🌟
## ✨ Features

- 🗨️ **Interactive Chat Mode** - Modern terminal UI based on BubbleTea with scrolling support
- ⏹️ **ESC Stop Streaming** - Press ESC to immediately stop AI output during streaming
- 📝 **Multi-line Input Support** - Paste multi-line code and text directly, format preserved
- ⚡ **Autonomous Quest Mode** - Automatically plan and execute multi-step programming tasks
- 🔍 **Code Review** - Security vulnerabilities, performance issues, code quality detection
- 🤖 **Sub-Agent System** - 6 specialized agents with enhanced capabilities
- 🌐 **Remote Development** - Connect to remote servers via SSH for remote coding
- 🎨 **Retro Terminal Splash** - Animated boot sequence with vintage CRT aesthetics
- 🧠 **Enhanced AI Memory** - Code semantics, learning history, error pattern recognition
- 🌐 **Advanced Context** - Code relationship graph, function tracking, project analysis
- 💻 **Programming Assistant** - 10 core capabilities including security & performance engineering
- 📝 **File Operations** - Read, write, edit files with precise replacements
- 💻 **Command Execution** - Execute shell commands with cross-platform support
- 🔒 **Safety First** - Confirmation for dangerous operations, command blacklist support
- 💾 **Session Persistence** - Automatic conversation history saving
- 🔄 **Real-Time Visibility** - Step-by-step display of AI thinking and tool execution
- 📊 **API Usage Stats** - Display token usage statistics on exit

## 🚀 Quick Install

### Linux / macOS

```bash
curl -fsSL https://raw.githubusercontent.com/acacMAX/accil/main/install.sh | bash
```

Or

```bash
git clone https://github.com/acacMAX/accil.git
cd accil
chmod +x install.sh
./install.sh
```

### Windows

**Option 1: PowerShell (Recommended)**

```powershell
# Download and run
[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.SecurityProtocolType]::Tls12
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/acacMAX/accil/main/install.ps1" -OutFile "$env:TEMP\accil-install.ps1"
& "$env:TEMP\accil-install.ps1"
```

**Option 2: Batch File**

Download [`install.bat`](https://raw.githubusercontent.com/acacMAX/accil/main/install.bat) and double-click to run

### Manual Installation

```bash
# Clone repository
git clone https://github.com/acacMAX/accil.git
cd accil

# Install dependencies
go mod tidy

# Build
go build -o accil .

# Run
./accil
```

**For Windows users:**
```powershell
# Build and install globally (recommended)
build.bat

# Or build manually
go build -o accil.exe .
.\accil.exe

# Or run directly (recommended for beginners)
go run .
```

## 📖 Usage

### Interactive Mode

```bash
# Start interactive session
accil

# Specify working directory
accil --workdir ./myproject

# Continue last session
accil --continue

# Auto-approve mode (skip confirmations)
accil --yolo
```

### Single-Shot Execution

```bash
# Execute a single task
accil "Read main.go and explain what it does"

# Create a file
accil "Create a hello world program in Python"

# Headless mode (for script integration)
accil --headless "Refactor this function"
```

### Built-in Commands

Type these in interactive mode:

| Command | Description |
|---------|-------------|
| `/help` | Show help message |
| `/clear` | Clear conversation |
| `/quit` | Exit program |
| `/chat` | Enter chat mode |
| `/quest` | Enter quest mode |
| `/review` | Enter review mode |
| `/agent` | Enter agent mode |
| `/remote` | Enter remote development mode |
| `/model <name>` | Change AI model |
| `/context` | Show current context |

### Keyboard Shortcuts

| Shortcut | Description |
|----------|-------------|
| `Ctrl+C` | Quit (press twice to confirm) |
| `Ctrl+L` | Clear screen |
| `Ctrl+S` / `F5` | Send message |
| `Ctrl+H` | Show/hide help |
| `ESC` | Stop current thinking/output |
| `Enter` | Insert newline (multi-line paste supported) |
| `PgUp/PgDn` | Page up/down |
| `Mouse Wheel` | Scroll messages |
| `Shift+Mouse Drag` | Select and copy text |

> **Tips**:
> - Paste multi-line text directly, format will be preserved automatically
> - Press `Ctrl+S` or `F5` to send your message

## 🔧 Configuration

### First Run

On first run, an interactive setup wizard will guide you through:
- API provider selection (OpenAI, DeepSeek, Anthropic, Ollama, etc.)
- API Key input
- Model selection

### Configuration File

Configuration is stored at `~/.accil/config.yaml`:

```yaml
api_key: "your-api-key"
base_url: "https://api.openai.com/v1"
model: "gpt-5.5"
max_tokens: 4096
auto_approve: false
block_list:
  - "rm -rf /"
  - "rm -rf /*"
  - "mkfs"
max_tool_calls: 30
command_timeout: 120

# Remote development configuration
remote:
  host: "your-server.com"
  port: "22"
  user: "username"
  key_path: "~/.ssh/id_rsa"
  workdir: "/home/user/project"
  use_agent: true
```

### Environment Variables

```bash
export AI_API_KEY="your-api-key"
export AI_BASE_URL="https://api.openai.com/v1"
```

## 🌐 Supported API Providers

| Provider | Base URL | Recommended Models |
|----------|----------|-------------------|
| OpenAI | `https://api.openai.com/v1` | gpt-5.5, gpt-5.4-mini, gpt-5.4-nano |
| DeepSeek | `https://api.deepseek.com/v1` | deepseek-v4-pro, deepseek-v4-flash |
| Anthropic | `https://api.anthropic.com/v1` | claude-sonnet-4-6, claude-opus-4-1 |
| Qwen | `https://dashscope.aliyuncs.com/compatible-mode/v1` | qwen3-max, qwen-plus, qwen-turbo |
| Zhipu AI | `https://open.bigmodel.cn/api/paas/v4` | glm-4.6, glm-4-plus |
| Ollama (Local) | `http://localhost:11434/v1` | qwen2.5-coder, llama3.3, mistral |

> DeepSeek note: `deepseek-chat` and `deepseek-reasoner` are legacy aliases scheduled for deprecation on `2026-07-24`.

## 🛠️ Tool System

AI can invoke the following tools:

| Tool | Description | Requires Confirmation |
|------|-------------|----------------------|
| `read_file` | Read file contents | No |
| `write_file` | Write/create files | Yes |
| `edit_file` | Precise content replacement | Yes |
| `run_command` | Execute shell commands | Yes |
| `list_dir` | List directory contents | No |
| `search_code` | Regex search in code | No |
| `glob` | File pattern matching | No |
| `web_search` | Search the web for information | No |
| `web_fetch` | Fetch content from a URL | No |

## 🌐 Remote Development

Connect to remote servers via SSH and develop directly on them:

### Quick Start

```bash
# Connect to remote server
accil remote user@hostname

# Or enter remote mode in interactive session
accil
/remote
/remote connect hostname
```

### Remote Tools

When connected to a remote server, all file operations work remotely:

| Tool | Description |
|------|-------------|
| `read_file` | Read remote file contents |
| `write_file` | Write to remote files |
| `edit_file` | Edit remote files |
| `run_command` | Execute commands on remote server |
| `list_dir` | List remote directory contents |
| `search_code` | Search code on remote server |
| `glob` | Match remote files |

### Authentication Methods

The remote client tries authentication in this order:
1. SSH Agent (if `use_agent: true`)
2. Private key file (specified by `key_path`)
3. Default SSH keys (`~/.ssh/id_rsa`, `~/.ssh/id_ed25519`)
4. Password (if configured)

## 🔒 Safety Mechanisms

- **Safe by Default**: All file writes and command executions require user confirmation
- **YOLO Mode**: Use `--yolo` flag to skip all confirmations (warning: use at your own risk)
- **Command Blacklist**: Dangerous commands are always blocked
- **Network Retry**: API calls automatically retry up to 3 times on failure

## 📁 Project Structure

```
accil/
├── cmd/                    # Command-line entry point
│   └── root.go
├── internal/
│   ├── ai/                 # AI client
│   ├── config/             # Configuration management
│   ├── context/            # Context memory
│   ├── memory/             # Project memory
│   ├── remote/             # Remote SSH client
│   ├── session/            # Session management
│   ├── tools/              # Tool system
│   ├── tui/                # Terminal UI
│   ├── quest/              # Autonomous quests
│   ├── agent/              # Sub-agents
│   └── review/             # Code review
├── main.go
├── go.mod
├── build.bat               # Windows build script
├── install.sh              # Linux/macOS installation script
├── install.bat             # Windows installation script
├── Makefile
├── LICENSE
├── README.md               # English Documentation
└── README_zh.md            # Chinese Documentation
```

## 🤝 Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## 📄 License

[MIT License](LICENSE)

---

<div align="center">

**If this project helps you, please give it a ⭐ Star!**

Made with ❤️ by the ACCIL Team

</div>
