# Provider 지원 현황

## 지원 Provider 목록

### ✅ 구현 예정

1. **Claude (Anthropic)**
   - API: Anthropic API
   - 토큰 계산: Claude tokenizer
   - 상태: ✅ M1 MVP 완료

2. **OpenAI**
   - API: OpenAI API
   - 토큰 계산: tiktoken
   - 상태: ✅ M1 MVP 완료

3. **Gemini (Google)**
   - API: Google Gemini API
   - 토큰 계산: Gemini tokenizer
   - 상태: ✅ M1 MVP 완료

4. **Cursor IDE/CLI**
   - 통합 방식: Cursor CLI + MCP 프로토콜
   - 특징: Cursor IDE와의 통합 지원
   - 상태: ✅ M1 MVP 완료

### ⚠️ 후순위 (M2 이후)

5. **Codex CLI**
   - 통합: AGENTS.md/config.toml 인식
   - 상태: M2에서 검토

6. **Ollama (로컬)**
   - 통합: 로컬 Ollama 서버
   - 상태: M2에서 검토

7. **AWS Bedrock / Google Vertex**
   - 통합: 각 플랫폼 API
   - 상태: M2 이후 검토

## Provider 추가 방법

Strategy 패턴을 사용하여 새로운 Provider를 쉽게 추가할 수 있습니다:

```go
type Provider interface {
    Evaluate(ctx context.Context, prompt string) (*Response, error)
    CalculateTokens(text string) (int, error)
    Name() string
}
```

### Gemini API 통합

- API 엔드포인트: `https://generativelanguage.googleapis.com/v1beta/models/{model}:generateContent`
- 인증: API 키 기반
- 지원 모델: `gemini-pro`, `gemini-pro-vision` 등

### Cursor CLI 통합

- Cursor CLI는 `cursor` 명령어로 접근
- MCP (Model Context Protocol) 지원
- 세션 캡처 및 분석 가능

## 설정 예시

```yaml
provider:
  name: gemini
  api_key: ${GEMINI_API_KEY}  # 환경 변수 사용
  model: gemini-pro
  timeout_sec: 30

# 또는 Cursor 사용
provider:
  name: cursor
  cli_path: /usr/local/bin/cursor
  mcp_enabled: true
```

---

**참고**: 각 Provider의 API 키는 환경 변수나 안전한 설정 파일에 저장해야 합니다.

