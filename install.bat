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
set "USER_PATH=%INSTALL_DIR%"

:: Check if already in PATH
echo %PATH% | findstr /C:"%INSTALL_DIR%" >nul
if %ERRORLEVEL% equ 0 (
    echo [OK] Already in PATH
) else (
    :: Add using setx
    for /f "2 tokens=3," %%a in ('reg query "HKCU\Environment" /v PATH') do set CURRENT_USER_PATH=%%b
    setx PATH "!CURRENT_USER_PATH!;%INSTALL_DIR%" >nul 2>&1
    if !ERRORLEVEL! equ 0 (
        echo [OK] Added to user PATH
    ) else (
        echo [WARNING] Failed to add to PATH automatically
        echo Please add manually: %INSTALL_DIR%
    )
)
echo.

:: Verify installation
echo [INFO] Verifying...
if exist "%INSTALL_DIR%\accil.exe" (
    echo [OK] accil.exe found at: %INSTALL_DIR%
) else (
    echo [ERROR] accil.exe not found!
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
echo   3. Type: accil --help
echo.
echo If 'accil' is still not recognized, run this command:
echo   setx PATH "%%PATH%%;%INSTALL_DIR%"
echo.
echo Then restart your command prompt again.
echo.

set /p run_now="Run ACCIL now? (y/n): "
if /i "%run_now%"=="y" (
    "%INSTALL_DIR%\accil.exe"
)

pause
