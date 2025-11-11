---
sidebar_position: 1
title: 아키텍처
---

# 아키텍처

curompt는 기능별 패키지로 나뉘어 있어 수집, 분석, 평가, 리포트 단계를 독립적으로 개선할 수 있습니다.

```
curompt/
├── cmd/curompt            # Cobra 진입점
├── internal/analyzer      # 정적 분석
├── internal/collector     # 로그/파일 수집
├── internal/evaluator     # LLM 공급자 호출
├── internal/provider      # Claude/OpenAI/Gemini/Local 어댑터
├── internal/reporter      # Markdown/JSON 렌더러
└── internal/scorer        # 점수 계산
```

## 데이터 흐름

1. **collector**가 프롬프트를 불러옵니다.
2. **analyzer**가 섹션 구조·금지어·토큰을 검사합니다.
3. **evaluator**가 선택한 LLM 공급자와 통신해 샘플을 만들고 스키마 적합성을 확인합니다.
4. **scorer**가 메트릭을 0–100 점수로 집계합니다.
5. **reporter**가 터미널/TUI 요약과 Markdown/JSON 파일을 생성합니다.

결과는 SQLite(`~/.curompt/db.sqlite`)에 저장되어 이후 비교나 대시보드에 활용할 수 있습니다.

## 공급자 확장

새로운 공급자를 추가하려면 `internal/provider/<name>`에 어댑터를 작성하고 CLI에 등록하면 됩니다. Provider 인터페이스는 간단합니다.

```go
type Provider interface {
  Name() string
  Evaluate(ctx context.Context, req EvalRequest) (EvalResult, error)
}
```

## 리포트 계층

1. 데이터 모델 생성 (`internal/reporter`)
2. Markdown/JSON 렌더링
3. CLI 플래그로 파일 경로나 stdout을 제어

덕분에 CI 스텝이나 외부 자동화에서 그대로 재사용할 수 있습니다.
