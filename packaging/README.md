# Packaging & Distribution Guide

This directory contains package definitions and manifests for publishing **Vidonex** to system package managers across macOS, Linux, and Windows.

---

## 1. Homebrew (macOS & Linux)

### Option A: Official Personal Tap (Instant Setup)
1. Create a GitHub repository named `homebrew-tap` (or `tap`) under `farshidrezaei`.
2. Copy `packaging/homebrew/vidonex.rb` into `Formula/vidonex.rb` in that repo.
3. Replace `REPLACE_WITH_*_SHA256` with the checksums from the GitHub Release `checksums.txt`.
4. Users can immediately run:
   ```bash
   brew install farshidrezaei/tap/vidonex
   ```

### Option B: Homebrew Core (`brew install vidonex`)
Once the repository gains ~50-75 GitHub stars:
1. Run `brew create https://github.com/farshidrezaei/vidonex/archive/refs/tags/v1.0.0.tar.gz`.
2. Submit a Pull Request to [Homebrew/homebrew-core](https://github.com/Homebrew/homebrew-core).

---

## 2. Windows Package Managers

### WinGet (Official Microsoft Store / CLI)
1. Fork [microsoft/winget-pkgs](https://github.com/microsoft/winget-pkgs).
2. Use the official `wingetcreate` CLI tool:
   ```powershell
   wingetcreate new https://github.com/farshidrezaei/vidonex/releases/download/v1.0.0/vidonex-cli-windows-amd64.exe
   ```
   Or submit `packaging/winget/Vidonex.Vidonex.yaml` to `manifests/v/Vidonex/Vidonex/1.0.0/`.
3. Submit a Pull Request. Once merged, anyone can install via:
   ```powershell
   winget install Vidonex.Vidonex
   ```

### Scoop
1. Create a repository named `scoop-bucket` on GitHub.
2. Place `packaging/scoop/vidonex.json` inside the `bucket/` folder.
3. Users install via:
   ```powershell
   scoop bucket add vidonex https://github.com/farshidrezaei/scoop-bucket
   scoop install vidonex
   ```

---

## 3. Debian / Ubuntu (APT) & Launchpad PPA

1. A standalone `.deb` package can be built using `nfpm` or `dpkg-deb`.
2. To offer `sudo apt-add-repository ppa:farshidrezaei/vidonex`:
   - Create a free account on [Launchpad.net](https://launchpad.net).
   - Create a PPA named `vidonex`.
   - Upload source packages (`debuild -S`).

---

## 4. Arch Linux (AUR)

1. Create a `PKGBUILD` pointing to the release tarball or binary:
   ```bash
   pkgname=vidonex-bin
   pkgver=1.0.0
   pkgrel=1
   arch=('x86_64')
   source=("https://github.com/farshidrezaei/vidonex/releases/download/v${pkgver}/vidonex-cli-linux-amd64")
   ```
2. Push to `aur.archlinux.org/vidonex-bin.git`.
3. Arch users can immediately install via:
   ```bash
   yay -S vidonex-bin
   ```
