@echo off
chcp 65001 >nul
setlocal EnableExtensions EnableDelayedExpansion

set "SCRIPT_DIR=%~dp0"
set "INSTALL_DIR=%USERPROFILE%\.accil\bin"
set "FALLBACK_INSTALL_DIR=%LOCALAPPDATA%\accil\bin"
set "LAST_RESORT_INSTALL_DIR=%TEMP%\accil\bin"
set "CACHE_DIR=%TEMP%\accil-cache"
set "GOCACHE=%CACHE_DIR%\go-build"
set "GOMODCACHE=%CACHE_DIR%\gomod"
set "REPO_URL=https://github.com/acacMAX/accil.git"
set "TEMP_DIR=%TEMP%\accil-install-%RANDOM%%RANDOM%"
set "SOURCE_DIR="
set "BUILD_OUTPUT=%TEMP%\accil-build-%RANDOM%%RANDOM%.exe"

echo.
echo ========================================
echo    ACCIL Installation Wizard
echo ========================================
echo.

where go >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Go is not installed or not in PATH.
    echo Please install Go from: https://go.dev/dl/
    pause
    exit /b 1
)
for /f "tokens=3" %%v in ('go version') do set "GO_VERSION=%%v"
echo [OK] Go version: %GO_VERSION%

if exist "%SCRIPT_DIR%go.mod" (
    set "SOURCE_DIR=%SCRIPT_DIR%"
    echo [INFO] Using local source tree: !SOURCE_DIR!
) else (
    where git >nul 2>&1
    if errorlevel 1 (
        echo [ERROR] Git is required when installing without a local source tree.
        pause
        exit /b 1
    )

    echo [INFO] Downloading ACCIL source package...
    echo [INFO] Repository: %REPO_URL%
    git clone --depth 1 "%REPO_URL%" "%TEMP_DIR%"
    if errorlevel 1 (
        echo [ERROR] Failed to download source package.
        pause
        exit /b 1
    )
    set "SOURCE_DIR=%TEMP_DIR%"
    echo [OK] Download completed
)

pushd "%SOURCE_DIR%" >nul || (
    echo [ERROR] Cannot access source directory.
    if exist "%TEMP_DIR%" rmdir /s /q "%TEMP_DIR%"
    pause
    exit /b 1
)

if not exist "!INSTALL_DIR!" mkdir "!INSTALL_DIR!" 2>nul
if not exist "!INSTALL_DIR!" (
    set "INSTALL_DIR=%FALLBACK_INSTALL_DIR%"
    if not exist "!INSTALL_DIR!" mkdir "!INSTALL_DIR!" 2>nul
)
if not exist "!INSTALL_DIR!" (
    set "INSTALL_DIR=%LAST_RESORT_INSTALL_DIR%"
    if not exist "!INSTALL_DIR!" mkdir "!INSTALL_DIR!" 2>nul
)
if not exist "!INSTALL_DIR!" (
    echo [ERROR] Failed to create install directory.
    echo [ERROR] Tried: %USERPROFILE%\.accil\bin
    echo [ERROR] Tried: %FALLBACK_INSTALL_DIR%
    echo [ERROR] Tried: %LAST_RESORT_INSTALL_DIR%
    popd >nul
    if exist "%TEMP_DIR%" rmdir /s /q "%TEMP_DIR%"
    pause
    exit /b 1
)

if not exist "%CACHE_DIR%" mkdir "%CACHE_DIR%"
if not exist "%GOCACHE%" mkdir "%GOCACHE%"
if not exist "%GOMODCACHE%" mkdir "%GOMODCACHE%"

set "ACCIL_VER=1.4.6"
if exist VERSION (
    for /f "usebackq delims=" %%a in ("VERSION") do set "ACCIL_VER=%%a"
)
set "ACCIL_VER=%ACCIL_VER: =%"

echo [INFO] Building ACCIL v%ACCIL_VER%...
go build -buildvcs=false -ldflags="-X github.com/accil/accil/cmd.Version=%ACCIL_VER%" -o "%BUILD_OUTPUT%" .
if errorlevel 1 (
    echo [ERROR] Build failed.
    popd >nul
    if exist "%TEMP_DIR%" rmdir /s /q "%TEMP_DIR%"
    pause
    exit /b 1
)
copy /y "%BUILD_OUTPUT%" "!INSTALL_DIR!\accil.exe" >nul 2>nul
if errorlevel 1 (
    if /I not "!INSTALL_DIR!"=="%FALLBACK_INSTALL_DIR%" (
        set "INSTALL_DIR=%FALLBACK_INSTALL_DIR%"
        if not exist "!INSTALL_DIR!" mkdir "!INSTALL_DIR!" 2>nul
        copy /y "%BUILD_OUTPUT%" "!INSTALL_DIR!\accil.exe" >nul 2>nul
    )
)
if errorlevel 1 (
    if /I not "!INSTALL_DIR!"=="%LAST_RESORT_INSTALL_DIR%" (
        set "INSTALL_DIR=%LAST_RESORT_INSTALL_DIR%"
        if not exist "!INSTALL_DIR!" mkdir "!INSTALL_DIR!" 2>nul
        copy /y "%BUILD_OUTPUT%" "!INSTALL_DIR!\accil.exe" >nul 2>nul
    )
)
if errorlevel 1 (
    echo [ERROR] Failed to copy accil.exe to an install directory.
    if exist "%BUILD_OUTPUT%" del /f /q "%BUILD_OUTPUT%" >nul 2>nul
    popd >nul
    if exist "%TEMP_DIR%" rmdir /s /q "%TEMP_DIR%"
    pause
    exit /b 1
)
if exist "%BUILD_OUTPUT%" del /f /q "%BUILD_OUTPUT%" >nul 2>nul
echo [OK] Build completed

echo %PATH% | findstr /I /C:"!INSTALL_DIR!" >nul
if errorlevel 1 (
    set "PATH=%PATH%;!INSTALL_DIR!"
    for /f "skip=2 tokens=2,*" %%a in ('reg query "HKCU\Environment" /v PATH 2^>nul') do set "CURRENT_USER_PATH=%%b"
    if defined CURRENT_USER_PATH (
        setx PATH "!CURRENT_USER_PATH!;!INSTALL_DIR!" >nul
    ) else (
        setx PATH "!INSTALL_DIR!" >nul
    )
    if errorlevel 1 (
        echo [WARNING] Could not add to PATH automatically.
        echo [WARNING] Add this folder manually: !INSTALL_DIR!
    ) else (
        echo [OK] Added to user PATH
    )
) else (
    echo [OK] Install directory already present in PATH
)

"!INSTALL_DIR!\accil.exe" version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Installation verification failed.
    popd >nul
    if exist "%TEMP_DIR%" rmdir /s /q "%TEMP_DIR%"
    pause
    exit /b 1
)

popd >nul
if exist "%TEMP_DIR%" rmdir /s /q "%TEMP_DIR%"

echo.
echo ========================================
echo    Installation Complete!
echo ========================================
echo.
echo Installed to: !INSTALL_DIR!\accil.exe
if /I "!INSTALL_DIR!"=="%LAST_RESORT_INSTALL_DIR%" (
    echo [WARNING] Primary install directories were not writable.
    echo [WARNING] Installed to a temp-backed fallback directory instead.
)
echo You can run: accil
echo If the command is not available in a new terminal yet, reopen the terminal.
echo.

set /p run_now="Run ACCIL now? (y/n): "
if /i "%run_now%"=="y" (
    "%INSTALL_DIR%\accil.exe"
)

pause
exit /b 0
