#!/usr/bin/env bash
# Vidonex Universal Installation Script
# Usage: curl -fsSL https://raw.githubusercontent.com/farshidrezaei/vidonex/main/install.sh | bash

set -e

OWNER="farshidrezaei"
REPO="vidonex"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Visual formatting
BOLD="$(tput bold 2>/dev/null || echo '')"
GREEN="$(tput setaf 2 2>/dev/null || echo '')"
CYAN="$(tput setaf 6 2>/dev/null || echo '')"
YELLOW="$(tput setaf 3 2>/dev/null || echo '')"
RED="$(tput setaf 1 2>/dev/null || echo '')"
RESET="$(tput sgr0 2>/dev/null || echo '')"

printf "${CYAN}${BOLD}🎬 Vidonex Installer${RESET}\n"
printf "Detecting operating system and platform architecture...\n"

OS="$(uname -s)"
ARCH="$(uname -m)"

case "${OS}" in
    Linux*)
        OS_TYPE="linux"
        ;;
    Darwin*)
        OS_TYPE="darwin"
        ;;
    CYGWIN*|MINGW*|MSYS*)
        OS_TYPE="windows"
        ;;
    *)
        printf "${RED}Error: Unsupported operating system: ${OS}${RESET}\n"
        exit 1
        ;;
esac

case "${ARCH}" in
    x86_64|amd64)
        ARCH_TYPE="amd64"
        ;;
    arm64|aarch64)
        ARCH_TYPE="arm64"
        ;;
    *)
        printf "${RED}Error: Unsupported CPU architecture: ${ARCH}${RESET}\n"
        exit 1
        ;;
esac

# Target artifact naming convention
if [ "${OS_TYPE}" = "darwin" ]; then
    CLI_BIN="vidonex-cli-darwin-universal"
elif [ "${OS_TYPE}" = "linux" ]; then
    CLI_BIN="vidonex-cli-linux-amd64"
elif [ "${OS_TYPE}" = "windows" ]; then
    CLI_BIN="vidonex-cli-windows-amd64.exe"
fi

DOWNLOAD_URL="https://github.com/${OWNER}/${REPO}/releases/latest/download/${CLI_BIN}"

printf "Fetching latest release from: ${BOLD}${DOWNLOAD_URL}${RESET}\n"

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "${TMP_DIR}"' EXIT

if command -v curl >/dev/null 2>&1; then
    curl -fSL "${DOWNLOAD_URL}" -o "${TMP_DIR}/vidonex"
elif command -v wget >/dev/null 2>&1; then
    wget -q "${DOWNLOAD_URL}" -O "${TMP_DIR}/vidonex"
else
    printf "${RED}Error: Neither curl nor wget is available.${RESET}\n"
    exit 1
fi

chmod +x "${TMP_DIR}/vidonex"

# Determine target install directory
if [ ! -w "${INSTALL_DIR}" ]; then
    printf "${YELLOW}Root permissions required to install to ${INSTALL_DIR}...${RESET}\n"
    sudo cp "${TMP_DIR}/vidonex" "${INSTALL_DIR}/vidonex"
else
    cp "${TMP_DIR}/vidonex" "${INSTALL_DIR}/vidonex"
fi

printf "${GREEN}${BOLD}✓ Vidonex successfully installed to ${INSTALL_DIR}/vidonex!${RESET}\n\n"
printf "Verify your installation with:\n"
printf "  ${BOLD}vidonex --help${RESET}\n"
printf "  ${BOLD}vidonex serve --port 8080${RESET}\n\n"
printf "Documentation: ${CYAN}https://farshidrezaei.github.io/vidonex/${RESET}\n"
