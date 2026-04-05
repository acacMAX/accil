@echo off
chcp 65001 >nul
echo.
echo ========================================
echo    ACCIL Quick Installation
echo ========================================
echo.

:: Check if we're in the right directory
if not exist "main.go" (
    echo [ERROR] Please run this script from the ACCIL project directory
    echo Current directory: %CD%
    pause
    exit /b 1
)

:: Build
echo [INFO] Building accil.exe...
go build -o accil.exe .
if %ERRORLEVEL% neq 0 (
    echo [ERROR] Build failed!
    pause
    exit /b 1
)
echo [OK] Build successful
echo.

:: Install directory
set INSTALL_DIR=%USERPROFILE%\.accil\bin
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"

:: Copy
echo [INFO] Installing to %INSTALL_DIR%...
copy /Y accil.exe "%INSTALL_DIR%\accil.exe" >nul
echo [OK] Installed
echo.

:: Add to PATH using setx
echo [INFO] Adding to PATH...
for /f "tokens=2,*" %%a in ('reg query "HKCU\Environment" /v PATH') do set USERPATH=%%b
setx PATH "%USERPATH%;%INSTALL_DIR%" >nul
echo [OK] PATH updated
echo.

echo ========================================
echo    Installation Complete!
echo ========================================
echo.
echo IMPORTANT NEXT STEPS:
echo.
echo 1. Close ALL command prompt and PowerShell windows
echo 2. Open a NEW command prompt
echo 3. Type: accil --help
echo.
echo If it still doesn't work, run this in a NEW cmd:
echo   setx PATH "%%PATH%%;%INSTALL_DIR%"
echo.
pause
