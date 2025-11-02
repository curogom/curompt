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
./bin/curo-prompt eval --file test-prompts/sample.md | head -20

echo ""
echo "=== 4. suggest 명령 테스트 ==="
./bin/curo-prompt suggest --file test-prompts/sample.md | head -15

echo ""
echo "=== 5. scan 명령 테스트 ==="
mkdir -p test-reports
./bin/curo-prompt scan --repo test-prompts --output test-reports 2>&1 | tail -10

echo ""
echo "=== 테스트 완료! ==="
echo "리포트 확인: ls -la test-reports/"
