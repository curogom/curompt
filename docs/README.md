# 문서 구조

이 디렉토리는 `curompt` 프로젝트의 문서를 포함합니다.

## 문서 분류

### 📖 공개 문서 (`public/`)

**외부 공개 문서**: 사용자, 기여자, 커뮤니티에게 공개되는 문서입니다.

- **ARCHITECTURE.md**: 시스템 아키텍처 및 모듈 구조
- **CONFIG.md**: 설정 파일 가이드
- **METRICS.md**: 점수화 메트릭 설명
- **PROVIDERS.md**: 지원 LLM Provider 목록
- **ROADMAP.md**: 프로젝트 로드맵 및 마일스톤
- **SECURITY.md**: 보안 및 프라이버시 정책
- **TECH_STACK.md**: 기술 스택 및 선택 이유

### 🔒 내부 문서 (`internal/`)

**내부 개발 문서**: 프로젝트 내부 개발 및 전략 문서입니다.

#### 개발 문서 (`internal/development/`)
- **DEVELOPMENT_PLAN.md**: 개발 계획, TDD 방법론, 디자인 패턴 적용
- **MVP_CLARIFICATION.md**: MVP 명세 및 목표
- **M1_COMPLETE_CHECKLIST.md**: M1 완료 체크리스트
- **M1_COMPLETION_SUMMARY.md**: M1 완료 요약
- **M1_STATUS.md**: M1 현재 상태

#### 전략 문서 (`internal/strategy/`)
- **BUSINESS_LICENSE_STRATEGY.md**: 라이선스 전략 검토
- **LICENSE_OPTIONS.md**: 라이선스 옵션 비교
- **SIMPLE_LICENSE_STRATEGY.md**: 단순 라이선스 전략
- **PUBLIC_REPO_CONSIDERATIONS.md**: Public 저장소 전환 고려사항

## 문서 작성 가이드

### 공개 문서 작성 시
- 사용자 관점에서 작성
- 기술적 정확성과 가독성 유지
- 예시와 코드 스니펫 포함
- 정기적으로 업데이트

### 내부 문서 작성 시
- 개발자 관점에서 작성
- 개발 프로세스, 의사결정 기록
- 미래 참조를 위한 컨텍스트 포함

## 문서 링크 업데이트

문서 구조 변경 시 다음 파일의 링크를 업데이트해야 합니다:
- `README.md` (루트)
- `docs/public/` 내 문서들 간 상호 참조
- 코드 내 주석

## Git 관리

### 공개 문서
- **버전 관리**: 모든 변경사항 추적
- **PR 필수**: 변경 시 PR을 통해 검토

### 내부 문서
- **버전 관리**: 개발 이력 추적 목적
- **자유로운 수정**: 내부 개발 프로세스 반영

---

**참고**: 공개 저장소로 전환 시 `internal/` 디렉토리는 포함하지 않을 수 있습니다.

