---
sidebar_position: 2
title: 설정
---

# 설정

curompt는 기본적으로 `curompt.yaml`을 읽어 공급자, 분석 옵션, 리포트 출력을 정의합니다.

## 예시

```yaml
storage:
  path: ~/.curompt/db.sqlite

providers:
  claude:
    api_key: $ANTHROPIC_API_KEY
    model: claude-3-7-sonnet

analysis:
  forbidden_terms:
    - "TBD"
evaluation:
  samples: 5
  schema: ./schemas/onboarding.yaml

report:
  output: reports
  single_output: summary.md
  format: markdown
```

### 주요 필드

- `storage.path`: 실행 이력을 저장할 SQLite 파일 경로
- `providers.*`: LLM 공급자별 API 키/모델 설정
- `analysis`: 토큰 한도, 금지 단어, 섹션 룰 등 정적 분석 옵션
- `evaluation`: 샘플 개수, 스키마 경로, 온도 등 동적 평가 파라미터
- `report`: Markdown/JSON 산출물의 경로와 파일명

## CLI 덮어쓰기

모든 설정값은 CLI 플래그로 덮어쓸 수 있습니다.

```bash
curompt scan --config ops/curompt.yaml --output reports
curompt eval --provider claude --samples 10
```
