# Metrics

## 상위 지표
- structure: 섹션 헤더·중복 규칙·금지어
- conciseness: 토큰/키워드 밀도, 압축 이득(LLM 요약 대비)
- instruction_following: JSONSchema/정규식 체크, LLM 보조 채점(옵션)
- self_consistency: 다중 샘플 간 키/포맷 일치율
- latency_cost: 입력/출력 토큰, p50/p95 지연, 캐시 적합도
- risk: 모호 표현, 민감 데이터 노출 가능성

## 기본 가중치
```yaml
weights:
  structure: 0.2
  conciseness: 0.15
  instruction_following: 0.3
  self_consistency: 0.15
  latency_cost: 0.1
  risk: 0.1
```

## 점수 계산
```text
score = Σ(weight[k] * normalize(metric[k]))
```
