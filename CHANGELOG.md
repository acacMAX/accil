# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
