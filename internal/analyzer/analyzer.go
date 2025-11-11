package analyzer

import (
	"github.com/curogom/curompt/internal/parser"
)

// staticAnalyzer implements Analyzer interface
type staticAnalyzer struct{}

// NewAnalyzer creates a new static analyzer
func NewAnalyzer() Analyzer {
	return &staticAnalyzer{}
}

const sectionConfidenceThreshold = 0.5

// Analyze performs static analysis on a prompt
func (a *staticAnalyzer) Analyze(prompt *parser.Prompt) *AnalysisResult {
	roleConfidence := clampConfidence(prompt.RoleConfidence)
	if roleConfidence == 0 && prompt.Role != "" {
		roleConfidence = 1.0
	}

	inputConfidence := clampConfidence(prompt.InputsConfidence)
	if inputConfidence == 0 && len(prompt.Inputs) > 0 {
		inputConfidence = 1.0
	}

	invariantConfidence := clampConfidence(prompt.InvariantsConfidence)
	if invariantConfidence == 0 && len(prompt.Invariants) > 0 {
		invariantConfidence = 1.0
	}

	outputConfidence := clampConfidence(prompt.OutputFormatConfidence)
	if outputConfidence == 0 && len(prompt.OutputFormat) > 0 {
		outputConfidence = 1.0
	}

	result := &AnalysisResult{
		DuplicateRules:         []string{},
		MissingSections:        []string{},
		RoleConfidence:         roleConfidence,
		InputsConfidence:       inputConfidence,
		InvariantsConfidence:   invariantConfidence,
		OutputFormatConfidence: outputConfidence,
	}

	// 섹션 존재 여부 체크 (확신도 기반)
	result.HasRole = result.RoleConfidence >= sectionConfidenceThreshold
	result.HasInputs = result.InputsConfidence >= sectionConfidenceThreshold
	result.HasInvariants = result.InvariantsConfidence >= sectionConfidenceThreshold
	result.HasOutputFormat = result.OutputFormatConfidence >= sectionConfidenceThreshold

	// 섹션 수 계산
	if result.HasRole {
		result.SectionCount++
	}
	if result.HasInputs {
		result.SectionCount++
	}
	if result.HasInvariants {
		result.SectionCount++
	}
	if result.HasOutputFormat {
		result.SectionCount++
	}

	// 필수 섹션 체크 (ROLE은 필수)
	if !result.HasRole {
		result.MissingSections = append(result.MissingSections, "ROLE")
	}

	// 중복 규칙 감지
	result.DuplicateRules = a.detectDuplicates(prompt.Invariants)

	return result
}

func clampConfidence(v float64) float64 {
	switch {
	case v < 0:
		return 0
	case v > 1:
		return 1
	default:
		return v
	}
}

// detectDuplicates finds duplicate items in a slice
func (a *staticAnalyzer) detectDuplicates(items []string) []string {
	seen := make(map[string]int)
	duplicates := []string{}

	for _, item := range items {
		seen[item]++
		if seen[item] == 2 {
			// 처음 중복이 발견될 때 추가
			duplicates = append(duplicates, item)
		}
	}

	return duplicates
}
