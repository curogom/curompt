# M1 완료 및 Initial Commit 체크리스트

## M1 완료 확인

### 핵심 기능 완료 여부
- [x] 프롬프트 파서 (Markdown 섹션 파싱)
- [x] 정적 분석기 (구조, 중복, 섹션 카운트)
- [x] 토큰 계산 (Claude, OpenAI)
- [x] 점수화 엔진 (0-100 점수)
- [x] Redaction (비밀값 마스킹)
- [x] Collector 인터페이스 및 기본 구현
- [x] SQLite 저장소
- [x] Evaluator (통합)
- [x] Reporter (Markdown 리포트)
- [x] CLI 명령 (eval 구현 완료)

### 테스트 커버리지
- [x] analyzer: 96.0%
- [x] parser: 87.9%
- [x] collector: 86.2%
- [x] evaluator: 90.0%
- [x] redactor: 100.0%
- [x] reporter: 88.5%
- [x] repository: 66.1%
- [x] scorer: 89.8%

**평균 커버리지**: ~85% ✅ (목표 80% 이상)

---

## Initial Commit 전 체크리스트

### 1. 코드 품질 ✅
- [x] `make ci` 통과 확인
- [x] `make lint` 통과 확인
- [x] `make test` 통과 확인
- [x] 테스트 커버리지 확인 (80% 이상)

### 2. 문서 ✅
- [x] README.md 완성
- [x] LICENSE (Apache 2.0) 업데이트
- [x] 개발 문서 완성 (docs/)
- [x] GitHub Description 준비

### 3. CI/CD 설정 ✅
- [x] `.github/workflows/ci.yml` - 테스트/린트/빌드
- [x] `.github/workflows/coverage.yml` - 커버리지 리포트
- [x] `.github/workflows/release.yml` - 릴리스 빌드
- [x] `.golangci.yml` - 린터 설정

### 4. 보안 점검 ✅
- [x] `.gitignore` 확인 (민감 정보 제외)
- [x] 하드코딩된 API 키/비밀값 없음
- [x] Redaction 기능 구현 확인

### 5. Git 설정
- [ ] `.git/` 디렉토리 확인
- [ ] 브랜치 전략 확인 (main vs develop)
- [ ] Initial commit 메시지 준비

---

## Initial Commit 절차

### 1. Git 초기화 (아직 안 했다면)
```bash
cd /Users/curogom/dev/curompt
git init
git branch -M main  # 또는 develop
```

### 2. 파일 추가 전 최종 점검
```bash
# CI 테스트
make ci

# 커버리지 확인
make coverage

# 빌드 확인
make build
```

### 3. Initial Commit
```bash
# 모든 파일 추가
git add .

# Initial commit
git commit -m "feat: M1 MVP 초기 구현

- 프롬프트 파서 (Markdown 섹션 파싱)
- 정적 분석기 (구조, 중복 규칙 감지)
- 토큰 계산 (Claude, OpenAI)
- 점수화 엔진 (0-100 종합 점수)
- Redaction (비밀값 마스킹)
- Collector 인터페이스 및 기본 구현
- SQLite 저장소 (Repository 패턴)
- Evaluator (수집 → 분석 → 점수화 통합)
- Reporter (Markdown 리포트 생성)
- CLI 명령 (eval 구현)
- CI/CD 파이프라인 설정
- 테스트 커버리지: ~85%

License: Apache License 2.0"

# 원격 저장소 연결 (GitHub)
git remote add origin https://github.com/curogom/curompt.git

# 첫 푸시
git push -u origin main
```

---

## Commit 메시지 템플릿 (선택)

### Conventional Commits 스타일
```
feat: M1 MVP 초기 구현

구현된 기능:
- 프롬프트 파서 및 정적 분석
- 점수화 엔진 (6개 메트릭)
- Collector 및 저장소 시스템
- CLI eval 명령

테스트 커버리지: ~85%
License: Apache License 2.0

Refs: #M1
```

---

## GitHub 저장소 설정 (Public 전환)

### 1. 저장소 생성 후
1. Settings → General → Visibility: **Public** 설정
2. Settings → General → Description 추가
   ```
   CLI 기반 LLM 프롬프트 분석·평가 도구 - Claude Code, Codex, Cursor 등에서 사용한 프롬프트를 자동 수집하고 정적 분석·점수화·리포트 생성 (Apache License 2.0)
   ```
3. About 섹션:
   - Topics 추가: `cli`, `llm`, `prompt-engineering`, `go`, `prompt-analysis`
   - Website (있으면)
   - License: Apache-2.0 (자동 인식)

### 2. README 개선 (선택)
- Badges 추가 (build, coverage, license)
- 사용 예시 스크린샷
- Contributing 가이드 링크

---

## 브랜치 전략 권장

### 옵션 1: 단순 전략
- `main`: 프로덕션/릴리스 브랜치
- Feature 브랜치: `feature/xxx`

### 옵션 2: Git Flow
- `main`: 릴리스
- `develop`: 개발
- `feature/xxx`: 기능 개발

**추천**: 옵션 1 (단순 전략) - MVP 단계에서는 충분

---

## 다음 단계 (M1 완료 후)

1. **피드백 수집**: Public 저장소에서 이슈/PR 확인
2. **문서 개선**: 실제 사용자 피드백 반영
3. **M2 시작**: 동적 평가 기능 추가 (필요 시)
4. **Collector 개선**: CLI 래퍼 세부 구현

---

## 빠른 체크 명령어

```bash
# 모든 검사 한 번에
make check

# CI 시뮬레이션
make ci

# 커버리지 확인
make coverage && open coverage.html
```

---

**준비 완료 시 Initial Commit 진행!** 🚀

