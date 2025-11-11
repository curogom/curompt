# curompt

CLI 기반 개발자의 LLM 프롬프트를 **분석·평가·최적화**하는 도구.

- 대상: Claude, OpenAI, Gemini, Cursor IDE/CLI, Codex CLI, Bedrock/Vertex/로컬 LLM을 CLI로 쓰는 워크플로우
- 목표: 정확도·재현성·비용을 지표화하고, 자동 리팩터 제안
- 기본 철학: **로컬 우선**, **계약 우선(OpenAPI/JSONSchema)**, **재현 가능 리포트**
- 라이선스: [Apache License 2.0](./LICENSE)

## 핵심 기능
- 프롬프트 정적 분석: 섹션 구조, 중복 규칙, 금지어, 스키마 유무
- 동적 평가(옵션): 다중 샘플 → 스키마 적합률·자체 일관성·지연·비용
- 점수화: 0–100 종합 점수 + 하위 지표
- 리팩터 제안: 토큰 절감, 규칙 분리, few-shot 축약, 캐시 최적화
- 리포트 출력: 터미널 요약 + `reports/*.md|json`

## 빠른 시작

### 1. 빌드
```bash
make build
./bin/curompt --help
```

### 2. 기본 테스트

```bash
# 샘플 프롬프트 평가
echo "# ROLE\nEngineer\n\n# INPUTS\n- task: string" | ./bin/curompt eval

# 파일로 평가
./bin/curompt eval --file prompts/dev_contract_v2.md

# 일괄 스캔 (reports/에 저장된 프롬프트 기준)
./bin/curompt scan --path prompts/

# 개선 제안 확인
./bin/curompt suggest --file prompts/dev_contract_v2.md
```

> `scan` 명령은 collect/eval로 저장된 프롬프트를 분석합니다.  
> 현재 경로에 해당하는 히스토리가 없으면 CLI가 Claude Code 또는 Codex 로그를 자동 수집할지 물어봅니다( Cursor 지원은 v1.1 예정 ).  
> 특정 프로젝트에 한정하려면 `--path /프로젝트/절대/경로` 옵션을 사용하세요.

### 3. 빠른 테스트 스크립트

```bash
# 자동 테스트 스크립트 실행 (빌드부터 테스트까지)
./test-quick.sh
```

이 스크립트는 다음을 자동으로 수행합니다:
- ✅ 빌드 확인
- ✅ 테스트 프롬프트 파일 생성
- ✅ eval, suggest, scan 명령 테스트
- ✅ 리포트 생성 확인

> **📖 상세한 실습 가이드**: [시작하기 가이드](./docs/public/GETTING_STARTED.md)를 참조하세요 (단계별 설명, 예시, 문제 해결 포함)

## 설치

### 요구사항
- Go ≥ 1.23
- macOS/Linux
- 선택: Python3 (플러그인 용, 향후)

### 설치 방법

#### 방법 1: Homebrew (macOS 권장) ⭐

macOS에서 `curompt`를 설치하는 가장 쉬운 방법:

```bash
# Tap 추가
brew tap curogom/curompt

# 설치
brew install curompt

# 설치 확인
curompt --version
```

**업데이트**:
```bash
brew upgrade curompt
```

> **✅ 사용 가능**: Homebrew tap이 준비되었습니다! 위 명령어로 `curompt`를 설치할 수 있습니다.

#### 방법 2: Go install

Go 모듈을 직접 설치하여 `$GOPATH/bin` 또는 `~/go/bin`에 설치합니다:

```bash
# 저장소 클론 (또는 이미 클론한 경우)
git clone https://github.com/curogom/curompt.git
cd curompt

# 설치
make install
# 또는 직접 설치
go install github.com/curogom/curompt/cmd/curompt@latest
```

**PATH 확인 및 설정**:
```bash
# Go bin 경로 확인
go env GOPATH
# 출력 예: /Users/username/go

# PATH에 추가 (이미 되어있을 수도 있음)
export PATH="$PATH:$(go env GOPATH)/bin"

# 확인
which curompt
curompt --help
```

**영구적으로 PATH 추가** (셸 설정 파일에 추가):
```bash
# zsh 사용자 (~/.zshrc)
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
source ~/.zshrc

# bash 사용자 (~/.bashrc 또는 ~/.bash_profile)
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc
source ~/.bashrc
```

#### 방법 2: 로컬 빌드 (개발용)

```bash
# 빌드
make build

# 실행 (PATH 없이)
./bin/curompt --help

# 또는 수동으로 PATH에 추가
sudo cp ./bin/curompt /usr/local/bin/
# 또는
cp ./bin/curompt ~/bin/
export PATH="$PATH:~/bin"
```

### 설치 확인

```bash
# 버전 확인
curompt --version

# 도움말 확인
curompt --help

# 기본 테스트
echo "# ROLE\nEngineer" | curompt eval
```

> **📖 상세 설치 가이드**: [설치 가이드](./docs/public/INSTALLATION.md) 참조 (PATH 설정, 문제 해결 포함)

### 개발용 빌드 및 테스트

```bash
# 빌드만 하기
make build

# 테스트 실행
make test

# 테스트 커버리지 확인
make coverage

# 모든 검사 (포맷, 린트, 테스트)
make check
```

## 통합 지원

### LLM Provider
- **Claude**: Anthropic API 통합 (토큰 계산, 메타데이터)
- **OpenAI**: OpenAI API 통합 (예정)
- **Gemini**: Google Gemini API 통합 (예정)

### 도구 통합
- **Cursor IDE/CLI**: CLI 래퍼를 통한 프롬프트 캡처
- **Codex CLI**: CLI 래퍼를 통한 프롬프트 캡처
- **Bedrock/Vertex/로컬**: 표준 입력/출력 래핑 지원

> **참고**: MVP에서는 LLM Provider는 메타데이터(토큰 계산, 비용 추정)만 제공하며, 실제 API 호출은 M2에서 구현됩니다.

## 보안 및 프라이버시

- **로컬 우선**: 기본적으로 로컬에서만 실행, 외부 전송 없음
- **자동 마스킹**: API 키/토큰/이메일/URL 쿼리/.env 참조 자동 마스킹
- **로컬 저장소**: SQLite로 로컬에만 저장 (`~/.curompt/db.sqlite`)

상세 내용은 [SECURITY.md](./docs/public/SECURITY.md) 참조.

## 아키텍처

모듈화된 아키텍처로 확장성과 유지보수성을 확보했습니다:

- **모듈 분리**: Collector, Parser, Analyzer, Scorer, Provider 등 명확한 책임 분리
- **디자인 패턴**: Strategy, Repository, Dependency Injection 패턴 적용
- **테스트 전략**: 단위 테스트(85% 커버리지) + 통합 테스트

상세 아키텍처는 [ARCHITECTURE.md](./docs/public/ARCHITECTURE.md) 참조.

## 테스트

### 단위 테스트
```bash
make test
```

### 통합 테스트
```bash
go test ./test/integration/... -v
```

### 테스트 커버리지
```bash
make coverage
```

**현재 커버리지**: 평균 85% (목표 80% 달성)

## 문서

### 사용자 문서 (공개)
- **[설치 가이드](./docs/public/INSTALLATION.md)**: 시스템 설치 및 PATH 설정 ⭐
- **[시작하기 가이드](./docs/public/GETTING_STARTED.md)**: 실습 테스트 가이드 ⭐
- [아키텍처](./docs/public/ARCHITECTURE.md): 모듈 구조 및 디자인 패턴
- [로드맵](./docs/public/ROADMAP.md): 마일스톤 및 기능 계획
- [기술 스택](./docs/public/TECH_STACK.md): 사용 기술 및 선택 이유
- [설정 가이드](./docs/public/CONFIG.md): YAML 설정 파일 가이드
- [메트릭](./docs/public/METRICS.md): 점수화 메트릭 설명
- [Provider 지원](./docs/public/PROVIDERS.md): 지원 LLM Provider 목록
- [보안](./docs/public/SECURITY.md): 보안 및 프라이버시 정책

> 모든 사용자 문서는 [`docs/public/`](./docs/public/)에서 확인할 수 있습니다.

## 기여

프로젝트는 Apache License 2.0 하에 배포됩니다. 기여를 환영합니다!

## 라이선스

Apache License 2.0

상세 내용은 [LICENSE](./LICENSE) 참조.
