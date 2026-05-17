# ACCIL

**AI 驱动的自主编程助手**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Version](https://img.shields.io/badge/release-1.4.6-blue.svg)](CHANGELOG.md)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey)](https://cli.acz.qzz.io/)
[![Website](https://img.shields.io/badge/Website-cli.acz.qzz.io-blue)](https://cli.acz.qzz.io/)

[English](README.md) | [中文](README_zh.md) | [官网](https://cli.acz.qzz.io/)

---

## 项目简介

ACCIL 是一个面向终端的 AI 编程助手，支持交互式对话、任务执行、代码审查、远程开发和多种 API 提供商配置。

## 1.4.6 更新内容

- 将运行时配置入口统一到 `/config`
- 支持通过 `/config provider`、`/config model`、`/config baseurl` 修改配置
- 保留旧的 `/provider`、`/model`、`/baseurl` 作为兼容别名
- 文档移除面向 GitHub 的使用说明，统一指向官网和本地发布流程

## 主要功能

- 交互式终端界面
- Quest 多步骤任务模式
- 代码审查模式
- 子代理系统
- SSH 远程开发
- 文件读写与命令执行
- 会话持久化
- API 用量统计

## 安装

### Linux / macOS

请到[官网](https://cli.acz.qzz.io/)获取最新安装包或安装说明。

### Windows

请到[官网](https://cli.acz.qzz.io/)获取最新安装器或安装说明。

### 手动安装

```bash
cd accil
go mod tidy
go build -o accil .
./accil
```

Windows:

```powershell
build.bat
```

## 使用

### 交互模式

```bash
accil
```

### 单次执行

```bash
accil "Read main.go and explain what it does"
accil --headless "Refactor this function"
```

## 内建命令

| 命令 | 说明 |
|------|------|
| `/help` | 显示帮助 |
| `/clear` | 清空对话 |
| `/quit` | 退出程序 |
| `/chat` | 进入对话模式 |
| `/quest` | 进入任务模式 |
| `/review` | 进入审查模式 |
| `/agent` | 进入代理模式 |
| `/remote` | 进入远程开发模式 |
| `/config` | 查看当前配置与可用配置命令 |
| `/context` | 显示当前上下文 |

## 配置

首次运行可以使用：

```bash
accil --setup
```

配置文件位置：

```text
~/.accil/config.yaml
```

环境变量：

```bash
export AI_API_KEY="your-api-key"
export AI_BASE_URL="https://api.openai.com/v1"
```

### 首次配置后修改服务商或模型

现在统一使用 `/config`：

```text
/config show
/config provider deepseek
/config model deepseek-v4-pro
/config baseurl https://api.deepseek.com/v1
```

旧命令 `/provider`、`/model`、`/baseurl` 仍可继续使用，但主入口已经改为 `/config`。

## 支持的 API 提供商

| 提供商 | Base URL | 推荐模型 |
|--------|----------|----------|
| OpenAI | `https://api.openai.com/v1` | gpt-5.5, gpt-5.4-mini, gpt-5.4-nano |
| DeepSeek | `https://api.deepseek.com/v1` | deepseek-v4-pro, deepseek-v4-flash |
| Anthropic | `https://api.anthropic.com/v1` | claude-sonnet-4-6, claude-opus-4-1 |
| Qwen | `https://dashscope.aliyuncs.com/compatible-mode/v1` | qwen3-max, qwen-plus, qwen-turbo |
| Zhipu AI | `https://open.bigmodel.cn/api/paas/v4` | glm-4.6, glm-4-plus |
| Ollama | `http://localhost:11434/v1` | qwen2.5-coder, llama3.3, mistral |

## 更多

- [快速开始](QUICKSTART.md)
- [变更日志](CHANGELOG.md)
- [贡献指南](CONTRIBUTING.md)

## 许可证

[MIT License](LICENSE)
