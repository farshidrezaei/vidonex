#!/usr/bin/env bash
# ==============================================================================
#  🎬 Vidonex Universal Installation Script
#  Declarative Video Composition & FFmpeg Filtergraph Engine
#
#  Usage:
#    curl -fsSL https://raw.githubusercontent.com/farshidrezaei/vidonex/main/install.sh | bash
#
#  Options (via environment variables or CLI flags):
#    VERSION=v1.0.0          Install specific version (default: latest)
#    INSTALL_DIR=~/.bin      Custom install directory
#    APP_TYPE=cli|desktop    Install headless CLI or Desktop Studio Workstation
#    NO_MODIFY_PATH=1        Do not automatically update shell configuration
#    QUIET=1                 Suppress decorative banners and verbose step output
# ==============================================================================

set -euo pipefail

OWNER="farshidrezaei"
REPO="vidonex"
REQUESTED_VERSION="${VERSION:-latest}"
APP_TYPE="${APP_TYPE:-cli}"
NO_MODIFY_PATH="${NO_MODIFY_PATH:-0}"
QUIET="${QUIET:-0}"

# --- Terminal Capabilities & Visual Styling -----------------------------------
setup_terminal_styles() {
    if [ -t 1 ] && [ "${QUIET}" != "1" ]; then
        BOLD="$(tput bold 2>/dev/null || printf '\033[1m')"
        DIM="$(tput dim 2>/dev/null || printf '\033[2m')"
        CYAN="$(tput setaf 6 2>/dev/null || printf '\033[36m')"
        GREEN="$(tput setaf 2 2>/dev/null || printf '\033[32m')"
        YELLOW="$(tput setaf 3 2>/dev/null || printf '\033[33m')"
        BLUE="$(tput setaf 4 2>/dev/null || printf '\033[34m')"
        MAGENTA="$(tput setaf 5 2>/dev/null || printf '\033[35m')"
        RED="$(tput setaf 1 2>/dev/null || printf '\033[31m')"
        WHITE="$(tput setaf 7 2>/dev/null || printf '\033[37m')"
        PURPLE="\033[38;2;168;85;247m"
        INDIGO="\033[38;2;99;102;241m"
        SKYBLUE="\033[38;2;56;189;248m"
        RESET="$(tput sgr0 2>/dev/null || printf '\033[0m')"
    else
        BOLD=""
        DIM=""
        CYAN=""
        GREEN=""
        YELLOW=""
        BLUE=""
        MAGENTA=""
        RED=""
        WHITE=""
        PURPLE=""
        INDIGO=""
        SKYBLUE=""
        RESET=""
    fi
}
setup_terminal_styles

# --- Parse CLI Arguments (if run as executable script) ------------------------
while [ $# -gt 0 ]; do
    case "$1" in
        --version|-v)
            REQUESTED_VERSION="$2"
            shift 2
            ;;
        --dir|-d)
            INSTALL_DIR="$2"
            shift 2
            ;;
        --desktop)
            APP_TYPE="desktop"
            shift
            ;;
        --cli)
            APP_TYPE="cli"
            shift
            ;;
        --no-modify-path)
            NO_MODIFY_PATH=1
            shift
            ;;
        --quiet|-q)
            QUIET=1
            setup_terminal_styles
            shift
            ;;
        --help|-h)
            printf "Usage: install.sh [options]\n\n"
            printf "Options:\n"
            printf "  --version, -v <ver>    Install specific release (e.g. v1.0.0)\n"
            printf "  --dir, -d <path>       Target install directory (default: ~/.local/bin or /usr/local/bin)\n"
            printf "  --desktop              Install Desktop Studio Workstation GUI binary\n"
            printf "  --cli                  Install Headless CLI binary (default)\n"
            printf "  --no-modify-path       Do not append install directory to shell PATH\n"
            printf "  --quiet, -q            Suppress decorative banner\n"
            printf "  --help, -h             Show this help screen\n"
            exit 0
            ;;
        *)
            printf "${RED}Unknown option: %s${RESET}\n" "$1"
            exit 1
            ;;
    esac
done

# --- UI Helpers ---------------------------------------------------------------
print_banner() {
    if [ "${QUIET}" = "1" ]; then return; fi
    printf "\n"
    printf "${PURPLE}${BOLD}   ██╗   ██╗██╗██████╗  ██████╗ ███╗   ██╗███████╗██╗  ██╗${RESET}\n"
    printf "${PURPLE}${BOLD}   ██║   ██║██║██╔══██╗██╔═══██╗████╗  ██║██╔════╝╚██╗██╔╝${RESET}\n"
    printf "${INDIGO}${BOLD}   ██║   ██║██║██║  ██║██║   ██║██╔██╗ ██║█████╗   ╚███╔╝ ${RESET}\n"
    printf "${INDIGO}${BOLD}   ╚██╗ ██╔╝██║██║  ██║██║   ██║██║╚██╗██║██╔══╝   ██╔██╗ ${RESET}\n"
    printf "${SKYBLUE}${BOLD}    ╚████╔╝ ██║██████╔╝╚██████╔╝██║ ╚████║███████╗██╔╝ ██╗${RESET}\n"
    printf "${SKYBLUE}${BOLD}     ╚═══╝  ╚═╝╚═════╝  ╚═════╝ ╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝${RESET}\n"
    printf "\n"
    printf "   ${WHITE}${BOLD}Vidonex Video Engine Universal Installer${RESET}\n"
    printf "   ${DIM}Declarative Video Composition & FFmpeg Filtergraph Compiler${RESET}\n"
    printf "   ${DIM}─────────────────────────────────────────────────────────────${RESET}\n\n"
}

step_badge() {
    local num="$1"
    local title="$2"
    printf "${CYAN}${BOLD}[%s/6]${RESET} ${WHITE}${BOLD}%s${RESET}\n" "${num}" "${title}"
}

success_msg() {
    printf "  ${GREEN}${BOLD}✓${RESET} %s\n" "$1"
}

info_msg() {
    printf "  ${CYAN}${BOLD}ℹ${RESET} %s\n" "$1"
}

warn_msg() {
    printf "  ${YELLOW}${BOLD}⚠${RESET} %s\n" "$1"
}

error_exit() {
    printf "\n  ${RED}${BOLD}✖ Error:${RESET} %s\n\n" "$1" >&2
    exit 1
}

# --- Cleanup Trap -------------------------------------------------------------
TMP_DIR="$(mktemp -d 2>/dev/null || mktemp -d -t 'vidonex-install.XXXXXX')"
cleanup() {
    rm -rf "${TMP_DIR}"
}
trap cleanup EXIT INT TERM

# ==============================================================================
# Pipeline Execution
# ==============================================================================

print_banner

# Step 1: Detect Operating System & Architecture Matrix
step_badge "1" "Detecting platform and architecture..."

RAW_OS="$(uname -s)"
RAW_ARCH="$(uname -m)"

case "${RAW_OS}" in
    Linux*)
        OS_TYPE="linux"
        OS_LABEL="Linux"
        ;;
    Darwin*)
        OS_TYPE="darwin"
        OS_LABEL="macOS"
        ;;
    CYGWIN*|MINGW*|MSYS*)
        OS_TYPE="windows"
        OS_LABEL="Windows (POSIX layer)"
        ;;
    *)
        error_exit "Unsupported operating system: ${RAW_OS}"
        ;;
esac

# Handle macOS Rosetta 2 emulation translation
if [ "${OS_TYPE}" = "darwin" ]; then
    if [ "${RAW_ARCH}" = "x86_64" ]; then
        IS_TRANSLATED="$(sysctl -in sysctl.proc_translated 2>/dev/null || echo '0')"
        if [ "${IS_TRANSLATED}" = "1" ]; then
            RAW_ARCH="arm64"
        fi
    fi
fi

case "${RAW_ARCH}" in
    x86_64|amd64)
        ARCH_TYPE="amd64"
        ARCH_LABEL="x86_64"
        ;;
    arm64*|aarch64*)
        ARCH_TYPE="arm64"
        ARCH_LABEL="Apple Silicon / ARM64"
        ;;
    *)
        error_exit "Unsupported CPU architecture: ${RAW_ARCH}"
        ;;
esac

success_msg "Platform verified: ${BOLD}${OS_LABEL} (${ARCH_LABEL})${RESET}"

# Step 2: Resolve Artifact & Target Binary
step_badge "2" "Resolving target package artifact..."

if [ "${APP_TYPE}" = "desktop" ]; then
    BIN_NAME="vidonex"
    if [ "${OS_TYPE}" = "darwin" ]; then
        ARTIFACT_NAME="vidonex-darwin-universal"
    elif [ "${OS_TYPE}" = "linux" ]; then
        ARTIFACT_NAME="vidonex-linux-amd64"
    elif [ "${OS_TYPE}" = "windows" ]; then
        ARTIFACT_NAME="vidonex-windows-amd64.exe"
        BIN_NAME="vidonex.exe"
    fi
    info_msg "Package selected: ${BOLD}Vidonex Studio (Desktop GUI Workstation)${RESET}"
else
    BIN_NAME="vidonex"
    if [ "${OS_TYPE}" = "darwin" ]; then
        ARTIFACT_NAME="vidonex-cli-darwin-universal"
    elif [ "${OS_TYPE}" = "linux" ]; then
        ARTIFACT_NAME="vidonex-cli-linux-amd64"
    elif [ "${OS_TYPE}" = "windows" ]; then
        ARTIFACT_NAME="vidonex-cli-windows-amd64.exe"
        BIN_NAME="vidonex.exe"
    fi
    info_msg "Package selected: ${BOLD}Vidonex CLI (Headless Composition Engine)${RESET}"
fi

if [ "${REQUESTED_VERSION}" = "latest" ]; then
    DOWNLOAD_BASE="https://github.com/${OWNER}/${REPO}/releases/latest/download"
else
    CLEAN_TAG="${REQUESTED_VERSION#v}"
    DOWNLOAD_BASE="https://github.com/${OWNER}/${REPO}/releases/download/v${CLEAN_TAG}"
fi

DOWNLOAD_URL="${DOWNLOAD_BASE}/${ARTIFACT_NAME}"
CHECKSUMS_URL="${DOWNLOAD_BASE}/checksums.txt?t=$(date +%s)"

info_msg "Target artifact: ${DIM}${ARTIFACT_NAME}${RESET}"

# Step 3: Download Artifact and SHA-256 Checksums
step_badge "3" "Downloading binary and checksums..."

download_file() {
    local url="$1"
    local dest="$2"

    if command -v curl >/dev/null 2>&1; then
        if [ -t 1 ] && [ "${QUIET}" != "1" ]; then
            curl -fSL# "${url}" -o "${dest}"
        else
            curl -fSL "${url}" -o "${dest}"
        fi
    elif command -v wget >/dev/null 2>&1; then
        if [ -t 1 ] && [ "${QUIET}" != "1" ]; then
            wget --progress=bar:force "${url}" -O "${dest}"
        else
            wget -q "${url}" -O "${dest}"
        fi
    else
        error_exit "Neither curl nor wget was found on your system. Please install one to continue."
    fi
}

info_msg "Fetching ${DOWNLOAD_URL}..."
if ! download_file "${DOWNLOAD_URL}" "${TMP_DIR}/${ARTIFACT_NAME}"; then
    error_exit "Failed downloading artifact from: ${DOWNLOAD_URL}"
fi

FILE_SIZE_BYTES="$(wc -c < "${TMP_DIR}/${ARTIFACT_NAME}" 2>/dev/null || stat -f%z "${TMP_DIR}/${ARTIFACT_NAME}" 2>/dev/null || echo 0)"
FILE_SIZE_MB="$(awk "BEGIN {printf \"%.2f\", ${FILE_SIZE_BYTES}/1048576}")"
success_msg "Artifact downloaded (${FILE_SIZE_MB} MB)"

info_msg "Fetching integrity checksums..."
if download_file "${CHECKSUMS_URL}" "${TMP_DIR}/checksums.txt" 2>/dev/null; then
    CHECKSUMS_AVAILABLE=1
else
    CHECKSUMS_AVAILABLE=0
    warn_msg "Checksum manifest unavailable for this release; skipping checksum validation"
fi

# Step 4: Cryptographic SHA-256 Verification
step_badge "4" "Verifying cryptographic integrity..."

calculate_sha256() {
    local target="$1"
    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum "${target}" | awk '{print $1}'
    elif command -v shasum >/dev/null 2>&1; then
        shasum -a 256 "${target}" | awk '{print $1}'
    elif command -v openssl >/dev/null 2>&1; then
        openssl dgst -sha256 "${target}" | awk '{print $NF}'
    elif command -v python3 >/dev/null 2>&1; then
        python3 -c "import hashlib, sys; print(hashlib.sha256(open(sys.argv[1], 'rb').read()).hexdigest())" "${target}"
    elif command -v python >/dev/null 2>&1; then
        python -c "import hashlib, sys; print(hashlib.sha256(open(sys.argv[1], 'rb').read()).hexdigest())" "${target}"
    else
        echo ""
    fi
}

if [ "${CHECKSUMS_AVAILABLE}" = "1" ]; then
    EXPECTED_HASH="$(tr -d '\r' < "${TMP_DIR}/checksums.txt" | awk -v name="${ARTIFACT_NAME}" '$2 == name || $2 == "*"name {print $1; exit}')"
    if [ -n "${EXPECTED_HASH}" ]; then
        ACTUAL_HASH="$(calculate_sha256 "${TMP_DIR}/${ARTIFACT_NAME}")"
        if [ -n "${ACTUAL_HASH}" ]; then
            if [ "$(echo "${ACTUAL_HASH}" | tr '[:upper:]' '[:lower:]')" != "$(echo "${EXPECTED_HASH}" | tr '[:upper:]' '[:lower:]')" ]; then
                printf "\n  ${RED}${BOLD}SECURITY ALERT: Checksum mismatch!${RESET}\n"
                printf "  Expected: %s\n" "${EXPECTED_HASH}"
                printf "  Actual:   %s\n" "${ACTUAL_HASH}"
                error_exit "Aborting installation to prevent compromised or corrupt execution."
            fi
            success_msg "SHA-256 integrity verified: ${DIM}${ACTUAL_HASH:0:16}...${RESET}"
        else
            warn_msg "No SHA-256 calculation utility found; skipped hash comparison"
        fi
    else
        warn_msg "Artifact not listed in checksums.txt; proceeding with size verification"
    fi
fi

# Step 5: Placement & Path Configuration (Zero-Sudo First)
step_badge "5" "Installing binary and configuring environment..."

chmod +x "${TMP_DIR}/${ARTIFACT_NAME}"

# Determine intelligent install directory
if [ -z "${INSTALL_DIR:-}" ]; then
    CURRENT_UID="$(id -u 2>/dev/null || echo 1000)"
    if [ "${CURRENT_UID}" -eq 0 ]; then
        INSTALL_DIR="/usr/local/bin"
    elif [ -w "/usr/local/bin" ]; then
        INSTALL_DIR="/usr/local/bin"
    else
        # Preferred user space directory (modern standard: XDG ~/.local/bin)
        INSTALL_DIR="${HOME}/.local/bin"
    fi
fi

mkdir -p "${INSTALL_DIR}" 2>/dev/null || true

# Copy file into target directory
TARGET_PATH="${INSTALL_DIR}/${BIN_NAME}"
if [ -w "${INSTALL_DIR}" ]; then
    cp -f "${TMP_DIR}/${ARTIFACT_NAME}" "${TARGET_PATH}"
else
    if command -v sudo >/dev/null 2>&1; then
        info_msg "Elevating with sudo to install into ${INSTALL_DIR}..."
        sudo cp -f "${TMP_DIR}/${ARTIFACT_NAME}" "${TARGET_PATH}"
        sudo chmod +x "${TARGET_PATH}"
    else
        FALLBACK_DIR="${HOME}/.local/bin"
        mkdir -p "${FALLBACK_DIR}"
        cp -f "${TMP_DIR}/${ARTIFACT_NAME}" "${FALLBACK_DIR}/${BIN_NAME}"
        INSTALL_DIR="${FALLBACK_DIR}"
        TARGET_PATH="${INSTALL_DIR}/${BIN_NAME}"
        warn_msg "Permissions denied for original directory. Installed to: ${TARGET_PATH}"
    fi
fi

chmod +x "${TARGET_PATH}" 2>/dev/null || true
success_msg "Binary installed to ${BOLD}${TARGET_PATH}${RESET}"

# Check whether INSTALL_DIR is in user's active PATH
PATH_UPDATED=0
PATH_DETECTED=0
case ":${PATH}:" in
    *:"${INSTALL_DIR}":*)
        PATH_DETECTED=1
        ;;
    *)
        PATH_DETECTED=0
        ;;
esac

TARGET_PROFILE=""
if [ "${PATH_DETECTED}" -eq 1 ]; then
    success_msg "Installation directory is already in your \${PATH}"
elif [ "${NO_MODIFY_PATH}" != "1" ]; then
    USER_SHELL="$(basename "${SHELL:-bash}")"

    case "${USER_SHELL}" in
        zsh)
            TARGET_PROFILE="${HOME}/.zshrc"
            [ ! -f "${TARGET_PROFILE}" ] && [ -f "${HOME}/.zprofile" ] && TARGET_PROFILE="${HOME}/.zprofile"
            ;;
        bash)
            if [ "${OS_TYPE}" = "darwin" ]; then
                TARGET_PROFILE="${HOME}/.bash_profile"
                [ ! -f "${TARGET_PROFILE}" ] && [ -f "${HOME}/.bashrc" ] && TARGET_PROFILE="${HOME}/.bashrc"
            else
                TARGET_PROFILE="${HOME}/.bashrc"
                [ ! -f "${TARGET_PROFILE}" ] && [ -f "${HOME}/.profile" ] && TARGET_PROFILE="${HOME}/.profile"
            fi
            ;;
        fish)
            TARGET_PROFILE="${HOME}/.config/fish/config.fish"
            ;;
        *)
            TARGET_PROFILE="${HOME}/.profile"
            ;;
    esac

    if [ -n "${TARGET_PROFILE}" ] && [ -f "${TARGET_PROFILE}" ] && [ -w "${TARGET_PROFILE}" ]; then
        LINE_TO_ADD="export PATH=\"${INSTALL_DIR}:\$PATH\""
        if [ "${USER_SHELL}" = "fish" ]; then
            LINE_TO_ADD="fish_add_path ${INSTALL_DIR}"
        fi

        if ! grep -Fq "${INSTALL_DIR}" "${TARGET_PROFILE}" 2>/dev/null; then
            printf "\n# Vidonex Video Engine\n%s\n" "${LINE_TO_ADD}" >> "${TARGET_PROFILE}"
            PATH_UPDATED=1
            success_msg "Added ${BOLD}${INSTALL_DIR}${RESET} to ${BOLD}${TARGET_PROFILE}${RESET}"
        fi
    fi
fi

# Step 6: System Readiness & FFmpeg Hardware Acceleration Probing
step_badge "6" "Probing video rendering engine & FFmpeg..."

if command -v ffmpeg >/dev/null 2>&1; then
    FFMPEG_RAW_VERSION="$(ffmpeg -version 2>/dev/null | head -n 1 | awk '{print $3}' || echo 'unknown')"
    FFMPEG_PATH="$(command -v ffmpeg)"
    success_msg "FFmpeg detected: ${BOLD}${FFMPEG_RAW_VERSION}${RESET} (${DIM}${FFMPEG_PATH}${RESET})"

    ENCODERS="$(ffmpeg -encoders 2>/dev/null || echo '')"
    ACCELERATORS=""

    if echo "${ENCODERS}" | grep -q "h264_nvenc"; then
        ACCELERATORS="${ACCELERATORS} NVENC"
    fi
    if echo "${ENCODERS}" | grep -q "videotoolbox"; then
        ACCELERATORS="${ACCELERATORS} VideoToolbox"
    fi
    if echo "${ENCODERS}" | grep -q "h264_vaapi"; then
        ACCELERATORS="${ACCELERATORS} VA-API"
    fi
    if echo "${ENCODERS}" | grep -q "h264_qsv"; then
        ACCELERATORS="${ACCELERATORS} QSV"
    fi

    if [ -n "${ACCELERATORS}" ]; then
        success_msg "Hardware accelerators available:${GREEN}${BOLD}${ACCELERATORS}${RESET}"
    else
        info_msg "Software rendering available (libx264/libx265)"
    fi
else
    warn_msg "FFmpeg not detected in PATH. Vidonex requires FFmpeg for video rendering."
    printf "\n  ${YELLOW}Quick install command for your system:${RESET}\n"
    if [ "${OS_TYPE}" = "darwin" ]; then
        printf "    ${BOLD}brew install ffmpeg${RESET}\n\n"
    elif [ "${OS_TYPE}" = "linux" ]; then
        if command -v apt-get >/dev/null 2>&1; then
            printf "    ${BOLD}sudo apt-get update && sudo apt-get install -y ffmpeg${RESET}\n\n"
        elif command -v pacman >/dev/null 2>&1; then
            printf "    ${BOLD}sudo pacman -S ffmpeg${RESET}\n\n"
        elif command -v dnf >/dev/null 2>&1; then
            printf "    ${BOLD}sudo dnf install -y ffmpeg${RESET}\n\n"
        elif command -v apk >/dev/null 2>&1; then
            printf "    ${BOLD}sudo apk add ffmpeg${RESET}\n\n"
        else
            printf "    ${BOLD}Please install ffmpeg through your system package manager.${RESET}\n\n"
        fi
    fi
fi

# --- Final Executive Summary Box ----------------------------------------------
INSTALLED_VERSION="$("${TARGET_PATH}" version 2>/dev/null | awk '{print $NF}' || echo '1.0.0')"

BOX_LINE="────────────────────────────────────────────────────────────"
printf "\n"
printf "${CYAN}┌%s┐${RESET}\n" "${BOX_LINE}"
printf "${CYAN}│${RESET}  ${GREEN}${BOLD}✨ Vidonex Successfully Installed!${RESET}\n"
printf "${CYAN}├%s┤${RESET}\n" "${BOX_LINE}"
printf "${CYAN}│${RESET}  ${WHITE}%-20s${RESET} : ${BOLD}%s${RESET}\n" "Binary Location" "${TARGET_PATH}"
printf "${CYAN}│${RESET}  ${WHITE}%-20s${RESET} : ${BOLD}%s${RESET}\n" "Engine Version" "v${INSTALLED_VERSION#v}"
printf "${CYAN}│${RESET}  ${WHITE}%-20s${RESET} : %s\n" "Platform Target" "${OS_LABEL} (${ARCH_LABEL})"
printf "${CYAN}│${RESET}  ${WHITE}%-20s${RESET} : %s\n" "Package Type" "${APP_TYPE}"
printf "${CYAN}├%s┤${RESET}\n" "${BOX_LINE}"
printf "${CYAN}│${RESET}  ${BOLD}Quick Start:${RESET}\n"

if [ "${PATH_DETECTED}" -eq 0 ] && [ "${PATH_UPDATED}" -eq 1 ]; then
    printf "${CYAN}│${RESET}    ${YELLOW}1. Reload your shell:${RESET} ${BOLD}source %s${RESET}\n" "${TARGET_PROFILE}"
    printf "${CYAN}│${RESET}    ${YELLOW}2. Start Web Studio:${RESET}  ${BOLD}vidonex serve --port 8080${RESET}\n"
elif [ "${PATH_DETECTED}" -eq 0 ]; then
    printf "${CYAN}│${RESET}    ${YELLOW}1. Add to PATH:${RESET}       ${BOLD}export PATH=\"%s:\$PATH\"${RESET}\n" "${INSTALL_DIR}"
    printf "${CYAN}│${RESET}    ${YELLOW}2. Start Web Studio:${RESET}  ${BOLD}vidonex serve --port 8080${RESET}\n"
else
    printf "${CYAN}│${RESET}    ${BOLD}vidonex serve --port 8080${RESET}       (Launch Web Studio)\n"
    printf "${CYAN}│${RESET}    ${BOLD}vidonex render project.yaml${RESET}     (Compile & render video)\n"
    printf "${CYAN}│${RESET}    ${BOLD}vidonex --help${RESET}                  (Explore CLI commands)\n"
fi

printf "${CYAN}│${RESET}\n"
printf "${CYAN}│${RESET}  ${DIM}Documentation: https://farshidrezaei.github.io/vidonex/${RESET}\n"
printf "${CYAN}└%s┘${RESET}\n\n" "${BOX_LINE}"
