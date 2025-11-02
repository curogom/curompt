# 프로젝트별 필터링 기능 설계

## 현재 상태

### 문제점
현재 `collect --from claude` 명령은:
- `~/.claude/history.jsonl` (전역 히스토리) 전체를 읽음
- 로컬에서 작업한 **모든 프로젝트**의 프롬프트를 수집
- 특정 레포/프로젝트만 필터링하는 기능 없음

### 확인된 데이터
- 각 항목에 `project` 필드 존재 (예: `/Users/curogom/dev/danggroom`)
- 여러 프로젝트의 프롬프트가 혼합되어 저장됨
- 프로젝트별 분리 필요

## 개선 방안

### 옵션 1: 필터링 옵션 추가 (권장)

```bash
# 특정 프로젝트만 수집
curo-prompt collect --from claude --project /Users/curogom/dev/danggroom

# 특정 레포만 수집 (경로 패턴)
curo-prompt collect --from claude --repo /Users/curogom/dev/*

# 제외 옵션
curo-prompt collect --from claude --exclude /Users/curogom/dev/devrock
```

### 옵션 2: 프로젝트별 히스토리 파일 지원

Claude Code는 프로젝트별 히스토리도 저장함:
- `~/.claude/projects/{project-path}/*.jsonl`
- 이 파일들을 직접 읽을 수도 있음

### 옵션 3: 수집 후 필터링

```bash
# 모든 프로젝트 수집 후
curo-prompt collect --from claude

# 특정 프로젝트만 조회
curo-prompt list --tool claude --project /Users/curogom/dev/danggroom
```

## 구현 계획

### 1단계: 필터링 옵션 추가
- `--project`: 특정 프로젝트 경로
- `--repo`: 레포 경로 패턴 (glob 지원)
- `--exclude`: 제외할 프로젝트 경로

### 2단계: 프로젝트별 히스토리 파일 지원
- 프로젝트별 `.jsonl` 파일 직접 읽기
- 세션 단위 수집

### 3단계: `list` 명령에 필터링 추가
- `--project`: 프로젝트별 조회
- `--repo`: 레포별 조회

## 구현 위치

1. `internal/cli/collect.go`: 필터 플래그 추가
2. `internal/collector/log_file_collector.go`: 필터링 로직 추가
3. `internal/cli/list.go`: 조회 필터 추가

