---
sidebar_position: 2
title: Configuration
---

# Configuration

curompt reads `curompt.yaml` (or the path passed via `--config`) to control providers, scoring thresholds, and report options.

## Example

```yaml
# ~/.config/curompt/curompt.yaml
storage:
  path: ~/.curompt/db.sqlite

providers:
  claude:
    api_key: $ANTHROPIC_API_KEY
    model: claude-3-7-sonnet
  openai:
    api_key: $OPENAI_API_KEY
    model: gpt-4.1-mini

analysis:
  max_tokens: 2000
  forbidden_terms:
    - "TBD"
    - "lorem ipsum"

evaluation:
  samples: 5
  schema: ./schemas/customer_onboarding.yaml
  temperature: 0.2

report:
  output: reports
  single_output: summary.md
  format: markdown
```

### Storage
- `storage.path`: SQLite database storing run history, scores, and prompt snapshots. Defaults to `~/.curompt/db.sqlite`.

### Providers
Configure any provider you need under `providers.<name>`. The CLI maps these to subcommands, so `curompt eval --provider claude` will use the section above.

### Analysis
- `max_tokens`: enforce prompt budget.
- `forbidden_terms`: list of phrases that should never appear (e.g., placeholders).

### Evaluation
- `samples`: number of generations per prompt.
- `schema`: JSON Schema / OpenAPI file that responses must satisfy.
- `temperature`: provider-specific generation temperature.

### Report
- `output`: directory for multi-file reports.
- `single_output`: optional all-in-one Markdown file.
- `format`: `markdown` or `json`.

## CLI overrides

Every config value can be overridden via flags:

```bash
curompt scan --config team/curompt.yaml --output reports/ --single-output all.md
curompt eval --provider claude --samples 10
```

Use `curompt config print` (coming soon) to render the resolved configuration for debugging.
