# Git User Profile Switcher (git-usr)

Quickly switch between different git user profiles with a simple command. Built with Go for speed and portability.

## ✨ Features

- 🚀 Fast profile switching with one command
- 💾 Store unlimited profiles in a central location
- 🔧 Easy to add/update/remove profiles
- 🌍 Cross-platform (macOS, Linux, Windows)
- 📝 Local or global scope support
- 🎯 Works as a native git subcommand
- ⚡ Single binary - no dependencies needed
- 🔄 Shell completion for Bash, Zsh, Fish, PowerShell

## 📋 Requirements

- Go 1.21 or later (only for building)
- Git

## 🚀 Quick Start

### Installation

#### Option 1: Download Pre-built Binary (Recommended)

Download the latest release for your platform from the [Releases](../../releases) page.

**Linux (AMD64):**
```bash
curl -LO https://github.com/amantham20/git-usr/releases/latest/download/git-usr-linux-amd64
chmod +x git-usr-linux-amd64
sudo mv git-usr-linux-amd64 /usr/local/bin/git-usr
```

**macOS (Apple Silicon):**
```bash
curl -LO https://github.com/amantham20/git-usr/releases/latest/download/git-usr-darwin-arm64
chmod +x git-usr-darwin-arm64
sudo mv git-usr-darwin-arm64 /usr/local/bin/git-usr
```

**macOS (Intel):**
```bash
curl -LO https://github.com/amantham20/git-usr/releases/latest/download/git-usr-darwin-amd64
chmod +x git-usr-darwin-amd64
sudo mv git-usr-darwin-amd64 /usr/local/bin/git-usr
```

**Windows (PowerShell):**
```powershell
Invoke-WebRequest -Uri "https://github.com/amantham20/git-usr/releases/latest/download/git-usr-windows-amd64.exe" -OutFile "git-usr.exe"
# Move git-usr.exe to a directory in your PATH
```

#### Option 2: Build from Source

#### macOS/Linux
```bash
./install.sh
```

#### Windows
```cmd
install.bat
```

The installer will build the binary and add it to your PATH.

### Setting up Tab Completion

For the best experience, set up tab completion for your shell. This enables auto-completion for profile names, commands, and flags.

#### Bash
```bash
git-usr completion bash | sudo tee /etc/bash_completion.d/git-usr
# Then restart your shell or run: source ~/.bashrc
```

#### Zsh
```bash
mkdir -p ~/.zsh/completions
git-usr completion zsh > ~/.zsh/completions/_git-usr
# Add to ~/.zshrc:
echo 'fpath=(~/.zsh/completions $fpath)' >> ~/.zshrc
echo 'autoload -U compinit && compinit' >> ~/.zshrc
# Then restart your shell or run: source ~/.zshrc
```

#### Fish
```bash
git-usr completion fish > ~/.config/fish/completions/git-usr.fish
# Restart your shell or run: source ~/.config/fish/config.fish
```

#### PowerShell
```powershell
git-usr completion powershell > git-usr-completion.ps1
# Add to your PowerShell profile ($PROFILE):
. path\to\git-usr-completion.ps1
```

### Development

```bash
# Build
make build

# Run tests
make test

# Run unit tests only
make test-unit

# Run integration tests (requires git)
make test-integration

# Generate coverage report
make coverage

# Install locally
make install
```

### First Time Setup

1. List the default profiles:
```bash
git-usr list
```

2. Add your profiles:
```bash
git-usr add work "John Doe" "john@work.com"
git-usr add personal "John Doe" "john@personal.com"
git-usr add company "Jane Smith" "jane@company.com"
```

3. Switch profiles:
```bash
git-usr work      # Local scope (current repo only)
git-usr personal  # Local scope
```

## 📖 Usage

Once installed, use as a native git subcommand:

### Switch Profiles
```bash
git-usr work              # Switch to work profile (local scope)
git-usr personal --global # Switch to personal profile (global scope)
```

### Manage Profiles
```bash
git-usr list                                    # List all profiles
git-usr add work "Name" "email@example.com"    # Add/update a profile
git-usr add freelance                           # Add profile (interactive)
git-usr remove oldprofile                       # Remove a profile
git-usr current                                 # Show current git config
```

### Shell Completion

Generate completion scripts for your shell:

```bash
# Bash
git-usr completion bash | sudo tee /etc/bash_completion.d/git-usr

# Zsh
mkdir -p ~/.zsh/completions
git-usr completion zsh > ~/.zsh/completions/_git-usr
# Add to ~/.zshrc: fpath=(~/.zsh/completions $fpath)

# Fish
git-usr completion fish > ~/.config/fish/completions/git-usr.fish

# PowerShell
git-usr completion powershell > git-usr-completion.ps1
# Add to $PROFILE: . path\to\git-usr-completion.ps1
```

That's it! Git automatically recognizes `git-usr` as a subcommand, so you can use `git-usr` directly.

## 📁 Configuration

Profiles are stored centrally in a home directory location:
- **macOS/Linux**: `~/.config/git-usr/profiles.json`
- **Windows**: `%APPDATA%\git-usr\profiles.json`

This central location ensures your profiles are accessible from any repository on your system.

You can manually edit this file if needed:
```json
{
  "work": {
    "name": "John Doe",
    "email": "john@work.com"
  },
  "personal": {
    "name": "John Doe",
    "email": "john@personal.com"
  },
  "opensource": {
    "name": "John Doe",
    "email": "john@opensource.dev"
  }
}
```

## 🎯 Use Cases

### Scenario 1: Work on different projects
```bash
cd ~/work/company-project
git-usr work

cd ~/personal/side-project
git-usr personal
```

### Scenario 2: Contributing to open source
```bash
git-usr add opensource "Your Name" "you@opensource.dev"
cd ~/projects/react
git-usr opensource
```

### Scenario 3: Freelance work
```bash
git-usr add client1 "Your Name" "you@client1.com"
git-usr add client2 "Your Name" "you@client2.com"
```

## 🛠️ Advanced

### Local vs Global Scope

- **Local** (default): Changes only affect the current repository
- **Global**: Changes affect all repositories on your system

```bash
git-usr work              # Local - only this repo
git-usr work --global     # Global - all repos
```

### Interactive Profile Creation

Simply omit the name and email to be prompted:
```bash
git-usr add newprofile
# Enter name: John Doe
# Enter email: john@example.com
```

### Tab Completion

After installing shell completion, you can tab-complete:
- Profile names
- Commands (list, add, remove, current, etc.)
- Shell types for the completion command
- Flags like `--global`

## 🤝 Contributing

Feel free to submit issues or pull requests!

## 📄 License

MIT License - feel free to use this however you want!
