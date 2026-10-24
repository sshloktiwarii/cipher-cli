#!/usr/bin/env bash
# ==============================================================================
# Cipher AI PaaS: Production One-Line Installer
# Usage: curl -fsSL https://raw.githubusercontent.com/sshloktiwarii/cipher-cli/master/scripts/install.sh | bash
# ==============================================================================
set -euo pipefail

CYAN='\033[0;36m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BOLD='\033[1m'
NC='\033[0m'

echo -e "${CYAN}${BOLD}"
echo "=================================================="
echo "🛡️  CIPHER AI PAAS — 5-SECOND ONE-LINE INSTALLER"
echo "The Cockpit & Airbag for Autonomous AI Agents"
echo "=================================================="
echo -e "${NC}"

# 1. Detect OS
OS_RAW="$(uname -s)"
case "${OS_RAW}" in
    Darwin) OS="darwin" ;;
    Linux)  OS="linux" ;;
    *)
        echo -e "${RED}Error: Unsupported operating system: ${OS_RAW}${NC}"
        exit 1
        ;;
esac

# 2. Detect Architecture
ARCH_RAW="$(uname -m)"
case "${ARCH_RAW}" in
    x86_64|amd64)   ARCH="amd64" ;;
    arm64|aarch64)  ARCH="arm64" ;;
    *)
        echo -e "${RED}Error: Unsupported hardware architecture: ${ARCH_RAW}${NC}"
        exit 1
        ;;
esac

TARGET="${OS}-${ARCH}"
echo -e "➔ Detected System Target: ${GREEN}${BOLD}${TARGET}${NC}"

# 3. Determine Installation Destination
INSTALL_DIR="/usr/local/bin"
USE_SUDO=""

if [ ! -w "${INSTALL_DIR}" ]; then
    if command -v sudo >/dev/null 2>&1; then
        USE_SUDO="sudo"
    else
        INSTALL_DIR="${HOME}/.local/bin"
        mkdir -p "${INSTALL_DIR}"
    fi
fi

TEMP_DIR="$(mktemp -d)"
trap 'rm -rf "${TEMP_DIR}"' EXIT

BINARY_URL="https://github.com/sshloktiwarii/cipher-cli/releases/latest/download/cipher-${TARGET}"
DEST_PATH="${INSTALL_DIR}/cipher"

echo -e "➔ Downloading pre-compiled static binary for ${TARGET}..."
if curl -fsSL --connect-timeout 5 --max-time 30 "${BINARY_URL}" -o "${TEMP_DIR}/cipher" 2>/dev/null; then
    echo -e "${GREEN}✓ Downloaded universal release binary${NC}"
else
    # Fallback to local go build if source tree is present or go is installed
    echo -e "${YELLOW}Notice: Release asset not yet published; checking for Go compiler fallback...${NC}"
    if command -v go >/dev/null 2>&1; then
        echo -e "➔ Compiling stripped static binary via Go..."
        CGO_ENABLED=0 go install -ldflags="-s -w" github.com/sshloktiwarii/cipher-cli@latest 2>/dev/null || \
        (cd "${TEMP_DIR}" && git clone --depth 1 https://github.com/sshloktiwarii/cipher-cli.git && cd cipher-cli && CGO_ENABLED=0 go build -ldflags="-s -w" -o "${TEMP_DIR}/cipher" .)
    else
        echo -e "${RED}Error: Could not retrieve release asset and Go compiler is unavailable.${NC}"
        exit 1
    fi
fi

chmod +x "${TEMP_DIR}/cipher"

echo -e "➔ Installing to ${GREEN}${DEST_PATH}${NC}..."
${USE_SUDO} mv "${TEMP_DIR}/cipher" "${DEST_PATH}"

# 4. Verification
if command -v cipher >/dev/null 2>&1 || [ -x "${DEST_PATH}" ]; then
    echo -e "\n${GREEN}${BOLD}=================================================="
    echo "🎉 CIPHER SUCCESSFULLY INSTALLED TO ${DEST_PATH}"
    echo "=================================================="
    echo -e "${NC}"
    echo "Quick Start Commands:"
    echo -e "  ${CYAN}cipher appliance${NC}   ➔ Boots full 8-service cluster in single supervised process (<500ms)"
    echo -e "  ${CYAN}cipher doctor${NC}      ➔ Runs pre-flight diagnostics across ports and memory"
    echo -e "  ${CYAN}cipher status${NC}      ➔ Displays real-time health across all services"
    echo -e "  ${CYAN}cipher test${NC}        ➔ Runs synthetic end-to-end integration tests"
    echo ""
    echo -e "Visit ${BOLD}http://localhost:8080/hud${NC} once booted for the Cyber-Cockpit War Room."
else
    echo -e "${RED}Installation failed. Please ensure ${INSTALL_DIR} is in your PATH.${NC}"
    exit 1
fi
