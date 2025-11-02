package analyzer

import (
	"github.com/curogom/curo-prompt/internal/parser"
)

// staticAnalyzer implements Analyzer interface
type staticAnalyzer struct{}

// NewAnalyzer creates a new static analyzer
func NewAnalyzer() Analyzer {
	return &staticAnalyzer{}
}

// Analyze performs static analysis on a prompt
func (a *staticAnalyzer) Analyze(prompt *parser.Prompt) *AnalysisResult {
	result := &AnalysisResult{
		DuplicateRules:  []string{},
		MissingSections: []string{},
	}

	// 섹션 존재 여부 체크
	result.HasRole = prompt.Role != ""
	result.HasInputs = len(prompt.Inputs) > 0
	result.HasInvariants = len(prompt.Invariants) > 0
	result.HasOutputFormat = len(prompt.OutputFormat) > 0

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
