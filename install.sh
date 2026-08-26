#!/usr/bin/env bash
# Terralings Installer Script
# Usage: curl -fsSL https://raw.githubusercontent.com/dnf0/terralings/main/install.sh | bash

set -euo pipefail

REPO="dnf0/terralings"
GITHUB_API="https://api.github.com/repos/${REPO}/releases/latest"

echo "==> Detecting OS and architecture..."
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "${ARCH}" in
  x86_64|amd64)
    ARCH="amd64"
    ;;
  arm64|aarch64)
    ARCH="arm64"
    ;;
  *)
    echo "Error: Unsupported architecture: ${ARCH}" >&2
    exit 1
    ;;
esac

case "${OS}" in
  linux|darwin)
    ;;
  *)
    echo "Error: Unsupported operating system: ${OS}" >&2
    exit 1
    ;;
esac

echo "==> Fetching latest release info for ${REPO}..."
LATEST_TAG=$(curl -sSL -H "Accept: application/vnd.github.v3+json" "${GITHUB_API}" | grep '"tag_name":' | head -1 | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "${LATEST_TAG}" ]; then
  LATEST_TAG="v0.1.0"
  echo "==> Fallback to tag: ${LATEST_TAG}"
else
  echo "==> Found release ${LATEST_TAG}"
fi

ARCHIVE_NAME="terralings-${LATEST_TAG}-${OS}-${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_TAG}/${ARCHIVE_NAME}"

TEMP_DIR="$(mktemp -d)"
trap 'rm -rf "${TEMP_DIR}"' EXIT

echo "==> Downloading ${DOWNLOAD_URL}..."
curl -fsSL "${DOWNLOAD_URL}" -o "${TEMP_DIR}/${ARCHIVE_NAME}"

echo "==> Extracting ${ARCHIVE_NAME}..."
tar -xzf "${TEMP_DIR}/${ARCHIVE_NAME}" -C "${TEMP_DIR}"

# Determine install location
INSTALL_DIR="/usr/local/bin"
if [ ! -w "${INSTALL_DIR}" ]; then
  INSTALL_DIR="${HOME}/.local/bin"
  mkdir -p "${INSTALL_DIR}"
fi

PKG_DIR="terralings-${LATEST_TAG}-${OS}-${ARCH}"
if [ -f "${TEMP_DIR}/${PKG_DIR}/terralings" ]; then
  cp "${TEMP_DIR}/${PKG_DIR}/terralings" "${INSTALL_DIR}/terralings"
elif [ -f "${TEMP_DIR}/terralings" ]; then
  cp "${TEMP_DIR}/terralings" "${INSTALL_DIR}/terralings"
else
  echo "Error: Could not locate terralings binary inside archive" >&2
  exit 1
fi

chmod +x "${INSTALL_DIR}/terralings"

echo "==> Terralings installed successfully to ${INSTALL_DIR}/terralings!"
echo ""
echo "To get started, run:"
echo "  terralings watch"
