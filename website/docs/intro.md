---
sidebar_position: 1
title: Overview
---

# curompt

curompt is a CLI for analyzing, evaluating, and optimizing LLM prompts. It helps teams measure quality (schema-fit, self-consistency), detect regressions, and produce reproducible prompt reports directly from a terminal or CI workflow.

## Why curompt?

- **Local-first & reproducible** – runs without sending prompts to third-party services unless a provider is explicitly selected.
- **Contract-first** – validates prompts against JSON Schema / OpenAPI definitions to keep expectations explicit.
- **Scoring & suggestions** – generates 0–100 health scores with actionable recommendations (token trimming, rule separation, few-shot cleanup).
- **Report automation** – produces Markdown/HTML summaries that slot into pull requests and runbooks.

## Core capabilities

| Area | Highlights |
| --- | --- |
| Static analysis | Section heuristics, forbidden rules, duplicate detection, token budget |
| Dynamic evaluation | Multi-sample schema fit, self-consistency, latency, cost (per provider) |
| Scoring | Overall + sub-metrics with history persistence (SQLite) |
| Suggestions | Automatic wording, structure, and formatting improvements |
| Reporting | Rich TUI output or file export via `--output` / `--single-output` |

## Typical workflow

```bash npm2yarn
make build
./bin/curompt scan --path prompts/ --output reports
./bin/curompt eval --file prompts/onboarding.md --provider claude
./bin/curompt suggest --file prompts/onboarding.md > suggestions.md
```

1. **Scan** a repository or prompt folder to generate baseline metrics.
2. **Evaluate** critical prompts with an LLM provider (Claude, OpenAI, Gemini, local).
3. **Suggest** improvements and track them in git.

> **Scan note:** `curompt scan` uses prompts already saved in the local DB. When no history exists for the requested path, the CLI offers to auto-collect from Claude Code or Codex logs (Cursor support lands in v1.1).

## Next steps

- ➡️ [Installation](getting-started/installation.md)
- ⚙️ [Configuration reference](getting-started/configuration.md)
- 🧠 [Architecture & concepts](concepts/architecture.md)
