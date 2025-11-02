# Config

기본 경로: `curo-prompt.yaml`

```yaml
provider:
  name: claude   # 또는 codex|openai|ollama
  endpoint: auto
  model: auto
  timeout_sec: 30
  samples: 5
  temperature_grid: [0.2, 0.7]

redaction:
  enable: true
  patterns:
    - '(?i)(sk-)[a-z0-9]+'
    - '(?i)bearer [a-z0-9\-_\.]+'
    - '(?i)api[_-]?key\s*[:=]\s*[a-z0-9]+'

metrics:
  use_schema_validation: true
  jsonschema_paths:
    - prompts/output_schema.json
  use_llm_judge: false
  weights:
    structure: 0.2
    conciseness: 0.15
    instruction_following: 0.3
    self_consistency: 0.15
    latency_cost: 0.1
    risk: 0.1

reports:
  dir: reports
  formats: [markdown, json]

suggest:
  max_token_reduction_ratio: 0.5
  keep_examples: 3

integrations:
  claude:
    detect_output_style: true
  codex:
    detect_agents_md: true
```
