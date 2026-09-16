# Installation & Downloads

Install **Vidonex** on macOS, Linux, or Windows via your favorite package manager or download standalone binaries.

---

## ⚡ Universal 1-Line Installer

The quickest way to install the latest `vidonex` binary on any system:

### macOS & Linux (Bash / Zsh)
```bash
curl -fsSL https://raw.githubusercontent.com/farshidrezaei/vidonex/main/install.sh | bash
```

### Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/farshidrezaei/vidonex/main/install.ps1 | iex
```

### ✨ What the Installer Does Automatically:
1. **Detects OS & Architecture Matrix**: Supports Linux (`x86_64`), macOS Universal (`Apple Silicon` + `Intel`), and Windows.
2. **Cryptographic SHA-256 Verification**: Verifies SHA-256 signatures before binary execution.
3. **Zero-Sudo Smart Pathing**: Prioritizes standard user space (`~/.local/bin`) or `/usr/local/bin` without prompting for sudo unless needed.
4. **Environment `$PATH` Configuration**: Automatically updates your shell configuration (`.bashrc`, `.zshrc`, `.config/fish/config.fish`, or Windows User PATH).
5. **System Readiness Check**: Probes for `ffmpeg` and detects available GPU hardware acceleration (NVENC, VideoToolbox, VA-API, QSV).

### 🛠️ Advanced Installer Options:
```bash
# Install a specific release version
VERSION=v1.0.0 curl -fsSL https://raw.githubusercontent.com/farshidrezaei/vidonex/main/install.sh | bash

# Install the Desktop Studio GUI Workstation instead of headless CLI
APP_TYPE=desktop curl -fsSL https://raw.githubusercontent.com/farshidrezaei/vidonex/main/install.sh | bash

# Install into a custom directory without touching shell profiles
INSTALL_DIR=~/.bin NO_MODIFY_PATH=1 curl -fsSL https://raw.githubusercontent.com/farshidrezaei/vidonex/main/install.sh | bash
```

---

## 🍺 Homebrew (macOS & Linux)

Install via the official Vidonex Homebrew tap:

```bash
brew tap farshidrezaei/tap
brew install vidonex
```

To upgrade in the future:
```bash
brew update && brew upgrade vidonex
```

---

## 🪟 Windows Package Managers

### WinGet (Official Microsoft Store / Windows 10 & 11)
```powershell
winget install Vidonex.Vidonex
```

### Scoop
```powershell
scoop bucket add vidonex https://github.com/farshidrezaei/scoop-bucket
scoop install vidonex
```

---

## 🐧 Linux (Debian / Ubuntu / Arch)

### Direct Debian / Ubuntu `.deb` Package
Download the latest Debian package from GitHub Releases and install with `dpkg`:

```bash
curl -fSL https://github.com/farshidrezaei/vidonex/releases/latest/download/vidonex_amd64.deb -o vidonex.deb
sudo dpkg -i vidonex.deb
```

### Arch Linux (AUR)
```bash
yay -S vidonex-bin
```

---

## 📦 Direct Binary Downloads (Latest Release)

These links always resolve directly to the **latest stable production binaries**:

### 🎬 Desktop Studio Workstation (GUI)

| Platform | Architecture | Download Link |
| :--- | :--- | :--- |
| **macOS** | Universal (Apple Silicon M1/M2/M3/M4 & Intel) | [Download `.dmg` / App](https://github.com/farshidrezaei/vidonex/releases/latest/download/vidonex-darwin-universal) |
| **Linux** | `x86_64` (glibc $\ge$ 2.31) | [Download Executable](https://github.com/farshidrezaei/vidonex/releases/latest/download/vidonex-linux-amd64) |
| **Windows** | `x86_64` (Windows 10/11) | [Download `.exe`](https://github.com/farshidrezaei/vidonex/releases/latest/download/vidonex-windows-amd64.exe) |

### 💻 Headless Automation CLI

| Platform | Architecture | Download Link |
| :--- | :--- | :--- |
| **Linux** | `x86_64` | [Download `vidonex-cli-linux-amd64`](https://github.com/farshidrezaei/vidonex/releases/latest/download/vidonex-cli-linux-amd64) |
| **macOS** | Universal (`arm64` + `amd64`) | [Download `vidonex-cli-darwin-universal`](https://github.com/farshidrezaei/vidonex/releases/latest/download/vidonex-cli-darwin-universal) |
| **Windows** | `x86_64` | [Download `vidonex-cli-windows-amd64.exe`](https://github.com/farshidrezaei/vidonex/releases/latest/download/vidonex-cli-windows-amd64.exe) |

---

## 🐹 Go Toolchain

If you have Go installed (`1.22+`):

```bash
go install github.com/farshidrezaei/vidonex/cmd/vidonex@latest
```

---

## 🔍 Verification & System Diagnostics
 
Once installed, verify that Vidonex is available in your `$PATH` and inspect system diagnostics:

```bash
# Check version
vidonex version

# Comprehensive system hardware specs, FFmpeg engine & GPU acceleration diagnostics
vidonex about

# Verify and check for updates
vidonex upgrade --check
```

---

## 🔄 Self-Upgrades

To keep Vidonex updated to the latest release without manually downloading archives:

```bash
# Check if an update is available
vidonex upgrade --check

# Upgrade to the latest release with atomic binary replacement
vidonex upgrade

# Force reinstall or update to latest release
vidonex upgrade --force
```

Additionally, when running the Web Studio (`vidonex serve`), you can click **"Install Update Automatically"** directly inside the **About** modal or top notification banner to perform a one-click in-app update with live progress streaming.

---

## ⚡ Shell Autocompletions

Vidonex includes native autocompletion for Bash, Zsh, and Fish:

### Zsh
```bash
# Generate completion into your Zsh completions directory:
mkdir -p ~/.zsh/completion
vidonex completion zsh > ~/.zsh/completion/_vidonex

# Add to ~/.zshrc if not already present:
# fpath=(~/.zsh/completion $fpath)
# autoload -Uz compinit && compinit
```

### Bash
```bash
vidonex completion bash | sudo tee /etc/bash_completion.d/vidonex > /dev/null
# Or user-local:
# vidonex completion bash >> ~/.bash_completion
```

### Fish
```fish
mkdir -p ~/.config/fish/completions
vidonex completion fish > ~/.config/fish/completions/vidonex.fish
```
