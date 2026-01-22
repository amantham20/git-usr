#!/bin/bash
# Installation script for git-usr (Unix/macOS/Linux)

echo "🔧 Installing git-usr..."

# Get the directory where the script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go from https://golang.org/dl/"
    exit 1
fi

# Build the binary
echo "Building git-usr..."
cd "$SCRIPT_DIR"
go build -o git-usr main.go

if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi

# Determine installation directory
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    INSTALL_DIR="/usr/local/bin"
else
    # Linux
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"
fi

# Copy binary to installation directory
cp "$SCRIPT_DIR/git-usr" "$INSTALL_DIR/git-usr"
chmod +x "$INSTALL_DIR/git-usr"

echo "✅ git-usr installed successfully!"
echo ""
echo "Usage (as a git subcommand):"
echo "  git usr work              # Switch to work profile"
echo "  git usr personal          # Switch to personal profile"
echo "  git usr list              # List all profiles"
echo "  git usr add <profile>     # Add a new profile"
echo ""
echo "First time setup:"
echo "  1. Run 'git usr list' to see default profiles"
echo "  2. Update profiles with: git usr add work \"Your Name\" \"email@example.com\""
echo "  3. Switch profiles with: git usr work"
echo ""
echo "Shell completion (optional):"
echo "  git usr completion bash > /etc/bash_completion.d/git-usr  # Bash"
echo "  git usr completion zsh > ~/.zsh/completions/_git-usr      # Zsh"
echo "  git usr completion fish > ~/.config/fish/completions/git-usr.fish  # Fish"
echo "  git usr completion powershell > git-usr-completion.ps1    # PowerShell"
echo ""
echo "Config stored in: ~/.config/git-usr/profiles.json"
