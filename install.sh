#!/bin/bash
# Dikogit Local Development Installer Script (without Homebrew)
set -e

BIN_DIR="${HOME}/.local/bin"
PROJECT_DIR="$(cd "$(dirname "$0")" && pwd)"
EXECUTABLE_SRC="${PROJECT_DIR}/bin/git-dikogit"

echo "=========================================="
echo "       Dikogit CLI Local Installer        "
echo "=========================================="

if [ ! -f "$EXECUTABLE_SRC" ]; then
    echo "Error: Executable $EXECUTABLE_SRC not found!"
    exit 1
fi

chmod +x "$EXECUTABLE_SRC"

mkdir -p "$BIN_DIR"
ln -sf "$EXECUTABLE_SRC" "$BIN_DIR/git-dikogit"
ln -sf "$EXECUTABLE_SRC" "$BIN_DIR/dikogit"

echo "✓ Linked binaries to $BIN_DIR"

if [[ ":$PATH:" != *":$BIN_DIR:"* ]]; then
    echo ""
    echo "Notice: $BIN_DIR is not in your \$PATH."
    echo "Add the following line to your ~/.zshrc or ~/.bash_profile:"
    echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
fi

echo ""
echo "Installation complete! Run:"
echo "  git dikogit generate"
