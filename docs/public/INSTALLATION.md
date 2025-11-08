# Installation Guide

> 🇰🇷 **Korean**: [한국어 버전](./INSTALLATION.ko.md)

## Overview

This guide helps you install `curompt` on your system so you can run the `curompt` command from anywhere.

## Installation Methods

### Method 1: Homebrew (Recommended for macOS) ⭐

The easiest way to install `curompt` on macOS. The tap repository is ready and available now!

#### Step 1: Add Tap

```bash
brew tap curogom/curompt
```

#### Step 2: Install

```bash
brew install curompt
```

#### Step 3: Verify Installation

```bash
# Check version
curompt --version

# Check help
curompt --help
```

#### Update

```bash
brew upgrade curompt
```

#### Uninstall

```bash
brew uninstall curompt
brew untap curogom/curompt  # Optional: Remove tap
```

> **✅ Available now**: The Homebrew tap is ready. You can install `curompt` using the commands above!

---

### Method 2: Go install (Recommended for Linux) ⭐

The simplest and standard method.

#### Step 1: Installation

```bash
# Method A: Using Makefile
git clone https://github.com/curogom/curompt.git
cd curompt
make install

# Method B: Direct go install
go install github.com/curogom/curompt/cmd/curompt@latest
```

Installation location:
- `$GOPATH/bin/curompt` (default: `~/go/bin/curompt`)

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
which curompt
# Output: /Users/username/go/bin/curompt

# Check version
curompt --version

# Check help
curompt --help
```

---

### Method 2: Local Build then Manual Installation

For development or when you want a specific version:

#### Step 1: Build

```bash
git clone https://github.com/curogom/curompt.git
cd curompt
make build
```

Binary location: `./bin/curompt`

#### Step 2: Copy to System Path

**Option A: System-wide installation (requires sudo)**

```bash
# macOS/Linux
sudo cp ./bin/curompt /usr/local/bin/
sudo chmod +x /usr/local/bin/curompt

# Verify
/usr/local/bin/curompt --version
```

**Option B: User-local installation**

```bash
# Create ~/bin directory
mkdir -p ~/bin

# Copy
cp ./bin/curompt ~/bin/

# Add to PATH (in shell config file)
echo 'export PATH="$PATH:$HOME/bin"' >> ~/.zshrc  # or ~/.bashrc
source ~/.zshrc

# Verify
which curompt
```

---

### Method 3: GitHub Releases (Future)

When official releases are available:

```bash
# Download (example)
wget https://github.com/curogom/curompt/releases/download/v1.0.0/curompt-linux-amd64

# Make executable
chmod +x curompt-linux-amd64

# Move to system path
sudo mv curompt-linux-amd64 /usr/local/bin/curompt

# Verify
curompt --version
```

---

## PATH Troubleshooting

### Issue: `command not found: curompt`

#### Solution 1: PATH Verification

```bash
# Check current PATH
echo $PATH

# Check Go bin path
go env GOPATH
# Example output: /Users/username/go

# Test by manually adding path
export PATH="$PATH:$(go env GOPATH)/bin"
curompt --help  # Should work now
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
$(go env GOPATH)/bin/curompt --help

# Or set alias
echo 'alias curompt="$(go env GOPATH)/bin/curompt"' >> ~/.zshrc
source ~/.zshrc
```

---

## Installation Location Verification

```bash
# If installed via go install
which curompt
# Output: /Users/username/go/bin/curompt

# Or
go env GOPATH
# Output: /Users/username/go
# Actual location: /Users/username/go/bin/curompt
```

---

## Updating

### If Using Go install

```bash
# Update to latest version
go install github.com/curogom/curompt/cmd/curompt@latest

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
sudo cp ./bin/curompt /usr/local/bin/  # or ~/bin/
```

---

## Uninstallation

### If Installed via Go install

```bash
rm $(go env GOPATH)/bin/curompt
```

### If Manually Installed

```bash
# Remove from system path
sudo rm /usr/local/bin/curompt

# Or remove from user local
rm ~/bin/curompt
```

---

## Installation Verification Checklist

- [ ] `curompt --version` command works
- [ ] `which curompt` shows the path
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
