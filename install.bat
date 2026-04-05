@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

:: ACCIL Installation Script for Windows (Batch version)

set REPO_URL=https://github.com/acacMAX/accil.git
set INSTALL_DIR=%USERPROFILE%\.accil\bin

echo.
echo ========================================
echo    ACCIL Installation Wizard
echo ========================================
echo.

:: Check if Go is installed
where go >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo [ERROR] Go is not installed!
    echo Please download and install Go from: https://go.dev/dl/
    pause
    exit /b 1
)

for /f "tokens=3" %%v in ('go version') do set GO_VERSION=%%v
echo [OK] Go version: %GO_VERSION%
echo.

:: Clone repository to temp directory
set TEMP_DIR=%TEMP%\accil-install-%RANDOM%
echo [INFO] Downloading ACCIL...
git clone %REPO_URL% "%TEMP_DIR%" >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo [ERROR] Failed to download from GitHub
    echo Please check your internet connection
    pause
    exit /b 1
)
echo [OK] Download completed
echo.

:: Build
echo [INFO] Building ACCIL...
cd /d "%TEMP_DIR%"
go mod tidy >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo [ERROR] Failed to install dependencies
    pause
    exit /b 1
)

if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"
go build -o "%INSTALL_DIR%\accil.exe" . >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo [ERROR] Build failed
    pause
    exit /b 1
)
echo [OK] Build completed
echo.

:: Cleanup
echo [INFO] Cleaning up...
cd /d "%INSTALL_DIR%"
rmdir /s /q "%TEMP_DIR%" 2>nul
echo [OK] Cleanup completed
echo.

:: Add to PATH
echo [INFO] Adding to PATH...
setx PATH "%PATH%;%INSTALL_DIR%" >nul 2>&1
echo [OK] Added to PATH
echo.

:: Success
echo ========================================
echo    Installation Complete!
echo ========================================
echo.
echo Install location: %INSTALL_DIR%\accil.exe
echo.
echo Usage:
echo   accil              Start interactive mode
echo   accil "hello"      Single execution
echo   accil --help       Show help
echo.
echo NOTE: Please restart your command prompt to use 'accil' command
echo.

set /p run_now="Run ACCIL now? (y/n): "
if /i "%run_now%"=="y" (
    "%INSTALL_DIR%\accil.exe"
)

pause
