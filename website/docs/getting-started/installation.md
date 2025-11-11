---
sidebar_position: 1
title: Installation
---

# Installation

curompt ships as a Go CLI, distributed via Homebrew or manual builds. Every release includes binaries for macOS and Linux.

## Homebrew (recommended)

```bash
brew tap curogom/curompt
brew install curompt
curompt --version
```

To upgrade:

```bash
brew update
brew upgrade curompt
```

## From source

```bash
git clone https://github.com/curogom/curompt.git
cd curompt
make build
./bin/curompt --help
```

The `make build` target injects version/build metadata. If you prefer `go install`, run:

```bash
go install ./cmd/curompt
```

## Binary download

Each GitHub release attaches pre-built artifacts (`curompt-linux-amd64`, `curompt-darwin-amd64`). Download, mark executable, and place them on your `$PATH`.

```bash
curl -L -o curompt https://github.com/curogom/curompt/releases/download/v0.2.0/curompt-darwin-amd64
chmod +x curompt
mv curompt /usr/local/bin/
```

## Verify

```bash
curompt --version
curompt scan --help
```

If the command runs successfully, proceed to [Configuration](./configuration.md).
