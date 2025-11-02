# Installation Guide

> 🇰🇷 **Korean**: [한국어 버전](./INSTALLATION.ko.md)

## Overview

This guide helps you install `curo-prompt` on your system so you can run the `curo-prompt` command from anywhere.

## Installation Methods

### Method 1: Homebrew (Recommended for macOS) ⭐

The easiest way to install `curo-prompt` on macOS. The tap repository is ready and available now!

#### Step 1: Add Tap

```bash
brew tap curogom/curo-prompt
```

#### Step 2: Install

```bash
brew install curo-prompt
```

#### Step 3: Verify Installation

```bash
# Check version
curo-prompt --version

# Check help
curo-prompt --help
```

#### Update

```bash
brew upgrade curo-prompt
```

#### Uninstall

```bash
brew uninstall curo-prompt
brew untap curogom/curo-prompt  # Optional: Remove tap
```

> **✅ Available now**: The Homebrew tap is ready. You can install `curo-prompt` using the commands above!

---

### Method 2: Go install (Recommended for Linux) ⭐

The simplest and standard method.

#### Step 1: Installation

```bash
# Method A: Using Makefile
git clone https://github.com/curogom/curo-prompt.git
cd curo-prompt
make install

# Method B: Direct go install
go install github.com/curogom/curo-prompt/cmd/curo-prompt@latest
```

Installation location:
- `$GOPATH/bin/curo-prompt` (default: `~/go/bin/curo-prompt`)

#### Step 2: PATH Verification

```bash
# Check Go bin path
go env GOPATH
# Example output: /Users/username/go

# Check if already in PATH
echo $PATH | grep -q "$(go env GOPATH)/bin" && echo "✅ Already in PATH" || echo "❌ Not in PATH"
```

#### Step 3: Add to PATH (if needed)

**macOS (zsh)**:
```bash
# Add to ~/.zshrc
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
source ~/.zshrc
```

**Linux (bash)**:
```bash
# Add to ~/.bashrc or ~/.bash_profile
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc
source ~/.bashrc
```

**fish**:
```fish
# Add to ~/.config/fish/config.fish
set -gx PATH $PATH (go env GOPATH)/bin
```

#### Step 4: Verification

```bash
# Check command
which curo-prompt
# Output: /Users/username/go/bin/curo-prompt

# Check version
curo-prompt --version

# Check help
curo-prompt --help
```

---

### Method 2: Local Build then Manual Installation

For development or when you want a specific version:

#### Step 1: Build

```bash
git clone https://github.com/curogom/curo-prompt.git
cd curo-prompt
make build
```

Binary location: `./bin/curo-prompt`

#### Step 2: Copy to System Path

**Option A: System-wide installation (requires sudo)**

```bash
# macOS/Linux
sudo cp ./bin/curo-prompt /usr/local/bin/
sudo chmod +x /usr/local/bin/curo-prompt

# Verify
/usr/local/bin/curo-prompt --version
```

**Option B: User-local installation**

```bash
# Create ~/bin directory
mkdir -p ~/bin

# Copy
cp ./bin/curo-prompt ~/bin/

# Add to PATH (in shell config file)
echo 'export PATH="$PATH:$HOME/bin"' >> ~/.zshrc  # or ~/.bashrc
source ~/.zshrc

# Verify
which curo-prompt
```

---

### Method 3: GitHub Releases (Future)

When official releases are available:

```bash
# Download (example)
wget https://github.com/curogom/curo-prompt/releases/download/v1.0.0/curo-prompt-linux-amd64

# Make executable
chmod +x curo-prompt-linux-amd64

# Move to system path
sudo mv curo-prompt-linux-amd64 /usr/local/bin/curo-prompt

# Verify
curo-prompt --version
```

---

## PATH Troubleshooting

### Issue: `command not found: curo-prompt`

#### Solution 1: PATH Verification

```bash
# Check current PATH
echo $PATH

# Check Go bin path
go env GOPATH
# Example output: /Users/username/go

# Test by manually adding path
export PATH="$PATH:$(go env GOPATH)/bin"
curo-prompt --help  # Should work now
```

#### Solution 2: Modify Shell Config File

**zsh users**:
```bash
# Check ~/.zshrc
cat ~/.zshrc | grep GOPATH

# Add if missing
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
source ~/.zshrc
```

**bash users**:
```bash
# Check ~/.bashrc or ~/.bash_profile
cat ~/.bashrc | grep GOPATH

# Add if missing
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc
source ~/.bashrc
```

#### Solution 3: Direct Path Specification

```bash
# Run with direct path
$(go env GOPATH)/bin/curo-prompt --help

# Or set alias
echo 'alias curo-prompt="$(go env GOPATH)/bin/curo-prompt"' >> ~/.zshrc
source ~/.zshrc
```

---

## Installation Location Verification

```bash
# If installed via go install
which curo-prompt
# Output: /Users/username/go/bin/curo-prompt

# Or
go env GOPATH
# Output: /Users/username/go
# Actual location: /Users/username/go/bin/curo-prompt
```

---

## Updating

### If Using Go install

```bash
# Update to latest version
go install github.com/curogom/curo-prompt/cmd/curo-prompt@latest

# Or update repository first
git pull
make install
```

### If Using Local Build

```bash
# Update repository
git pull

# Rebuild and install
make build
sudo cp ./bin/curo-prompt /usr/local/bin/  # or ~/bin/
```

---

## Uninstallation

### If Installed via Go install

```bash
rm $(go env GOPATH)/bin/curo-prompt
```

### If Manually Installed

```bash
# Remove from system path
sudo rm /usr/local/bin/curo-prompt

# Or remove from user local
rm ~/bin/curo-prompt
```

---

## Installation Verification Checklist

- [ ] `curo-prompt --version` command works
- [ ] `which curo-prompt` shows the path
- [ ] PATH contains the correct path
- [ ] Works in new terminals (verify after shell restart)

---

## Next Steps

Once installation is complete:

1. Follow [Getting Started Guide](./GETTING_STARTED.md) for hands-on testing
2. Configure settings using [Config Guide](./CONFIG.md)
3. Read [Architecture](./ARCHITECTURE.md) to understand the system better

---

**Issues?**: Please report via GitHub Issues or check the documentation.
