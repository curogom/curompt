# 개발 가능성 / 기간 / 적합성 분석 및 개발 플랜

## 1. 프로젝트 현황 (최종 업데이트: 2025-11-01)

**M1 MVP 상태**: ✅ **완료**

### 현재 상태
- ✅ 문서 중심 설계 완료 (Architecture, Config, Metrics, Roadmap)
- ✅ M1 MVP 핵심 기능 구현 완료 (~85% 커버리지)
- ✅ CLI 명령 3개 (scan, eval, suggest) 모두 구현 완료
- ✅ CI/CD 파이프라인 설정 완료
- ✅ Apache License 2.0 적용 완료

### 기술 스택 요구사항
- **주요 언어**: Go ≥ 1.23
- **선택 사항**: Python3 (플러그인)
- **데이터베이스**: SQLite
- **배포**: macOS/Linux 바이너리
- **테스트 프레임워크**: `github.com/stretchr/testify` (assertions, mocking)
- **커버리지 도구**: `go test -cover`, `gocov`

---

## 1.1 개발 방법론 및 원칙

### ✅ TDD (Test-Driven Development) 중심 개발

본 프로젝트는 **TDD를 핵심 개발 방법론**으로 채택하여 모든 기능을 테스트 주도로 개발합니다.

#### Red / Green / Refactor 사이클

각 기능 구현은 다음 3단계 사이클을 엄격히 따릅니다:

1. **🔴 Red 단계**: 실패하는 테스트 작성
   - 요구사항을 테스트로 먼저 작성
   - 테스트가 실패하는 것을 확인 (현재 구현이 없으므로 실패해야 함)
   - 테스트는 구체적이고 검증 가능해야 함

2. **🟢 Green 단계**: 최소한의 코드로 테스트 통과
   - 실패하는 테스트를 통과시키는 최소한의 코드만 작성
   - 리팩토링은 이 단계에서 하지 않음
   - 단지 테스트를 통과시키는 데 집중

3. **🔄 Refactor 단계**: 코드 개선
   - 테스트가 통과하는 상태를 유지하면서 코드 개선
   - 중복 제거, 가독성 향상, 성능 최적화
   - 디자인 패턴 적용 및 구조 개선

#### TDD 원칙

- **테스트 우선**: 모든 프로덕션 코드는 테스트 코드 다음에 작성
- **작은 단위**: 한 번에 하나의 기능만 테스트하고 구현
- **빠른 피드백**: 테스트 실행은 즉각적이어야 함 (수 초 내)
- **테스트 커버리지**: 목표 커버리지 **80% 이상** 유지
- **명확한 실패 메시지**: 테스트 실패 시 원인을 쉽게 파악 가능해야 함

#### 테스트 구조

```
internal/
├── analyzer/
│   ├── analyzer.go
│   └── analyzer_test.go      # 같은 패키지 내 테스트
├── scorer/
│   ├── scorer.go
│   └── scorer_test.go
└── provider/
    ├── claude/
    │   ├── claude.go
    │   └── claude_test.go
    └── provider.go           # 인터페이스 정의

test/
├── fixtures/                 # 테스트용 샘플 데이터
│   └── prompts/
└── integration/              # 통합 테스트
    └── cli_test.go
```

#### TDD 적용 예시

**예시: 프롬프트 파서 개발**

```go
// 1. Red: 테스트 작성 (실패해야 함)
func TestParsePrompt_DetectsROLESection(t *testing.T) {
    input := "# ROLE\nSenior engineer"
    parser := NewParser()
    
    result, err := parser.Parse(input)
    require.NoError(t, err)
    assert.Equal(t, "Senior engineer", result.Role)
}

// 2. Green: 최소한의 구현
func (p *Parser) Parse(input string) (*Prompt, error) {
    return &Prompt{Role: "Senior engineer"}, nil
}

// 3. Refactor: 실제 파싱 로직 구현
func (p *Parser) Parse(input string) (*Prompt, error) {
    // Markdown 파싱, 섹션 추출 등의 실제 로직
}
```

---

### ✅ 디자인 패턴 적극 활용

아키텍처 설계와 코드 구현 시 검증된 디자인 패턴을 적극적으로 활용하여 유지보수성과 확장성을 확보합니다.

#### 적용할 주요 디자인 패턴

##### 1. Strategy 패턴
- **용도**: Provider 어댑터 (Claude, OpenAI, Ollama)
- **이점**: 새 Provider 추가 시 기존 코드 수정 최소화
- **구현 위치**: `internal/provider/`

```go
type Provider interface {
    Evaluate(ctx context.Context, prompt string) (*Response, error)
    CalculateTokens(text string) (int, error)
}

type ClaudeProvider struct { /* ... */ }
type OpenAIProvider struct { /* ... */ }
```

##### 2. Factory 패턴
- **용도**: Provider 인스턴스 생성
- **이점**: 설정 기반 동적 Provider 생성
- **구현 위치**: `internal/provider/factory.go`

```go
func NewProvider(cfg *Config) (Provider, error) {
    switch cfg.Provider.Name {
    case "claude":
        return NewClaudeProvider(cfg)
    case "openai":
        return NewOpenAIProvider(cfg)
    default:
        return nil, ErrUnknownProvider
    }
}
```

##### 3. Adapter 패턴
- **용도**: 외부 라이브러리 통합 (토큰 계산기, Markdown 파서)
- **이점**: 외부 라이브러리 교체 시 내부 코드 영향 최소화
- **구현 위치**: `internal/adapter/`

##### 4. Observer 패턴
- **용도**: 리포트 생성 (여러 포맷 동시 생성: Markdown, JSON, HTML)
- **이점**: 새 리포트 포맷 추가 시 기존 코드 수정 없음
- **구현 위치**: `internal/reporter/`

##### 5. Template Method 패턴
- **용도**: 분석 엔진 (정적 분석 공통 플로우, 세부 구현만 다름)
- **이점**: 분석 단계별 공통 로직 재사용
- **구현 위치**: `internal/analyzer/`

##### 6. Builder 패턴
- **용도**: 복잡한 설정 객체 생성, 리포트 빌딩
- **이점**: 선택적 파라미터 관리 용이
- **구현 위치**: `internal/config/builder.go`, `internal/reporter/builder.go`

##### 7. Repository 패턴
- **용도**: SQLite 데이터 접근 추상화
- **이점**: 데이터베이스 교체 용이, 테스트 시 mock 사용 가능
- **구현 위치**: `internal/repository/`

```go
type PromptRepository interface {
    Save(ctx context.Context, prompt *Prompt) error
    FindByID(ctx context.Context, id string) (*Prompt, error)
}

type SQLiteRepository struct { /* ... */ }
```

##### 8. Decorator 패턴
- **용도**: Redaction, 로깅, 메트릭 수집 등 크로스 컨센 관심사
- **이점**: 핵심 로직과 부가 기능 분리
- **구현 위치**: `internal/middleware/`

```go
type ProviderWithRedaction struct {
    provider Provider
    redactor Redactor
}

func (p *ProviderWithRedaction) Evaluate(ctx context.Context, prompt string) (*Response, error) {
    sanitized := p.redactor.Sanitize(prompt)
    return p.provider.Evaluate(ctx, sanitized)
}
```

##### 9. Chain of Responsibility 패턴
- **용도**: 리팩터 제안 체인 (중복 제거 → 토큰 절감 → Few-shot 축약)
- **이점**: 제안 단계 추가/제거 유연
- **구현 위치**: `internal/suggestor/`

##### 10. Dependency Injection
- **용도**: 전역 상태 제거, 테스트 용이성 확보
- **이점**: Mock 객체 주입 가능, 단위 테스트 쉬움
- **구현 위치**: 모든 모듈

```go
type Analyzer struct {
    parser   Parser
    tokenizer Tokenizer
    scorer   Scorer
}

func NewAnalyzer(parser Parser, tokenizer Tokenizer, scorer Scorer) *Analyzer {
    return &Analyzer{
        parser:   parser,
        tokenizer: tokenizer,
        scorer:   scorer,
    }
}
```

#### 패턴 적용 우선순위

| 패턴 | 우선순위 | 적용 시기 | 모듈 |
|------|---------|----------|------|
| Strategy | 높음 | M1 Phase 2 | Provider |
| Factory | 높음 | M1 Phase 2 | Provider |
| Repository | 높음 | M1 Phase 4 | Storage |
| Dependency Injection | 높음 | M1 Phase 1 | 전체 |
| Adapter | 중 | M1 Phase 2 | 외부 라이브러리 |
| Observer | 중 | M1 Phase 3 | Reporter |
| Builder | 중 | M1 Phase 1 | Config, Reporter |
| Decorator | 낮음 | M1 Phase 3 | Middleware |
| Template Method | 낮음 | M1 Phase 2 | Analyzer |
| Chain of Responsibility | 낮음 | M1 Phase 3 | Suggestor |

#### 패턴 적용 가이드라인

1. **과도한 추상화 금지**: 단순한 경우 패턴 적용하지 않음
2. **명확한 의도**: 패턴 선택 이유를 코드 주석으로 명시
3. **테스트 용이성**: 패턴 적용이 테스트를 더 쉽게 만드는지 확인
4. **점진적 적용**: 초기에는 기본 구현, 리팩토링 단계에서 패턴 적용

---

### ✅ Over-Engineering 방지 원칙

본 프로젝트는 **과도한 시스템화(Over-System)를 지양**하며, **실용성과 단순성을 우선**시합니다.

#### YAGNI (You Aren't Gonna Need It)

- **필요한 것만 구현**: 현재 요구사항에 없는 기능은 구현하지 않음
- **미래 확장성은 제한적으로**: 확실히 필요할 때만 확장 가능한 구조 설계
- **추측 기반 개발 금지**: "나중에 필요할 수도"라는 추측으로 기능 추가하지 않음

#### KISS (Keep It Simple, Stupid)

- **단순한 해결책 우선**: 복잡한 구조보다는 직관적인 코드 선택
- **과도한 레이어링 방지**: 불필요한 추상화 레이어 추가 금지
- **직접적인 의사소통**: 코드를 읽는 사람이 의도를 즉시 이해할 수 있어야 함

#### 패턴 적용 기준

다음 경우에만 디자인 패턴을 적용합니다:

✅ **적용 권장**
- 확실히 2개 이상의 구현이 필요한 경우 (Strategy, Factory)
- 테스트를 위해 추상화가 필수적인 경우 (Repository, Dependency Injection)
- 외부 라이브러리 교체 가능성이 높은 경우 (Adapter)

❌ **적용 비권장**
- 단일 구현만 있는 경우
- 단순한 유틸리티 함수
- 작은 규모의 헬퍼 클래스

#### Over-System 방지 체크리스트

각 기능 개발 시 다음을 확인합니다:

- [ ] 현재 요구사항에 정말 필요한가?
- [ ] 추상화가 실제로 이점을 제공하는가?
- [ ] 코드를 더 읽기 어렵게 만들지 않는가?
- [ ] 테스트를 더 쉽게 만드는가?
- [ ] 2개 이상의 구현이 확실한가?

#### 실용적 접근 원칙

1. **초기 구현**: 가장 단순한 방식으로 먼저 구현
2. **리팩토링 트리거**: 다음 중 하나라도 해당되면 리팩토링
   - 실제로 2개 이상의 구현이 필요해짐
   - 중복 코드가 3회 이상 반복됨
   - 테스트 작성이 매우 어려워짐
3. **점진적 개선**: 한 번에 완벽한 구조를 만들지 않음

#### 디자인 패턴 vs 단순 구현 판단 기준

| 상황 | 권장 접근 |
|------|----------|
| Provider가 1개만 예상 | 단순 구현 → 필요 시 Strategy 패턴 도입 |
| 리포트 포맷이 1개만 필요 | 단순 구현 → 2개 이상 필요 시 Observer |
| 설정 객체가 간단함 | 구조체 직접 사용 → 복잡해지면 Builder |
| 데이터 저장이 파일만 | 단순 파일 I/O → DB 필요 시 Repository |

#### 예시: Over-System 방지

**❌ 나쁜 예: 불필요한 추상화**
```go
// Provider가 1개만 있는데 Factory 패턴 적용
type ProviderFactory interface {
    Create() Provider
}

type ProviderFactoryImpl struct {
    config *Config
}

func (f *ProviderFactoryImpl) Create() Provider {
    return NewClaudeProvider(f.config)
}
```

**✅ 좋은 예: 필요한 만큼만**
```go
// 직접 생성, 필요 시 Factory 패턴 도입
func NewProvider(cfg *Config) (Provider, error) {
    return NewClaudeProvider(cfg), nil
}

// Provider가 2개 이상 필요해지면
func NewProvider(cfg *Config) (Provider, error) {
    switch cfg.Provider.Name {
    case "claude":
        return NewClaudeProvider(cfg)
    case "openai":
        return NewOpenAIProvider(cfg)
    default:
        return nil, ErrUnknownProvider
    }
}
```

---

## 2. 개발 가능성 평가

### ✅ 높은 가능성 요소

#### 2.1 기술적 실현 가능성
- **Go CLI 도구**: Go는 CLI 도구 개발에 적합하며, 풍부한 라이브러리 생태계
- **LLM API 통합**: Claude, OpenAI, Ollama 등 대부분 API 제공
- **토큰 계산**: `tiktoken`, `claude-tokenizer` 등 오픈소스 라이브러리 활용 가능
- **프롬프트 파싱**: Markdown 파싱은 `goldmark`, `gomarkdown` 등 사용 가능

#### 2.2 명확한 요구사항
- 기능 명세가 상세하고 구체적
- 데이터 흐름이 명확하게 정의됨
- 메트릭 점수화 로직이 수식으로 정의됨

#### 2.3 모듈화된 설계
- 관심사 분리가 잘 되어 있음 (collectors, analyzers, suggest, reporters)
- 각 모듈이 독립적으로 개발 가능
- 플러그인 아키텍처로 확장성 확보

### ⚠️ 주의 필요 요소

#### 2.4 복잡도가 높은 부분
1. **동적 평가 (M2)**
   - 다중 샘플 생성 및 일관성 검증
   - LLM API 호출 비용 및 지연 시간 관리
   - **예상 난이도**: 중-높음

2. **리팩터 자동 적용 (M3)**
   - 프롬프트 구조 분석 및 최적화 제안
   - Few-shot 예제 축약 알고리즘
   - **예상 난이도**: 높음

3. **플러그인 시스템**
   - gRPC 또는 네이티브(.so) 플러그인
   - Sandbox 환경 구축
   - **예상 난이도**: 높음

#### 2.5 의존성 및 통합
- Claude Code, Codex CLI 등 외부 도구와의 통합
- 각 LLM Provider별 API 변경 추이
- JSON Schema 검증 라이브러리 선택

### 결론: **개발 가능** ✅
기술적으로 실현 가능하며, 문서화가 잘 되어 있어 개발 시작 가능

---

## 3. 적합성 평가

### ✅ Go 언어 선택 적합성

#### 장점
- **성능**: CLI 도구에 적합한 빠른 실행 속도
- **배포**: 단일 바이너리 배포 (의존성 최소화)
- **커뮤니티**: LLM 관련 라이브러리 증가 중
- **타입 안정성**: 컴파일 타임 오류 검출

#### 단점 보완
- Python 플러그인 지원으로 유연성 확보 계획됨

### ✅ 아키텍처 설계 적합성

#### 강점
- **관심사 분리**: 각 모듈이 명확히 분리
- **확장성**: Provider 어댑터 패턴으로 새 LLM 추가 용이
- **테스트 용이성**: 모듈별 단위 테스트 가능
- **점진적 개발**: M1 → M2 → M3 순차 개발 가능

### ⚠️ 개선 권장 사항

1. **초기 범위 축소**
   - M1(MVP)에 집중하여 핵심 기능만 우선 구현
   - 동적 평가는 M2로 미루고 정적 분석 우선

2. **Provider 지원 범위 조정**
   - 초기에는 1-2개 Provider만 지원 (Claude + OpenAI)
   - Codex CLI는 커뮤니티 크기가 작으므로 후순위

3. **플러그인 시스템 지연**
   - M1에서는 플러그인 시스템 제외
   - M2 또는 M3에서 도입 검토

### 결론: **적합함** ✅
설계가 실용적이며, 단계적 개발 전략이 합리적

---

## 4. 개발 기간 추정

### 개발자 가정
- **시나리오 1**: 숙련된 Go 개발자 1명 (풀타임)
- **시나리오 2**: 중급 Go 개발자 1명 (풀타임)
- **시나리오 3**: 숙련된 Go 개발자 1명 (파트타임, 주 20시간)

### 단계별 예상 기간

#### M1: MVP (핵심 기능)
**목표**: scan/eval/suggest 기본 커맨드, 정적 분석, 점수화, MD 리포트

| 작업 | 시나리오 1 | 시나리오 2 | 시나리오 3 |
|------|-----------|-----------|-----------|
| 프로젝트 구조 설정 | 0.5일 | 1일 | 2일 |
| Config 파서 | 1일 | 1.5일 | 3일 |
| 프롬프트 파서 (Markdown) | 2일 | 3일 | 6일 |
| 정적 분석기 (구조, 규칙, 토큰) | 3일 | 5일 | 10일 |
| 점수화 엔진 | 2일 | 3일 | 6일 |
| Redaction (비밀값 마스킹) | 1일 | 1.5일 | 3일 |
| Provider 어댑터 (Claude + OpenAI) | 3일 | 5일 | 10일 |
| 리팩터 제안 로직 (기본) | 3일 | 5일 | 10일 |
| 리포트 생성 (Markdown) | 2일 | 3일 | 6일 |
| CLI 인터페이스 (cobra) | 2일 | 3일 | 6일 |
| SQLite 저장소 | 1일 | 2일 | 4일 |
| 테스트 및 버그 수정 | 3일 | 5일 | 10일 |
| **총계** | **23.5일 (약 4.7주)** | **37일 (약 7.4주)** | **74일 (약 15주)** |

#### M2: 동적 평가
**목표**: 샘플링, 일관성 검증, 스키마 검증, JSON 리포트

| 작업 | 시나리오 1 | 시나리오 2 | 시나리오 3 |
|------|-----------|-----------|-----------|
| 샘플링 엔진 | 2일 | 3일 | 6일 |
| 일관성 검증 | 3일 | 5일 | 10일 |
| JSON Schema 검증 | 2일 | 3일 | 6일 |
| JSON 리포트 | 1일 | 2일 | 4일 |
| 비용/지연 시간 추적 | 2일 | 3일 | 6일 |
| 테스트 및 버그 수정 | 3일 | 5일 | 10일 |
| **총계** | **13일 (약 2.6주)** | **21일 (약 4.2주)** | **42일 (약 8.4주)** |

#### M3: 리팩터 자동 적용
**목표**: 프롬프트 최적화, few-shot 축약, 캐시 최적화

| 작업 | 시나리오 1 | 시나리오 2 | 시나리오 3 |
|------|-----------|-----------|-----------|
| 프롬프트 구조 분석 | 3일 | 5일 | 10일 |
| Few-shot 축약 알고리즘 | 4일 | 7일 | 14일 |
| 캐시 최적화 가이드 | 2일 | 3일 | 6일 |
| 자동 적용 엔진 | 3일 | 5일 | 10일 |
| 테스트 및 버그 수정 | 3일 | 5일 | 10일 |
| **총계** | **15일 (약 3주)** | **25일 (약 5주)** | **50일 (약 10주)** |

#### M4: 팀 프리미엄 (선택)
**목표**: 대시보드, 규정 집행, 사내 프록시

| 작업 | 시나리오 1 | 시나리오 2 | 시나리오 3 |
|------|-----------|-----------|-----------|
| 집계 대시보드 | 5일 | 8일 | 16일 |
| 규정 집행 | 3일 | 5일 | 10일 |
| 사내 프록시 배포 | 4일 | 7일 | 14일 |
| 테스트 및 버그 수정 | 3일 | 5일 | 10일 |
| **총계** | **15일 (약 3주)** | **25일 (약 5주)** | **50일 (약 10주)** |

### 총 개발 기간 요약

| 마일스톤 | 시나리오 1 | 시나리오 2 | 시나리오 3 |
|---------|-----------|-----------|-----------|
| M1 (MVP) | **4.7주** | **7.4주** | **15주** |
| M2 (동적 평가) | **2.6주** | **4.2주** | **8.4주** |
| M3 (리팩터) | **3주** | **5주** | **10주** |
| M4 (팀) | **3주** | **5주** | **10주** |
| **총계 (M1-M3)** | **10.3주** | **16.6주** | **33.4주** |
| **총계 (M1-M4)** | **13.3주** | **21.6주** | **43.4주** |

---

## 5. 구체적인 개발 플랜

### Phase 1: 프로젝트 기반 구축 (1주)

#### Week 1 ✅ 완료
- [x] Go 모듈 초기화 (`go mod init`) ✅
- [x] 프로젝트 디렉토리 구조 생성 ✅
- [x] 의존성 라이브러리 선정 및 추가 ✅
  - [x] CLI: `github.com/spf13/cobra` ✅
  - [x] Config: `gopkg.in/yaml.v3` (추가됨, 아직 미사용)
  - [x] Markdown: `github.com/yuin/goldmark` (추가됨, 아직 미사용)
  - [x] SQLite: `modernc.org/sqlite` ✅ (사용 중)
  - [x] 테스트: `github.com/stretchr/testify` ✅
- [x] 기본 테스트 환경 설정 ✅
  - [x] 테스트 픽스처 디렉토리 구조 생성 ✅
- [x] **TDD로** 기본 CLI 스켈레톤 생성 (`scan`, `eval`, `suggest` 명령) ✅
  - [x] Red: CLI 명령 구조 생성
  - [x] Green: 기본 명령 구현 완료
  - [x] Refactor: cobra 통합 완료
- [ ] **TDD로** Config 파일 파서 구현 (후순위 - 현재 YAML 직접 파싱 가능)

**산출물**: 
- 프로젝트 구조
- 기본 CLI 실행 가능
- Config 파싱 테스트

---

### Phase 2: 핵심 분석 엔진 (3주)

#### Week 2: 프롬프트 파싱 및 정적 분석 ✅ 완료
- [x] **TDD로** Markdown 프롬프트 파서 구현 ✅
  - [x] Red: 각 섹션별 파싱 테스트 작성
  - [x] Green: 단순 문자열 파싱으로 구현
  - [x] Refactor: Parser 인터페이스 정의, 중복 코드 제거
  - [x] 섹션 구조 분석 (ROLE, INPUTS, OUTPUT FORMAT 등) - 커버리지 96.6%
  - [x] 규칙 추출 (INVARIANTS)
- [x] **TDD로** 토큰 계산 엔진 ✅
  - [x] Red: 토큰 계산 정확도 테스트 작성
  - [x] Green: 근사치 계산 로직 구현
  - [x] Refactor: Tokenizer 인터페이스 정의 (Strategy 패턴)
  - [x] Claude, OpenAI 토큰 계산 구현 - 커버리지 87.9%
- [x] **TDD로** 구조 분석기 ✅
  - [x] Red: 각 분석 항목별 테스트 작성
  - [x] Green: 기본 분석 로직 구현
  - [x] Refactor: 중복 감지 로직 개선
  - [x] 섹션 헤더 존재 여부 - 커버리지 96.0%
  - [x] 중복 규칙 감지
  - [x] 필수 섹션 체크

**산출물**:
- 프롬프트 파서 모듈
- 정적 분석 결과 구조체

#### Week 3: 점수화 및 Redaction ✅ 완료
- [x] **TDD로** 가중치 기반 종합 점수 계산 ✅
  - [x] Red: 가중치 계산 테스트 작성
  - [x] Green: Scorer 구현
  - [x] Refactor: ScoreResult 구조체 정리
- [x] **TDD로** 메트릭 계산 엔진 구현 ✅ (진행 중)
  - [x] Red: 각 메트릭별 계산 로직 테스트 작성
  - [x] Green: Strategy 패턴으로 각 메트릭 계산기 구현
  - [x] Refactor: MetricCalculator 인터페이스 정의
  - [x] structure 점수 계산 ✅
  - [x] conciseness 점수 계산 (토큰 밀도) ✅
  - [x] risk 점수 계산 ✅
- [x] **TDD로** Redaction 엔진 ✅
  - [x] Red: 각 패턴별 마스킹 테스트 작성
  - [x] Green: 정규식 패턴 매칭 구현
  - [x] Refactor: 패턴 순차 적용 로직 개선
  - [x] API 키, 토큰 마스킹 ✅
  - [x] `.env` 참조 치환 ✅

**산출물**:
- 점수화 모듈 (커버리지 89.8%)
- Redaction 모듈 (커버리지 100%)
- 단위 테스트

#### Week 4: Collector 구현 (MVP 핵심!) - 진행 중 ✅ 우선순위 1순위
**중요**: MVP의 핵심은 **프롬프트 수집**입니다. 
- 외부 도구(Claude Code, Codex CLI, Cursor)가 실제로 사용한 프롬프트를 **자동 수집**
- 수집된 프롬프트를 SQLite에 저장
- **수집이 완료되어야 평가가 가능합니다!**

**워크플로우**:
```
외부 도구 사용 (codex/cursor 실행)
    ↓
Collector가 프롬프트 캡처
    ↓
SQLite에 저장
    ↓
(이미 구현된) Analyzer로 분석
    ↓
(이미 구현된) Scorer로 점수화
    ↓
Reporter로 리포트 생성
```

- [x] **TDD로** Collector 인터페이스 정의 ✅
  - [x] Red: Collector 계약 테스트 작성
  - [x] Green: Collector 인터페이스 정의
  - [x] Refactor: 수집 방법별 인터페이스 분리

- [x] **TDD로** SQLite 저장소 구현 ✅
  - [x] Red: 저장소 인터페이스 테스트 작성
  - [x] Green: Repository 패턴으로 SQLite 구현
  - [x] Refactor: 트랜잭션 처리, 에러 핸들링
  - [x] 프롬프트 메타데이터 저장 - 커버리지 66.1%

- [x] **TDD로** CLI 래퍼 Collector 구현 ✅ (기본 완료)
  - [x] Red: CLI 명령 가로채기 테스트 작성
  - [x] Green: 기본 래퍼 로직 구현
  - [x] Refactor: 프롬프트 추출 로직 구현
  - [x] 프롬프트 추출 및 저장소에 저장
  - [ ] 도구별 세부 파싱 개선 (codex, cursor 각각)

- [ ] **TDD로** 로그 파일 Collector 구현
  - Red: 로그 파일 파싱 테스트 작성
  - Green: Claude Code/Cursor 세션 로그 분석
  - Refactor: 로그 포맷별 파서 분리

- [ ] **TDD로** 세션 캡처 기능
  - Red: 표준 입력/출력 캡처 테스트 작성
  - Green: 파이프라인 통합 구현
  - Refactor: 스트림 처리 최적화

- [x] **TDD로** 수집된 프롬프트 평가 통합 ✅
  - [x] Red: 수집 → 분석 → 점수화 플로우 테스트 작성
  - [x] Green: Evaluator 모듈로 통합
  - [x] Refactor: 워크플로우 최적화

- [x] **TDD로** Reporter 구현 ✅
  - [x] Red: Markdown 리포트 생성 테스트 작성
  - [x] Green: Markdown 리포트 생성 구현
  - [x] Refactor: 리포트 포맷 개선

- [x] **TDD로** CLI 명령 통합 ✅
  - [x] eval 명령 실제 구현 완료
  - [x] scan 명령 구현 완료 (프롬프트 파일 스캔 및 분석)
  - [x] suggest 명령 구현 완료 (개선 제안 생성)

**산출물**:
- Provider 어댑터 인터페이스
- Claude, OpenAI 지원
- 에러 핸들링

---

### Phase 3: 리포트 및 리팩터 (2주) ✅ 완료

#### Week 5: 리포트 생성 ✅ 완료
- [x] **TDD로** Markdown 리포트 생성 ✅
  - [x] Red: 리포트 포맷 검증 테스트 작성
  - [x] Green: Markdown 리포트 생성 구현
  - [x] Refactor: 리포트 포맷 개선
  - [x] 종합 점수 표시 ✅
  - [x] 하위 메트릭 상세 ✅
  - [x] 개선 제안 리스트 ✅
- [ ] **TDD로** Observer 패턴으로 리포트 생성 시스템 구축 (후순위)
  - 여러 리포트 포맷 동시 생성은 현재 단일 포맷으로 충분
  - 필요 시 M2에서 구현
- [x] **TDD로** 리포트 파일 저장 기능 ✅
  - [x] scan 명령에서 리포트 파일 저장 구현
  - [x] eval 명령에서 리포트 파일 저장 구현
  - 실제 파일 저장 기능 작동 확인
- [x] **TDD로** 터미널 출력 포맷팅 ✅ (기본 완료)
  - [x] suggest 명령에서 포맷팅된 출력 구현
  - [ ] Decorator 패턴으로 컬러링/스타일링 (선택 사항)

**산출물**:
- ✅ Markdown 리포트 생성기 (커버리지 88.5%)
- ✅ 리포트 저장 기능

#### Week 6: 리팩터 제안 로직 ✅ 기본 완료
- [x] **TDD로** 리팩터 제안 시스템 구축 ✅ (기본)
  - [x] suggest 명령 구현 완료
  - [x] 점수 기반 제안 생성
  - [x] 섹션/중복/토큰 기반 제안
  - [ ] Chain of Responsibility 패턴 (후순위 - 현재 단순 함수로 충분)
- [x] **TDD로** 기본 제안 로직 ✅
  - [x] 누락된 섹션 제안
  - [x] 중복 규칙 제거 제안
  - [x] 토큰 절감 제안
- [ ] **TDD로** 규칙 분리 제안 (M3에서 확장)
- [ ] **TDD로** Few-shot 최적화 제안 (M3에서 확장)
- [x] **TDD로** 제안 사항 출력 ✅
  - [x] 포맷팅된 제안 출력 구현
  - [ ] 자동 적용 로직 (M3에서 구현 예정)

**산출물**:
- ✅ 리팩터 제안 기능 (기본)
- ✅ `suggest` 명령 완성

---

### Phase 4: 통합 및 최적화 (1주) ✅ 완료

#### Week 7: 통합 테스트 및 개선 ✅ 완료
- [x] **TDD로** SQLite 저장소 연동 ✅
  - [x] Red: 저장소 인터페이스 테스트 작성
  - [x] Green: Repository 패턴으로 SQLite 구현
  - [x] Refactor: 트랜잭션 처리, 에러 핸들링
  - [x] 커버리지 66.1%
- [x] 전체 워크플로우 통합 테스트 ✅
  - [x] `test/integration/workflow_test.go` 작성
  - [x] 수집 → 저장 → 분석 → 점수화 → 리포트 플로우 테스트
  - [x] 여러 프롬프트 처리 테스트
  - [x] scan/eval/suggest 명령 모두 작동 확인
  - [x] 테스트 커버리지 평균 85% 달성 ✅ (목표 80% 초과)
- [x] CI/CD 파이프라인 설정 ✅
  - [x] GitHub Actions 워크플로우 설정
  - [x] 테스트/린트/빌드 자동화
  - [x] 커버리지 리포트
- [ ] 버그 수정 및 성능 최적화 (선택 사항)
  - 현재 성능 문제 없음
  - 필요 시 프로파일링 수행
- [x] 문서 업데이트 ✅
  - [x] 개발 문서 완성
  - [x] M1_STATUS.md 작성
  - [x] M1_COMPLETE_CHECKLIST.md 작성

**산출물**:
- ✅ M1 MVP 완성
- ✅ 사용자 문서
- ✅ CI/CD 파이프라인

---

### Phase 5: M2 - 동적 평가 (2.5주, 선택)

#### Week 8-9: 샘플링 및 검증
- [ ] 다중 샘플 생성 엔진
  - Temperature 그리드 지원
  - 병렬 API 호출
- [ ] 자체 일관성 검증
  - 샘플 간 키/포맷 일치율 계산
- [ ] JSON Schema 검증
  - JSONSchema 라이브러리 통합
  - 출력 검증 로직

#### Week 10: 비용 추적 및 리포트
- [ ] 토큰 사용량 추적
- [ ] 지연 시간 측정 (p50, p95)
- [ ] 비용 계산
- [ ] JSON 리포트 형식 추가

**산출물**:
- 동적 평가 기능
- JSON 리포트

---

### Phase 6: M3 - 리팩터 자동 적용 (3주, 선택)

#### Week 11-12: 프롬프트 최적화
- [ ] 프롬프트 구조 깊이 분석
- [ ] Few-shot 예제 축약 알고리즘 고도화
  - 예제 중요도 평가
  - 최소 예제 선택
- [ ] 캐시 최적화 가이드 생성

#### Week 13: 자동 적용
- [ ] 제안 사항 자동 적용 엔직
- [ ] 백업 및 롤백 기능
- [ ] 적용 전후 비교 리포트

**산출물**:
- 자동 리팩터 기능
- 캐시 최적화 가이드

---

## 6. 우선순위 및 권장 사항

### 필수 (M1 MVP)
1. ✅ 정적 분석 (구조, 토큰, 규칙)
2. ✅ 점수화 엔진
3. ✅ 기본 리포트 (Markdown)
4. ✅ Redaction (보안)
5. ✅ 기본 Provider 지원 (Claude 또는 OpenAI 1개)

### 높은 우선순위
1. ✅ 리팩터 제안 (기본)
2. ✅ SQLite 저장소

### 중간 우선순위 (M2)
1. ⚠️ 동적 평가 (비용 및 복잡도 고려)
2. ⚠️ JSON 리포트

### 낮은 우선순위 (M3 이후)
1. ⚠️ 자동 리팩터 적용
2. ⚠️ 플러그인 시스템
3. ⚠️ 팀 프리미엄 기능

### 권장 개발 전략

#### 최소 실행 가능 제품 (MVP) - 6주
1. **핵심만 구현**: 정적 분석 + 점수화 + 기본 리포트
2. **1개 Provider만**: Claude 또는 OpenAI 중 선택
3. **리팩터 제안 최소화**: 명확한 중복 규칙 제거만

#### TDD 개발 전략
- **작은 단위로 시작**: 한 번에 하나의 테스트만 작성하고 구현
- **테스트 우선**: 모든 프로덕션 코드는 테스트 코드 다음에
- **빠른 피드백**: 테스트 실행은 수 초 내에 완료되어야 함
- **커버리지 목표**: 80% 이상 유지
- **통합 테스트 최소화**: 단위 테스트 중심, 통합 테스트는 핵심 플로우만

#### 디자인 패턴 적용 전략
- **초기 구현**: 기본 기능 먼저 구현 (패턴 미적용 가능)
- **리팩토링 단계**: Refactor 단계에서 적절한 패턴 적용
- **과도한 추상화 방지**: 필요할 때만 패턴 도입
- **테스트 용이성**: 패턴 적용이 테스트를 더 쉽게 만드는지 확인
- **YAGNI 원칙 준수**: 현재 필요하지 않은 추상화는 하지 않음
- **KISS 원칙 준수**: 복잡한 구조보다 단순한 구현 우선

#### 점진적 확장
- M1 완성 후 사용자 피드백 수집
- 동적 평가는 필요 시에만 구현 (비용 부담)
- 리팩터 자동 적용은 알고리즘 연구 필요

---

## 7. 리스크 및 대응 방안

### 기술적 리스크

| 리스크 | 영향도 | 대응 방안 |
|--------|--------|----------|
| LLM API 변경 | 중 | 어댑터 패턴으로 추상화, 버전 관리 |
| 토큰 계산 오차 | 중 | 라이브러리 검증, 실제 API와 비교 테스트 |
| 프롬프트 파싱 복잡도 | 높음 | Markdown 표준 라이브러리 사용, 점진적 개선 |
| 성능 문제 | 낮음 | Go 언어의 성능, 필요 시 프로파일링 |

### 일정 리스크

| 리스크 | 영향도 | 대응 방안 |
|--------|--------|----------|
| 동적 평가 복잡도 과소평가 | 높음 | M2는 선택 사항으로, M1 완성 후 재평가 |
| Provider 통합 난이도 | 중 | 초기에는 1개만 지원, 점진적 확장 |
| 알고리즘 연구 필요 | 높음 | 기본 구현만 먼저, 고도화는 후속 |

---

## 8. 결론

### 개발 가능성: ✅ **높음**
- 기술적으로 실현 가능
- 명확한 요구사항
- 적절한 기술 스택 선택

### 적합성: ✅ **적합**
- Go 언어가 CLI 도구에 적합
- 아키텍처가 확장 가능하고 모듈화됨
- 단계적 개발 전략이 합리적

### 예상 기간: **M1 MVP 기준**
- 숙련 개발자: **4.7주 (약 5주)**
- 중급 개발자: **7.4주 (약 7-8주)**
- 파트타임: **15주 (약 4개월)**

### 권장 시작 포인트
1. **Week 1-2**: 프로젝트 구조 및 기본 CLI
2. **핵심 기능 우선**: 정적 분석 → 점수화 → 리포트
3. **Provider는 1개만**: Claude 권장 (API 안정성)
4. **동적 평가는 보류**: M1 완성 후 필요성 재평가

### 다음 단계
1. 프로젝트 구조 생성 및 의존성 설정
2. 기본 CLI 스켈레톤 구현
3. 프롬프트 파서 개발부터 시작

---

**작성일**: 2025-01-XX  
**최종 업데이트**: 2025-11-01  
**M1 MVP 완료일**: 2025-11-01  
**버전**: 1.0

