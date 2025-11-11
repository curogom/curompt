---
sidebar_position: 1
title: Architecture
---

# Architecture

curompt is organized as a modular Go project so that each concern (collecting prompts, evaluating, reporting) can evolve independently.

```
curompt/
├── cmd/curompt           # Cobra entrypoint
├── internal/analyzer     # Static analysis (sections, tokens, lint rules)
├── internal/collector    # Log/import pipelines
├── internal/evaluator    # Provider orchestration & schema fit
├── internal/provider     # Claude / Gemini / OpenAI / local adapters
├── internal/reporter     # Markdown + structured reports
├── internal/repository   # SQLite persistence
└── internal/scorer       # Composite scoring engine
```

## Data flow

1. **Collector** ingests prompt sources (local files, repos, CLI pipes).
2. **Analyzer** runs static heuristics (section detection, duplication, forbidden terms).
3. **Evaluator** optionally calls LLM providers to generate samples, then validates them against schemas.
4. **Scorer** aggregates metrics into headline scores and flags regressions vs. previous runs.
5. **Reporter** renders TUI summaries plus Markdown/JSON outputs for CI or docs.

All results are persisted in SQLite (`~/.curompt/db.sqlite` by default) so historical comparisons and dashboards can be built externally.

## Providers

Each provider lives under `internal/provider/<name>` and implements a small interface:

```go
type Provider interface {
  Name() string
  Evaluate(ctx context.Context, req EvalRequest) (EvalResult, error)
}
```

Adding a new provider typically means:

1. Creating an adapter inside `internal/provider`.
2. Registering it with the CLI via `internal/cli`.
3. Adding config schema and documentation.

## Reports

Reports are generated in layers:

1. `internal/reporter` builds data models (scores, suggestions, token diff).
2. Renderers output Markdown (default), JSON, or future HTML dashboards.
3. CLI flags route the output—standard out for humans, files for CI artifacts.

This structure keeps the CLI thin and allows reuse from other entrypoints (e.g., API server or automation scripts) in the future.
