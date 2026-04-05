# ACCIL PowerShell Installation Script
# This script downloads and installs ACCIL on Windows

param(
    [string]$Version = "latest",
    [switch]$Help
)

if ($Help) {
    Write-Host "ACCIL Installer for Windows"
    Write-Host ""
    Write-Host "Usage:"
    Write-Host "  .\install.ps1 [-Version <version>] [-Help]"
    Write-Host ""
    Write-Host "Options:"
    Write-Host "  -Version    Version to install (default: latest)"
    Write-Host "  -Help       Show this help message"
    exit 0
}

$ErrorActionPreference = "Stop"

# Colors
function Write-Success {
    param([string]$Message)
    Write-Host "[SUCCESS] " -ForegroundColor Green -NoNewline
    Write-Host $Message
}

function Write-Error-Custom {
    param([string]$Message)
    Write-Host "[ERROR] " -ForegroundColor Red -NoNewline
    Write-Host $Message
}

function Write-Info {
    param([string]$Message)
    Write-Host "[INFO] " -ForegroundColor Cyan -NoNewline
    Write-Host $Message
}

Write-Host ""
Write-Host "╔═══════════════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║           ACCIL Installation Wizard                  ║" -ForegroundColor Cyan
Write-Host "╚═══════════════════════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host ""

# Check if Go is installed
Write-Info "Checking Go installation..."
try {
    $goVersion = go version 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Go is installed: $goVersion"
    } else {
        throw "Go not found"
    }
} catch {
    Write-Error-Custom "Go is not installed or not in PATH"
    Write-Host ""
    Write-Host "Please install Go from https://golang.org/dl/" -ForegroundColor Yellow
    Write-Host "After installation, restart your terminal and run this script again." -ForegroundColor Yellow
    exit 1
}

# Determine installation directory
$InstallDir = "$env:USERPROFILE\.accil\bin"
if (!(Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
    Write-Info "Created installation directory: $InstallDir"
}

# Clone or download the project
$TempDir = Join-Path $env:TEMP "accil-install-$((Get-Date).ToString('yyyyMMddHHmmss'))"
$RepoUrl = "https://github.com/acacMAX/accil.git"

Write-Info "Downloading ACCIL from GitHub..."
try {
    git clone $RepoUrl $TempDir 2>&1 | Out-Null
    if ($LASTEXITCODE -ne 0) {
        throw "Git clone failed"
    }
    Write-Success "Download completed"
} catch {
    Write-Error-Custom "Failed to download from GitHub"
    Write-Host ""
    Write-Host "Please check your internet connection and try again." -ForegroundColor Yellow
    exit 1
}

# Build the project
Write-Info "Building ACCIL..."
Set-Location $TempDir

try {
    # Install dependencies
    Write-Info "Installing dependencies..."
    go mod tidy 2>&1 | Out-Null
    if ($LASTEXITCODE -ne 0) {
        throw "Failed to install dependencies"
    }

    # Build
    Write-Info "Compiling..."
    go build -o "$InstallDir\accil.exe" . 2>&1 | Out-Null
    if ($LASTEXITCODE -ne 0) {
        throw "Build failed"
    }

    Write-Success "Build completed successfully!"
} catch {
    Write-Error-Custom $_.Exception.Message
    Write-Host ""
    Write-Host "Build failed. Cleaning up..." -ForegroundColor Yellow
    
    # Cleanup on failure
    if (Test-Path $TempDir) {
        Remove-Item -Recurse -Force $TempDir -ErrorAction SilentlyContinue
    }
    exit 1
}

# Cleanup temporary files
Write-Info "Cleaning up temporary files..."
Set-Location $InstallDir
if (Test-Path $TempDir) {
    Remove-Item -Recurse -Force $TempDir -ErrorAction SilentlyContinue
}
Write-Success "Cleanup completed"

# Add to PATH if not already present
Write-Info "Checking PATH..."
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*$InstallDir*") {
    Write-Info "Adding $InstallDir to PATH..."
    $newPath = "$currentPath;$InstallDir"
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    Write-Success "Added to PATH"

    # Refresh current session PATH
    $env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
} else {
    Write-Success "PATH already configured"
}

# Verify installation
Write-Info "Verifying installation..."
try {
    & "$InstallDir\accil.exe" version
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Installation verified!"
    }
} catch {
    Write-Error-Custom "Verification failed"
    exit 1
}

Write-Host ""
Write-Host "╔═══════════════════════════════════════════════════════╗" -ForegroundColor Green
Write-Host "║              Installation Complete!                    ║" -ForegroundColor Green
Write-Host "╚═══════════════════════════════════════════════════════╝" -ForegroundColor Green
Write-Host ""
Write-Host "You can now run ACCIL by typing:" -ForegroundColor Cyan
Write-Host "  accil" -ForegroundColor White
Write-Host ""
Write-Host "For first-time setup, run:" -ForegroundColor Cyan
Write-Host "  accil --setup" -ForegroundColor White
Write-Host ""
Write-Host "Would you like to run ACCIL now? (Y/N): " -ForegroundColor Yellow -NoNewline
$response = Read-Host

if ($response -eq "Y" -or $response -eq "y") {
    Write-Host ""
    & "$InstallDir\accil.exe"
}
