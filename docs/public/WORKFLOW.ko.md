# 워크플로우 & 프롬프트 수집

> 🇬🇧 **English**: [English Version](./WORKFLOW.md)

## 개요

이 문서는 `curompt`가 프롬프트를 수집, 저장, 관리하는 방법을 설명합니다.

> **참고:** `curompt scan`은 DB에 저장된 프롬프트를 분석합니다.  
> 지정한 경로에 데이터가 없으면 CLI가 Claude Code 또는 Codex 로그를 자동 수집할지를 물어보고, Cursor 지원은 v1.1에서 추가됩니다.

## 현재 워크플로우

### 1. 초기 설정

처음 `curompt`를 실행하면 데이터베이스가 자동으로 초기화됩니다:

```bash
# 첫 실행 - ~/.curompt/db.sqlite에 DB 생성됨
curompt scan --path .
```

데이터베이스 위치: `~/.curompt/db.sqlite`

### 2. 프롬프트 수집 방법

#### 방법 1: 파일 스캔 (현재 - MVP)

레포지토리 내 기존 프롬프트 파일을 스캔:

```bash
# 현재 디렉토리에서 프롬프트 파일 스캔
curompt scan --path .

# 특정 디렉토리 스캔
curompt scan --path ./prompts

# 커스텀 파일 패턴
curompt scan --path . --patterns "*.md" --patterns "*.txt"
```

**작동 방식 (한 번에 모든 작업 수행):**
1. 패턴에 맞는 모든 프롬프트 파일 찾기
2. **수집**: 각 프롬프트를 CollectedPrompt로 생성
3. **평가**: 각 프롬프트 평가 (Evaluate 호출)
4. **저장**: 데이터베이스에 저장 (`~/.curompt/db.sqlite`)
5. **리포트 생성**: `reports/` 디렉토리에 리포트 생성

**참고**: `scan` 명령은 한 번에 모든 작업을 수행합니다: 수집 → 평가 → 저장 → 리포트 생성

#### 방법 2: 직접 평가 (DB에 저장 안 함)

저장하지 않고 단일 프롬프트 파일 평가:

```bash
# 저장 없이 평가
curompt eval --file prompt.md

# 평가 후 리포트 저장
curompt eval --file prompt.md --output reports/
```

**참고**: `eval` 명령은 현재 **데이터베이스에 저장하지 않습니다**. 저장하려면 `scan` 명령을 사용하세요.

#### 방법 3: 로그 파일 수집 ✅

도구의 히스토리/로그 파일에서 프롬프트 수집:

```bash
# Claude Code에서 수집 (현재 프로젝트만)
cd /path/to/project
curompt collect --from claude

# 모든 프로젝트의 프롬프트 수집
curompt collect --from claude --all

# Codex에서 수집 (현재 디렉토리를 프로젝트로 사용)
cd /path/to/project
curompt collect --from codex

# 모든 Codex 프롬프트 수집
curompt collect --from codex --all

# 수집 후 자동 평가
curompt collect --from claude --eval
```

**작동 방식:**
1. `~/.claude/history.jsonl` 또는 `~/.codex/history.jsonl` 읽기
2. 로그 항목에서 프롬프트 파싱
3. 프로젝트 정보 추출 (Codex는 session 파일에서)
4. 프로젝트별 필터링 (`--all` 미사용 시)
5. 데이터베이스에 자동 저장

**프로젝트 필터링:**

**Claude Code:**
- 프로젝트 루트에 `CLAUDE.md` 또는 `Claude.md` 파일 필요
- `--all` 없이: 현재 프로젝트만 수집
- `--all` 사용: 모든 프로젝트 수집

**Codex:**
- 현재 디렉토리를 프로젝트 경로로 사용 (`CLAUDE.md` 불필요)
- Session 파일의 `cwd` 필드와 매칭
- `--all` 없이: 현재 디렉토리만 수집
- `--all` 사용: 모든 프로젝트 수집

**예시:**

```bash
# Claude Code: 프로젝트별 수집
cd ~/projects/my-app  # CLAUDE.md 파일 필요
curompt collect --from claude
# → ~/projects/my-app 프로젝트 프롬프트만 수집

# Codex: 프로젝트별 수집
cd ~/projects/my-app  # CLAUDE.md 불필요
curompt collect --from codex
# → ~/projects/my-app 프로젝트 프롬프트만 수집

# 모든 프로젝트 수집
curompt collect --from claude --all
curompt collect --from codex --all
# → 히스토리의 모든 프로젝트 프롬프트 수집
```

#### 방법 4: CLI 래퍼 수집 ⚠️

외부 CLI 도구를 래핑하여 프롬프트 캡처:

```bash
# codex/cursor 명령에서 프롬프트 캡처
curompt wrap codex exec "TASK: 기능 추가"
curompt wrap cursor chat "로그인 구현"
```

**상태**: Collector 인프라는 있으나 도구별 파싱 구현이 필요함.

### 3. 저장된 프롬프트 조회

`list` 명령으로 저장된 프롬프트 조회:

```bash
# 최근 프롬프트 10개 조회
curompt list

# 더 많이 조회
curompt list --limit 20

# 특정 도구로 수집된 프롬프트만 조회
curompt list --tool scan
curompt list --tool codex

# 조회 후 재평가
curompt list --eval
```

**출력 예시:**
```
저장된 프롬프트: 5개

[1] ID: a1b2c3d4
    도구: scan
    시간: 2025-01-15 14:30:22
    명령: scan --path ./prompts
    경로: /Users/user/project
    프롬프트: # ROLE\nSenior Engineer\n\n# INPUTS\n- task: string...
```

### 4. 저장된 프롬프트 재평가

데이터베이스에 저장된 프롬프트를 재평가:

```bash
# 최근 프롬프트 모두 조회 후 평가
curompt list --eval --limit 10

# 특정 도구의 프롬프트 평가
curompt list --tool scan --eval --output reports/
```

## 데이터베이스 스키마

저장되는 정보:

- `id`: 고유 프롬프트 ID
- `tool`: 수집 도구 (scan, codex, cursor 등)
- `raw_prompt`: 원본 프롬프트 텍스트
- `role`, `inputs`, `invariants`, `output_format`: 파싱된 섹션
- `timestamp`: 수집 시간 (Unix timestamp)
- `command`: 수집에 사용된 명령어
- `working_dir`: 수집 시 작업 디렉토리
- `metadata`: 추가 메타데이터 (JSON)

## 현재 제한사항

### ❌ 아직 구현되지 않음

1. **자동 CLI 래퍼 수집**
   - Collector는 있으나 도구별 파싱이 미완성
   - codex, cursor 명령 파싱 구현 필요

2. **Eval 명령이 저장하지 않음**
   - `eval` 명령은 평가하지만 데이터베이스에 저장하지 않음
   - 저장하려면 `scan` 명령 사용

3. **Cursor IDE 수집**
   - Cursor 로그 파일 파싱 아직 구현되지 않음
   - 예정: Cursor 워크스페이스 로그 파싱

4. **세션 캡처**
   - 아직 구현되지 않음
   - 예정: CLI 도구의 stdin/stdout 캡처

### ✅ 현재 작동 중

1. **파일 스캔 & 저장**
   - `scan` 명령이 파일을 찾아 평가하고 데이터베이스에 저장
   
2. **로그 파일 수집**
   - `collect` 명령이 Claude Code와 Codex의 히스토리 파일 파싱
   - 프로젝트별 필터링 지원
   - Session 파일에서 자동 프로젝트 감지 (Codex)

3. **프롬프트 목록**
   - `list` 명령으로 저장된 프롬프트 표시
   - 도구별 필터링, 결과 제한

4. **재평가**
   - `list --eval`로 저장된 프롬프트 재평가

## 권장 워크플로우

### 신규 사용자

1. **초기 설정**
   ```bash
   # 기존 프롬프트 파일 스캔
   curompt scan --path ./prompts
   ```

2. **수집된 프롬프트 확인**
   ```bash
   # 수집된 내용 확인
   curompt list
   ```

3. **리포트 검토**
   ```bash
   # 생성된 리포트 확인
   ls reports/
   cat reports/your_prompt_report.md
   ```

### 일반적인 사용

1. **새 프롬프트 작성 후**
   ```bash
   # 새 프롬프트 스캔
   curompt scan --path ./prompts
   ```

2. **빠른 평가 (저장 없이)**
   ```bash
   # 저장 없이 평가만
   curompt eval --file new_prompt.md
   ```

3. **저장된 프롬프트 재평가**
   ```bash
   # 최근 프롬프트 모두 재평가
   curompt list --eval
   ```

## 향후 개선 사항 (M2+)

- CLI 래퍼를 통한 자동 수집
- 로그 파일 파싱
- 세션 캡처
- Git hook 통합
- 저장된 프롬프트 배치 작업

---

**질문?** [시작하기 가이드](./GETTING_STARTED.md) 또는 [아키텍처](./ARCHITECTURE.md)를 확인하세요.
