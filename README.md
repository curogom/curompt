# curompt

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
- **Report Output**: Terminal rich summary by default (<=100 lines). Optional file output via `--output` and `--single-output`.

## Documentation

The new Docusaurus site lives under [`website/`](website) and mirrors the English/Korean guides from `README*.md` and `docs/public/**`.

Run it locally:

```bash
cd website
npm install
npm run start
```

CI builds the site via [`.github/workflows/docs.yml`](.github/workflows/docs.yml). Adjust `docusaurus.config.ts` once you select the final production URL (e.g., GitHub Pages).

## Quick Start

### 1. Build

```bash
make build
./bin/curompt --help
```

### 2. Basic Testing (scan 중심)

```bash
# 기본: 프로젝트 스캔 → 배치 수집 → 병렬 평가 → 리치 요약(stdout)
./bin/curompt scan

# 하위 점수 Top-20만 보고 싶을 때
./bin/curompt scan --top 20

# 전체 리포트를 콘솔로 확인
./bin/curompt scan --full

# 파일 저장(단일 파일 병합)
./bin/curompt scan --output reports --single-output all_reports.md

# 병렬 작업 수 조정(예: 8)
./bin/curompt scan --concurrency 8
```

> **Note**: `scan` analyzes prompts already stored in the local DB (saved via `curompt collect …` or `curompt eval …`).  
> If the selected path has no history yet, the CLI will offer to auto-collect from Claude Code or Codex logs (Cursor support ships in v1.1).  
> Use `--path /absolute/project/path` to filter by project directory.

### 3. Quick Test Script

```bash
# Automated test script (build to test)
./test-quick.sh
```

This script automatically:
- ✅ Checks build
- ✅ Creates test prompt files
- ✅ Tests scan command and verifies rich summary
- ✅ Verifies optional report generation

> **📖 Detailed guide**: See [Getting Started Guide](./docs/public/GETTING_STARTED.md) (step-by-step instructions, examples, troubleshooting)

## Installation

### Requirements
- Go ≥ 1.23
- macOS/Linux
- Optional: Python3 (for plugins, future)

### Installation Methods

#### Method 1: Homebrew (Recommended for macOS) ⭐

The easiest way to install `curompt` on macOS:

```bash
# Add tap
brew tap curogom/curompt

# Install
brew install curompt

# Verify installation
curompt --version
```

**Update**:
```bash
brew upgrade curompt
```

> **✅ Available now**: Homebrew tap is ready! You can install `curompt` using the commands above.

#### Method 2: Go install

Install directly using Go modules to `$GOPATH/bin` or `~/go/bin`:

```bash
# Clone repository (or if already cloned)
git clone https://github.com/curogom/curompt.git
cd curompt

# Install
make install
# Or install directly
go install github.com/curogom/curompt/cmd/curompt@latest
```

**PATH verification and setup**:
```bash
# Check Go bin path
go env GOPATH
# Output example: /Users/username/go

# Add to PATH (may already be there)
export PATH="$PATH:$(go env GOPATH)/bin"

# Verify
which curompt
curompt --help
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
./bin/curompt --help

# Or manually add to PATH
sudo cp ./bin/curompt /usr/local/bin/
# Or
cp ./bin/curompt ~/bin/
export PATH="$PATH:~/bin"
```

### Installation Verification

```bash
# Check version
curompt --version

# Check help
curompt --help

# Basic test (rich summary)
curompt scan --top 5
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
- **Local storage**: SQLite storage only (`~/.curompt/db.sqlite`)

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
