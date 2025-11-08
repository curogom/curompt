# MVP 목표 명확화

## 핵심 목표 재정의

### 주요 목표: 프롬프트 **수집** 및 **평가**

MVP의 핵심은:
1. **프롬프트 수집 (Collector)**
   - Claude Code, Codex CLI, Cursor IDE/CLI 등 외부 도구를 사용할 때
   - 실제로 전송된 프롬프트를 **자동으로 수집**
   - 수집된 프롬프트를 SQLite에 저장

2. **프롬프트 평가 (Analyzer + Scorer)**
   - 수집된 프롬프트를 분석
   - 정적 분석 (구조, 토큰, 규칙)
   - 점수화 (0-100점)
   - 리포트 생성

### ❌ 잘못된 이해 (수정 필요)

**Provider 직접 호출**: 
- 우리가 직접 Claude/OpenAI API를 호출하는 것이 **아님**
- 외부 도구가 이미 사용한 프롬프트를 **수집**하는 것

### ✅ 올바른 구조

```
외부 도구 (Claude Code/Codex/Cursor)
    ↓ (프롬프트 전송)
Collector (가로채기/캡처)
    ↓ (수집된 프롬프트)
Analyzer (정적 분석)
    ↓ (분석 결과)
Scorer (점수화)
    ↓ (점수)
Reporter (리포트 생성)
```

## Collector 구현 방법

### 1. CLI 래퍼 방식
```bash
# 사용자가 원래 실행하려던 명령
codex -C . exec "TASK: SSE 추가"

# Collector가 래핑
curompt wrap codex -C . exec "TASK: SSE 추가"
# → codex 실행하면서 프롬프트 캡처
# → SQLite에 저장
# → 분석 및 리포트 생성
```

### 2. 로그 파일 분석
- Claude Code, Cursor 등이 남긴 세션 로그 분석
- `.cursor/`, `.codex/` 등 디렉토리 스캔

### 3. 표준 입력/출력 캡처
- 파이프라인 통합
- `codex ... | curompt capture`

### 4. Git Hook
- 프롬프트 파일이 커밋될 때 자동 분석

## Provider의 역할 재정의

Provider는:
- ❌ **직접 LLM을 호출하는 용도가 아님**
- ✅ **수집된 프롬프트를 분석하기 위한 메타데이터 제공**
  - 토큰 계산 (어떤 Provider였는지에 따라 다른 토큰 계산)
  - 비용 추정
  - Provider별 특성 분석

## 수정된 우선순위

### MVP 핵심 (우선순위 높음)
1. ✅ **Collector 구현** ← **가장 중요!**
   - CLI 래퍼 (codex, cursor 명령 가로채기)
   - 로그 파일 파싱
   - 세션 캡처
2. ✅ Analyzer (구현 완료)
3. ✅ Scorer (구현 완료)
4. ✅ Reporter (구현 예정)

### 후순위
- Provider 직접 호출 (동적 평가는 M2에서)

---

**결론**: Collector 모듈이 MVP의 핵심이며, 현재 빠져있습니다. 이를 최우선으로 구현해야 합니다.

