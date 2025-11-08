# 프롬프트 수집 방법 조사 결과

> **목적**: Claude Code, Codex, Cursor IDE에서 프롬프트를 자동으로 수집하는 방법 조사

## 조사 결과 요약

### ✅ 확실한 방법

1. **터미널 세션 로깅** (모든 CLI 도구)
   - `script` 명령어 활용
   - 세션의 모든 입력/출력 기록
   - 이후 로그 파일 파싱

2. **API 호출 로깅** (Codex 등 API 기반)
   - HTTP 프록시 또는 미들웨어
   - API 요청/응답 로깅

### ⚠️ 확인 필요

1. **Claude Code 세션 히스토리**
   - 로컬 저장 위치 불명확
   - Anthropic 계정과 동기화 가능성
   - 추가 조사 필요

2. **Cursor IDE 대화 히스토리**
   - IDE 내부 저장소 위치 불명확
   - MCP 프로토콜을 통한 접근 가능성
   - 추가 조사 필요

3. **Codex CLI**
   - 실제 CLI 도구 존재 여부 확인 필요
   - vs OpenAI Codex API

## 구현 가능한 방법

### Method 1: 터미널 세션 로깅 (우선 구현) ✅

```bash
# script 명령어로 세션 기록
script -a session.log
claude
# (대화 진행)
exit

# 로그 파일 파싱
curompt collect --from-log session.log --tool claude
```

**장점**: 모든 CLI 도구에 적용 가능  
**단점**: 사용자가 script 명령어 사용 필요

### Method 2: 로그 파일 모니터링 ✅

일반적인 로그 저장 위치 모니터링:

```
~/.claude/
~/.config/claude/
~/.cursor/
~/.config/cursor/
~/.codex/
```

**장점**: 자동 모니터링 가능  
**단점**: 로그 위치와 형식 확인 필요

### Method 3: 환경 변수/설정 파일 기반 ✅

각 도구의 설정 파일에서 히스토리 위치 파악:

```
~/.claude.json
~/.cursor/config.json
~/.codex/config.json
```

### Method 4: 실시간 프로세스 모니터링 (고급) ⚠️

프로세스의 stdin/stdout 실시간 캡처:

```go
// 프로세스 실행 및 I/O 모니터링
cmd := exec.Command("claude")
// stdin/stdout/stderr 캡처
```

**장점**: 완전 자동화  
**단점**: 복잡한 구현 필요

## 구현 우선순위

### Phase 1: 즉시 구현 가능 (M1)

1. **로그 파일 Collector** ✅
   - 일반적인 로그 위치 스캔
   - 파일 형식 파싱
   - 프롬프트 추출

2. **터미널 세션 로그 파서** ✅
   - `script` 명령어 출력 파싱
   - 입력/출력 구분
   - 프롬프트 추출

3. **설정 파일 기반 히스토리 위치 탐색** ✅
   - 각 도구 설정 파일 확인
   - 히스토리 경로 추출

### Phase 2: 추가 조사 필요 (M1.5)

1. **Claude Code 세션 히스토리**
   - Anthropic API 확인
   - 로컬 저장소 위치 조사

2. **Cursor IDE 히스토리**
   - IDE 내부 저장소 구조 확인
   - MCP 프로토콜 활용

3. **실시간 프로세스 모니터링**
   - ptrace 또는 유사 도구 활용
   - 복잡도 고려

## 불확실한 부분 (별도 AI Agent 리서치 필요)

1. **Claude Code 세션 히스토리 저장 위치**
   - Anthropic 공식 문서 확인 필요
   - 로컬 저장소 구조 분석

2. **Cursor IDE 대화 히스토리 접근 방법**
   - IDE 확장 API 확인
   - MCP 프로토콜 상세 사양

3. **Codex CLI 실제 존재 여부**
   - vs OpenAI Codex API
   - CLI 래퍼 도구 확인

## 권장 구현 계획

### M1에 포함할 내용:

1. ✅ **로그 파일 Collector**
   - 일반적인 위치 스캔
   - 기본 파싱 로직

2. ✅ **터미널 세션 로그 파서**
   - script 출력 파싱
   - 간단한 입력/출력 구분

3. ✅ **설정 파일 기반 탐색**
   - 일반적인 설정 위치 확인
   - 히스토리 경로 추출 시도

### M1.5 (추가 조사 후):

1. Claude Code 세션 히스토리 자동 수집
2. Cursor IDE 히스토리 통합
3. 실시간 프로세스 모니터링

---

**다음 단계**: 실제 시스템에서 로그 위치와 파일 형식 확인 후 구현 시작

