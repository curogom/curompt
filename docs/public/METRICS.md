# Metrics

> 🇰🇷 **Korean**: [한국어 버전](./METRICS.ko.md)

## Top-level Metrics
- **structure**: Section headers, duplicate rules, forbidden words
- **conciseness**: Token/keyword density, compression gain (compared to LLM summarization)
- **instruction_following**: JSONSchema/regex checks, LLM-assisted scoring (optional)
- **self_consistency**: Key/format match rate across multiple samples
- **latency_cost**: Input/output tokens, p50/p95 latency, cache suitability
- **risk**: Ambiguous expressions, sensitive data exposure potential

## Default Weights
```yaml
weights:
  structure: 0.2
  conciseness: 0.15
  instruction_following: 0.3
  self_consistency: 0.15
  latency_cost: 0.1
  risk: 0.1
```

## Score Calculation
```text
score = Σ(weight[k] * normalize(metric[k]))
```
