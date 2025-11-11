## M1.5 제로-LLM 구조 채점 보완 계획

### 목표
- 헤더 키워드 의존도를 낮추고, 오탈자/다국어/자유형식 프롬프트에서도 안정적으로 구조 점수 산출
- 현 점수식·리포트 출력과의 호환 유지, 성능/의존성 증대 없이 적용

### 범위
- 파서: 섹션 신호 기반 추출 + 퍼지 매칭 + 별칭 설정 주입 + 확신도 산출
- 분석기: 이진 필드 유지 + 확신도 필드 추가 활용
- 스코어러: Structure 점수에 부분 점수(확신도 가중) 반영
- 구성: 섹션 별칭·임계치 YAML 추가
- 테스트: 단위/통합/회귀

### 피드백 반영(최종 보완)
- 성능 최적화
  - 퍼지 매칭 비교 범위를 “섹션 수 × 별칭 후보 수”로 제한(헤더 후보만 비교)
  - 별칭·키워드 사전은 사전 정규화(소문자/공백 제거) 캐시 사용
  - 코드블록 파싱 및 JSON/YAML 유효성 검사는 첫 N개 블록(예: 3개)까지만 검사
  - 목표: 5k 토큰 입력 기준 < 50ms(단일 스레드) 유지
- 설정 관리 전략
  - 기본값(영/한 별칭, 임계치)은 코드에 내장
  - `configs/sections.yaml`은 “override 병합” 방식(없으면 기본 사용)
  - 프로젝트 레벨 파일이 있을 때만 덮어쓰기, 키 단위 부분 병합 지원
- 확신도 결합 설계
  - 초기 가중치(문서화/튜닝 계획 포함):
    - 헤더(직접/퍼지) 0.6, 키워드 0.2, 패턴(목록/JSON 등) 0.2
  - 결합 방식: 정규화된 가중합 후 [0,1]로 클램프, 임계치(예: 0.3) 미만은 0 처리
  - 튜닝 계획: 테스트 코퍼스 기반 Grid/BO 간단 탐색, 주 단위 점검
- 테스트 범위 구체화
  - 언어: 한국어/영어/혼용
  - 오탈자: 1–2자 편집 거리, 자모 분해 케이스
  - 포맷: 헤더 유/무, 불릿/번호 목록, 코드블록(JSON/YAML), 표 형태
  - 노이즈: 이모지, 마크다운 장식, 긴 컨텍스트
  - 샘플 구성: `test/fixtures/prompts/m15/*.md`에 케이스별 최소 3개씩
- Reporter/CLI 토글
  - 기본: 확신도 비표시(현행과 동일)
  - CLI 플래그: `--show-structure-confidence` 추가(또는 설정 `reporter.show_confidence: true`)
  - 표시는 소수 1자리 백분율(예: Role 78.0%)

### 변경사항 설계
- parser (`internal/parser/parser.go`)
  - 섹션 추출 파이프라인 확장: 기존 헤더 파싱 + 신호 기반 감지(불릿/번호/키워드/코드블록 JSON 등) + 퍼지 매칭(Jaro-Winkler/n-gram)
  - 결과 모델에 확신도 추가:
    - `Prompt`: `RoleConfidence float64`, `InputsConfidence float64`, `InvariantsConfidence float64`, `OutputFormatConfidence float64`
  - 설정 로딩: `configs/sections.yaml`에서 별칭·임계치 로드
- analyzer (`internal/analyzer/analyzer.go`)
  - `AnalysisResult` 확장: 위 4개 `…Confidence` 반영
  - 중복 검출 유지(추가로 유사 문장 중복은 후속)
- scorer (`internal/scorer/metrics_calculator.go`)
  - Structure 계산식에 부분 점수 반영(0~1 확신도 기반):
    - `50*RoleConfidence + 10*InputsConfidence + 10*InvariantsConfidence + 10*OutputFormatConfidence - 5*DuplicateCount` (0~100 클램프)
  - 기존 이진 로직은 Confidence=0/1로 자동 호환
- reporter (`internal/reporter/markdown_reporter.go`)
  - 선택: 확신도(%) 표시 한 줄 추가. 기본 비표시로 호환성 유지
  - CLI 플래그/설정으로 노출 제어(`--show-structure-confidence`/`reporter.show_confidence`)
- 구성 파일
  - `configs/sections.yaml` (예시)
    ```
    sections:
      role:        ["role", "역할", "you are", "as a", "당신은"]
      inputs:      ["inputs", "입력", "parameters", "요구 입력"]
      invariants:  ["invariants", "불변", "제약", "금지", "constraints"]
      output:      ["output format", "출력 형식", "응답 형식", "schema", "format"]
    thresholds:
      fuzzy: 0.85
    ```

### 알고리듬 개요
- 텍스트 정규화: 소문자화/자모 분해(한글)/공백·구두점 정리/코드블록 분리
- 퍼지 매칭: Damerau-Levenshtein, Jaro-Winkler, n-gram Jaccard로 별칭 후보 스코어링(임계치 ≥ 0.85)
- 구조 신호:
  - Inputs/Invariant: 불릿(`- `)·번호(`1.`) 목록, “입력/parameters/제약/must/금지” 등의 키워드
  - Output Format: 코드블록 내 JSON/YAML/XML 유효성, 표 헤더 패턴, “format/schema/출력 예시”
  - Role: 선언적 역할 문장(“You are…/As a…/당신은 … 역할”)
- 확신도 결합: 헤더/키워드/패턴/퍼지 스코어를 앙상블해 섹션별 `Confidence ∈ [0,1]` 산출
- 점수화: 부분 점수(가산점 × 확신도) + 중복 감점, 최종 클램프

### 성능/의존성
- 표준 라이브러리 + 경량 퍼지 구현(또는 단일 소형 라이브러리)만 사용
- LLM 비사용(비용/지연 증가 없음)

### 단계별 일정(예시: 3–4일)
- D1
  - 설정 로더·데이터 구조 추가
  - 퍼지 매칭 유틸, 텍스트 정규화(소문자/공백/구두점/코드블록 분리)
- D2
  - 파서 신호 기반 섹션 감지 + 퍼지 결합 → 확신도 산출
  - 단위 테스트(한/영/오탈자 케이스)
- D3
  - analyzer 확신도 전달, Structure 부분 점수 반영
  - 회귀 테스트(기존 헤더 기반 프롬프트 점수 동일성)
- D4
  - 리포트(선택적 확신도 표시), 문서/예제 업데이트
  - 테스트 코퍼스 확대 및 임계치 튜닝

### 수용 기준(AC)
- 헤더 없는 한국어 프롬프트에서 `Structure ≥ 20` 이상(내용 충실 시)
- 오탈자 1–2자 내 섹션 라벨에서 헤더 인식률 ≥ 90%
- 기존 헤더 완비 영문 프롬프트의 종합 점수 변화 ±1점 이내
- 성능: 5k 토큰 프롬프트 처리 < 50ms(퍼지 + 정규식 경로, 단일 스레드 기준)

### 리스크·완화
- 과검출: 다중 신호(키워드+패턴) 동시 충족 시 확신도 상향, 단일 신호는 상한 제한
- 언어별 편차: 언어 감지 후 별칭 목록 가중치 적용, 임계치 언어별 분기
- 설정 스프롤: 프로젝트 로컬 YAML 병합, 기본값 보장

### Out of Scope (M1.5)
- LLM 보조 파싱/라벨링
- 의미 중복(문장 임베딩 기반) 감점 고도화
- SelfConsistency/LatencyCost 동적 구현


