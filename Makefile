.PHONY: build test lint clean install run help ci coverage

# 변수
BINARY_NAME=curo-prompt
BINARY_PATH=./bin/$(BINARY_NAME)
MAIN_PATH=./cmd/curo-prompt
VERSION?=dev
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 빌드 플래그
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# 기본 타겟
.DEFAULT_GOAL := help

help: ## 도움말 출력
	@echo "사용 가능한 명령어:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## 바이너리 빌드
	@echo "빌드 중..."
	@mkdir -p bin
	go build $(LDFLAGS) -o $(BINARY_PATH) $(MAIN_PATH)
	@if [ ! -f $(BINARY_PATH) ]; then \
		echo "Error: Build failed - binary not found"; \
		exit 1; \
	fi
	@echo "빌드 완료: $(BINARY_PATH)"

install: ## 시스템에 설치 (GOPATH/bin 또는 ~/go/bin)
	go install $(LDFLAGS) $(MAIN_PATH)
	@echo "설치 완료: $$(go env GOPATH)/bin/$(BINARY_NAME)"

test: ## 테스트 실행
	go test -v -cover ./...

test-coverage: ## 테스트 커버리지 리포트
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "커버리지 리포트 생성: coverage.html"

lint: ## 코드 린트 검사 (CI와 동일한 방식)
	@echo "🔍 Running lint checks (same as CI)..."
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run --timeout=5m || (echo "❌ Lint failed! Fix issues before pushing." && exit 1)

fmt: ## 코드 포맷팅
	go fmt ./...
	@echo "포맷팅 완료"

clean: ## 빌드 결과물 삭제
	rm -rf bin/
	rm -f coverage.out coverage.html
	@echo "정리 완료"

run: build ## 빌드 후 실행
	$(BINARY_PATH) --help

deps: ## 의존성 업데이트
	go get -u ./...
	go mod tidy

deps-download: ## 의존성 다운로드
	go mod download

vet: ## go vet 실행
	go vet ./...

check: fmt vet lint test ## 모든 검사 실행 (포맷, vet, 린트, 테스트)

release-prepare: ## 릴리스 준비 (태그 생성 및 Formula 업데이트)
	@echo "릴리스 준비:"
	@echo "1. 버전 태그 생성 예: git tag -a v0.1.0 -m 'Release v0.1.0'"
	@echo "2. 태그 푸시: git push origin v0.1.0"
	@echo "3. GitHub Release가 자동 생성됩니다"
	@echo "4. Formula 업데이트: ./scripts/update-homebrew-formula.sh v0.1.0"

.PHONY: release-prepare

ci: fmt vet lint test ## CI 파이프라인 실행 (로컬에서 CI 테스트)
	@echo "✅ CI 검사 완료"

coverage: test-coverage ## 테스트 커버리지 (별칭)

