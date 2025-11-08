package analyzer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/curogom/curompt/internal/parser"
)

func TestAnalyzer_DetectsRoleSection(t *testing.T) {
	// Red: 실패하는 테스트 작성
	prompt := &parser.Prompt{
		Role: "Senior engineer",
	}
	analyzer := NewAnalyzer()

	result := analyzer.Analyze(prompt)

	assert.True(t, result.HasRole)
	assert.False(t, result.HasInputs)
}

func TestAnalyzer_DetectsAllSections(t *testing.T) {
	prompt := &parser.Prompt{
		Role:         "Engineer",
		Inputs:       []string{"task: string"},
		Invariants:   []string{"rule1"},
		OutputFormat: []string{"format1"},
	}
	analyzer := NewAnalyzer()

	result := analyzer.Analyze(prompt)

	assert.True(t, result.HasRole)
	assert.True(t, result.HasInputs)
	assert.True(t, result.HasInvariants)
	assert.True(t, result.HasOutputFormat)
	assert.Equal(t, 4, result.SectionCount)
}

func TestAnalyzer_DetectsMissingSections(t *testing.T) {
	prompt := &parser.Prompt{
		Role: "Engineer",
		// INPUTS, INVARIANTS, OUTPUT FORMAT 없음
	}
	analyzer := NewAnalyzer()

	result := analyzer.Analyze(prompt)

	assert.True(t, result.HasRole)
	assert.False(t, result.HasInputs)
	assert.False(t, result.HasInvariants)
	assert.False(t, result.HasOutputFormat)
	assert.Equal(t, 1, result.SectionCount)

	// 필수 섹션 체크 (ROLE은 필수이므로 없으면 MissingSections에 포함)
	// 여기서는 ROLE이 있으므로 MissingSections는 비어있어야 함
	assert.Equal(t, 0, len(result.MissingSections))
}

func TestAnalyzer_DetectsDuplicateRules(t *testing.T) {
	prompt := &parser.Prompt{
		Role: "Engineer",
		Invariants: []string{
			"계약 우선",
			"allowed_packages 안에서만 선택",
			"계약 우선", // 중복
		},
	}
	analyzer := NewAnalyzer()

	result := analyzer.Analyze(prompt)

	assert.Equal(t, 1, len(result.DuplicateRules))
	assert.Contains(t, result.DuplicateRules, "계약 우선")
}

func TestAnalyzer_DetectsMultipleDuplicates(t *testing.T) {
	prompt := &parser.Prompt{
		Role: "Engineer",
		Invariants: []string{
			"rule1",
			"rule2",
			"rule1", // 중복
			"rule3",
			"rule2", // 중복
		},
	}
	analyzer := NewAnalyzer()

	result := analyzer.Analyze(prompt)

	assert.Equal(t, 2, len(result.DuplicateRules))
	assert.Contains(t, result.DuplicateRules, "rule1")
	assert.Contains(t, result.DuplicateRules, "rule2")
}

func TestAnalyzer_NoDuplicatesWhenAllUnique(t *testing.T) {
	prompt := &parser.Prompt{
		Role: "Engineer",
		Invariants: []string{
			"rule1",
			"rule2",
			"rule3",
		},
	}
	analyzer := NewAnalyzer()

	result := analyzer.Analyze(prompt)

	assert.Equal(t, 0, len(result.DuplicateRules))
}

func TestAnalyzer_CountsSectionsCorrectly(t *testing.T) {
	prompt := &parser.Prompt{
		Role:         "Engineer",
		Inputs:       []string{"input1"},
		Invariants:   []string{"rule1"},
		OutputFormat: []string{"format1"},
	}
	analyzer := NewAnalyzer()

	result := analyzer.Analyze(prompt)

	assert.Equal(t, 4, result.SectionCount)
}
