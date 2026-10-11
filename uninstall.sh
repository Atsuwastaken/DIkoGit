#!/bin/bash
# Dikogit Local Uninstaller Script
set -e

BIN_DIR="${HOME}/.local/bin"

echo "Uninstalling Dikogit from $BIN_DIR..."
for bin_name in git-dikogit dikogit git-dikogits dikogits; do
    if [ -f "$BIN_DIR/$bin_name" ]; then
        rm -f "$BIN_DIR/$bin_name"
    fi
done

# Clean up any leftover legacy hooks in current repository
if [ -d ".git/hooks" ]; then
    rm -f .git/hooks/pre-commit .git/hooks/prepare-commit-msg
fi

echo "✓ Dikogit uninstalled successfully."
