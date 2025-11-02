# 프롬프트 수집 테스트 가이드

> 🇬🇧 **English**: [English Version](./TESTING_COLLECTION.md)

## 개요

이 가이드는 Claude Code, Codex CLI, Cursor CLI 같은 외부 도구를 사용할 때 `curo-prompt`가 프롬프트를 성공적으로 수집하는지 테스트하는 방법을 안내합니다.

## 중요 사항

### Claude Code CLI

- **Claude Code**: CLI 도구 - `wrap` 명령으로 래핑 가능
- **Codex CLI / Cursor CLI**: 명령줄 도구 - `wrap` 명령으로 래핑 가능

## 테스트 방법

### 방법 1: 파일 스캔 (현재 작동 중) ✅

가장 간단하고 안정적인 방법:

```bash
# 1. 테스트 프롬프트 파일 생성
cat > test-prompt.md << 'EOF'
# ROLE
Senior Software Engineer

# INPUTS
- task: string
- context: string

# INVARIANTS
- Code must be testable
- Error handling required

# OUTPUT FORMAT
JSON format
EOF

# 2. 스캔 및 수집
curo-prompt scan --repo .

# 3. 수집 확인
curo-prompt list
```

**예상 결과:**
- 프롬프트가 데이터베이스에 저장됨
- `reports/` 디렉토리에 리포트 생성됨
- `list` 명령으로 프롬프트 확인 가능

### 방법 2: CLI 래퍼 (CLI 도구용) ⚠️

Codex CLI나 Cursor CLI 같은 명령줄 도구용:

```bash
# CLI 명령 래핑
curo-prompt wrap codex exec "TASK: 인증 기능 추가"

# 수집 확인
curo-prompt list --tool codex
```

**제한사항:**
- CLI 도구에만 작동, IDE에는 작동 안 함
- 도구별 프롬프트 추출 로직 필요
- 현재 기본 구현 상태

### 방법 3: Claude Code CLI 래퍼 ✅

Claude Code CLI용:

```bash
# Claude Code CLI 명령 래핑
curo-prompt wrap claude-code exec "TASK: 인증 기능 추가"

# 또는 다른 Claude Code 명령
curo-prompt wrap claude-code chat "로그인 구현"

# 수집 확인
curo-prompt list --tool claude-code
```

**작동 방식:**
1. 래핑된 명령 실행
2. 명령 인수에서 프롬프트 추출
3. 데이터베이스에 저장
4. 원래 명령 출력은 그대로 유지됨

## 단계별 테스트 절차

### 테스트 1: 기본 수집 테스트

```bash
# 1. 기존 데이터 삭제 (선택 사항)
rm -rf ~/.curo-prompt/db.sqlite

# 2. 테스트 프롬프트 생성
echo "# ROLE\nEngineer" > test.md

# 3. 스캔 및 수집
curo-prompt scan --repo .

# 4. 저장 확인
curo-prompt list

# 예상: 1개의 프롬프트가 표시되어야 함
```

### 테스트 2: 여러 프롬프트 수집

```bash
# 1. 여러 프롬프트 생성
mkdir -p test-prompts
echo "# ROLE\nEngineer1" > test-prompts/p1.md
echo "# ROLE\nEngineer2" > test-prompts/p2.md

# 2. 모두 스캔
curo-prompt scan --repo test-prompts

# 3. 개수 확인
curo-prompt list --limit 10

# 예상: 2개의 프롬프트가 표시되어야 함
```

### 테스트 3: 실제 사용 후 수집

```bash
# 1. Claude Code IDE를 정상적으로 사용
# (프롬프트 생성, 요청 전송 등)

# 2. Claude Code 사용 후 수집 시도:
# 옵션 A: 로그 파싱이 구현되어 있다면
curo-prompt collect --from claude-code

# 옵션 B: 프롬프트를 파일로 저장했다면
curo-prompt scan --repo ./my-prompts

# 3. 수집 확인
curo-prompt list
```

## 확인 체크리스트

수집 실행 후:

- [ ] 데이터베이스 파일 존재: `~/.curo-prompt/db.sqlite`
- [ ] `list` 명령이 수집된 프롬프트 표시
- [ ] 프롬프트에 올바른 도구 식별자 포함
- [ ] 타임스탬프 정확
- [ ] 리포트 생성됨 (scan 명령의 경우)

## 문제 해결

### 문제: 프롬프트가 수집되지 않음

**확인:**
1. 데이터베이스 존재: `ls -la ~/.curo-prompt/db.sqlite`
2. 도구 식별자 확인: `curo-prompt list --tool scan`
3. scan 명령의 파일 경로 확인

**해결:**
```bash
# 상세 출력과 함께 재스캔
curo-prompt scan --repo . --patterns "*.md"
curo-prompt list
```

### 문제: 래퍼가 프롬프트를 캡처하지 않음

**CLI 도구의 경우:**
- 명령 형식 확인
- 도구별 파서 필요 여부 확인
- 따옴표로 감싼 프롬프트로 시도: `curo-prompt wrap codex exec "PROMPT: ..."`

### 문제: Claude Code 프롬프트가 수집되지 않음

**현재 상태**: 아직 구현되지 않음

**임시 해결책:**
1. 프롬프트를 수동으로 파일로 저장
2. 해당 파일들에 `scan` 명령 사용

## 구현을 위한 다음 단계

Claude Code를 완전히 지원하려면:

1. **로그 파일 Collector**
   - 로그 파일 파서 구현
   - Claude Code 세션 로그 위치 찾기
   - 로그에서 프롬프트 추출

2. **세션 캡처**
   - IDE 프로세스에 후크
   - API 호출 모니터링
   - 실시간 프롬프트 캡처

3. **대안: 수동 내보내기**
   - Claude Code에서 프롬프트 내보내기
   - 파일로 저장
   - `scan` 명령 사용

---

**M1 MVP를 위해**: 파일 스캔(`scan` 명령)이 완전히 작동하며 프롬프트 수집 및 평가 테스트에 충분합니다.

