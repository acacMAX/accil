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
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Windows%20|%20macOS%20|%20Linux-lightgrey)](https://github.com/acacMAX/accil)

[English](README.md) | [中文](README_zh.md)

</div>

---

## ✨ Features

- 🗨️ **Interactive Chat Mode** - Modern terminal UI based on BubbleTea with scrolling support
- ⚡ **Autonomous Quest Mode** - Automatically plan and execute multi-step programming tasks
- 🔍 **Code Review** - Security vulnerabilities, performance issues, code quality detection
- 🤖 **Sub-Agent System** - Specialized agents: Coder, Reviewer, Architect, Tester, Debugger
- 📝 **File Operations** - Read, write, edit files with precise replacements
- 💻 **Command Execution** - Execute shell commands with cross-platform support
- 🧠 **Context Memory** - Project-aware assistance with automatic memory management
- 🔒 **Safety First** - Confirmation for dangerous operations, command blacklist support
- 💾 **Session Persistence** - Automatic conversation history saving
- 🔄 **Real-Time Visibility** - Step-by-step display of AI thinking and tool execution

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
| `/quest` | Enter quest mode |
| `/review` | Enter review mode |
| `/agent` | Enter agent mode |
| `/model <name>` | Change AI model |

### Keyboard Shortcuts

| Shortcut | Description |
|----------|-------------|
| `Ctrl+C` | Quit |
| `Ctrl+L` | Clear screen |
| `↑/↓` | Browse history / Scroll messages |
| `PgUp/PgDn` | Page up/down |
| `Mouse Wheel` | Scroll messages |
| `Shift+Mouse Drag` | Select and copy text |

> **Tip**: Since mouse scrolling is enabled, direct mouse dragging is captured by the terminal. Hold `Shift` while dragging to select text for copying.

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
model: "gpt-4o"
max_tokens: 4096
auto_approve: false
block_list:
  - "rm -rf /"
  - "rm -rf /*"
  - "mkfs"
max_tool_calls: 30
```

### Environment Variables

```bash
export AI_API_KEY="your-api-key"
export AI_BASE_URL="https://api.openai.com/v1"
```

## 🌐 Supported API Providers

| Provider | Base URL | Recommended Models |
|----------|----------|-------------------|
| OpenAI | `https://api.openai.com/v1` | gpt-4o, gpt-4-turbo |
| DeepSeek | `https://api.deepseek.com/v1` | deepseek-chat, deepseek-coder |
| Anthropic | `https://api.anthropic.com/v1` | claude-3-opus, claude-3-sonnet |
| Qwen | `https://dashscope.aliyuncs.com/compatible-mode/v1` | qwen-turbo, qwen-max |
| Zhipu AI | `https://open.bigmodel.cn/api/paas/v4` | glm-4 |
| Ollama (Local) | `http://localhost:11434/v1` | llama3, codellama, mistral |

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
│   ├── session/            # Session management
│   ├── tools/              # Tool system
│   ├── tui/                # Terminal UI
│   ├── quest/              # Autonomous quests
│   ├── agent/              # Sub-agents
│   └── review/             # Code review
├── main.go
├── go.mod
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
