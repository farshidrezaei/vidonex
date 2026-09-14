# Installation & Downloads

Install **Vidonex** on macOS, Linux, or Windows via your favorite package manager or download standalone binaries.

---

## ⚡ Universal 1-Line Installer (macOS & Linux)

The quickest way to install the latest `vidonex` CLI on Linux and macOS:

```bash
curl -fsSL https://raw.githubusercontent.com/farshidrezaei/vidonex/main/install.sh | bash
```

This script automatically detects your operating system and CPU architecture, downloads the latest binary from GitHub Releases, and puts it in `/usr/local/bin/vidonex`.

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

## 🔍 Verification

Once installed, verify that Vidonex is available in your `$PATH`:

```bash
vidonex version
vidonex --help
```
