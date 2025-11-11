---
sidebar_position: 1
title: 개요
---

# curompt

curompt는 개발자가 CLI에서 프롬프트를 분석·평가·최적화할 수 있도록 설계된 도구입니다. JSON Schema/OpenAPI 기반 검증, 점수화, 리포트 자동화를 통해 팀이 프롬프트 품질을 꾸준히 측정할 수 있게 돕습니다.

## 핵심 가치

- **로컬 우선**: 기본적으로 로컬에서 실행되며, 명시적으로 선택한 공급자에만 요청을 보냅니다.
- **계약 우선**: JSON Schema/OpenAPI로 기대 결과를 선언하고, curompt가 일치 여부를 검사합니다.
- **점수 + 제안**: 0–100 점수와 함께 토큰 절감, 규칙 분리, few-shot 축약 등의 개선책을 제공합니다.
- **재현 가능한 리포트**: 터미널 요약 + Markdown/JSON 산출물을 CI나 문서화 파이프라인에 쉽게 넣을 수 있습니다.

## 대표 시나리오

```bash
make build
./bin/curompt scan --path prompts/ --output reports
./bin/curompt eval --provider claude --file prompts/onboarding.md
./bin/curompt suggest --file prompts/onboarding.md > suggestions.md
```

1. **scan**: 저장소 전체를 분석하여 기초 점수 확보
2. **eval**: 중요 프롬프트를 다중 샘플로 검증
3. **suggest**: 자동 리팩터링 제안을 수집

> **참고:** `curompt scan`은 DB에 저장된 프롬프트를 기반으로 하며, 지정한 경로에 데이터가 없으면 Claude Code 또는 Codex 로그를 자동 수집할지를 묻습니다. (Cursor 지원은 v1.1 예정)

## 다음 단계

- ➡️ [설치](getting-started/installation.md)
- ⚙️ [설정](getting-started/configuration.md)
- 🧠 [아키텍처](concepts/architecture.md)
