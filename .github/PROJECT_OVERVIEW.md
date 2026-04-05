# ACCIL Project Overview

## 🎯 Project Vision

ACCIL (AI-powered Coding Assistant CLI) is a powerful terminal-based AI programming assistant that helps developers write, review, and manage code through natural language interaction.

## 🏗️ Architecture

### Core Components

```
accil/
├── cmd/                    # Command-line interface entry
│   └── root.go            # Main command routing and application logic
│
├── internal/              # Internal packages (not exported)
│   ├── ai/                # AI client for API communication
│   │   └── client.go      # OpenAI-compatible API client with retry logic
│   │
│   ├── config/            # Configuration management
│   │   ├── config.go      # Config loading/saving
│   │   └── wizard.go      # Interactive setup wizard
│   │
│   ├── context/           # Project context awareness
│   │   └── context.go     # Context collection and relevance calculation
│   │
│   ├── memory/            # Long-term project memory
│   │   └── memory.go      # AGENTS.md management
│   │
│   ├── session/           # Conversation persistence
│   │   └── session.go     # Session save/load functionality
│   │
│   ├── tools/             # Tool execution system
│   │   └── tools.go       # 7 tool implementations
│   │
│   ├── tui/               # Terminal user interface
│   │   ├── app.go         # BubbleTea TUI model
│   │   └── markdown.go    # Markdown rendering
│   │
│   ├── quest/             # Autonomous task execution
│   │   └── quest.go       # Multi-step planning and execution
│   │
│   ├── agent/             # Specialized sub-agents
│   │   └── agent.go       # 6 agent types (coder, reviewer, etc.)
│   │
│   └── review/            # Code review system
│       └── review.go      # Security and quality analysis
│
└── main.go                # Application entry point
```

## 🔄 Data Flow

1. **User Input** → TUI captures message
2. **Session Manager** → Adds to conversation history
3. **Context Manager** → Enriches with project context
4. **AI Client** → Sends to LLM with tools
5. **Tool Executor** → Executes tool calls
6. **Message Handler** → Streams results back to TUI
7. **Display Update** → Real-time UI refresh

## 🛠️ Technology Stack

| Component | Technology | Purpose |
|-----------|-----------|---------|
| Language | Go 1.21+ | Type-safe, concurrent, fast |
| CLI Framework | Cobra | Command routing |
| TUI Library | BubbleTea | Terminal UI framework |
| Viewport | bubbles/viewport | Scrollable message display |
| Text Input | bubbles/textinput | User input field |
| Styling | lipgloss | Terminal styling |
| YAML | gopkg.in/yaml.v3 | Configuration parsing |
| HTTP | net/http | API communication |
| Markdown | Custom renderer | Message formatting |

## 🔑 Key Features Implementation

### Real-Time Execution Display
- **Goroutines**: Async message processing
- **Channels**: Stream updates from backend to UI
- **Message Types**: ProcessingUpdate, ToolCallMessage, ToolResultMessage
- **Immediate Feedback**: Each step visible as it happens

### Tool System
1. `read_file` - Read file contents
2. `write_file` - Create/write files
3. `edit_file` - Precise content replacement
4. `run_command` - Execute shell commands
5. `list_dir` - List directory contents
6. `search_code` - Regex search in code
7. `glob` - File pattern matching

### Safety Mechanisms
- Confirmation required for write/command operations
- Command blacklist checking
- YOLO mode (`--yolo`) to skip confirmations
- Network retry with exponential backoff

### Multi-Provider Support
- OpenAI-compatible API interface
- Configurable base URL
- Environment variable overrides
- Proxy support via environment variables

## 📊 Development Status

### Completed ✅
- Core TUI with real-time updates
- Tool execution system
- Session management
- Configuration wizard
- Quest mode
- Sub-agent system
- Code review
- Context memory
- Installation scripts
- Documentation (EN/ZH)
- CI/CD workflows

### Planned 🚧
- Plugin system
- Web dashboard
- Team collaboration
- Advanced analytics
- Model fine-tuning support

## 🤝 Contributing

We welcome contributions! Please see:
- [CONTRIBUTING.md](../CONTRIBUTING.md) - How to contribute
- [CODE_OF_CONDUCT.md](../CODE_OF_CONDUCT.md) - Community guidelines
- [SECURITY.md](../SECURITY.md) - Security policy

## 📄 License

MIT License - See [LICENSE](../LICENSE) for details

---

<div align="center">

Built with ❤️ by the ACCIL Team

</div>
