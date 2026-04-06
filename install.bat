@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

:: ACCIL Installation Script for Windows

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

:: Check if git is installed
where git >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo [ERROR] Git is not installed!
    echo Please download and install Git from: https://git-scm.com/download/win
    echo After installing Git, run this script again.
    pause
    exit /b 1
)

echo [OK] Git is installed
echo.

:: Clone repository to temp directory
set TEMP_DIR=%TEMP%\accil-install-%RANDOM%
echo [INFO] Downloading ACCIL from GitHub...
echo [INFO] Repository: %REPO_URL%
echo.

git clone --depth 1 %REPO_URL% "%TEMP_DIR%"
if %ERRORLEVEL% neq 0 (
    echo.
    echo [ERROR] Failed to download from GitHub
    echo.
    echo Possible solutions:
    echo   1. Check your internet connection
    echo   2. If using proxy, configure git:
    echo      git config --global http.proxy http://127.0.0.1:7890
    echo   3. Or download manually from:
    echo      https://github.com/acacMAX/accil/archive/refs/heads/main.zip
    pause
    exit /b 1
)
echo [OK] Download completed
echo.

:: Build
echo [INFO] Building ACCIL...
if not exist "%TEMP_DIR%" (
    echo [ERROR] Download directory not found
    pause
    exit /b 1
)

cd /d "%TEMP_DIR%"
if %ERRORLEVEL% neq 0 (
    echo [ERROR] Cannot access download directory
    pause
    exit /b 1
)

echo [INFO] Installing dependencies...
go mod tidy
if %ERRORLEVEL% neq 0 (
    echo [ERROR] Failed to install dependencies
    echo.
    echo Try setting Go proxy:
    echo   go env -w GOPROXY=https://goproxy.cn,direct
    cd /d "%USERPROFILE%"
    pause
    exit /b 1
)

if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"

echo [INFO] Compiling...
go build -o "%INSTALL_DIR%\accil.exe" .
if %ERRORLEVEL% neq 0 (
    echo [ERROR] Build failed
    cd /d "%USERPROFILE%"
    pause
    exit /b 1
)
echo [OK] Build completed

cd /d "%USERPROFILE%"

:: Cleanup
echo [INFO] Cleaning up...
if exist "%TEMP_DIR%" rmdir /s /q "%TEMP_DIR%"
echo [OK] Cleanup completed
echo.

:: Add to PATH
echo [INFO] Configuring PATH...

:: Check if already in PATH
echo %PATH% | findstr /C:"%INSTALL_DIR%" >nul
if %ERRORLEVEL% equ 0 (
    echo [OK] Already in PATH
) else (
    :: Add to current session immediately
    set "PATH=%PATH%;%INSTALL_DIR%"

    :: Get current user PATH from registry
    for /f "skip=2 tokens=2,*" %%a in ('reg query "HKCU\Environment" /v PATH 2^>nul') do set CURRENT_USER_PATH=%%b
    if "!CURRENT_USER_PATH!"=="" (
        setx PATH "%INSTALL_DIR%" >nul
    ) else (
        setx PATH "!CURRENT_USER_PATH!;%INSTALL_DIR%" >nul
    )
    if !ERRORLEVEL! equ 0 (
        echo [OK] Added to user PATH
        echo     You can now use 'accil' command globally
    ) else (
        echo [WARNING] Could not add to PATH automatically
        echo Please add manually: %INSTALL_DIR%
    )
)
echo.

:: Verify installation
echo [INFO] Verifying installation...
if exist "%INSTALL_DIR%\accil.exe" (
    echo [OK] Installation verified
    echo     Location: %INSTALL_DIR%\accil.exe
) else (
    echo [ERROR] Installation verification failed
    pause
    exit /b 1
)
echo.

:: Success
echo ========================================
echo    Installation Complete!
echo ========================================
echo.
echo Install location: %INSTALL_DIR%\accil.exe
echo.
echo USAGE:
echo   - You can use 'accil' command NOW in this window
echo   - For NEW command prompts, simply type: accil
echo.
echo If 'accil' is not recognized in new windows:
echo   1. Open System Properties ^> Environment Variables
echo   2. Add to PATH: %INSTALL_DIR%
echo   Or run: setx PATH "%%PATH%%;%INSTALL_DIR%"
echo.

set /p run_now="Run ACCIL now? (y/n): "
if /i "%run_now%"=="y" (
    "%INSTALL_DIR%\accil.exe"
)

pause
