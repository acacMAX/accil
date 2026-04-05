# Contributing to ACCIL

[English](#english) | [中文](#chinese)

---

<a name="english"></a>
## 🤝 How to Contribute

Thank you for your interest in contributing to ACCIL! This document provides guidelines and instructions for contributing.

### 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [How to Contribute](#how-to-contribute)
- [Development Setup](#development-setup)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Commit Message Guidelines](#commit-message-guidelines)
- [Reporting Bugs](#reporting-bugs)
- [Feature Requests](#feature-requests)

### Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code.

### Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/accil.git
   cd accil
   ```
3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/accil/accil.git
   ```
4. **Create a branch** for your changes:
   ```bash
   git checkout -b feature/your-feature-name
   ```

### How to Contribute

#### Reporting Bugs

Before creating bug reports, please check existing issues. When creating a bug report, include:

- **Clear title and description**
- **Steps to reproduce** the behavior
- **Expected vs actual behavior**
- **Screenshots or logs** if applicable
- **Environment information**:
  - OS (Windows/macOS/Linux)
  - Go version (`go version`)
  - ACCIL version

Example:
```markdown
**Describe the bug**
Tool execution results are not displayed in real-time.

**To Reproduce**
1. Run `accil`
2. Enter "Create a hello.txt file"
3. Observe that output only appears after completion

**Expected behavior**
Should see each tool call and result as it happens.

**Environment:**
- OS: Windows 11
- Go: 1.21.5
- ACCIL: 0.1.0
```

#### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

- **Use a clear title**
- **Describe the current behavior** and what you'd like to see
- **Explain why this would be useful**
- **Provide examples** if possible

#### Your First Code Contribution

Look for issues labeled [`good first issue`](https://github.com/accil/accil/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22) or [`help wanted`](https://github.com/accil/accil/issues?q=is%3Aissue+is%3Aopen+label%3A%22help+wanted%22) to get started.

### Development Setup

#### Prerequisites

- **Go 1.21+** ([Download](https://golang.org/dl/))
- **Git** ([Download](https://git-scm.com/downloads))
- A code editor (VS Code recommended)

#### Building from Source

```bash
# Clone the repository
git clone https://github.com/accil/accil.git
cd accil

# Install dependencies
go mod tidy

# Build
go build -o accil .

# Run tests
go test ./...

# Run with verbose logging
go run . --help
```

#### Project Structure

```
accil/
├── cmd/              # CLI entry point
├── internal/         # Internal packages
│   ├── ai/           # AI client
│   ├── config/       # Configuration
│   ├── context/      # Context management
│   ├── memory/       # Memory system
│   ├── session/      # Session handling
│   ├── tools/        # Tool implementations
│   ├── tui/          # Terminal UI
│   ├── quest/        # Quest mode
│   ├── agent/        # Sub-agents
│   └── review/       # Code review
└── main.go
```

### Pull Request Process

1. **Update your fork** with the latest changes from upstream:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Make your changes** following our coding standards

3. **Test your changes**:
   ```bash
   go test ./...
   go vet ./...
   go fmt ./...
   ```

4. **Commit your changes** following commit message guidelines

5. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

6. **Open a Pull Request**:
   - Use a clear title
   - Reference related issues (e.g., "Fixes #123")
   - Describe your changes in detail
   - Include screenshots for UI changes
   - List any breaking changes

### Coding Standards

#### Go Code Style

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use `gofmt` to format your code:
  ```bash
  gofmt -w .
  ```
- Run `go vet` to catch common mistakes:
  ```bash
  go vet ./...
  ```

#### Naming Conventions

- Use **camelCase** for local variables and functions
- Use **PascalCase** for exported names
- Use **snake_case** for file names
- Be descriptive with names

#### Comments

- Comment exported functions, types, and variables
- Write comments in English
- Keep comments clear and concise
- Use `//` for single-line comments

#### Error Handling

- Always check errors
- Return errors with context:
  ```go
  if err != nil {
      return fmt.Errorf("failed to load config: %w", err)
  }
  ```

### Commit Message Guidelines

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

#### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, semicolons, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

#### Examples

```
feat(tui): add real-time tool execution display

fix(config): handle missing API key gracefully

docs(readme): update installation instructions

refactor(ai): simplify chat message handling
```

### Review Process

All submissions require review. Reviewers will check for:

- Functionality and correctness
- Code quality and style
- Test coverage
- Documentation updates
- No breaking changes (unless justified)

---

<a name="chinese"></a>
## 🤝 如何贡献

感谢您对 ACCIL 的贡献！本文档提供了贡献的指南和说明。

### 📋 目录

- [行为准则](#行为准则)
- [开始使用](#开始使用)
- [如何贡献](#如何贡献-1)
- [开发环境设置](#开发环境设置)
- [Pull Request 流程](#pull-request-流程)
- [代码规范](#代码规范)
- [提交信息指南](#提交信息指南)
- [报告 Bug](#报告-bug)
- [功能请求](#功能请求)

### 行为准则

本项目及所有参与者都受我们的行为准则约束。参与即表示您应遵守此准则。

### 开始使用

1. **在 GitHub 上 Fork 仓库**
2. **克隆您的 fork** 到本地：
   ```bash
   git clone https://github.com/YOUR_USERNAME/accil.git
   cd accil
   ```
3. **添加上游远程仓库**：
   ```bash
   git remote add upstream https://github.com/accil/accil.git
   ```
4. **为您的更改创建分支**：
   ```bash
   git checkout -b feature/your-feature-name
   ```

### 如何贡献

#### 报告 Bug

在创建 Bug 报告之前，请检查是否已有相关问题。创建 Bug 报告时，请包含：

- **清晰的标题和描述**
- **重现步骤**
- **期望行为与实际行为**
- **截图或日志**（如果适用）
- **环境信息**：
  - 操作系统（Windows/macOS/Linux）
  - Go 版本（`go version`）
  - ACCIL 版本

示例：
```markdown
**描述 Bug**
工具执行结果未实时显示。

**重现步骤**
1. 运行 `accil`
2. 输入 "创建一个 hello.txt 文件"
3. 观察到输出仅在完成后出现

**期望行为**
应该在每次工具调用时立即看到结果。

**环境：**
- 操作系统: Windows 11
- Go: 1.21.5
- ACCIL: 0.1.0
```

#### 建议增强功能

增强功能建议作为 GitHub Issue 跟踪。创建建议时：

- **使用清晰的标题**
- **描述当前行为**和您希望看到的行为
- **解释为什么这会有用**
- **提供示例**（如果可能）

#### 您的第一次代码贡献

查找标记为 [`good first issue`](https://github.com/accil/accil/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22) 或 [`help wanted`](https://github.com/accil/accil/issues?q=is%3Aissue+is%3Aopen+label%3A%22help+wanted%22) 的问题以开始。

### 开发环境设置

#### 前置要求

- **Go 1.21+** ([下载](https://golang.org/dl/))
- **Git** ([下载](https://git-scm.com/downloads))
- 代码编辑器（推荐 VS Code）

#### 从源码构建

```bash
# 克隆仓库
git clone https://github.com/accil/accil.git
cd accil

# 安装依赖
go mod tidy

# 构建
go build -o accil .

# 运行测试
go test ./...

# 详细日志运行
go run . --help
```

#### 项目结构

```
accil/
├── cmd/              # CLI 入口点
├── internal/         # 内部包
│   ├── ai/           # AI 客户端
│   ├── config/       # 配置管理
│   ├── context/      # 上下文管理
│   ├── memory/       # 记忆系统
│   ├── session/      # 会话处理
│   ├── tools/        # 工具实现
│   ├── tui/          # 终端 UI
│   ├── quest/        # 任务模式
│   ├── agent/        # 子代理
│   └── review/       # 代码审查
└── main.go
```

### Pull Request 流程

1. **更新您的 fork**，获取上游最新更改：
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **进行更改**，遵循我们的代码规范

3. **测试您的更改**：
   ```bash
   go test ./...
   go vet ./...
   go fmt ./...
   ```

4. **提交更改**，遵循提交信息指南

5. **推送到您的 fork**：
   ```bash
   git push origin feature/your-feature-name
   ```

6. **打开 Pull Request**：
   - 使用清晰的标题
   - 引用相关问题（例如，"修复 #123"）
   - 详细描述您的更改
   - UI 更改请包含截图
   - 列出任何破坏性更改

### 代码规范

#### Go 代码风格

- 遵循 [Effective Go](https://go.dev/doc/effective_go)
- 使用 `gofmt` 格式化代码：
  ```bash
  gofmt -w .
  ```
- 运行 `go vet` 捕获常见错误：
  ```bash
  go vet ./...
  ```

#### 命名约定

- 局部变量和函数使用 **camelCase**
- 导出的名称使用 **PascalCase**
- 文件名使用 **snake_case**
- 名称应具有描述性

#### 注释

- 为导出的函数、类型和变量添加注释
- 使用英文编写注释
- 保持注释清晰简洁
- 使用 `//` 进行单行注释

#### 错误处理

- 始终检查错误
- 返回带有上下文的错误：
  ```go
  if err != nil {
      return fmt.Errorf("加载配置失败: %w", err)
  }
  ```

### 提交信息指南

我们遵循 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```
<类型>(<范围>): <描述>

[可选正文]

[可选脚注]
```

#### 类型

- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更改
- `style`: 代码风格更改（格式、分号等）
- `refactor`: 代码重构
- `test`: 添加或更新测试
- `chore`: 维护任务

#### 示例

```
feat(tui): 添加工具执行实时显示

fix(config): 优雅处理缺失的 API 密钥

docs(readme): 更新安装说明

refactor(ai): 简化聊天消息处理
```

### 审查流程

所有提交都需要审查。审查者将检查：

- 功能正确性
- 代码质量和风格
- 测试覆盖率
- 文档更新
- 无破坏性更改（除非有正当理由）

---

<div align="center">

**再次感谢您的贡献！** 🎉

</div>
