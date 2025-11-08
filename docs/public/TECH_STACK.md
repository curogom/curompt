# Technology Stack & Framework Decisions

> 🇰🇷 **Korean**: [한국어 버전](./TECH_STACK.ko.md)

## Finalized Technology Stack

### Core Language & Runtime

#### ✅ Go (Golang)
- **Version**: Go ≥ 1.23
- **Reasons**:
  - Optimized for CLI tool development
  - Single binary deployment (minimal dependencies)
  - Fast execution speed
  - Strong type system
  - Rich standard library
  - Cross-platform compilation (macOS/Linux)
- **Target Platforms**: macOS, Linux

#### Optional: Python3
- **Purpose**: Plugin system (after M2)
- **Reason**: Optional support for plugin extensibility

---

## Key Libraries & Frameworks

### CLI Framework

#### ✅ Cobra
- **Package**: `github.com/spf13/cobra`
- **Purpose**: CLI command structure (`scan`, `eval`, `suggest`)
- **Selection Reasons**:
  - Standard CLI framework in Go ecosystem
  - Supports subcommands, flags, automatic help generation
  - Integrates with Viper (configuration file management)

### Configuration File Processing

#### ✅ YAML Parser
- **Package**: `gopkg.in/yaml.v3`
- **Purpose**: Parsing `curompt.yaml` configuration file
- **Selection Reasons**:
  - Most widely used YAML library in Go
  - Stable and well-documented
  - Suitable for configuration file format

### Prompt Parsing

#### ✅ Goldmark
- **Package**: `github.com/yuin/goldmark`
- **Purpose**: Markdown prompt parsing
- **Selection Reasons**:
  - Fast Markdown parser written in Go
  - Extensible architecture
  - CommonMark compatible

**Alternative Considerations**:
- `gomarkdown`: More features but heavier and slower

### Database

#### ✅ SQLite
- **Package**: `modernc.org/sqlite` (pure Go implementation)
- **Purpose**: Store prompt analysis results (`~/.curompt/db.sqlite`)
- **Selection Reasons**:
  - File-based, no setup required
  - Pure Go implementation without CGO dependency
  - Aligns with local-only philosophy

**Alternative Considerations**:
- `github.com/mattn/go-sqlite3`: Requires CGO, faster but dependency complexity

### Testing Framework

#### ✅ Testify
- **Package**: `github.com/stretchr/testify`
- **Purpose**: Assertions, Mocking, test helpers
- **Selection Reasons**:
  - Most popular testing library in Go
  - Good assertion readability
  - Easy mock creation and usage
  - Supports Suite pattern

#### Coverage Tools
- **Standard**: `go test -cover`
- **Advanced**: `gocov` (optional)

### HTTP Client

#### ✅ Standard Library or HTTP Client Library
- **Option 1**: `net/http` (standard library)
- **Option 2**: `github.com/go-resty/resty` (optional)
- **Purpose**: LLM Provider API calls (Claude, OpenAI)
- **Selection Criteria**: 
  - Use standard library by default
  - Consider resty for complex request/response handling

### Token Calculation Libraries

#### Claude Token Calculation
- **Option 1**: `github.com/anthropics/anthropic-sdk-go` (official SDK)
- **Option 2**: `github.com/pkoukk/tiktoken-go` (compatible library)
- **Purpose**: Claude API token calculation
- **Selection**: Abstracted via Adapter pattern for easy replacement

#### OpenAI Token Calculation
- **Package**: `github.com/pkoukk/tiktoken-go`
- **Purpose**: OpenAI API token calculation
- **Selection Reason**: Go port of tiktoken Python implementation

### JSON Schema Validation (M2)

#### ✅ JSON Schema Validator
- **Option 1**: `github.com/xeipuuv/gojsonschema`
- **Option 2**: `github.com/santhosh-tekuri/jsonschema/v5`
- **Purpose**: Output schema validation (M2)
- **Selection Criteria**: Compare performance and features in M2 before deciding

---

## Development Tools

### Linter & Formatter

#### ✅ golangci-lint
- **Purpose**: Code quality checking
- **Configuration**: Manage rules via `.golangci.yml` file
- **Recommended Rules**:
  - `gofmt`, `govet`: Basic formatting
  - `golint`, `errcheck`: Code quality
  - `gosec`: Security checking

### Build Tools

#### Make
- **Purpose**: Build, test, release automation
- **File**: `Makefile`

---

## Project Structure

```
curompt/
├── cmd/
│   └── curompt/          # CLI entry point
│       └── main.go
├── internal/                 # Internal packages (no external access)
│   ├── analyzer/             # Static/dynamic analysis
│   ├── scorer/               # Scoring engine
│   ├── provider/            # LLM Provider adapters
│   │   ├── claude/
│   │   ├── openai/
│   │   └── provider.go       # Interface
│   ├── parser/               # Prompt parser
│   ├── reporter/             # Report generation
│   ├── suggestor/            # Refactoring suggestions
│   ├── redactor/             # Secret masking
│   ├── repository/           # Data storage
│   └── config/               # Configuration management
├── pkg/                      # Public packages (optional)
├── test/                     # Test fixtures
│   ├── fixtures/
│   └── integration/
├── docs/                     # Documentation
├── examples/                 # Examples
├── prompts/                  # Prompt templates
├── go.mod
├── go.sum
├── Makefile
├── .golangci.yml            # Linter configuration
└── README.md
```

---

## Dependency List (Expected go.mod)

```go
module github.com/curogom/curompt

go 1.23

require (
    // CLI
    github.com/spf13/cobra v1.8.0
    
    // Configuration
    gopkg.in/yaml.v3 v3.0.1
    
    // Markdown parsing
    github.com/yuin/goldmark v1.7.0
    
    // Database
    modernc.org/sqlite v1.29.0
    
    // Testing
    github.com/stretchr/testify v1.9.0
    
    // HTTP (if needed)
    github.com/go-resty/resty/v2 v2.12.0
    
    // JSON Schema (M2)
    github.com/xeipuuv/gojsonschema v1.2.0
    
    // Token calculation
    github.com/pkoukk/tiktoken-go v0.1.6
    github.com/anthropics/anthropic-sdk-go v1.0.0 // or latest version
)
```

---

## GUI / Web Interface

### ❌ GUI Not Supported (Current)

**Decision**: Develop as **CLI-only**

**Reasons**:
1. **CLI Tool Nature**: Optimized for terminal-based workflows
2. **Simplicity**: GUI development increases complexity and maintenance costs
3. **Easy Integration**: CLI can be directly integrated into pipelines, scripts, CI/CD
4. **YAGNI Principle**: Don't implement features not currently needed

**Output Methods**:
- Terminal output: Summary information and progress
- File reports: `reports/*.md`, `reports/*.json`
- Pipeline integration: Leverage standard input/output

### Future Considerations (M4 Team Premium)

**Web-based Dashboard** (optional):
- **Purpose**: Team-wide statistics and trend analysis
- **Technology Options** (if needed):
  - Simple HTTP server (`net/http`)
  - Static HTML + JavaScript chart library
  - Or separate as web service
- **Current Status**: Excluded from M1-M3, decide after re-evaluating necessity

**Desktop GUI**:
- **Decision**: Not implemented
- **Reason**: Can be replaced by web dashboard, high development complexity

---

## Technologies Not Selected (Reasons)

### ❌ Rust
- **Reason**: Go alone is fast enough, ecosystem is more mature
- **Additional Learning Curve**: Rust requires memory safety learning

### ❌ Python (Main Language)
- **Reason**: 
  - Increased deployment complexity (virtual environments, dependency management)
  - Slow execution speed
  - Difficult single binary deployment

### ❌ Node.js / TypeScript
- **Reason**: 
  - Requires runtime dependency
  - High memory usage
  - More suitable for web applications than CLI tools

### ❌ GUI Frameworks
- **Examples**: Fyne, Wails, Electron
- **Reason**:
  - Doesn't align with CLI tool nature
  - Increases development complexity
  - Complex deployment and dependency management
  - Terminal-based workflow is more efficient

### ❌ PostgreSQL / MySQL
- **Reason**: 
  - Requires setup and management
  - Doesn't align with local-only philosophy
  - SQLite is sufficient

---

## Technology Stack Verification Checklist

- [x] Suitable for CLI tool development? → Go + Cobra
- [x] Single binary deployment possible? → Go compilation
- [x] Cross-platform support possible? → Go cross-compilation
- [x] Easy to write tests? → Testify
- [x] Simple dependency management? → Go Modules
- [x] Performance sufficient? → Go is fast enough
- [x] Good community support? → Go ecosystem is active
- [x] Implementable without excessive abstraction? → Go is practical
- [x] Implementable without GUI? → CLI-only, file reports

---

## Final Recommended Technology Stack Summary

| Category | Selected Technology | Version/Package |
|----------|-------------------|----------------|
| **Language** | Go | ≥ 1.23 |
| **CLI Framework** | Cobra | `github.com/spf13/cobra` |
| **Config File** | YAML | `gopkg.in/yaml.v3` |
| **Markdown Parsing** | Goldmark | `github.com/yuin/goldmark` |
| **Database** | SQLite | `modernc.org/sqlite` |
| **Testing** | Testify | `github.com/stretchr/testify` |
| **HTTP Client** | net/http | Standard library |
| **Token Calculation** | tiktoken-go | `github.com/pkoukk/tiktoken-go` |
| **JSON Schema** | gojsonschema | `github.com/xeipuuv/gojsonschema` (M2) |
| **Linter** | golangci-lint | - |

---

**Final Confirmation**: Proceed with development using this technology stack.
