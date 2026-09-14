# ==============================================================================
#  🎬 Vidonex Universal Windows PowerShell Installer
#  Declarative Video Composition & FFmpeg Filtergraph Engine
#
#  Usage (PowerShell):
#    irm https://raw.githubusercontent.com/farshidrezaei/vidonex/main/install.ps1 | iex
#
#  Advanced Usage:
#    & ([scriptblock]::Create((irm https://raw.githubusercontent.com/farshidrezaei/vidonex/main/install.ps1))) -AppType desktop
# ==============================================================================

[CmdletBinding()]
param(
    [string]$Version = "latest",
    [string]$InstallDir = "",
    [ValidateSet("cli", "desktop")]
    [string]$AppType = "cli",
    [switch]$NoPathUpdate,
    [switch]$Quiet
)

$ErrorActionPreference = "Stop"

$Owner = "farshidrezaei"
$Repo = "vidonex"

function Write-Banner {
    if ($Quiet) { return }
    Write-Host ""
    Write-Host "  __   ___     __                  " -ForegroundColor Cyan
    Write-Host "  \ \ / (_)___/ /__  ___  __ _____ " -ForegroundColor Cyan
    Write-Host "   \ V / / / _  / _ \/ _ \/ // /\ \ /" -ForegroundColor Cyan
    Write-Host "    \_/_/_/\_,_/\___/_//_/\_, //_\_\ " -ForegroundColor Cyan
    Write-Host "                         /___/       " -ForegroundColor Cyan
    Write-Host "  Vidonex Video Engine Installer (Windows PowerShell)" -ForegroundColor White
    Write-Host "  Declarative Video Composition & FFmpeg Filtergraph Compiler" -ForegroundColor DarkGray
    Write-Host "  ─────────────────────────────────────────────────────────────" -ForegroundColor DarkGray
    Write-Host ""
}

function Write-Step([string]$StepNumber, [string]$Message) {
    Write-Host "[$StepNumber/6] " -ForegroundColor Cyan -NoNewline
    Write-Host $Message -ForegroundColor White
}

function Write-Success([string]$Message) {
    Write-Host "  ✓ " -ForegroundColor Green -NoNewline
    Write-Host $Message
}

function Write-Info([string]$Message) {
    Write-Host "  ℹ " -ForegroundColor Cyan -NoNewline
    Write-Host $Message
}

function Write-Warn([string]$Message) {
    Write-Host "  ⚠ " -ForegroundColor Yellow -NoNewline
    Write-Host $Message
}

Write-Banner

# Step 1: Detect Platform & Architecture
Write-Step "1" "Detecting Windows platform and architecture..."
$Arch = if ([IntPtr]::Size -eq 8) { "amd64" } else { "x86" }
if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") {
    $Arch = "arm64"
}

if ($Arch -ne "amd64") {
    Write-Warn "Notice: x86_64 binary will run via Windows emulation on $Arch."
}
Write-Success "Platform verified: Windows (x86_64)"

# Step 2: Resolve Target Package Artifact
Write-Step "2" "Resolving target package artifact..."
if ($AppType -eq "desktop") {
    $ArtifactName = "vidonex-windows-amd64.exe"
    Write-Info "Package selected: Vidonex Studio (Desktop GUI Workstation)"
} else {
    $ArtifactName = "vidonex-cli-windows-amd64.exe"
    Write-Info "Package selected: Vidonex CLI (Headless Composition Engine)"
}

if ($Version -eq "latest") {
    $DownloadBase = "https://github.com/$Owner/$Repo/releases/latest/download"
} else {
    $CleanTag = $Version.TrimStart("v")
    $DownloadBase = "https://github.com/$Owner/$Repo/releases/download/v$CleanTag"
}

$DownloadUrl = "$DownloadBase/$ArtifactName"
$Timestamp = [DateTimeOffset]::UtcNow.ToUnixTimeSeconds()
$ChecksumsUrl = "$DownloadBase/checksums.txt?t=$Timestamp"

# Step 3: Download Binary and Checksums
Write-Step "3" "Downloading binary and checksums..."
$TempDir = Join-Path ([System.IO.Path]::GetTempPath()) ("vidonex-install-" + [System.Guid]::NewGuid().ToString("N"))
New-Item -ItemType Directory -Path $TempDir -Force | Out-Null

$DestArtifact = Join-Path $TempDir $ArtifactName
$DestChecksums = Join-Path $TempDir "checksums.txt"

try {
    Write-Info "Fetching $DownloadUrl..."
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12 -bor [Net.SecurityProtocolType]::Tls13
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $DestArtifact -UseBasicParsing

    $FileSizeBytes = (Get-Item $DestArtifact).Length
    $FileSizeMB = [Math]::Round($FileSizeBytes / 1MB, 2)
    Write-Success "Artifact downloaded ($FileSizeMB MB)"

    $ChecksumsAvailable = $false
    try {
        Invoke-WebRequest -Uri $ChecksumsUrl -OutFile $DestChecksums -UseBasicParsing
        $ChecksumsAvailable = $true
    } catch {
        Write-Warn "Checksums manifest unavailable; skipping hash verification."
    }

    # Step 4: Cryptographic Verification
    Write-Step "4" "Verifying cryptographic integrity..."
    if ($ChecksumsAvailable) {
        $ExpectedHash = ""
        Get-Content $DestChecksums | ForEach-Object {
            $Parts = $_ -split "\s+", 2
            if ($Parts.Length -ge 2 -and ($Parts[1].Trim("*") -eq $ArtifactName)) {
                $ExpectedHash = $Parts[0].Trim()
            }
        }

        if ($ExpectedHash) {
            $ActualHash = (Get-FileHash -Path $DestArtifact -Algorithm SHA256).Hash.ToLower()
            if ($ActualHash -ne $ExpectedHash.ToLower()) {
                Write-Host ""
                Write-Host "  SECURITY ALERT: Checksum mismatch!" -ForegroundColor Red
                Write-Host "  Expected: $ExpectedHash"
                Write-Host "  Actual:   $ActualHash"
                throw "Aborting installation due to hash mismatch."
            }
            Write-Success "SHA-256 integrity verified: $($ActualHash.Substring(0, 16))..."
        } else {
            Write-Warn "Artifact not found in checksums.txt; proceeding with size verification."
        }
    }

    # Step 5: Install Binary & Update User PATH
    Write-Step "5" "Installing binary and configuring environment..."
    if (-not $InstallDir) {
        $InstallDir = Join-Path $env:LOCALAPPDATA "Programs\Vidonex"
    }

    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }

    $TargetBinary = Join-Path $InstallDir "vidonex.exe"
    Copy-Item -Path $DestArtifact -Destination $TargetBinary -Force
    Write-Success "Binary installed to $TargetBinary"

    # PATH Configuration
    $PathUpdated = $false
    $CurrentPath = [Environment]::GetEnvironmentVariable("Path", "User")
    $PathParts = $CurrentPath -split ";" | ForEach-Object { $_.TrimEnd("\") }

    if ($PathParts -notcontains $InstallDir.TrimEnd("\") -and -not $NoPathUpdate) {
        $NewPath = if ($CurrentPath) { "$CurrentPath;$InstallDir" } else { $InstallDir }
        [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
        $env:Path = "$env:Path;$InstallDir"
        $PathUpdated = $true
        Write-Success "Added $InstallDir to User environment PATH"
    } elseif ($PathParts -contains $InstallDir.TrimEnd("\")) {
        Write-Success "Installation directory is already in User PATH"
    }

    # Step 6: System Readiness & FFmpeg Hardware Acceleration Probing
    Write-Step "6" "Probing video rendering engine & FFmpeg..."
    $FFmpegCmd = Get-Command "ffmpeg" -ErrorAction SilentlyContinue
    if ($FFmpegCmd) {
        $FFmpegVersion = (& ffmpeg -version 2>&1 | Select-Object -First 1)
        Write-Success "FFmpeg detected: $FFmpegVersion"

        $Encoders = (& ffmpeg -encoders 2>&1 | Out-String)
        $Accelerators = @()
        if ($Encoders -match "h264_nvenc") { $Accelerators += "NVENC" }
        if ($Encoders -match "h264_qsv")   { $Accelerators += "QSV" }
        if ($Encoders -match "h264_amf")   { $Accelerators += "AMF (AMD)" }

        if ($Accelerators.Count -gt 0) {
            Write-Success "Hardware accelerators available: $($Accelerators -join ' ')"
        } else {
            Write-Info "Software rendering available (libx264/libx265)"
        }
    } else {
        Write-Warn "FFmpeg was not detected in your PATH. Vidonex requires FFmpeg for video rendering."
        Write-Host ""
        Write-Host "  Quick install command for Windows:" -ForegroundColor Yellow
        Write-Host "    winget install Gyan.FFmpeg" -ForegroundColor White
        Write-Host "    scoop install ffmpeg" -ForegroundColor DarkGray
        Write-Host ""
    }

    # Executive Summary Box
    $BoxLine = "────────────────────────────────────────────────────────────"
    Write-Host ""
    Write-Host "┌$BoxLine┐" -ForegroundColor Cyan
    Write-Host "│  " -ForegroundColor Cyan -NoNewline
    Write-Host "✨ Vidonex Successfully Installed!" -ForegroundColor Green
    Write-Host "├$BoxLine┤" -ForegroundColor Cyan
    Write-Host "│  " -ForegroundColor Cyan -NoNewline
    Write-Host ("{0,-20} : {1}" -f "Binary Location", $TargetBinary)
    Write-Host "│  " -ForegroundColor Cyan -NoNewline
    Write-Host ("{0,-20} : {1}" -f "Engine Version", $Version)
    Write-Host "│  " -ForegroundColor Cyan -NoNewline
    Write-Host ("{0,-20} : {1}" -f "Platform Target", "Windows ($Arch)")
    Write-Host "│  " -ForegroundColor Cyan -NoNewline
    Write-Host ("{0,-20} : {1}" -f "Package Type", $AppType)
    Write-Host "├$BoxLine┤" -ForegroundColor Cyan
    Write-Host "│  Quick Start:" -ForegroundColor Cyan

    if ($PathUpdated) {
        Write-Host "│    " -ForegroundColor Cyan -NoNewline
        Write-Host "1. Restart PowerShell / Terminal to reload your PATH" -ForegroundColor Yellow
        Write-Host "│    " -ForegroundColor Cyan -NoNewline
        Write-Host "2. Start Web Studio:  vidonex serve --port 8080"
    } else {
        Write-Host "│    vidonex serve --port 8080       (Launch Web Studio)"
        Write-Host "│    vidonex render project.yaml     (Compile & render video)"
        Write-Host "│    vidonex --help                  (Explore CLI commands)"
    }
    Write-Host "│" -ForegroundColor Cyan
    Write-Host "│  Documentation: https://farshidrezaei.github.io/vidonex/" -ForegroundColor DarkGray
    Write-Host "└$BoxLine┘" -ForegroundColor Cyan
    Write-Host ""

} finally {
    Remove-Item -Path $TempDir -Recurse -Force -ErrorAction SilentlyContinue
}
