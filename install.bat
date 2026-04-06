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

for /f "tokens=3" %%v in ('go version 2^>^&1') do set GO_VERSION=%%v
echo [OK] Go version: %GO_VERSION%

:: Check if git is installed
where git >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo [ERROR] Git is not installed!
    echo Please download and install Git from: https://git-scm.com/download/win
    pause
    exit /b 1
)

for /f "tokens=3" %%v in ('git --version 2^>^&1') do set GIT_VERSION=%%v
echo [OK] Git version: %GIT_VERSION%
echo.

:: Clone repository to temp directory
set TEMP_DIR=%TEMP%\accil-install-%RANDOM%
echo [INFO] Downloading ACCIL from GitHub...

git clone --depth 1 %REPO_URL% "%TEMP_DIR%"
if %ERRORLEVEL% neq 0 (
    echo [ERROR] Failed to download from GitHub
    echo Please check:
    echo   1. Your internet connection
    echo   2. If you need a proxy, set it first:
    echo      git config --global http.proxy http://127.0.0.1:7890
    pause
    exit /b 1
)
echo [OK] Download completed
echo.

:: Build
echo [INFO] Building ACCIL...
pushd "%TEMP_DIR%"
if %ERRORLEVEL% neq 0 (
    echo [ERROR] Failed to enter build directory
    pause
    exit /b 1
)

echo [INFO] Installing dependencies...
go mod tidy
if %ERRORLEVEL% neq 0 (
    echo [ERROR] Failed to install dependencies
    echo This might be a network issue. Try setting Go proxy:
    echo   go env -w GOPROXY=https://goproxy.cn,direct
    popd
    pause
    exit /b 1
)

if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"

echo [INFO] Compiling...
go build -o "%INSTALL_DIR%\accil.exe" .
if %ERRORLEVEL% neq 0 (
    echo [ERROR] Build failed
    popd
    pause
    exit /b 1
)
echo [OK] Build completed
popd

:: Cleanup
echo [INFO] Cleaning up...
rmdir /s /q "%TEMP_DIR%" 2>nul
echo [OK] Cleanup completed
echo.

:: Add to PATH
echo [INFO] Adding to PATH...

:: Check if already in PATH
echo %PATH% | findstr /C:"%INSTALL_DIR%" >nul
if %ERRORLEVEL% equ 0 (
    echo [OK] Already in PATH
) else (
    :: Get current user PATH
    for /f "skip=2 tokens=2,*" %%a in ('reg query "HKCU\Environment" /v PATH 2^>nul') do set CURRENT_USER_PATH=%%b
    if "!CURRENT_USER_PATH!"=="" (
        setx PATH "%INSTALL_DIR%" >nul 2>&1
    ) else (
        setx PATH "!CURRENT_USER_PATH!;%INSTALL_DIR%" >nul 2>&1
    )
    if !ERRORLEVEL! equ 0 (
        echo [OK] Added to user PATH
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
echo IMPORTANT: Please follow these steps:
echo   1. Close ALL command prompt windows
echo   2. Open a NEW command prompt
echo   3. Type: accil
echo.
echo If 'accil' is not recognized, run:
echo   setx PATH "%%PATH%%;%INSTALL_DIR%"
echo Then restart your command prompt.
echo.

set /p run_now="Run ACCIL now? (y/n): "
if /i "%run_now%"=="y" (
    "%INSTALL_DIR%\accil.exe"
)

pause
