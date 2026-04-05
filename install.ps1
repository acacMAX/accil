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

# Build the project
Write-Info "Building ACCIL..."
$BuildDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $BuildDir

try {
    # Install dependencies
    Write-Info "Installing dependencies..."
    go mod tidy
    if ($LASTEXITCODE -ne 0) {
        throw "Failed to install dependencies"
    }

    # Build
    Write-Info "Compiling..."
    go build -o "$InstallDir\accil.exe" .
    if ($LASTEXITCODE -ne 0) {
        throw "Build failed"
    }

    Write-Success "Build completed successfully!"
} catch {
    Write-Error-Custom $_.Exception.Message
    exit 1
}

# Add to PATH if not already present
Write-Info "Checking PATH..."
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*$InstallDir*") {
    Write-Info "Adding $InstallDir to PATH..."
    [Environment]::SetEnvironmentVariable(
        "Path",
        "$currentPath;$InstallDir",
        "User"
    )
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
