@echo off
setlocal EnableExtensions EnableDelayedExpansion

set "SCRIPT_DIR=%~dp0"
pushd "%SCRIPT_DIR%" >nul || (
    echo [ERROR] Failed to enter script directory.
    exit /b 1
)

where go >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Go is not installed or not in PATH.
    popd >nul
    exit /b 1
)

set "ACCIL_VER=1.4.6"
if exist VERSION (
    for /f "usebackq delims=" %%a in ("VERSION") do set "ACCIL_VER=%%a"
)
set "ACCIL_VER=%ACCIL_VER: =%"

set "INSTALL_DIR=%USERPROFILE%\.accil\bin"
set "FALLBACK_INSTALL_DIR=%LOCALAPPDATA%\accil\bin"
set "LAST_RESORT_INSTALL_DIR=%TEMP%\accil\bin"
set "CACHE_DIR=%TEMP%\accil-cache"
set "GOCACHE=%CACHE_DIR%\go-build"
set "GOMODCACHE=%CACHE_DIR%\gomod"
set "OUTPUT_EXE=%SCRIPT_DIR%accil.exe"

if not exist "%CACHE_DIR%" mkdir "%CACHE_DIR%"
if not exist "%GOCACHE%" mkdir "%GOCACHE%"
if not exist "%GOMODCACHE%" mkdir "%GOMODCACHE%"

echo [INFO] Building ACCIL v%ACCIL_VER%...
go build -buildvcs=false -ldflags="-X github.com/accil/accil/cmd.Version=%ACCIL_VER%" -o "%OUTPUT_EXE%" .
if errorlevel 1 (
    echo [ERROR] Build failed.
    popd >nul
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
    exit /b 1
)

copy /y "%OUTPUT_EXE%" "!INSTALL_DIR!\accil.exe" >nul 2>nul
if errorlevel 1 (
    if /I not "!INSTALL_DIR!"=="%FALLBACK_INSTALL_DIR%" (
        set "INSTALL_DIR=%FALLBACK_INSTALL_DIR%"
        if not exist "!INSTALL_DIR!" mkdir "!INSTALL_DIR!" 2>nul
        copy /y "%OUTPUT_EXE%" "!INSTALL_DIR!\accil.exe" >nul 2>nul
    )
)
if errorlevel 1 (
    if /I not "!INSTALL_DIR!"=="%LAST_RESORT_INSTALL_DIR%" (
        set "INSTALL_DIR=%LAST_RESORT_INSTALL_DIR%"
        if not exist "!INSTALL_DIR!" mkdir "!INSTALL_DIR!" 2>nul
        copy /y "%OUTPUT_EXE%" "!INSTALL_DIR!\accil.exe" >nul 2>nul
    )
)
if errorlevel 1 (
    echo [ERROR] Failed to copy accil.exe to an install directory.
    popd >nul
    exit /b 1
)

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
        echo [WARNING] Installed, but failed to update user PATH automatically.
        echo [WARNING] Add this folder manually: !INSTALL_DIR!
    ) else (
        echo [OK] Added !INSTALL_DIR! to user PATH.
    )
) else (
    echo [OK] Install directory already present in PATH.
)

"!INSTALL_DIR!\accil.exe" version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Installed binary verification failed.
    popd >nul
    exit /b 1
)

echo [OK] Build and global install complete.
echo [OK] Binary: !INSTALL_DIR!\accil.exe
if /I "!INSTALL_DIR!"=="%LAST_RESORT_INSTALL_DIR%" (
    echo [WARNING] Primary install directories were not writable.
    echo [WARNING] Installed to a temp-backed fallback directory instead.
)
echo [INFO] Open a new terminal if the ^`accil^` command is not available yet.

popd >nul
exit /b 0
