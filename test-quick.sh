#!/bin/bash
# 빠른 테스트 스크립트

echo "=== 1. 빌드 확인 ==="
make build || exit 1

echo ""
echo "=== 2. 테스트 프롬프트 생성 ==="
mkdir -p test-prompts
cat > test-prompts/sample.md << 'PROMPT'
# ROLE
Senior Engineer

# INPUTS
- task: string
- context: string

# INVARIANTS
- Must be testable
- Must be documented

# OUTPUT FORMAT
JSON
PROMPT

echo "✅ test-prompts/sample.md 생성 완료"

echo ""
echo "=== 3. eval 명령 테스트 ==="
./bin/curompt eval --file test-prompts/sample.md | head -20

echo ""
echo "=== 4. suggest 명령 테스트 ==="
./bin/curompt suggest --file test-prompts/sample.md | head -15

echo ""
echo "=== 5. 참고: scan은 히스토리 기반 ==="
echo "현재 스크립트에서는 eval/suggest만 확인합니다."
echo "실제 환경에서는 'curompt collect --from claude' 후 'curompt scan --path <project>'를 실행하세요."

echo ""
echo "=== 테스트 완료! ==="
