# Architecture

## 개요

`curompt`는 CLI 기반 LLM 프롬프트 분석·평가·최적화 도구입니다. 모듈화된 아키텍처를 통해 확장성과 유지보수성을 확보했습니다.

## 모듈 구성

### 1. Collectors (`internal/collector/`)

프롬프트를 수집하는 모듈입니다. 다양한 소스에서 프롬프트를 캡처합니다.

#### 구현
- **CLI Wrapper Collector**: 외부 CLI 도구(Codex, Cursor 등)의 실행을 래핑하여 프롬프트 추출
- **Log File Collector**: 로그 파일에서 프롬프트 파싱 (예정)
- **Session Capture**: 표준 입력/출력 캡처 (예정)

#### 인터페이스
```go
type Collector interface {
    Collect(ctx context.Context) (*model.CollectedPrompt, error)
    Name() string
}
```

#### 데이터 흐름
```
외부 도구 실행 → CLI Wrapper → 프롬프트 추출 → Repository 저장
```

---

### 2. Parsers (`internal/parser/`)

프롬프트 텍스트를 구조화된 데이터로 파싱합니다.

#### 기능
- **Markdown 섹션 파싱**: ROLE, INPUTS, INVARIANTS, OUTPUT FORMAT 추출
- **토큰 계산**: Provider별 토큰 수 계산 (Strategy 패턴)

#### Strategy 패턴 적용
```go
type Tokenizer interface {
    CountTokens(text string) (int, error)
}

// 구현체
- ClaudeTokenizer: Claude 모델용 토큰 계산
- OpenAITokenizer: OpenAI 모델용 토큰 계산
```

#### 파싱 결과 구조
```go
type Prompt struct {
    Role         string
    Inputs       []string
    Invariants   []string
    OutputFormat string
    Raw          string
}
```

---

### 3. Analyzers (`internal/analyzer/`)

파싱된 프롬프트에 대한 정적 분석을 수행합니다.

#### 분석 항목
- **섹션 존재 여부**: 필수 섹션(ROLE, INPUTS 등) 존재 확인
- **중복 규칙 감지**: INVARIANTS 내 중복된 규칙 탐지
- **섹션 개수 계산**: 각 섹션의 항목 개수 계산
- **누락된 섹션**: 필수 섹션 누락 여부

#### 분석 결과
```go
type AnalysisResult struct {
    HasRole         bool
    HasInputs       bool
    HasInvariants   bool
    HasOutputFormat bool
    DuplicateRules  []string
    SectionCounts   map[string]int
    MissingSections []string
}
```

---

### 4. Scorers (`internal/scorer/`)

프롬프트 품질을 점수화합니다.

#### 메트릭 계산기 (Strategy 패턴)
- **StructureMetricCalculator**: 구조 점수 (섹션 존재, 중복 제거)
- **ConcisenessMetricCalculator**: 간결성 점수 (토큰 밀도 기반)
- **RiskMetricCalculator**: 위험 점수 (모호한 표현 감지)

#### 점수 계산
- 가중치 기반 종합 점수 (0-100)
- 각 메트릭별 점수 제공
- 사용자 정의 가중치 지원

```go
type ScoreResult struct {
    OverallScore float64
    Metrics     struct {
        Structure float64
        Conciseness float64
        Risk       float64
    }
}
```

---

### 5. Providers (`internal/provider/`)

LLM Provider별 메타데이터 및 기능을 제공합니다.

#### Strategy 패턴 적용
```go
type Provider interface {
    Evaluate(ctx context.Context, prompt string) (*Response, error)
    CalculateTokens(text string) (int, error)
    Name() string
}
```

#### 구현체
- **ClaudeProvider**: Anthropic Claude API 지원
- **OpenAIProvider**: OpenAI API 지원 (예정)
- **GeminiProvider**: Google Gemini API 지원 (예정)

#### 역할 (MVP)
- **중요**: MVP에서는 **직접 LLM API 호출이 아닌 메타데이터 제공**에 집중
- 토큰 계산, 비용 추정, 모델 정보 제공
- 실제 평가는 M2에서 구현

---

### 6. Evaluators (`internal/evaluator/`)

전체 평가 워크플로우를 조율합니다.

#### 워크플로우
```
CollectedPrompt → Parse → Analyze → Score → EvaluationResult
```

#### 역할
- Parser, Analyzer, Scorer, Provider를 조합
- 종합적인 평가 결과 생성

```go
type EvaluationResult struct {
    Score       *scorer.ScoreResult
    Analysis    *analyzer.AnalysisResult
    TokenCount  int
    ParsedPrompt *parser.Prompt
}
```

---

### 7. Reporters (`internal/reporter/`)

평가 결과를 다양한 포맷으로 출력합니다.

#### 구현
- **MarkdownReporter**: Markdown 포맷 리포트 생성
- JSON Reporter (M2 예정)
- HTML Reporter (M2 예정)

#### 리포트 내용
- 종합 점수 및 메트릭별 점수
- 정적 분석 결과
- 토큰 정보
- 개선 제안

---

### 8. Repository (`internal/repository/`)

프롬프트 데이터의 영구 저장을 담당합니다.

#### Repository 패턴 적용
```go
type PromptRepository interface {
    Save(ctx context.Context, prompt *model.CollectedPrompt) error
    FindByID(ctx context.Context, id string) (*model.CollectedPrompt, error)
    FindByTool(ctx context.Context, tool string) ([]*model.CollectedPrompt, error)
    FindRecent(ctx context.Context, limit int) ([]*model.CollectedPrompt, error)
    Close() error
}
```

#### 구현
- **SQLiteRepository**: SQLite 데이터베이스 구현
- 저장 위치: `~/.curompt/db.sqlite`

---

### 9. Redactor (`internal/redactor/`)

프롬프트 내 민감 정보를 마스킹합니다.

#### 마스킹 대상
- API 키 (형식: `sk-*`, `api_key:*`)
- Bearer 토큰
- 환경 변수 참조 (`.env`, `$VAR`)
- 이메일 주소
- URL 쿼리 파라미터

#### 사용 위치
- 리포트 생성 시
- 로그 출력 시
- 외부 전송 전 (향후)

---

### 10. CLI (`internal/cli/`)

사용자 인터페이스를 제공합니다.

#### 명령어
1. **eval**: 단일 프롬프트 평가 및 점수화
2. **scan**: 레포지토리 내 프롬프트 파일 일괄 스캔 및 분석
3. **suggest**: 프롬프트 개선 제안 생성

#### 의존성 주입
- Evaluator, Reporter, Repository를 CLI 레벨에서 조합
- 테스트 용이성 확보

---

## 데이터 흐름

### 전체 워크플로우

```
1. 수집 (Collection)
   └─> 외부 도구 사용 (Codex, Cursor 등)
       └─> Collector가 프롬프트 캡처
           └─> Repository에 저장

2. 평가 (Evaluation)
   └─> CollectedPrompt
       └─> Parser: 텍스트 → 구조화된 Prompt
           └─> Analyzer: 정적 분석 수행
               └─> Provider: 토큰 계산
                   └─> Scorer: 점수 계산
                       └─> EvaluationResult

3. 리포트 (Reporting)
   └─> EvaluationResult
       └─> Reporter: 포맷 변환
           └─> 파일 또는 터미널 출력
```

### CLI 명령별 흐름

#### `eval` 명령
```
파일/Stdin → CollectedPrompt 생성 → Evaluator.Evaluate() → Reporter.Generate() → 출력
```

#### `scan` 명령
```
레포지토리 스캔 → 파일 목록 수집 → 각 파일 평가 → 리포트 생성 → 저장소 저장 (선택)
```

#### `suggest` 명령
```
파일/Stdin → Evaluator.Evaluate() → 제안 생성 로직 → 포맷팅된 출력
```

---

## 디자인 패턴 적용

### 1. Strategy 패턴 ✅

**위치**: `internal/parser/tokenizer.go`, `internal/provider/types.go`

**용도**: 
- Provider별 토큰 계산 전략 분리
- LLM Provider별 동작 분리

**이점**: 
- 새 Provider 추가 시 기존 코드 수정 없음
- 테스트 시 Mock 구현체 주입 가능

### 2. Repository 패턴 ✅

**위치**: `internal/repository/`

**용도**: 데이터 저장소 추상화

**이점**:
- 데이터베이스 교체 용이 (SQLite → PostgreSQL 등)
- 테스트 시 In-Memory 구현체 사용 가능

### 3. Dependency Injection ✅

**위치**: 모든 모듈

**용도**: 전역 상태 제거, 테스트 용이성

**예시**:
```go
// Bad: 전역 상태
var parser = NewParser()

// Good: DI
func NewEvaluator(provider string) *Evaluator {
    return &Evaluator{
        parser: parser.NewParser(),
        provider: providerFactory.New(provider),
    }
}
```

### 4. Factory 패턴 (부분적)

**위치**: `internal/evaluator/evaluator.go`

**용도**: Provider 인스턴스 생성

**현재 구현**:
```go
func NewEvaluator(provider string) *Evaluator {
    // 단순 switch문, 필요 시 Factory 패턴으로 전환
}
```

---

## 테스트 전략

### 단위 테스트
- **위치**: 각 모듈의 `*_test.go` 파일
- **커버리지 목표**: 80% 이상
- **현재 커버리지**: 평균 85%

### 통합 테스트
- **위치**: `test/integration/`
- **테스트 항목**:
  - 전체 워크플로우 (수집 → 분석 → 점수화 → 리포트)
  - 여러 프롬프트 처리
  - CLI 명령 E2E 테스트

---

## 확장성

### 새 Provider 추가
1. `internal/provider/{provider}/` 디렉토리 생성
2. `Provider` 인터페이스 구현
3. `evaluator.NewEvaluator()`에 등록

### 새 리포트 포맷 추가
1. `internal/reporter/`에 새 Reporter 구현
2. `Reporter` 인터페이스 구현
3. CLI에서 선택 옵션 추가

### 새 Collector 추가
1. `internal/collector/`에 새 Collector 구현
2. `Collector` 인터페이스 구현
3. Collector 팩토리에 등록

---

## 저장소 구조

```
curompt/
├── cmd/curompt/       # CLI 진입점
├── internal/
│   ├── cli/               # CLI 명령어 구현
│   ├── collector/         # 프롬프트 수집
│   ├── parser/            # 프롬프트 파싱
│   ├── analyzer/          # 정적 분석
│   ├── scorer/            # 점수화
│   ├── provider/          # LLM Provider
│   ├── evaluator/         # 평가 조율
│   ├── reporter/          # 리포트 생성
│   ├── repository/        # 데이터 저장
│   ├── redactor/          # 민감 정보 마스킹
│   └── model/             # 공유 데이터 모델
├── test/
│   ├── fixtures/          # 테스트 픽스처
│   └── integration/       # 통합 테스트
└── docs/                  # 문서
```

---

## 보안 고려사항

### 로컬 우선 설계
- 기본적으로 로컬에서만 실행
- 외부 전송 기능 비활성화 (기본값)

### 민감 정보 마스킹
- Redactor 모듈로 자동 마스킹
- 리포트 생성 시 적용
- 로그 출력 시 적용

### 데이터 저장
- SQLite 로컬 저장
- 사용자 홈 디렉토리 (`~/.curompt/`)
- 원격 전송 없음

---

## 성능 최적화

### 현재 상태
- 단일 스레드 처리 (CLI 도구 특성상 충분)
- 파일 I/O 최소화
- 메모리 효율적 구조체 사용

### 향후 개선
- 병렬 처리 (M2: 동적 평가 시)
- 캐싱 메커니즘 (M3)
- 인덱싱 최적화 (대용량 저장소)

---

## 라이선스

Apache License 2.0

상세 내용은 [LICENSE](../../LICENSE) 참조.
