@echo off
chcp 65001 >nul
setlocal enabledelayexpansion

:: ACCIL 一键安装脚本 (Windows)

set REPO_URL=https://github.com/acacMAX/accil.git
set INSTALL_DIR=%USERPROFILE%\.accil\bin
set BINARY_NAME=accil.exe

:: 颜色设置（需要Windows 10+）
for /F %%a in ('echo prompt $E^| cmd') do set "ESC=%%a"
set "RED=!ESC![91m"
set "GREEN=!ESC![92m"
set "YELLOW=!ESC![93m"
set "BLUE=!ESC![94m"
set "NC=!ESC![0m"

echo.
echo !BLUE!╔══════════════════════════════════════════════════════════════╗!NC!
echo !BLUE!║                                                              ║!NC!
echo !BLUE!║   █████╗ ██████╗ ██████╗  ██████╗██╗  ██╗██╗     ███████╗   ║!NC!
echo !BLUE!║  ██╔══██╗██╔══██╗██╔══██╗██╔════╝██║  ██║██║     ██╔════╝   ║!NC!
echo !BLUE!║  ███████║██████╔╝██████╔╝██║     ███████║██║     █████╗     ║!NC!
echo !BLUE!║  ██╔══██║██╔══██╗██╔══██╗██║     ██╔══██║██║     ██╔══╝     ║!NC!
echo !BLUE!║  ██║  ██║██████╔╝██████╔╝╚██████╗██║  ██║███████╗███████╗   ║!NC!
echo !BLUE!║  ╚═╝  ╚═╝╚═════╝ ╚═════╝  ╚═════╝╚═╝  ╚═╝╚══════╝╚══════╝   ║!NC!
echo !BLUE!║                                                              ║!NC!
echo !BLUE!║           AI驱动的自主编程助手 - 安装程序                    ║!NC!
echo !BLUE!║                                                              ║!NC!
echo !BLUE!╚══════════════════════════════════════════════════════════════╝!NC!
echo.

:: 检查 Go
where go >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo !RED!错误: 未安装 Go!NC!
    echo !YELLOW!请从 https://go.dev/dl/ 下载并安装 Go!NC!
    pause
    exit /b 1
)

for /f "tokens=3" %%v in ('go version') do set GO_VERSION=%%v
echo !GREEN!✓ 检测到 Go %GO_VERSION%!NC!

:: 设置 Go 代理
set GOPROXY=https://goproxy.cn,direct

:: 编译
echo !BLUE!正在编译...!NC!
go mod tidy
go build -ldflags="-s -w" -o %BINARY_NAME% .
if %ERRORLEVEL% neq 0 (
    echo !RED!编译失败!NC!
    pause
    exit /b 1
)

:: 创建安装目录
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"

:: 移动文件
move /Y %BINARY_NAME% "%INSTALL_DIR%\" >nul

:: 添加到 PATH（需要管理员权限或手动添加）
echo.
echo !GREEN!✓ 编译完成!NC!
echo.
echo 安装位置: %INSTALL_DIR%\accil.exe
echo.

:: 检查配置
if not exist "%USERPROFILE%\.accil\config.yaml" (
    echo.
    echo !BLUE!首次使用配置!NC!
    echo.
    echo 请选择 API 提供商:
    echo   1) OpenAI
    echo   2) DeepSeek
    echo   3) 本地 Ollama
    echo   4) 自定义
    echo.
    set /p provider_choice="请选择 [1-4]: "
    
    if "!provider_choice!"=="1" (
        set BASE_URL=https://api.openai.com/v1
        set DEFAULT_MODEL=gpt-4o
    )
    if "!provider_choice!"=="2" (
        set BASE_URL=https://api.deepseek.com/v1
        set DEFAULT_MODEL=deepseek-chat
    )
    if "!provider_choice!"=="3" (
        set BASE_URL=http://localhost:11434/v1
        set DEFAULT_MODEL=llama3
    )
    if "!provider_choice!"=="4" (
        set /p BASE_URL="输入 API URL: "
        set /p DEFAULT_MODEL="输入模型名称: "
    )
    
    set /p API_KEY="输入 API Key: "
    
    :: 创建配置目录
    if not exist "%USERPROFILE%\.accil" mkdir "%USERPROFILE%\.accil"
    
    :: 写入配置文件
    (
        echo api_key: "!API_KEY!"
        echo base_url: "!BASE_URL!"
        echo model: "!DEFAULT_MODEL!"
        echo max_tokens: 4096
        echo auto_approve: false
    ) > "%USERPROFILE%\.accil\config.yaml"
    
    echo !GREEN!✓ 配置已保存!NC!
)

:: 添加到用户 PATH
echo.
echo !YELLOW!正在添加到 PATH...!NC!
setx PATH "%PATH%;%INSTALL_DIR%" >nul 2>&1

echo.
echo !GREEN!╔══════════════════════════════════════════════════════════════╗!NC!
echo !GREEN!║                    安装成功!                                  ║!NC!
echo !GREEN!╚══════════════════════════════════════════════════════════════╝!NC!
echo.
echo 使用方法:
echo   accil              # 启动交互模式
echo   accil "你好"       # 单次执行
echo   accil --help       # 查看帮助
echo.
echo !YELLOW!注意: 请重新打开命令提示符窗口使 PATH 生效!NC!
echo.

:: 快速运行选项
set /p run_now="是否立即运行? (y/n): "
if /i "!run_now!"=="y" (
    "%INSTALL_DIR%\accil.exe"
)

pause
