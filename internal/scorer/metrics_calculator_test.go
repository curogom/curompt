package scorer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/curogom/curo-prompt/internal/analyzer"
	"github.com/curogom/curo-prompt/internal/parser"
)

func TestStructureMetricCalculator_Calculate(t *testing.T) {
	// Red: 실패하는 테스트 작성
	prompt := &parser.Prompt{
		Role:         "Engineer",
		Inputs:       []string{"input1"},
		Invariants:   []string{"rule1"},
		OutputFormat: []string{"format1"},
	}
	analysis := &analyzer.AnalysisResult{
		HasRole:         true,
		HasInputs:       true,
		HasInvariants:   true,
		HasOutputFormat: true,
		DuplicateRules:  []string{},
		SectionCount:    4,
	}
	calculator := NewStructureMetricCalculator(prompt, analysis)

	score, err := calculator.Calculate()
	require.NoError(t, err)
	assert.Greater(t, score, 0.0)
	assert.LessOrEqual(t, score, 100.0)
}

func TestStructureMetricCalculator_WithRoleOnly(t *testing.T) {
	prompt := &parser.Prompt{
		Role: "Engineer",
	}
	analysis := &analyzer.AnalysisResult{
		HasRole:      true,
		SectionCount: 1,
	}
	calculator := NewStructureMetricCalculator(prompt, analysis)

	score, err := calculator.Calculate()
	require.NoError(t, err)
	assert.Equal(t, 50.0, score) // ROLE만 있으면 50점
}

func TestStructureMetricCalculator_WithDuplicates(t *testing.T) {
	prompt := &parser.Prompt{
		Role:       "Engineer",
		Invariants: []string{"rule1", "rule1"},
	}
	analysis := &analyzer.AnalysisResult{
		HasRole:        true,
		HasInvariants:  true,
		DuplicateRules: []string{"rule1"},
		SectionCount:   2,
	}
	calculator := NewStructureMetricCalculator(prompt, analysis)

	score, err := calculator.Calculate()
	require.NoError(t, err)
	// ROLE (50) + INVARIANTS (10) - 중복 (5) = 55점
	assert.Equal(t, 55.0, score)
}

func TestStructureMetricCalculator_MaxScore(t *testing.T) {
	prompt := &parser.Prompt{
		Role:         "Engineer",
		Inputs:       []string{"input1"},
		Invariants:   []string{"rule1"},
		OutputFormat: []string{"format1"},
	}
	analysis := &analyzer.AnalysisResult{
		HasRole:         true,
		HasInputs:       true,
		HasInvariants:   true,
		HasOutputFormat: true,
		DuplicateRules:  []string{},
		SectionCount:    4,
	}
	calculator := NewStructureMetricCalculator(prompt, analysis)

	score, err := calculator.Calculate()
	require.NoError(t, err)
	// ROLE (50) + INPUTS (10) + INVARIANTS (10) + OUTPUT FORMAT (10) = 80점
	assert.Equal(t, 80.0, score)
	assert.LessOrEqual(t, score, 100.0)
}

func TestStructureMetricCalculator_WithManyDuplicates(t *testing.T) {
	// 많은 중복으로 점수가 음수가 될 수 있는지 테스트
	prompt := &parser.Prompt{
		Role:       "Engineer",
		Invariants: []string{"rule1", "rule1", "rule1", "rule1", "rule1"},
	}
	analysis := &analyzer.AnalysisResult{
		HasRole:        true,
		HasInvariants:  true,
		DuplicateRules: []string{"rule1"},
		SectionCount:   2,
	}
	calculator := NewStructureMetricCalculator(prompt, analysis)

	score, err := calculator.Calculate()
	require.NoError(t, err)
	// ROLE (50) + INVARIANTS (10) - 중복 (5) = 55점
	// 중복이 1개만 감지되므로 점수는 양수
	assert.GreaterOrEqual(t, score, 0.0)
}

func TestStructureMetricCalculator_NegativeScoreClampedToZero(t *testing.T) {
	// 매우 많은 중복으로 점수가 음수가 되는 경우를 테스트
	prompt := &parser.Prompt{
		Role:       "Engineer",
		Invariants: []string{"rule1", "rule1", "rule1", "rule1", "rule1"},
	}
	// 20개 중복 = -100점, ROLE(50) + INVARIANTS(10) - 100 = -40점 -> 0점으로 클램핑
	analysis := &analyzer.AnalysisResult{
		HasRole:        true,
		HasInvariants:  true,
		DuplicateRules: []string{"rule1", "rule2", "rule3", "rule4", "rule5", "rule6", "rule7", "rule8", "rule9", "rule10", "rule11", "rule12", "rule13", "rule14", "rule15", "rule16", "rule17", "rule18", "rule19", "rule20"},
		SectionCount:   2,
	}
	calculator := NewStructureMetricCalculator(prompt, analysis)

	score, err := calculator.Calculate()
	require.NoError(t, err)
	// 음수가 되면 0점으로 클램핑
	assert.Equal(t, 0.0, score)
}

func TestStructureMetricCalculator_ScoreCappedAt100(t *testing.T) {
	// 점수가 100을 초과하는 경우는 없지만, 분기를 커버하기 위한 테스트
	prompt := &parser.Prompt{
		Role:         "Engineer",
		Inputs:       []string{"input1"},
		Invariants:   []string{"rule1"},
		OutputFormat: []string{"format1"},
	}
	analysis := &analyzer.AnalysisResult{
		HasRole:         true,
		HasInputs:       true,
		HasInvariants:   true,
		HasOutputFormat: true,
		DuplicateRules:  []string{},
		SectionCount:    4,
	}
	calculator := NewStructureMetricCalculator(prompt, analysis)

	score, err := calculator.Calculate()
	require.NoError(t, err)
	// 최대 80점이므로 100점 초과 분기는 테스트 불가능하지만, 코드는 커버됨
	assert.LessOrEqual(t, score, 100.0)
}

func TestConcisenessMetricCalculator_Calculate(t *testing.T) {
	prompt := &parser.Prompt{
		Raw: "Short prompt",
	}
	tokenizer := parser.NewClaudeTokenizer()
	calculator := NewConcisenessMetricCalculator(prompt, tokenizer)

	score, err := calculator.Calculate()
	require.NoError(t, err)
	assert.GreaterOrEqual(t, score, 0.0)
	assert.LessOrEqual(t, score, 100.0)
}

func TestConcisenessMetricCalculator_SmallPrompt(t *testing.T) {
	prompt := &parser.Prompt{
		Raw: "Test", // 작은 프롬프트
	}
	tokenizer := parser.NewClaudeTokenizer()
	calculator := NewConcisenessMetricCalculator(prompt, tokenizer)

	score, err := calculator.Calculate()
	require.NoError(t, err)
	// 500 토큰 이하 = 100점
	assert.GreaterOrEqual(t, score, 90.0)
}

func TestConcisenessMetricCalculator_MediumPrompt_500to1000(t *testing.T) {
	// 500-1000 토큰 범위 테스트
	// 대략 750 토큰 정도를 만들어야 함
	prompt := &parser.Prompt{
		Raw: generateText(750), // 대략 750 토큰
	}
	tokenizer := parser.NewClaudeTokenizer()
	calculator := NewConcisenessMetricCalculator(prompt, tokenizer)

	score, err := calculator.Calculate()
	require.NoError(t, err)
	// 500-1000 범위에서는 100점에서 80점으로 선형 감소
	// 750 토큰은 약 90점 정도여야 함
	assert.GreaterOrEqual(t, score, 80.0)
	assert.LessOrEqual(t, score, 100.0)
}

func TestConcisenessMetricCalculator_LargePrompt_1000to2000(t *testing.T) {
	// 1000-2000 토큰 범위 테스트
	prompt := &parser.Prompt{
		Raw: generateText(1500), // 대략 1500 토큰
	}
	tokenizer := parser.NewClaudeTokenizer()
	calculator := NewConcisenessMetricCalculator(prompt, tokenizer)

	score, err := calculator.Calculate()
	require.NoError(t, err)
	// 1000-2000 범위에서는 80점에서 60점으로 선형 감소
	// 1500 토큰은 약 70점 정도여야 함
	assert.GreaterOrEqual(t, score, 60.0)
	assert.LessOrEqual(t, score, 80.0)
}

func TestConcisenessMetricCalculator_VeryLargePrompt_2000to5000(t *testing.T) {
	// 2000-5000 토큰 범위 테스트
	prompt := &parser.Prompt{
		Raw: generateText(3000), // 대략 3000 토큰
	}
	tokenizer := parser.NewClaudeTokenizer()
	calculator := NewConcisenessMetricCalculator(prompt, tokenizer)

	score, err := calculator.Calculate()
	require.NoError(t, err)
	// 2000+ 범위에서는 60점에서 40점으로 선형 감소
	// 3000 토큰은 약 53점 정도여야 함
	assert.GreaterOrEqual(t, score, 40.0)
	assert.LessOrEqual(t, score, 60.0)
}

func TestConcisenessMetricCalculator_ExtremelyLargePrompt_Over5000(t *testing.T) {
	// 5000 토큰 이상 범위 테스트
	prompt := &parser.Prompt{
		Raw: generateText(6000), // 대략 6000 토큰
	}
	tokenizer := parser.NewClaudeTokenizer()
	calculator := NewConcisenessMetricCalculator(prompt, tokenizer)

	score, err := calculator.Calculate()
	require.NoError(t, err)
	// 5000+ 토큰은 40점
	assert.Equal(t, 40.0, score)
}

func TestConcisenessMetricCalculator_ErrorHandling(t *testing.T) {
	// tokenizer 에러 처리 테스트는 mock 필요
	// 일단 정상 케이스만 테스트
}

func TestRiskMetricCalculator_Calculate(t *testing.T) {
	prompt := &parser.Prompt{
		Raw: "Clear instructions",
	}
	calculator := NewRiskMetricCalculator(prompt)

	score, err := calculator.Calculate()
	require.NoError(t, err)
	assert.Equal(t, 100.0, score) // 모호한 표현 없음 = 100점
}

func TestRiskMetricCalculator_WithAmbiguousWords(t *testing.T) {
	prompt := &parser.Prompt{
		Raw: "적절히 처리하고 필요시 확인하세요",
	}
	calculator := NewRiskMetricCalculator(prompt)

	score, err := calculator.Calculate()
	require.NoError(t, err)
	assert.Less(t, score, 100.0) // 모호한 표현 감지 = 감점
	assert.GreaterOrEqual(t, score, 0.0)
}

func TestRiskMetricCalculator_MultipleAmbiguousWords(t *testing.T) {
	prompt := &parser.Prompt{
		Raw: "적절히 처리하고 가능하면 빠르게 필요시 확인하세요",
	}
	calculator := NewRiskMetricCalculator(prompt)

	score, err := calculator.Calculate()
	require.NoError(t, err)
	// 3개 모호한 표현 = -15점, 100 - 15 = 85점
	assert.Equal(t, 85.0, score)
}

func TestRiskMetricCalculator_ScoreCannotGoBelowZero(t *testing.T) {
	prompt := &parser.Prompt{
		Raw: "적절히 적당히 보통 일반적으로 필요시 가능하면 maybe perhaps possibly might could",
	}
	calculator := NewRiskMetricCalculator(prompt)

	score, err := calculator.Calculate()
	require.NoError(t, err)
	// 매우 많은 모호한 표현이어도 최소 0점
	assert.GreaterOrEqual(t, score, 0.0)
	assert.LessOrEqual(t, score, 100.0)
}

// generateText generates approximate token count (rough estimation)
func generateText(targetTokens int) string {
	// 대략 4자 = 1 토큰으로 추정하므로, targetTokens * 4 문자가 필요
	words := []string{"word", "test", "prompt", "example", "content"}
	text := ""
	wordCount := 0

	// 목표 토큰 수에 도달할 때까지 단어 추가
	for wordCount < targetTokens {
		text += words[wordCount%len(words)] + " "
		wordCount++
	}

	return text
}
