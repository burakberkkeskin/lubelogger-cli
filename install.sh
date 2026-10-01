#!/bin/sh
set -e

REPO="burakberkkeskin/lubelogger-cli"
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

case "$OS" in
  linux|darwin) ;;
  *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

BINARY="lubelog-${OS}-${ARCH}"
URL="https://github.com/${REPO}/releases/latest/download/${BINARY}"
DEST_DIR="/usr/local/bin"

if [ ! -w "$DEST_DIR" ]; then
  SUDO="sudo"
else
  SUDO=""
fi

echo "==> Downloading ${BINARY}..."
curl -fsSL "$URL" -o /tmp/lubelog
chmod +x /tmp/lubelog
$SUDO mv /tmp/lubelog "${DEST_DIR}/lubelog"

echo "==> Installed successfully to ${DEST_DIR}/lubelog"
"${DEST_DIR}/lubelog" --version
