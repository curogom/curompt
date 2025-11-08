# 문서 구조 가이드

## 개요

`curompt` 프로젝트의 문서는 **외부 공개**와 **내부 개발**로 명확히 구분되어 있습니다.

## 디렉토리 구조

```
docs/
├── README.md                    # 문서 구조 설명
├── public/                      # 📖 공개 문서
│   ├── README.md
│   ├── ARCHITECTURE.md         # 시스템 아키텍처
│   ├── CONFIG.md               # 설정 가이드
│   ├── METRICS.md              # 점수화 메트릭
│   ├── PROVIDERS.md            # Provider 지원 목록
│   ├── ROADMAP.md              # 프로젝트 로드맵
│   ├── SECURITY.md             # 보안 정책
│   └── TECH_STACK.md           # 기술 스택
│
└── internal/                    # 🔒 내부 문서
    ├── README.md
    ├── development/            # 개발 문서
    │   ├── DEVELOPMENT_PLAN.md
    │   ├── MVP_CLARIFICATION.md
    │   ├── M1_COMPLETE_CHECKLIST.md
    │   ├── M1_COMPLETION_SUMMARY.md
    │   └── M1_STATUS.md
    │
    └── strategy/                # 전략 문서
        ├── BUSINESS_LICENSE_STRATEGY.md
        ├── LICENSE_OPTIONS.md
        ├── SIMPLE_LICENSE_STRATEGY.md
        └── PUBLIC_REPO_CONSIDERATIONS.md
```

## 문서 분류 기준

### 📖 공개 문서 (`docs/public/`)

**외부에 공개되는 문서**

- ✅ 사용자 가이드
- ✅ 기여자 가이드
- ✅ 프로젝트 소개 및 로드맵
- ✅ 기술 문서 (아키텍처, 설정 등)

**특징:**
- 모든 사용자가 읽을 수 있음
- GitHub Public 저장소에 포함
- 지속적인 유지보수 필요
- PR을 통한 검토 권장

### 🔒 내부 문서 (`docs/internal/`)

**프로젝트 내부용 문서**

- 개발 프로세스 기록
- 의사결정 과정
- 전략 검토 문서
- 내부 체크리스트

**특징:**
- 개발 팀 내부 참고용
- Public 저장소 전환 시 제외 가능
- 자유로운 수정 가능

## 문서 작성 가이드

### 공개 문서 작성 시

1. **사용자 관점**: 기술적으로 정확하면서도 이해하기 쉬운 표현
2. **예시 포함**: 실제 사용 예시와 코드 스니펫
3. **정기 업데이트**: 기능 변경 시 문서 동기화
4. **링크 관리**: 상호 참조 링크 정확성 유지

### 내부 문서 작성 시

1. **의사결정 기록**: 왜 그렇게 결정했는지 배경 기록
2. **컨텍스트 보존**: 미래 참조를 위한 충분한 정보
3. **프로세스 문서화**: 개발 프로세스와 방법론 기록

## 링크 규칙

### 공개 문서 간 링크
```markdown
[문서명](./파일명.md)  # 같은 디렉토리 내
```

### 루트에서 공개 문서로
```markdown
[docs/public/파일명.md](./docs/public/파일명.md)
```

### 내부 문서 참조
- 공개 문서에서는 내부 문서를 직접 링크하지 않음
- 필요 시 일반적인 설명으로 대체

## Git 관리 전략

### 공개 문서
- 모든 변경사항 버전 관리
- PR 필수 (변경 검토)
- 릴리스 노트에 주요 변경사항 포함

### 내부 문서
- 개발 이력 추적 목적의 버전 관리
- 자유로운 수정 가능
- Public 저장소 전환 시 제외 옵션

## 공개 저장소 전환 시

### 포함할 디렉토리
- ✅ `docs/public/`
- ✅ 루트 `README.md`
- ✅ `LICENSE`

### 선택적 포함
- `docs/internal/` (원하는 경우 포함 가능, 비공개로 명시)

### 제외 가능
- `docs/internal/strategy/` (비즈니스 전략 문서)
- 일부 내부 개발 문서 (민감 정보 포함 시)

## 문서 업데이트 체크리스트

새 문서 작성 또는 이동 시:

- [ ] 적절한 디렉토리에 배치 (public vs internal)
- [ ] README.md에 문서 목록 업데이트
- [ ] 상호 참조 링크 확인 및 수정
- [ ] 루트 README.md 링크 업데이트 (공개 문서인 경우)
- [ ] 문서 작성 가이드 준수 확인

---

**최종 업데이트**: 2025-11-01

