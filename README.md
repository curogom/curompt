# curo-prompt

CLI tool for **analyzing, evaluating, and optimizing** LLM prompts used by developers.

- **Target**: Claude, OpenAI, Gemini, Cursor IDE/CLI, Codex CLI, Bedrock/Vertex/Local LLM CLI workflows
- **Goal**: Measure accuracy, reproducibility, and cost, with automatic refactoring suggestions
- **Philosophy**: **Local-first**, **Contract-first (OpenAPI/JSONSchema)**, **Reproducible reports**
- **License**: [Apache License 2.0](./LICENSE)

> 🇰🇷 **Korean Documentation**: [Korean README](./README.ko.md) | [Korean User Docs](./docs/public/README.ko.md)

## Features

- **Static Prompt Analysis**: Section structure, duplicate rules, forbidden words, schema presence
- **Dynamic Evaluation** (Optional): Multi-sample → Schema fit rate, self-consistency, latency, cost
- **Scoring**: 0–100 overall score + sub-metrics
- **Refactoring Suggestions**: Token reduction, rule separation, few-shot summarization, cache optimization
- **Report Output**: Terminal summary + `reports/*.md|json`

## Quick Start

### 1. Build

```bash
make build
./bin/curo-prompt --help
```

### 2. Basic Testing

```bash
# Evaluate sample prompt
echo "# ROLE\nEngineer\n\n# INPUTS\n- task: string" | ./bin/curo-prompt eval

# Evaluate from file
./bin/curo-prompt eval --file prompts/dev_contract_v2.md

# Batch scan
./bin/curo-prompt scan --repo prompts/

# Get improvement suggestions
./bin/curo-prompt suggest --file prompts/dev_contract_v2.md
```

### 3. Quick Test Script

```bash
# Automated test script (build to test)
./test-quick.sh
```

This script automatically:
- ✅ Checks build
- ✅ Creates test prompt files
- ✅ Tests eval, suggest, scan commands
- ✅ Verifies report generation

> **📖 Detailed guide**: See [Getting Started Guide](./docs/public/GETTING_STARTED.md) (step-by-step instructions, examples, troubleshooting)

## Installation

### Requirements
- Go ≥ 1.23
- macOS/Linux
- Optional: Python3 (for plugins, future)

### Installation Methods

#### Method 1: Homebrew (Recommended for macOS) ⭐

The easiest way to install `curo-prompt` on macOS:

```bash
# Add tap
brew tap curogom/curo-prompt

# Install
brew install curo-prompt

# Verify installation
curo-prompt --version
```

**Update**:
```bash
brew upgrade curo-prompt
```

> **✅ Available now**: Homebrew tap is ready! You can install `curo-prompt` using the commands above.

#### Method 2: Go install

Install directly using Go modules to `$GOPATH/bin` or `~/go/bin`:

```bash
# Clone repository (or if already cloned)
git clone https://github.com/curogom/curo-prompt.git
cd curo-prompt

# Install
make install
# Or install directly
go install github.com/curogom/curo-prompt/cmd/curo-prompt@latest
```

**PATH verification and setup**:
```bash
# Check Go bin path
go env GOPATH
# Output example: /Users/username/go

# Add to PATH (may already be there)
export PATH="$PATH:$(go env GOPATH)/bin"

# Verify
which curo-prompt
curo-prompt --help
```

**Permanently add to PATH** (add to shell config file):
```bash
# zsh users (~/.zshrc)
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
source ~/.zshrc

# bash users (~/.bashrc or ~/.bash_profile)
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc
source ~/.bashrc
```

#### Method 2: Local Build (Development)

```bash
# Build
make build

# Run (without PATH)
./bin/curo-prompt --help

# Or manually add to PATH
sudo cp ./bin/curo-prompt /usr/local/bin/
# Or
cp ./bin/curo-prompt ~/bin/
export PATH="$PATH:~/bin"
```

### Installation Verification

```bash
# Check version
curo-prompt --version

# Check help
curo-prompt --help

# Basic test
echo "# ROLE\nEngineer" | curo-prompt eval
```

> **📖 Detailed installation guide**: See [Installation Guide](./docs/public/INSTALLATION.md) (PATH setup, troubleshooting included)

### Development Build and Testing

```bash
# Build only
make build

# Run tests
make test

# Check test coverage
make coverage

# All checks (format, lint, test)
make check
```

## Integration Support

### LLM Providers
- **Claude**: Anthropic API integration (token counting, metadata)
- **OpenAI**: OpenAI API integration (planned)
- **Gemini**: Google Gemini API integration (planned)

### Tool Integration
- **Cursor IDE/CLI**: Prompt capture via CLI wrapper
- **Codex CLI**: Prompt capture via CLI wrapper
- **Bedrock/Vertex/Local**: Standard input/output wrapping support

> **Note**: In MVP, LLM Providers only provide metadata (token counting, cost estimation). Actual API calls will be implemented in M2.

## Security and Privacy

- **Local-first**: Runs locally by default, no external transmission
- **Auto masking**: Automatic masking of API keys/tokens/emails/URL queries/.env references
- **Local storage**: SQLite storage only (`~/.curo-prompt/db.sqlite`)

For details, see [SECURITY.md](./docs/public/SECURITY.md).

## Architecture

Modular architecture ensures scalability and maintainability:

- **Module separation**: Clear separation of responsibilities (Collector, Parser, Analyzer, Scorer, Provider, etc.)
- **Design patterns**: Strategy, Repository, Dependency Injection patterns applied
- **Test strategy**: Unit tests (85% coverage) + Integration tests

For detailed architecture, see [ARCHITECTURE.md](./docs/public/ARCHITECTURE.md).

## Testing

### Unit Tests
```bash
make test
```

### Integration Tests
```bash
go test ./test/integration/... -v
```

### Test Coverage
```bash
make coverage
```

**Current coverage**: Average 85% (Target 80% achieved)

## Documentation

### User Documentation (Public)
- **[Installation Guide](./docs/public/INSTALLATION.md)**: System installation and PATH setup ⭐
- **[Getting Started Guide](./docs/public/GETTING_STARTED.md)**: Hands-on testing guide ⭐
- [Architecture](./docs/public/ARCHITECTURE.md): Module structure and design patterns
- [Roadmap](./docs/public/ROADMAP.md): Milestones and feature plans
- [Tech Stack](./docs/public/TECH_STACK.md): Technologies used and selection rationale
- [Config Guide](./docs/public/CONFIG.md): YAML configuration file guide
- [Metrics](./docs/public/METRICS.md): Scoring metrics explanation
- [Provider Support](./docs/public/PROVIDERS.md): Supported LLM Provider list
- [Security](./docs/public/SECURITY.md): Security and privacy policy

> All user documentation is available in [`docs/public/`](./docs/public/).

## Contributing

This project is distributed under the Apache License 2.0. Contributions are welcome!

## License

Apache License 2.0

For details, see [LICENSE](./LICENSE).
