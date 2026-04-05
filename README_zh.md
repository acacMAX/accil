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

**AI驱动的自主编程助手**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Windows%20|%20macOS%20|%20Linux-lightgrey)](https://github.com/accil/accil)

[English](README.md) | [中文](README_zh.md)

</div>

---

## ✨ 功能特性

- 🗨️ **交互对话模式** - 基于 BubbleTea 的现代化终端界面，支持滚动和中文
- ⚡ **自主任务模式 (Quest)** - 自动规划并执行多步骤编程任务
- 🔍 **代码审查** - 安全漏洞、性能问题、代码质量检测
- 🤖 **子代理系统** - 专业化代理：编码、审查、架构、测试、调试
- 📝 **文件操作** - 读取、写入、编辑文件，支持精确替换
- 💻 **命令执行** - 执行 Shell 命令，自动处理跨平台差异
- 🧠 **上下文记忆** - 项目感知，自动生成和维护项目记忆
- 🔒 **安全确认** - 危险操作需要确认，支持命令黑名单
- 💾 **会话持久化** - 自动保存对话历史
- 🔄 **实时执行可见性** - 逐步显示AI思考和工具调用过程

## 🚀 一键安装

### Linux / macOS

```bash
curl -fsSL https://raw.githubusercontent.com/accil/accil/main/install.sh | bash
```

或

```bash
git clone https://github.com/accil/accil.git
cd accil
chmod +x install.sh
./install.sh
```

### Windows

```powershell
# PowerShell
irm https://raw.githubusercontent.com/accil/accil/main/install.ps1 | iex
```

或下载后双击运行 `install.bat`

### 手动安装

```bash
# 克隆仓库
git clone https://github.com/accil/accil.git
cd accil

# 安装依赖
go mod tidy

# 编译
go build -o accil .

# 运行
./accil
```

## 📖 使用方法

### 交互模式

```bash
# 启动交互式会话
accil

# 指定工作目录
accil --workdir ./myproject

# 继续上次会话
accil --continue

# 自动批准模式（跳过确认）
accil --yolo
```

### 单次执行

```bash
# 执行单个任务
accil "读取 main.go 并解释其功能"

# 创建文件
accil "在桌面创建一个介绍你的txt文件"

# 无头模式（用于脚本集成）
accil --headless "重构这个函数"
```

### 内建命令

在交互模式下输入：

| 命令 | 说明 |
|------|------|
| `/help` | 显示帮助 |
| `/clear` | 清除对话 |
| `/quit` | 退出程序 |
| `/quest` | 进入任务模式 |
| `/review` | 进入审查模式 |
| `/agent` | 进入代理模式 |
| `/model <名称>` | 更改模型 |

### 快捷键

| 快捷键 | 说明 |
|--------|------|
| `Ctrl+C` | 退出 |
| `Ctrl+L` | 清屏 |
| `↑/↓` | 浏览历史 / 滚动消息 |
| `PgUp/PgDn` | 翻页 |
| `鼠标滚轮` | 滚动消息 |

## 🔧 配置

### 首次运行

首次运行会自动启动配置向导，引导设置：
- API 提供商选择（OpenAI、DeepSeek、Anthropic、Ollama 等）
- API Key 输入
- 模型选择

### 配置文件

配置文件位于 `~/.accil/config.yaml`：

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

### 环境变量

```bash
export AI_API_KEY="your-api-key"
export AI_BASE_URL="https://api.openai.com/v1"
```

## 🌐 支持的 API 提供商

| 提供商 | Base URL | 推荐模型 |
|--------|----------|----------|
| OpenAI | `https://api.openai.com/v1` | gpt-4o, gpt-4-turbo |
| DeepSeek | `https://api.deepseek.com/v1` | deepseek-chat, deepseek-coder |
| Anthropic | `https://api.anthropic.com/v1` | claude-3-opus, claude-3-sonnet |
| 通义千问 | `https://dashscope.aliyuncs.com/compatible-mode/v1` | qwen-turbo, qwen-max |
| 智谱 AI | `https://open.bigmodel.cn/api/paas/v4` | glm-4 |
| Ollama (本地) | `http://localhost:11434/v1` | llama3, codellama, mistral |

## 🛠️ 工具系统

AI 可以调用以下工具：

| 工具 | 说明 | 需要确认 |
|------|------|----------|
| `read_file` | 读取文件内容 | 否 |
| `write_file` | 写入/创建文件 | 是 |
| `edit_file` | 精确替换文件内容 | 是 |
| `run_command` | 执行 Shell 命令 | 是 |
| `list_dir` | 列出目录内容 | 否 |
| `search_code` | 正则搜索代码 | 否 |
| `glob` | 文件模式匹配 | 否 |

## 🔒 安全机制

- **默认安全**：所有写文件和执行命令操作需要用户确认
- **YOLO 模式**：使用 `--yolo` 参数跳过所有确认（警告：使用风险自负）
- **命令黑名单**：危险命令始终被阻止
- **网络重试**：API调用失败时自动重试最多3次

## 📁 项目结构

```
accil/
├── cmd/                    # 命令行入口
│   └── root.go
├── internal/
│   ├── ai/                 # AI 客户端
│   ├── config/             # 配置管理
│   ├── context/            # 上下文记忆
│   ├── memory/             # 项目记忆
│   ├── session/            # 会话管理
│   ├── tools/              # 工具系统
│   ├── tui/                # 终端 UI
│   ├── quest/              # 自主任务
│   ├── agent/              # 子代理
│   └── review/             # 代码审查
├── main.go
├── go.mod
├── install.sh              # Linux/macOS 安装脚本
├── install.bat             # Windows 安装脚本
├── Makefile
├── LICENSE
├── README.md               # English Documentation
└── README_zh.md            # Chinese Documentation
```

## 🤝 贡献

欢迎贡献！请查看 [贡献指南](CONTRIBUTING.md)。

## 📄 许可证

[MIT License](LICENSE)

---

<div align="center">

**如果这个项目对你有帮助，请给一个 ⭐ Star！**

Made with ❤️ by the ACCIL Team

</div>
