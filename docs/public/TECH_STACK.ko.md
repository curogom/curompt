# 기술 스택 및 프레임워크 결정

## 최종 결정된 기술 스택

### 핵심 언어 및 런타임

#### ✅ Go (Golang)
- **버전**: Go ≥ 1.23
- **이유**:
  - CLI 도구 개발에 최적화
  - 단일 바이너리 배포 (의존성 최소화)
  - 빠른 실행 속도
  - 강력한 타입 시스템
  - 풍부한 표준 라이브러리
  - 크로스 플랫폼 컴파일 (macOS/Linux)
- **대상 플랫폼**: macOS, Linux

#### 선택 사항: Python3
- **용도**: 플러그인 시스템 (M2 이후)
- **이유**: 플러그인 확장성을 위한 선택적 지원

---

## 주요 라이브러리 및 프레임워크

### CLI 프레임워크

#### ✅ Cobra
- **패키지**: `github.com/spf13/cobra`
- **용도**: CLI 명령어 구조 (`scan`, `eval`, `suggest`)
- **선택 이유**:
  - Go 생태계에서 표준 CLI 프레임워크
  - 서브커맨드, 플래그, 도움말 자동 생성 지원
  - Viper과 통합 가능 (설정 파일 관리)

### 설정 파일 처리

#### ✅ YAML 파서
- **패키지**: `gopkg.in/yaml.v3`
- **용도**: `curompt.yaml` 설정 파일 파싱
- **선택 이유**:
  - Go에서 가장 널리 사용되는 YAML 라이브러리
  - 안정적이고 문서화가 잘 됨
  - 설정 파일 형식에 적합

### 프롬프트 파싱

#### ✅ Goldmark
- **패키지**: `github.com/yuin/goldmark`
- **용도**: Markdown 프롬프트 파싱
- **선택 이유**:
  - Go로 작성된 빠른 Markdown 파서
  - 확장 가능한 아키텍처
  - CommonMark 호환

**대안 고려사항**:
- `gomarkdown`: 더 많은 기능이지만 무겁고 느림

### 데이터베이스

#### ✅ SQLite
- **패키지**: `modernc.org/sqlite` (순수 Go 구현)
- **용도**: 프롬프트 분석 결과 저장 (`~/.curompt/db.sqlite`)
- **선택 이유**:
  - 파일 기반, 설정 불필요
  - CGO 의존성 없는 순수 Go 구현
  - 로컬 전용 철학과 부합

**대안 고려사항**:
- `github.com/mattn/go-sqlite3`: CGO 필요, 더 빠르지만 의존성 복잡

### 테스트 프레임워크

#### ✅ Testify
- **패키지**: `github.com/stretchr/testify`
- **용도**: Assertions, Mocking, 테스트 헬퍼
- **선택 이유**:
  - Go에서 가장 인기 있는 테스트 라이브러리
  - Assertions 가독성 좋음
  - Mock 생성 및 사용 간편
  - Suite 패턴 지원

#### 커버리지 도구
- **표준**: `go test -cover`
- **고급**: `gocov` (선택 사항)

### HTTP 클라이언트

#### ✅ 표준 라이브러리 또는 HTTP 클라이언트 라이브러리
- **옵션 1**: `net/http` (표준 라이브러리)
- **옵션 2**: `github.com/go-resty/resty` (선택 사항)
- **용도**: LLM Provider API 호출 (Claude, OpenAI)
- **선택 기준**: 
  - 기본은 표준 라이브러리 사용
  - 복잡한 요청/응답 처리 필요 시 resty 고려

### 토큰 계산 라이브러리

#### Claude 토큰 계산
- **옵션 1**: `github.com/anthropics/anthropic-sdk-go` (공식 SDK)
- **옵션 2**: `github.com/pkoukk/tiktoken-go` (호환 라이브러리)
- **용도**: Claude API 토큰 계산
- **선택**: Adapter 패턴으로 추상화하여 교체 가능하게 구성

#### OpenAI 토큰 계산
- **패키지**: `github.com/pkoukk/tiktoken-go`
- **용도**: OpenAI API 토큰 계산
- **선택 이유**: tiktoken Python 구현의 Go 포트

### JSON Schema 검증 (M2)

#### ✅ JSON Schema Validator
- **옵션 1**: `github.com/xeipuuv/gojsonschema`
- **옵션 2**: `github.com/santhosh-tekuri/jsonschema/v5`
- **용도**: 출력 스키마 검증 (M2)
- **선택 기준**: M2에서 성능 및 기능 비교 후 결정

---

## 개발 도구

### 린터 및 포맷터

#### ✅ golangci-lint
- **용도**: 코드 품질 검사
- **설정**: `.golangci.yml` 파일로 규칙 관리
- **권장 규칙**:
  - `gofmt`, `govet`: 기본 포맷팅
  - `golint`, `errcheck`: 코드 품질
  - `gosec`: 보안 검사

### 빌드 도구

#### Make
- **용도**: 빌드, 테스트, 릴리스 자동화
- **파일**: `Makefile`

```makefile
.PHONY: build test lint clean

build:
	go build -o bin/curompt ./cmd/curompt

test:
	go test -v -cover ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/
```

---

## 프로젝트 구조

```
curompt/
├── cmd/
│   └── curompt/          # CLI 진입점
│       └── main.go
├── internal/                 # 내부 패키지 (외부 접근 불가)
│   ├── analyzer/             # 정적/동적 분석
│   ├── scorer/               # 점수화 엔진
│   ├── provider/            # LLM Provider 어댑터
│   │   ├── claude/
│   │   ├── openai/
│   │   └── provider.go       # 인터페이스
│   ├── parser/               # 프롬프트 파서
│   ├── reporter/             # 리포트 생성
│   ├── suggestor/            # 리팩터 제안
│   ├── redactor/             # 비밀값 마스킹
│   ├── repository/           # 데이터 저장소
│   └── config/               # 설정 관리
├── pkg/                      # 공개 패키지 (선택 사항)
├── test/                     # 테스트 픽스처
│   ├── fixtures/
│   └── integration/
├── docs/                     # 문서
├── examples/                 # 예제
├── prompts/                  # 프롬프트 템플릿
├── go.mod
├── go.sum
├── Makefile
├── .golangci.yml            # 린터 설정
└── README.md
```

---

## 의존성 목록 (go.mod 예상)

```go
module github.com/curogom/curompt

go 1.23

require (
    // CLI
    github.com/spf13/cobra v1.8.0
    
    // 설정
    gopkg.in/yaml.v3 v3.0.1
    
    // Markdown 파싱
    github.com/yuin/goldmark v1.7.0
    
    // 데이터베이스
    modernc.org/sqlite v1.29.0
    
    // 테스트
    github.com/stretchr/testify v1.9.0
    
    // HTTP (필요 시)
    github.com/go-resty/resty/v2 v2.12.0
    
    // JSON Schema (M2)
    github.com/xeipuuv/gojsonschema v1.2.0
    
    // 토큰 계산
    github.com/pkoukk/tiktoken-go v0.1.6
    github.com/anthropics/anthropic-sdk-go v1.0.0 // 또는 최신 버전
)
```

---

## GUI / 웹 인터페이스

### ❌ GUI 미지원 (현재)

**결정**: **CLI 전용**으로 개발

**이유**:
1. **CLI 도구의 본질**: 터미널 기반 워크플로우에 최적화
2. **단순성**: GUI 개발 복잡도와 유지보수 비용 증가
3. **통합 용이성**: CLI는 파이프라인, 스크립트, CI/CD에 직접 통합 가능
4. **YAGNI 원칙**: 현재 필요하지 않은 기능은 구현하지 않음

**출력 방식**:
- 터미널 출력: 요약 정보 및 진행 상황
- 파일 리포트: `reports/*.md`, `reports/*.json`
- 파이프라인 통합: 표준 입력/출력 활용

### 미래 고려사항 (M4 팀 프리미엄)

**웹 기반 대시보드** (선택 사항):
- **용도**: 팀 단위 통계 및 트렌드 분석
- **기술 옵션** (필요 시):
  - 간단한 HTTP 서버 (`net/http`)
  - 정적 HTML + JavaScript 차트 라이브러리
  - 또는 별도 웹 서비스로 분리
- **현재 상태**: M1-M3에서 제외, 필요성 재평가 후 결정

**데스크톱 GUI**:
- **결정**: 구현하지 않음
- **이유**: 웹 대시보드로 대체 가능하며, 개발 복잡도가 높음

---

## 선택하지 않은 기술 (이유)

### ❌ Rust
- **이유**: Go만으로도 충분히 빠르고, 생태계가 더 성숙함
- **추가 학습 곡선**: Rust는 메모리 안전성 학습 필요

### ❌ Python (메인 언어)
- **이유**: 
  - 배포 복잡도 증가 (가상환경, 의존성 관리)
  - 실행 속도 느림
  - 단일 바이너리 배포 어려움

### ❌ Node.js / TypeScript
- **이유**: 
  - 런타임 의존성 필요
  - 메모리 사용량 높음
  - CLI 도구보다 웹 애플리케이션에 적합

### ❌ GUI 프레임워크
- **예시**: Fyne, Wails, Electron
- **이유**:
  - CLI 도구의 본질과 부합하지 않음
  - 개발 복잡도 증가
  - 배포 및 의존성 관리 복잡
  - 터미널 기반 워크플로우가 더 효율적

### ❌ PostgreSQL / MySQL
- **이유**: 
  - 설정 및 관리 필요
  - 로컬 전용 철학과 부합하지 않음
  - SQLite로 충분함

---

## 기술 스택 검증 체크리스트

- [x] CLI 도구 개발에 적합한가? → Go + Cobra
- [x] 단일 바이너리 배포 가능한가? → Go 컴파일
- [x] 크로스 플랫폼 지원 가능한가? → Go 크로스 컴파일
- [x] 테스트 작성이 쉬운가? → Testify
- [x] 의존성 관리가 간단한가? → Go Modules
- [x] 성능이 충분한가? → Go는 충분히 빠름
- [x] 커뮤니티 지원이 좋은가? → Go 생태계 활발
- [x] 과도한 추상화 없이 구현 가능한가? → Go는 실용적
- [x] GUI 없이 구현 가능한가? → CLI 전용, 파일 리포트

---

## 최종 권장 기술 스택 요약

| 카테고리 | 선택 기술 | 버전/패키지 |
|---------|---------|------------|
| **언어** | Go | ≥ 1.23 |
| **CLI 프레임워크** | Cobra | `github.com/spf13/cobra` |
| **설정 파일** | YAML | `gopkg.in/yaml.v3` |
| **Markdown 파싱** | Goldmark | `github.com/yuin/goldmark` |
| **데이터베이스** | SQLite | `modernc.org/sqlite` |
| **테스트** | Testify | `github.com/stretchr/testify` |
| **HTTP 클라이언트** | net/http | 표준 라이브러리 |
| **토큰 계산** | tiktoken-go | `github.com/pkoukk/tiktoken-go` |
| **JSON Schema** | gojsonschema | `github.com/xeipuuv/gojsonschema` (M2) |
| **린터** | golangci-lint | - |

---

**최종 확인**: 이 기술 스택으로 개발을 진행합니다.

**변경 이력**:
- 2025-01-XX: 초기 기술 스택 결정

