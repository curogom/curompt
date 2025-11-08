# 시작하기 가이드

이 문서는 `curompt`를 처음 사용하는 사용자를 위한 실습 가이드입니다.

## 전제 조건

- Go ≥ 1.23
- macOS 또는 Linux
- 프로젝트 빌드 완료

## 1. 빌드 및 설치

```bash
# 프로젝트 디렉토리로 이동
cd /Users/curogom/dev/curompt

# 빌드
make build

# 바이너리 확인
./bin/curompt --help
```

예상 출력:
```
curompt는 CLI 기반 개발자의 LLM 프롬프트를 분석·평가·최적화하는 도구입니다.

Usage:
  curompt [command]

Available Commands:
  eval        프롬프트 평가 및 점수화
  scan        레포지토리 내 프롬프트 파일 스캔 및 분석
  suggest     프롬프트 개선 제안
```

## 2. 테스트용 프롬프트 파일 생성

먼저 테스트할 프롬프트 파일을 만들어봅시다:

```bash
# 테스트 디렉토리 생성
mkdir -p test-prompts

# 샘플 프롬프트 파일 생성
cat > test-prompts/sample.md << 'EOF'
# ROLE
Senior Software Engineer

# INPUTS
- task: string (개발할 기능 설명)
- context: string (프로젝트 컨텍스트)
- requirements: array (요구사항 목록)

# INVARIANTS
- 코드는 반드시 테스트 가능해야 함
- 모든 공개 함수는 문서화되어야 함
- 에러 처리는 명시적으로 처리해야 함

# OUTPUT FORMAT
JSON 형식으로 다음 구조를 반환:
{
  "code": "string",
  "tests": "string",
  "documentation": "string"
}
EOF
```

또는 더 복잡한 예시:

```bash
cat > test-prompts/complex.md << 'EOF'
# ROLE
Full-stack Developer

# INPUTS
- feature_description: string
- tech_stack: array
- deadline: date

# INVARIANTS
- 프론트엔드와 백엔드 모두 구현
- RESTful API 설계 원칙 준수
- 반응형 디자인 필수

# OUTPUT FORMAT
Markdown 형식의 구현 계획서
EOF
```

## 3. 기본 테스트: eval 명령

### 3.1 파일에서 평가

```bash
./bin/curompt eval --file test-prompts/sample.md
```

예상 출력:
- 프롬프트 분석 결과
- 종합 점수 (0-100)
- 메트릭별 점수
- 토큰 수
- 개선 제안

### 3.2 표준 입력에서 평가

```bash
cat test-prompts/sample.md | ./bin/curompt eval
```

또는 직접 입력:

```bash
echo "# ROLE\nEngineer\n\n# INPUTS\n- task: string" | ./bin/curompt eval
```

### 3.3 리포트 파일로 저장

```bash
./bin/curompt eval --file test-prompts/sample.md --output reports/sample_report.md
```

리포트 파일이 `reports/sample_report.md`에 저장됩니다.

### 3.4 Provider 변경 (토큰 계산용)

```bash
# Claude 사용 (기본값)
./bin/curompt eval --file test-prompts/sample.md --provider claude

# OpenAI 사용
./bin/curompt eval --file test-prompts/sample.md --provider openai
```

## 4. 일괄 스캔: scan 명령

### 4.1 현재 디렉토리 스캔

```bash
# test-prompts 디렉토리 스캔
./bin/curompt scan --repo test-prompts
```

예상 출력:
```
발견된 프롬프트 파일: 2개

[1/2] 분석 중: test-prompts/sample.md
  ✅ 완료 - 점수: 85.3/100, 리포트: reports/sample.md_report.md

[2/2] 분석 중: test-prompts/complex.md
  ✅ 완료 - 점수: 72.1/100, 리포트: reports/complex.md_report.md

총 2개 파일 분석 완료
리포트 저장 위치: reports
```

### 4.2 커스텀 출력 디렉토리

```bash
./bin/curompt scan --repo test-prompts --output my-reports
```

### 4.3 파일 패턴 지정

```bash
# .md 파일만 스캔
./bin/curompt scan --repo test-prompts --patterns "*.md"

# 여러 패턴
./bin/curompt scan --repo test-prompts --patterns "*.md" --patterns "*.txt"
```

## 5. 개선 제안: suggest 명령

### 5.1 기본 제안 확인

```bash
./bin/curompt suggest --file test-prompts/sample.md
```

예상 출력:
```
📋 개선 제안:
==================================================

1. ✅ 프롬프트가 이미 잘 구성되어 있습니다!

==================================================

현재 점수: 85.3 / 100
```

### 5.2 개선이 필요한 프롬프트 테스트

개선이 필요한 샘플 파일 생성:

```bash
cat > test-prompts/needs-improvement.md << 'EOF'
# ROLE
Engineer
EOF
```

이 파일로 제안 확인:

```bash
./bin/curompt suggest --file test-prompts/needs-improvement.md
```

예상 출력:
```
📋 개선 제안:
==================================================

1. 🟡 INPUTS 섹션 추가 권장 - 입력값을 명시하면 프롬프트가 더 명확해집니다
2. 🟡 INVARIANTS 섹션 추가 권장 - 규칙과 제약사항을 명시하세요
3. 🟡 OUTPUT FORMAT 섹션 추가 권장 - 출력 형식을 명시하면 결과 품질이 향상됩니다
4. 🟠 종합 점수 개선 여지가 있습니다 - 제안 사항을 적용하여 더 높은 점수를 목표로 하세요

==================================================

현재 점수: 45.2 / 100
```

### 5.3 표준 입력에서 제안 확인

```bash
cat test-prompts/sample.md | ./bin/curompt suggest
```

## 6. 실제 워크플로우 테스트

### 6.1 전체 워크플로우

```bash
# 1. 프롬프트 파일 준비
cat > my-prompt.md << 'EOF'
# ROLE
Developer

# INPUTS
- feature: string

# OUTPUT FORMAT
Code in Go
EOF

# 2. 평가 및 점수 확인
./bin/curompt eval --file my-prompt.md

# 3. 개선 제안 확인
./bin/curompt suggest --file my-prompt.md

# 4. 리포트 저장
./bin/curompt eval --file my-prompt.md --output my-report.md

# 5. 리포트 확인
cat my-report.md
```

### 6.2 여러 프롬프트 일괄 평가

```bash
# prompts 디렉토리에 여러 파일 준비
mkdir -p prompts
# ... 여러 프롬프트 파일 생성 ...

# 일괄 스캔 및 분석
./bin/curompt scan --repo prompts --output reports

# 리포트 확인
ls -la reports/
```

## 7. 예상 결과 해석

### 7.1 점수 해석

- **90-100점**: 매우 우수한 프롬프트
- **80-89점**: 양호한 프롬프트
- **70-79점**: 개선 여지 있음
- **60점 미만**: 전면 개선 필요

### 7.2 메트릭 이해

**Structure (구조)**
- 섹션 존재 여부
- 중복 규칙 제거
- 필수 섹션 포함

**Conciseness (간결성)**
- 토큰 밀도
- 불필요한 설명 제거

**Risk (위험)**
- 모호한 표현 감지
- 민감 정보 노출 가능성

### 7.3 리포트 파일 내용

리포트 파일(`.md`)에는 다음이 포함됩니다:
- 프롬프트 메타데이터
- 종합 점수 및 메트릭별 점수
- 정적 분석 결과
- 토큰 정보
- 개선 제안

## 8. 문제 해결

### 8.1 바이너리를 찾을 수 없음

```bash
# 빌드 확인
make build

# 바이너리 존재 확인
ls -la ./bin/curompt
```

### 8.2 파일을 읽을 수 없음

```bash
# 파일 경로 확인 (절대 경로 사용)
./bin/curompt eval --file /full/path/to/prompt.md

# 파일 권한 확인
ls -la test-prompts/sample.md
```

### 8.3 점수가 0점 또는 이상한 결과

- 프롬프트 형식 확인 (Markdown 섹션 구조)
- 파일 인코딩 확인 (UTF-8)
- 파일 내용 확인 (`cat`으로 출력)

### 8.4 리포트 파일이 생성되지 않음

```bash
# 출력 디렉토리 권한 확인
mkdir -p reports
chmod 755 reports

# 다시 시도
./bin/curompt scan --repo test-prompts --output reports
```

## 9. 고급 사용법

### 9.1 파이프라인 사용

```bash
# 프롬프트 생성 → 평가 → 제안
echo "# ROLE\nEngineer" | \
  ./bin/curompt eval | \
  tee eval-output.txt | \
  ./bin/curompt suggest
```

### 9.2 배치 처리 스크립트

```bash
#!/bin/bash
# evaluate-all.sh

for file in prompts/*.md; do
    echo "Evaluating: $file"
    ./bin/curompt eval --file "$file" --output "reports/$(basename $file .md)_report.md"
done
```

## 10. 다음 단계

테스트가 성공적으로 완료되었다면:

1. **실제 프롬프트로 테스트**: 자신의 프로젝트 프롬프트 사용
2. **리포트 분석**: 생성된 리포트 확인 및 개선 사항 적용
3. **문서 탐색**: [아키텍처](./ARCHITECTURE.md), [설정 가이드](./CONFIG.md) 참조
4. **기여**: 이슈 리포트나 개선 제안 환영

---

**문제가 발생하면**: GitHub Issues에 리포트하거나 문서를 확인하세요.

