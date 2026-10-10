#!/bin/bash
# Dikogits macOS Installer Script
set -e

BIN_DIR="${HOME}/.local/bin"
PROJECT_DIR="$(cd "$(dirname "$0")" && pwd)"
EXECUTABLE_SRC="${PROJECT_DIR}/bin/git-dikogits"

echo "=========================================="
echo "      Dikogits macOS Installer            "
echo "=========================================="

if [ ! -f "$EXECUTABLE_SRC" ]; then
    echo "Error: Executable $EXECUTABLE_SRC not found!"
    exit 1
fi

chmod +x "$EXECUTABLE_SRC"

mkdir -p "$BIN_DIR"
ln -sf "$EXECUTABLE_SRC" "$BIN_DIR/git-dikogits"
ln -sf "$EXECUTABLE_SRC" "$BIN_DIR/dikogits"

echo "✓ Installed binaries to $BIN_DIR"

# Check if BIN_DIR is in PATH
if [[ ":$PATH:" != *":$BIN_DIR:"* ]]; then
    echo ""
    echo "Notice: $BIN_DIR is not in your \$PATH."
    echo "Add the following line to your ~/.zshrc or ~/.bash_profile:"
    echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
fi

echo ""
echo "Installing Dikogits Git hooks for current repository..."
"$BIN_DIR/git-dikogits" install

echo ""
echo "Setup complete! You can now run:"
echo "  git dikogits generate"
echo "  git commit"
