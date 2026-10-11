#!/bin/bash
# Dikogits Uninstaller Script
set -e

BIN_DIR="${HOME}/.local/bin"

echo "Uninstalling Dikogits..."
if [ -f "$BIN_DIR/git-dikogit" ]; then
    rm -f "$BIN_DIR/git-dikogit"
fi
if [ -f "$BIN_DIR/dikogit" ]; then
    rm -f "$BIN_DIR/dikogit"
fi
if [ -f "$BIN_DIR/git-dikogits" ]; then
    rm -f "$BIN_DIR/git-dikogits"
fi
if [ -f "$BIN_DIR/dikogits" ]; then
    rm -f "$BIN_DIR/dikogits"
fi

if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    if command -v git-dikogit >/dev/null 2>&1; then
        git-dikogit uninstall || true
    fi
fi

echo "✓ Dikogits uninstalled successfully."
