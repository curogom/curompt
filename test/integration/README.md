# 통합 테스트

## 개요

통합 테스트는 전체 워크플로우와 CLI 명령어의 End-to-End 동작을 검증합니다.

## 테스트 구조

### 워크플로우 테스트 (`workflow_test.go`)

**목적**: 핵심 모듈 간 통합 및 데이터 흐름 검증

#### 테스트 목록
1. **TestWorkflow_CollectAnalyzeScoreReport**
   - 전체 워크플로우: 수집 → 저장 → 분석 → 점수화 → 리포트
   - 검증 항목:
     - SQLite 저장소 저장/조회
     - Evaluator 평가 결과
     - 점수 범위 (0-100)
     - 토큰 수 계산

2. **TestWorkflow_MultiplePrompts**
   - 여러 프롬프트 일괄 처리
   - 검증 항목:
     - 여러 프롬프트 저장/평가
     - 최근 프롬프트 조회

3. **TestWorkflow_ScanEvaluateReport**
   - scan 명령 워크플로우 시뮬레이션
   - 검증 항목:
     - 파일 읽기 → 수집 → 평가 → 리포트 생성

### CLI E2E 테스트 (`cli_e2e_test.go`)

**목적**: 실제 바이너리 실행을 통한 End-to-End 검증

#### 테스트 목록
1. **TestCLI_EvalCommand**
   - `eval` 명령 실행
   - 검증 항목:
     - 명령 성공
     - 출력에 점수 포함
     - 파싱된 섹션 포함

2. **TestCLI_ScanCommand**
   - `scan` 명령 실행
   - 검증 항목:
     - 명령 성공
     - 리포트 파일 생성
     - 분석 진행 메시지

3. **TestCLI_SuggestCommand**
   - `suggest` 명령 실행
   - 검증 항목:
     - 명령 성공
     - 개선 제안 출력
     - 점수 출력

4. **TestCLI_StdinInput**
   - 표준 입력 처리
   - 검증 항목:
     - stdin에서 프롬프트 읽기
     - 평가 결과 출력

5. **TestWorkflow_EvaluatorIntegration**
   - Evaluator 모듈 통합 테스트
   - 검증 항목:
     - 평가 결과 구조
     - 분석 결과 정확성

## 실행 방법

### 모든 통합 테스트 실행
```bash
go test ./test/integration/... -v
```

### 특정 테스트 실행
```bash
go test ./test/integration/... -v -run TestCLI_EvalCommand
```

### 커버리지 포함 실행
```bash
go test ./test/integration/... -v -coverprofile=coverage.out -coverpkg=./...
go tool cover -func=coverage.out
```

### 바이너리 필요 테스트 (CLI E2E)
CLI E2E 테스트는 `./bin/curo-prompt` 바이너리가 필요합니다.
바이너리가 없으면 해당 테스트는 자동으로 스킵됩니다.

```bash
# 빌드 먼저 수행
make build

# 그 다음 테스트 실행
go test ./test/integration/... -v
```

## 테스트 전제 조건

1. **바이너리 빌드**: CLI E2E 테스트는 `./bin/curo-prompt` 필요
2. **임시 디렉토리**: `t.TempDir()` 사용으로 자동 정리
3. **SQLite**: 임시 데이터베이스 파일 생성

## 현재 상태

- ✅ **8개 테스트 모두 통과**
- ✅ **평균 실행 시간**: ~0.6초
- ✅ **커버리지**: 62% (통합 테스트 범위)

## 추가 테스트 고려 사항

### 향후 추가 가능한 테스트
1. **에러 처리 테스트**
   - 잘못된 파일 경로
   - 잘못된 프롬프트 형식
   - 저장소 오류 처리

2. **성능 테스트**
   - 대용량 프롬프트 처리
   - 여러 프롬프트 동시 처리

3. **Provider별 테스트**
   - 다양한 Provider (Claude, OpenAI 등) 테스트
   - Provider별 토큰 계산 검증

4. **Reporter 테스트**
   - 다양한 리포트 포맷 생성
   - 리포트 파일 저장 검증

---

**최종 업데이트**: 2025-11-01

