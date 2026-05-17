# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.4.6] - 2026-05-16

### Added

- Unified interactive runtime configuration under `/config`.
- `/config show`, `/config provider`, `/config model`, and `/config baseurl` are now documented as the primary config workflow.

### Changed

- Updated app version, release metadata, and installer defaults to `v1.4.6`.
- User documentation now points to the official website instead of GitHub-facing install and community links.
- Chinese README was refreshed to match the current command set and release messaging.

### Deprecated

- `/provider`, `/model`, and `/baseurl` remain available as compatibility aliases, but `/config ...` is now the main entry point.

## [1.4.0] - 2026-05-10

### Added

- In-app config commands for changing provider, model, base URL, and editing config without rerunning the first-time wizard.
- Interactive `/config` display for checking the active provider, base URL, model, and whether tool calling is enabled.

### Changed

- Refreshed documented model recommendations, including the latest DeepSeek model names and legacy-alias notes.
- Updated release metadata, installers, and version badges to `v1.4.0`.
- Security policy now clearly documents the currently supported release lines in both English and Chinese.
- README and quick-start docs now explain how to change provider and model after the initial setup.

### Fixed

- Improved request flow for DeepSeek follow-up responses that previously failed with HTTP 400 in multi-turn use.
- Restored TUI loading feedback while the model is thinking before output starts.
- Fixed UI state regressions where the interface could disappear after a task completed.
- Improved task execution stability around post-task state transitions.

## [1.3.5] - 2026-05-05

### Security

- SSH remote: enforce host key verification via `known_hosts` by default; unknown keys require an explicit `knownhosts` line (no silent `InsecureIgnoreHostKey` fallback).
- Remote `run_command`: apply the same configurable `block_list` as local tools; add regex-based dangerous-command heuristics in the remote executor.

### Added

- Root `VERSION` file for reproducible builds and release tagging.
- `--headless` with no CLI args reads the prompt from **stdin** (script/pipe friendly).

### Changed

- Session storage directory unified to `~/.accil/sessions`.
- Default and documented API model identifiers refreshed (OpenAI / Anthropic / Qwen / Zhipu / Ollama examples).
- TUI splash line shows **`v1.3.5`**.
- Markdown rendering width follows the terminal width.
- `ReviewProject` collects source files in pure Go (no Unix-only `find | head` on Windows).
- Makefile, `build.bat`, `install.bat`, `install.ps1`, and `install.sh` embed version with `-X github.com/accil/accil/cmd.Version=...` and use `-buildvcs=false` where appropriate.

### Fixed

- Quest CLI progress printed incorrect step counters.
- Remote file read error path used a non-constant `fmt.Errorf` format string (vet/build hygiene).

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

[1.4.6]: CHANGELOG.md
[1.4.0]: CHANGELOG.md
[1.3.5]: CHANGELOG.md
[1.3.0]: CHANGELOG.md
[0.1.0]: CHANGELOG.md
