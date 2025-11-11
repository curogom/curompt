package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsePrompt_DetectsROLESection(t *testing.T) {
	// Red: 실패하는 테스트 작성
	input := "# ROLE\nSenior engineer"
	parser := NewParser()

	result, err := parser.Parse(input)
	require.NoError(t, err)
	assert.Equal(t, "Senior engineer", result.Role)
	assert.InDelta(t, 1.0, result.RoleConfidence, 0.0001)
}

func TestParsePrompt_DetectsINPUTSSection(t *testing.T) {
	input := `# ROLE
Engineer

# INPUTS
- stack_profile: YAML
- task: string`
	parser := NewParser()

	result, err := parser.Parse(input)
	require.NoError(t, err)
	assert.Equal(t, 2, len(result.Inputs))
	assert.Contains(t, result.Inputs, "stack_profile: YAML")
	assert.Contains(t, result.Inputs, "task: string")
	assert.InDelta(t, 1.0, result.InputsConfidence, 0.0001)
}

func TestParsePrompt_DetectsINVARIANTSection(t *testing.T) {
	input := `# ROLE
Engineer

# INVARIANTS
- 계약 우선
- allowed_packages 안에서만 선택`
	parser := NewParser()

	result, err := parser.Parse(input)
	require.NoError(t, err)
	assert.Equal(t, 2, len(result.Invariants))
	assert.Contains(t, result.Invariants, "계약 우선")
	assert.Contains(t, result.Invariants, "allowed_packages 안에서만 선택")
	assert.InDelta(t, 1.0, result.InvariantsConfidence, 0.0001)
}

func TestParsePrompt_DetectsOUTPUTFORMATSection(t *testing.T) {
	input := `# ROLE
Engineer

# OUTPUT FORMAT
1) PLAN — 변경 요약
2) DIFF — 경로별 unified diff`
	parser := NewParser()

	result, err := parser.Parse(input)
	require.NoError(t, err)
	assert.Equal(t, 2, len(result.OutputFormat))
	assert.Contains(t, result.OutputFormat, "1) PLAN — 변경 요약")
	assert.InDelta(t, 1.0, result.OutputFormatConfidence, 0.0001)
}

func TestParsePrompt_HandlesCompletePrompt(t *testing.T) {
	input := `# ROLE
Senior full-stack engineer. Framework-agnostic. 모르면 추정 금지.

# INPUTS
- stack_profile: YAML
- task: 한 줄

# INVARIANTS
- 계약 우선(OpenAPI/gRPC/DB 스키마)
- allowed_packages 안에서만 선택
- 숨은 전제·임의 파일 생성 금지

# OUTPUT FORMAT
1) PLAN — 변경 요약
2) DIFF — 경로별 unified diff
3) RUN — 재현·검증 명령어(선택된 스택별)
4) ROLLBACK — 되돌리기
5) CHECKS — 테스트·린트·보안 점검`
	parser := NewParser()

	result, err := parser.Parse(input)
	require.NoError(t, err)

	assert.Equal(t, "Senior full-stack engineer. Framework-agnostic. 모르면 추정 금지.", result.Role)
	assert.Equal(t, 2, len(result.Inputs))
	assert.Equal(t, 3, len(result.Invariants))
	assert.Equal(t, 5, len(result.OutputFormat))
	assert.InDelta(t, 1.0, result.RoleConfidence, 0.0001)
	assert.InDelta(t, 1.0, result.InputsConfidence, 0.0001)
	assert.InDelta(t, 1.0, result.InvariantsConfidence, 0.0001)
	assert.InDelta(t, 1.0, result.OutputFormatConfidence, 0.0001)
}

func TestParsePrompt_PreservesRawContent(t *testing.T) {
	input := "# ROLE\nTest"
	parser := NewParser()

	result, err := parser.Parse(input)
	require.NoError(t, err)
	assert.Equal(t, input, result.Raw)
}

func TestParsePrompt_FuzzyMatchesTypos(t *testing.T) {
	input := `# ROEL
Engineer

# INVARIANTZ
- keep safety rules

# OUTPUT FORMTA
1) Summary

# INPTUS
- task: string`

	parser := NewParser()

	result, err := parser.Parse(input)
	require.NoError(t, err)

	assert.InDelta(t, 0.9, result.RoleConfidence, 0.1)
	assert.InDelta(t, 0.9, result.InputsConfidence, 0.1)
	assert.InDelta(t, 0.9, result.InvariantsConfidence, 0.1)
	assert.InDelta(t, 0.9, result.OutputFormatConfidence, 0.1)
}

func TestParsePrompt_DetectsKoreanAliasesWithoutHeaders(t *testing.T) {
	input := `당신은 시니어 엔지니어입니다.

입력:
- stack_profile: YAML
- task: string

제약:
- 외부 라이브러리 금지

출력 형식:
- plan: 한 줄 요약
- diff: 변경 사항`

	parser := NewParser()

	result, err := parser.Parse(input)
	require.NoError(t, err)

	assert.True(t, result.RoleConfidence > 0.5)
	assert.True(t, result.InputsConfidence > 0.5)
	assert.True(t, result.InvariantsConfidence > 0.5)
	assert.True(t, result.OutputFormatConfidence > 0.5)
	assert.Contains(t, result.Inputs, "stack_profile: YAML")
	assert.Contains(t, result.Invariants, "외부 라이브러리 금지")
	assert.Len(t, result.OutputFormat, 2)
}

func TestParsePrompt_FromFile(t *testing.T) {
	// 실제 파일로 테스트
	filePath := filepath.Join("..", "..", "test", "fixtures", "prompts", "sample.md")
	content, err := os.ReadFile(filePath)
	require.NoError(t, err)

	parser := NewParser()
	result, err := parser.Parse(string(content))
	require.NoError(t, err)

	assert.Equal(t, "Senior full-stack engineer. Framework-agnostic. 모르면 추정 금지.", result.Role)
	assert.Equal(t, 2, len(result.Inputs))
	assert.Equal(t, 3, len(result.Invariants))
	assert.Equal(t, 5, len(result.OutputFormat))
}

func TestParsePrompt_HandlesEmptySections(t *testing.T) {
	input := `# ROLE
Engineer

# INPUTS

# INVARIANTS
- rule1`
	parser := NewParser()

	result, err := parser.Parse(input)
	require.NoError(t, err)
	assert.Equal(t, "Engineer", result.Role)
	assert.Equal(t, 0, len(result.Inputs))
	assert.Equal(t, 1, len(result.Invariants))
}
