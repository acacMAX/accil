#!/bin/bash

# ACCIL 一键安装脚本
# 支持: Linux, macOS, Windows (Git Bash)

set -e

REPO_URL="https://github.com/accil/accil.git"
INSTALL_DIR="$HOME/.accil/bin"
BINARY_NAME="accil"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_logo() {
    echo -e "${BLUE}"
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║                                                              ║"
    echo "║   █████╗ ██████╗ ██████╗  ██████╗██╗  ██╗██╗     ███████╗   ║"
    echo "║  ██╔══██╗██╔══██╗██╔══██╗██╔════╝██║  ██║██║     ██╔════╝   ║"
    echo "║  ███████║██████╔╝██████╔╝██║     ███████║██║     █████╗     ║"
    echo "║  ██╔══██║██╔══██╗██╔══██╗██║     ██╔══██║██║     ██╔══╝     ║"
    echo "║  ██║  ██║██████╔╝██████╔╝╚██████╗██║  ██║███████╗███████╗   ║"
    echo "║  ╚═╝  ╚═╝╚═════╝ ╚═════╝  ╚═════╝╚═╝  ╚═╝╚══════╝╚══════╝   ║"
    echo "║                                                              ║"
    echo "║           AI驱动的自主编程助手 - 安装程序                    ║"
    echo "║                                                              ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

check_go() {
    if ! command -v go &> /dev/null; then
        echo -e "${RED}错误: 未安装 Go${NC}"
        echo -e "${YELLOW}请先安装 Go 1.21 或更高版本:${NC}"
        echo "  - macOS:   brew install go"
        echo "  - Linux:   sudo apt install golang-go 或 sudo dnf install golang"
        echo "  - Windows: 从 https://go.dev/dl/ 下载安装"
        exit 1
    fi

    GO_VERSION=$(go version | grep -oP 'go\K[0-9.]+' | head -1)
    echo -e "${GREEN}✓ 检测到 Go $GO_VERSION${NC}"
}

install_from_source() {
    echo -e "${BLUE}正在从源码安装...${NC}"
    
    # 创建临时目录
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"
    
    # 克隆仓库
    echo -e "${YELLOW}→ 克隆仓库...${NC}"
    git clone "$REPO_URL" 2>/dev/null || {
        # 如果仓库不存在，使用当前目录
        echo -e "${YELLOW}→ 使用本地源码...${NC}"
        SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
        cd "$SCRIPT_DIR"
    }
    
    cd accil 2>/dev/null || true
    
    # 设置 Go 代理
    export GOPROXY=https://goproxy.cn,direct
    
    # 编译
    echo -e "${YELLOW}→ 编译中...${NC}"
    go mod tidy
    go build -ldflags="-s -w" -o "$BINARY_NAME" .
    
    # 创建安装目录
    mkdir -p "$INSTALL_DIR"
    
    # 移动二进制文件
    mv "$BINARY_NAME" "$INSTALL_DIR/"
    
    # 清理
    cd -
    rm -rf "$TEMP_DIR"
    
    echo -e "${GREEN}✓ 编译完成${NC}"
}

add_to_path() {
    echo -e "${BLUE}配置 PATH...${NC}"
    
    # 检查是否已在 PATH 中
    if [[ ":$PATH:" == *":$INSTALL_DIR:"* ]]; then
        echo -e "${GREEN}✓ PATH 已配置${NC}"
        return
    fi
    
    # 添加到 shell 配置文件
    SHELL_RC=""
    if [ -n "$ZSH_VERSION" ]; then
        SHELL_RC="$HOME/.zshrc"
    elif [ -n "$BASH_VERSION" ]; then
        SHELL_RC="$HOME/.bashrc"
    fi
    
    if [ -n "$SHELL_RC" ]; then
        echo "" >> "$SHELL_RC"
        echo "# ACCIL" >> "$SHELL_RC"
        echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$SHELL_RC"
        echo -e "${GREEN}✓ 已添加到 $SHELL_RC${NC}"
    fi
    
    # 当前会话生效
    export PATH="$PATH:$INSTALL_DIR"
}

run_setup() {
    echo -e "${BLUE}"
    echo "╭───────────────────────────────────────╮"
    echo "│         首次使用配置                   │"
    echo "╰───────────────────────────────────────╯"
    echo -e "${NC}"
    
    echo "请选择 API 提供商:"
    echo "  1) OpenAI"
    echo "  2) DeepSeek"
    echo "  3) Anthropic"
    echo "  4) 本地 Ollama"
    echo "  5) 自定义"
    echo ""
    read -p "请选择 [1-5]: " provider_choice
    
    case $provider_choice in
        1)
            BASE_URL="https://api.openai.com/v1"
            DEFAULT_MODEL="gpt-4o"
            ;;
        2)
            BASE_URL="https://api.deepseek.com/v1"
            DEFAULT_MODEL="deepseek-chat"
            ;;
        3)
            BASE_URL="https://api.anthropic.com/v1"
            DEFAULT_MODEL="claude-3-opus"
            ;;
        4)
            BASE_URL="http://localhost:11434/v1"
            DEFAULT_MODEL="llama3"
            ;;
        5)
            read -p "输入 API URL: " BASE_URL
            read -p "输入模型名称: " DEFAULT_MODEL
            ;;
        *)
            BASE_URL="https://api.openai.com/v1"
            DEFAULT_MODEL="gpt-4o"
            ;;
    esac
    
    read -p "输入 API Key: " API_KEY
    
    # 创建配置目录
    CONFIG_DIR="$HOME/.accil"
    mkdir -p "$CONFIG_DIR"
    
    # 写入配置文件
    cat > "$CONFIG_DIR/config.yaml" << EOF
api_key: "$API_KEY"
base_url: "$BASE_URL"
model: "$DEFAULT_MODEL"
max_tokens: 4096
auto_approve: false
EOF
    
    echo -e "${GREEN}✓ 配置已保存到 $CONFIG_DIR/config.yaml${NC}"
}

print_success() {
    echo ""
    echo -e "${GREEN}╔══════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║                    安装成功!                                  ║${NC}"
    echo -e "${GREEN}╚══════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "安装位置: ${YELLOW}$INSTALL_DIR/accil${NC}"
    echo ""
    echo "使用方法:"
    echo "  accil              # 启动交互模式"
    echo "  accil '你好'       # 单次执行"
    echo "  accil --help       # 查看帮助"
    echo "  accil --setup      # 重新配置"
    echo ""
    echo -e "${YELLOW}注意: 请重新打开终端或运行 'source ~/.bashrc' 使 PATH 生效${NC}"
    echo ""
}

# Windows 特定检测
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
    BINARY_NAME="accil.exe"
    INSTALL_DIR="$HOME/.accil/bin"
fi

# 主流程
print_logo
check_go
install_from_source
add_to_path

# 检查是否需要配置
if [ ! -f "$HOME/.accil/config.yaml" ]; then
    run_setup
fi

print_success
