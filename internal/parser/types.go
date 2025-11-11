package parser

// Prompt represents a parsed prompt structure
type Prompt struct {
	Role                   string   // ROLE 섹션 내용
	Inputs                 []string // INPUTS 섹션 항목들
	Invariants             []string // INVARIANTS 섹션 규칙들
	OutputFormat           []string // OUTPUT FORMAT 섹션 항목들
	RoleConfidence         float64  // ROLE 섹션 확신도 (0-1)
	InputsConfidence       float64  // INPUTS 섹션 확신도 (0-1)
	InvariantsConfidence   float64  // INVARIANTS 섹션 확신도 (0-1)
	OutputFormatConfidence float64  // OUTPUT FORMAT 섹션 확신도 (0-1)
	Raw                    string   // 원본 프롬프트 텍스트
}

// Parser interface for parsing prompts
type Parser interface {
	Parse(content string) (*Prompt, error)
}
