package scorer

import (
	"strings"

	"github.com/curogom/curompt/internal/analyzer"
	"github.com/curogom/curompt/internal/parser"
)

// StructureMetricCalculator calculates structure score
type StructureMetricCalculator struct {
	prompt   *parser.Prompt
	analysis *analyzer.AnalysisResult
}

// NewStructureMetricCalculator creates a new structure metric calculator
func NewStructureMetricCalculator(prompt *parser.Prompt, analysis *analyzer.AnalysisResult) MetricCalculator {
	return &StructureMetricCalculator{
		prompt:   prompt,
		analysis: analysis,
	}
}

// Calculate calculates structure score (0-100)
// 점수 기준:
// - 필수 섹션 존재: ROLE (50점)
// - 추가 섹션 존재: INPUTS, INVARIANTS, OUTPUT FORMAT (각 10점)
// - 중복 규칙 감점: 중복당 -5점
func (c *StructureMetricCalculator) Calculate() (float64, error) {
	score := 0.0

	// 필수 섹션 체크
	if c.analysis.HasRole {
		score += 50.0
	}

	// 추가 섹션 체크
	if c.analysis.HasInputs {
		score += 10.0
	}
	if c.analysis.HasInvariants {
		score += 10.0
	}
	if c.analysis.HasOutputFormat {
		score += 10.0
	}

	// 중복 규칙 감점
	deduction := float64(len(c.analysis.DuplicateRules)) * 5.0
	score -= deduction

	// 최소 0점, 최대 100점
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score, nil
}

// ConcisenessMetricCalculator calculates conciseness score
type ConcisenessMetricCalculator struct {
	prompt    *parser.Prompt
	tokenizer parser.Tokenizer
}

// NewConcisenessMetricCalculator creates a new conciseness metric calculator
func NewConcisenessMetricCalculator(prompt *parser.Prompt, tokenizer parser.Tokenizer) MetricCalculator {
	return &ConcisenessMetricCalculator{
		prompt:    prompt,
		tokenizer: tokenizer,
	}
}

// Calculate calculates conciseness score (0-100)
// 점수 기준:
// - 토큰 수가 적을수록 높은 점수
// - 기준: 0-500 토큰 = 100점, 500-1000 = 80점, 1000-2000 = 60점, 2000+ = 40점
func (c *ConcisenessMetricCalculator) Calculate() (float64, error) {
	tokenCount, err := c.tokenizer.CountTokens(c.prompt.Raw)
	if err != nil {
		return 0, err
	}

	score := 0.0
	switch {
	case tokenCount <= 500:
		score = 100.0
	case tokenCount <= 1000:
		// 500-1000: 100점에서 80점으로 선형 감소
		score = 100.0 - ((float64(tokenCount)-500.0)/500.0)*20.0
	case tokenCount <= 2000:
		// 1000-2000: 80점에서 60점으로 선형 감소
		score = 80.0 - ((float64(tokenCount)-1000.0)/1000.0)*20.0
	default:
		// 2000+: 60점에서 40점으로 선형 감소 (최대 5000 토큰 기준)
		if tokenCount > 5000 {
			score = 40.0
		} else {
			score = 60.0 - ((float64(tokenCount)-2000.0)/3000.0)*20.0
		}
	}

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score, nil
}

// RiskMetricCalculator calculates risk score (higher is better, lower risk = higher score)
type RiskMetricCalculator struct {
	prompt *parser.Prompt
}

// NewRiskMetricCalculator creates a new risk metric calculator
func NewRiskMetricCalculator(prompt *parser.Prompt) MetricCalculator {
	return &RiskMetricCalculator{
		prompt: prompt,
	}
}

// Calculate calculates risk score (0-100)
// 점수 기준:
// - 기본 100점에서 시작
// - 모호한 표현 감지 시 감점
// - 민감 데이터 패턴 감지 시 감점
func (c *RiskMetricCalculator) Calculate() (float64, error) {
	score := 100.0

	// 모호한 표현 감지 (예: "적절히", "필요시", "가능하면")
	ambiguousPatterns := []string{
		"적절히", "필요시", "가능하면", "적당히", "보통", "일반적으로",
		"maybe", "perhaps", "possibly", "might", "could",
	}

	text := c.prompt.Raw
	for _, pattern := range ambiguousPatterns {
		if contains(text, pattern) {
			score -= 5.0 // 각 모호한 표현마다 -5점
		}
	}

	// 최소 0점
	if score < 0 {
		score = 0
	}

	return score, nil
}

// contains checks if text contains substring (case-insensitive)
func contains(text, substr string) bool {
	return strings.Contains(strings.ToLower(text), strings.ToLower(substr))
}
