# ACCIL PowerShell installation script

param(
    [string]$Version = "latest",
    [switch]$Help
)

if ($Help) {
    Write-Host "ACCIL Installer for Windows"
    Write-Host ""
    Write-Host "Usage:"
    Write-Host "  .\install.ps1 [-Version <version>] [-Help]"
    exit 0
}

$ErrorActionPreference = "Stop"

function Write-Success {
    param([string]$Message)
    Write-Host "[SUCCESS] " -ForegroundColor Green -NoNewline
    Write-Host $Message
}

function Write-Info {
    param([string]$Message)
    Write-Host "[INFO] " -ForegroundColor Cyan -NoNewline
    Write-Host $Message
}

function Write-Error-Custom {
    param([string]$Message)
    Write-Host "[ERROR] " -ForegroundColor Red -NoNewline
    Write-Host $Message
}

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$InstallDir = Join-Path $env:USERPROFILE ".accil\bin"
$FallbackInstallDir = Join-Path $env:LOCALAPPDATA "accil\bin"
$LastResortInstallDir = Join-Path $env:TEMP "accil\bin"
$CacheDir = Join-Path $env:TEMP "accil-cache"
$env:GOCACHE = Join-Path $CacheDir "go-build"
$env:GOMODCACHE = Join-Path $CacheDir "gomod"
$RepoUrl = "https://github.com/acacMAX/accil.git"
$TempDir = Join-Path $env:TEMP ("accil-install-" + [guid]::NewGuid().ToString("N"))
$BuildOutput = Join-Path $env:TEMP ("accil-build-" + [guid]::NewGuid().ToString("N") + ".exe")
$SourceDir = $null

Write-Host ""
Write-Host "========================================"
Write-Host "   ACCIL Installation Wizard"
Write-Host "========================================"
Write-Host ""

Write-Info "Checking Go installation..."
try {
    $goVersion = & go version 2>&1
    if ($LASTEXITCODE -ne 0) {
        throw "Go not found"
    }
    Write-Success "Go is installed: $goVersion"
} catch {
    Write-Error-Custom "Go is not installed or not in PATH."
    exit 1
}

if (Test-Path (Join-Path $ScriptDir "go.mod")) {
    $SourceDir = $ScriptDir
    Write-Info "Using local source tree: $SourceDir"
} else {
    try {
        $null = Get-Command git -ErrorAction Stop
    } catch {
        Write-Error-Custom "Git is required when installing without a local source tree."
        exit 1
    }

    Write-Info "Downloading ACCIL source package..."
    & git clone --depth 1 $RepoUrl $TempDir 2>&1 | Out-Null
    if ($LASTEXITCODE -ne 0) {
        Write-Error-Custom "Failed to download source package."
        exit 1
    }
    $SourceDir = $TempDir
    Write-Success "Download completed"
}

if (!(Test-Path $InstallDir)) {
    try {
        New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
    } catch {
        $InstallDir = $FallbackInstallDir
        try {
            New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
        } catch {
            $InstallDir = $LastResortInstallDir
            New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
        }
    }
}
if (!(Test-Path $env:GOCACHE)) {
    New-Item -ItemType Directory -Force -Path $env:GOCACHE | Out-Null
}
if (!(Test-Path $env:GOMODCACHE)) {
    New-Item -ItemType Directory -Force -Path $env:GOMODCACHE | Out-Null
}

Push-Location $SourceDir
try {
    $accilVer = "1.4.6"
    $verFile = Join-Path $SourceDir "VERSION"
    if (Test-Path $verFile) {
        $accilVer = (Get-Content $verFile -TotalCount 1).Trim()
    }
    if ([string]::IsNullOrWhiteSpace($accilVer)) {
        $accilVer = "1.4.6"
    }

    Write-Info "Building ACCIL v$accilVer..."
    & go build -buildvcs=false "-ldflags=-X github.com/accil/accil/cmd.Version=$accilVer" "-o=$BuildOutput" . 2>&1 | Out-Null
    if ($LASTEXITCODE -ne 0) {
        throw "Build failed"
    }
    try {
        Copy-Item -Force $BuildOutput (Join-Path $InstallDir "accil.exe")
    } catch {
        $InstallDir = $FallbackInstallDir
        try {
            if (!(Test-Path $InstallDir)) {
                New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
            }
            Copy-Item -Force $BuildOutput (Join-Path $InstallDir "accil.exe")
        } catch {
            $InstallDir = $LastResortInstallDir
            if (!(Test-Path $InstallDir)) {
                New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
            }
            Copy-Item -Force $BuildOutput (Join-Path $InstallDir "accil.exe")
        }
    }
    Write-Success "Build completed"
} catch {
    Write-Error-Custom $_.Exception.Message
    Pop-Location
    if (Test-Path $BuildOutput) {
        Remove-Item -Force $BuildOutput -ErrorAction SilentlyContinue
    }
    if ((Test-Path $TempDir) -and ($SourceDir -eq $TempDir)) {
        Remove-Item -Recurse -Force $TempDir -ErrorAction SilentlyContinue
    }
    exit 1
}
Pop-Location
if (Test-Path $BuildOutput) {
    Remove-Item -Force $BuildOutput -ErrorAction SilentlyContinue
}

$currentUserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ([string]::IsNullOrWhiteSpace($currentUserPath)) {
    $currentUserPath = ""
}
if ($currentUserPath -notlike "*$InstallDir*") {
    $newUserPath = if ([string]::IsNullOrWhiteSpace($currentUserPath)) { $InstallDir } else { "$currentUserPath;$InstallDir" }
    [Environment]::SetEnvironmentVariable("Path", $newUserPath, "User")
    $env:Path = [Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + $newUserPath
    Write-Success "Added $InstallDir to user PATH"
} else {
    Write-Success "Install directory already present in PATH"
}

try {
    & "$InstallDir\accil.exe" version 2>&1 | Out-Null
    if ($LASTEXITCODE -ne 0) {
        throw "Verification failed"
    }
    Write-Success "Installation verified"
} catch {
    Write-Error-Custom $_.Exception.Message
    exit 1
}

if ((Test-Path $TempDir) -and ($SourceDir -eq $TempDir)) {
    Remove-Item -Recurse -Force $TempDir -ErrorAction SilentlyContinue
}

Write-Host ""
Write-Host "========================================"
Write-Host "   Installation Complete!"
Write-Host "========================================"
Write-Host ""
Write-Host "Installed to: $InstallDir\accil.exe"
if ($InstallDir -eq $LastResortInstallDir) {
    Write-Host "[WARNING] Primary install directories were not writable." -ForegroundColor Yellow
    Write-Host "[WARNING] Installed to a temp-backed fallback directory instead." -ForegroundColor Yellow
}
Write-Host "You can now run: accil"
Write-Host ""
$response = Read-Host "Run ACCIL now? (y/n)"
if ($response -match '^(?i)y(es)?$') {
    & "$InstallDir\accil.exe"
}
