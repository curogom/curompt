package analyzer

import "github.com/curogom/curompt/internal/parser"

// AnalysisResult contains the result of static analysis
type AnalysisResult struct {
	HasRole               bool     // ROLE 섹션 존재 여부
	HasInputs             bool     // INPUTS 섹션 존재 여부
	HasInvariants         bool     // INVARIANTS 섹션 존재 여부
	HasOutputFormat       bool     // OUTPUT FORMAT 섹션 존재 여부
	RoleConfidence        float64  // ROLE 섹션 확신도 (0-1)
	InputsConfidence      float64  // INPUTS 섹션 확신도 (0-1)
	InvariantsConfidence  float64  // INVARIANTS 섹션 확신도 (0-1)
	OutputFormatConfidence float64 // OUTPUT FORMAT 섹션 확신도 (0-1)
	DuplicateRules        []string // 중복된 규칙 목록
	MissingSections       []string // 누락된 필수 섹션 목록
	SectionCount          int      // 전체 섹션 수
}

// Analyzer performs static analysis on prompts
type Analyzer interface {
	Analyze(prompt *parser.Prompt) *AnalysisResult
}
