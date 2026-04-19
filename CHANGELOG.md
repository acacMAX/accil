# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.3.0] - 2026-04-12

### Added
- **ESC Key Stop Functionality** - Press ESC to immediately stop AI streaming output
- **Enhanced AI Programming Capabilities** - 10 core programming capabilities including security, performance engineering, and modern DevOps practices
- **Upgraded Agent System Prompts** - All 6 sub-agents (coder, reviewer, architect, tester, debugger, researcher) now have detailed, professional capability descriptions
- **Keyboard Event Priority Handling** - Fixed keyboard shortcuts (Ctrl+C, ESC, scroll) during AI output streaming

### Fixed
- ESC key not working during AI streaming output
- Keyboard shortcuts becoming unresponsive during streaming
- Channel blocking issues in message processing loop

## [1.2.5] - 2026-04-12

### Added
- **Enhanced AI Memory System** - Code semantics memory, learning history tracking, error pattern recognition
- **Advanced Context Management** - Code relationship graph, function info tracking, intelligent project analysis
- **Upgraded Programming Capabilities** - Code analysis, architecture design, refactoring, debugging support
- **Retro Terminal Splash Screen** - Animated boot sequence with CRT scanline effects
- **Interactive Remote Login Form** - Form-based SSH connection setup with visual feedback

### Changed
- Improved AI system prompts with detailed programming capability descriptions
- Enhanced context module with code graph and dependency tracking
- Memory module now supports semantic code understanding

## [0.3.0] - 2026-04-06

### Added
- **Multi-line Input Support** - Paste multi-line code and text directly with format preserved
- `build.bat` script for Windows - build and install globally with one command
- Textarea component replaces textinput for better multi-line editing
- **Command Timeout** - Commands now timeout after 120 seconds by default (configurable)
- `command_timeout` config option to set custom timeout

### Changed
- **Keyboard shortcuts updated**:
  - `Enter` now inserts newline (for multi-line input)
  - `Ctrl+S` or `F5` sends the message
- Improved input field with 2-line default height
- Updated documentation with new keyboard shortcuts

## [0.2.0] - 2026-04-06

### Added
- **Web Search Tool** - AI can now search the web for information using DuckDuckGo
- **Web Fetch Tool** - AI can fetch and read content from URLs
- Documentation for Shift+drag to select/copy text in terminal

### Changed
- Security response time updated to 7×24 hours
- Security contact email updated to acac74151@gmail.com

## [Unreleased]

### Added
- Real-time tool execution visibility with step-by-step display
- Streaming message updates in interactive TUI mode
- Tool call logging with emoji indicators (🔧 ✅ ❌)
- Processing status messages showing current AI operation round
- Enhanced error handling with retry mechanism (up to 3 attempts)
- HTTP proxy support via environment variables
- Installation scripts for Windows (install.bat, install.ps1) and Linux/macOS (install.sh)
- GitHub Actions CI/CD workflows
- Comprehensive documentation in English and Chinese
- CODE_OF_CONDUCT.md and SECURITY.md
- CONTRIBUTING.md with bilingual guidelines

### Changed
- Increased default max tool calls from 10 to 30
- Made tool call limit configurable in config.yaml
- Improved TUI rendering with minimum bounds checking
- Refactored message handling for better real-time updates
- Updated README with installation instructions

### Fixed
- Messages displaying twice in UI
- TUI viewport panic with slice index out of range
- Network timeout issues with extended timeout (300s)
- Unused variable compilation errors
- Missing import statements

## [0.1.0] - 2026-04-05

### Added
- Initial release of ACCIL
- Interactive chat mode with BubbleTea TUI
- Autonomous Quest mode for multi-step tasks
- Code review functionality
- Sub-agent system (coder, reviewer, architect, tester, debugger, researcher)
- Context memory management
- Session persistence
- First-time setup wizard
- Support for multiple API providers (OpenAI, DeepSeek, Anthropic, Ollama, etc.)
- Tool system with 7 tool types
- Command blacklist for safety
- Cross-platform support (Windows, macOS, Linux)

[Unreleased]: https://github.com/accil/accil/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/accil/accil/releases/tag/v0.1.0
