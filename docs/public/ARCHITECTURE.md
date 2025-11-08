# Architecture

> 🇰🇷 **Korean**: [한국어 버전](./ARCHITECTURE.ko.md)

## Overview

`curompt` is a CLI-based tool for analyzing, evaluating, and optimizing LLM prompts. It uses a modular architecture to ensure scalability and maintainability.

## Module Structure

### 1. Collectors (`internal/collector/`)

Module for collecting prompts from various sources.

#### Implementation
- **CLI Wrapper Collector**: Wraps external CLI tools (Codex, Cursor, etc.) to extract prompts
- **Log File Collector**: Parses prompts from log files (planned)
- **Session Capture**: Captures standard input/output (planned)

#### Interface
```go
type Collector interface {
    Collect(ctx context.Context) (*model.CollectedPrompt, error)
    Name() string
}
```

#### Data Flow
```
External Tool → CLI Wrapper → Prompt Extraction → Repository Storage
```

---

### 2. Parsers (`internal/parser/`)

Parses prompt text into structured data.

#### Features
- **Markdown Section Parsing**: Extracts ROLE, INPUTS, INVARIANTS, OUTPUT FORMAT
- **Token Calculation**: Calculates tokens per provider (Strategy pattern)

#### Strategy Pattern Applied
```go
type Tokenizer interface {
    CountTokens(text string) (int, error)
}

// Implementations
- ClaudeTokenizer: Token calculation for Claude models
- OpenAITokenizer: Token calculation for OpenAI models
```

#### Parsing Result Structure
```go
type Prompt struct {
    Role         string
    Inputs       []string
    Invariants   []string
    OutputFormat string
    Raw          string
}
```

---

### 3. Analyzers (`internal/analyzer/`)

Performs static analysis on parsed prompts.

#### Analysis Items
- **Section Presence**: Checks for required sections (ROLE, INPUTS, etc.)
- **Duplicate Rule Detection**: Detects duplicate rules in INVARIANTS
- **Section Count**: Counts items in each section
- **Missing Sections**: Identifies missing required sections

#### Analysis Result
```go
type AnalysisResult struct {
    HasRole         bool
    HasInputs       bool
    HasInvariants   bool
    HasOutputFormat bool
    DuplicateRules  []string
    SectionCounts   map[string]int
    MissingSections []string
}
```

---

### 4. Scorers (`internal/scorer/`)

Scores prompt quality.

#### Metric Calculators (Strategy Pattern)
- **StructureMetricCalculator**: Structure score (section presence, duplicate removal)
- **ConcisenessMetricCalculator**: Conciseness score (token density based)
- **RiskMetricCalculator**: Risk score (ambiguous expression detection)

#### Score Calculation
- Weighted overall score (0-100)
- Per-metric scores
- Custom weight support

```go
type ScoreResult struct {
    OverallScore float64
    Metrics     struct {
        Structure float64
        Conciseness float64
        Risk       float64
    }
}
```

---

### 5. Providers (`internal/provider/`)

Provides LLM Provider-specific metadata and functionality.

#### Strategy Pattern Applied
```go
type Provider interface {
    Evaluate(ctx context.Context, prompt string) (*Response, error)
    CalculateTokens(text string) (int, error)
    Name() string
}
```

#### Implementations
- **ClaudeProvider**: Anthropic Claude API support
- **OpenAIProvider**: OpenAI API support (planned)
- **GeminiProvider**: Google Gemini API support (planned)

#### Role (MVP)
- **Important**: In MVP, focuses on **metadata provision rather than direct LLM API calls**
- Token calculation, cost estimation, model information
- Actual evaluation will be implemented in M2

---

### 6. Evaluators (`internal/evaluator/`)

Orchestrates the entire evaluation workflow.

#### Workflow
```
CollectedPrompt → Parse → Analyze → Score → EvaluationResult
```

#### Role
- Combines Parser, Analyzer, Scorer, Provider
- Generates comprehensive evaluation results

```go
type EvaluationResult struct {
    Score       *scorer.ScoreResult
    Analysis    *analyzer.AnalysisResult
    TokenCount  int
    ParsedPrompt *parser.Prompt
}
```

---

### 7. Reporters (`internal/reporter/`)

Outputs evaluation results in various formats.

#### Implementation
- **MarkdownReporter**: Generates Markdown format reports
- JSON Reporter (M2 planned)
- HTML Reporter (M2 planned)

#### Report Content
- Overall score and per-metric scores
- Static analysis results
- Token information
- Improvement suggestions

---

### 8. Repository (`internal/repository/`)

Handles persistent storage of prompt data.

#### Repository Pattern Applied
```go
type PromptRepository interface {
    Save(ctx context.Context, prompt *model.CollectedPrompt) error
    FindByID(ctx context.Context, id string) (*model.CollectedPrompt, error)
    FindByTool(ctx context.Context, tool string) ([]*model.CollectedPrompt, error)
    FindRecent(ctx context.Context, limit int) ([]*model.CollectedPrompt, error)
    Close() error
}
```

#### Implementation
- **SQLiteRepository**: SQLite database implementation
- Storage location: `~/.curompt/db.sqlite`

---

### 9. Redactor (`internal/redactor/`)

Masks sensitive information in prompts.

#### Masking Targets
- API keys (format: `sk-*`, `api_key:*`)
- Bearer tokens
- Environment variable references (`.env`, `$VAR`)
- Email addresses
- URL query parameters

#### Usage Locations
- During report generation
- During log output
- Before external transmission (future)

---

### 10. CLI (`internal/cli/`)

Provides user interface.

#### Commands
1. **eval**: Single prompt evaluation and scoring
2. **scan**: Batch scan and analysis of prompt files in repository
3. **suggest**: Generate prompt improvement suggestions

#### Dependency Injection
- Combines Evaluator, Reporter, Repository at CLI level
- Ensures testability

---

## Data Flow

### Complete Workflow

```
1. Collection
   └─> External tool usage (Codex, Cursor, etc.)
       └─> Collector captures prompt
           └─> Store in Repository

2. Evaluation
   └─> CollectedPrompt
       └─> Parser: Text → Structured Prompt
           └─> Analyzer: Perform static analysis
               └─> Provider: Calculate tokens
                   └─> Scorer: Calculate score
                       └─> EvaluationResult

3. Reporting
   └─> EvaluationResult
       └─> Reporter: Format conversion
           └─> File or terminal output
```

### CLI Command Flow

#### `eval` Command
```
File/Stdin → Create CollectedPrompt → Evaluator.Evaluate() → Reporter.Generate() → Output
```

#### `scan` Command
```
Repository Scan → Collect File List → Evaluate Each File → Generate Reports → Store in Repository (optional)
```

#### `suggest` Command
```
File/Stdin → Evaluator.Evaluate() → Generate Suggestions Logic → Formatted Output
```

---

## Design Patterns Applied

### 1. Strategy Pattern ✅

**Location**: `internal/parser/tokenizer.go`, `internal/provider/types.go`

**Purpose**: 
- Separates token calculation strategy per provider
- Separates behavior per LLM provider

**Benefits**: 
- No code changes needed when adding new providers
- Mock implementations can be injected for testing

### 2. Repository Pattern ✅

**Location**: `internal/repository/`

**Purpose**: Data storage abstraction

**Benefits**:
- Easy database replacement (SQLite → PostgreSQL, etc.)
- In-Memory implementations can be used for testing

### 3. Dependency Injection ✅

**Location**: All modules

**Purpose**: Remove global state, improve testability

**Example**:
```go
// Bad: Global state
var parser = NewParser()

// Good: DI
func NewEvaluator(provider string) *Evaluator {
    return &Evaluator{
        parser: parser.NewParser(),
        provider: providerFactory.New(provider),
    }
}
```

### 4. Factory Pattern (Partial)

**Location**: `internal/evaluator/evaluator.go`

**Purpose**: Provider instance creation

**Current Implementation**:
```go
func NewEvaluator(provider string) *Evaluator {
    // Simple switch statement, can convert to Factory pattern if needed
}
```

---

## Test Strategy

### Unit Tests
- **Location**: `*_test.go` files in each module
- **Coverage Target**: 80% or above
- **Current Coverage**: Average 85%

### Integration Tests
- **Location**: `test/integration/`
- **Test Items**:
  - Complete workflow (collection → analysis → scoring → reporting)
  - Multiple prompt processing
  - CLI command E2E tests

---

## Extensibility

### Adding New Provider
1. Create `internal/provider/{provider}/` directory
2. Implement `Provider` interface
3. Register in `evaluator.NewEvaluator()`

### Adding New Report Format
1. Implement new Reporter in `internal/reporter/`
2. Implement `Reporter` interface
3. Add selection option in CLI

### Adding New Collector
1. Implement new Collector in `internal/collector/`
2. Implement `Collector` interface
3. Register in Collector factory

---

## Repository Structure

```
curompt/
├── cmd/curompt/       # CLI entry point
├── internal/
│   ├── cli/               # CLI command implementation
│   ├── collector/         # Prompt collection
│   ├── parser/            # Prompt parsing
│   ├── analyzer/          # Static analysis
│   ├── scorer/            # Scoring
│   ├── provider/          # LLM Provider
│   ├── evaluator/         # Evaluation orchestration
│   ├── reporter/          # Report generation
│   ├── repository/        # Data storage
│   ├── redactor/          # Sensitive information masking
│   └── model/             # Shared data models
├── test/
│   ├── fixtures/          # Test fixtures
│   └── integration/       # Integration tests
└── docs/                  # Documentation
```

---

## Security Considerations

### Local-First Design
- Runs locally by default
- External transmission disabled (default)

### Sensitive Information Masking
- Automatic masking via Redactor module
- Applied during report generation
- Applied during log output

### Data Storage
- SQLite local storage
- User home directory (`~/.curompt/`)
- No remote transmission

---

## Performance Optimization

### Current Status
- Single-threaded processing (sufficient for CLI tool characteristics)
- Minimized file I/O
- Memory-efficient structures

### Future Improvements
- Parallel processing (M2: during dynamic evaluation)
- Caching mechanism (M3)
- Indexing optimization (large-scale storage)

---

## License

Apache License 2.0

For details, see [LICENSE](../../LICENSE).
