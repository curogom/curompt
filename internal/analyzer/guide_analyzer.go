package analyzer

import (
	"fmt"
	"strings"

	"github.com/curogom/curo-prompt/internal/parser"
)

// GuideComplianceResult represents compliance with prompt engineering guides
type GuideComplianceResult struct {
	ClarityScore     float64  // 명확성 점수 (0-100)
	ContextScore     float64  // 컨텍스트 제공 점수 (0-100)
	ExampleScore     float64  // 예시 제공 점수 (0-100)
	StructureScore   float64  // 구조화 점수 (0-100)
	ConstraintScore  float64  // 제약 명시 점수 (0-100)
	Suggestions      []string // 구체적인 개선 제안
	AmbiguousPhrases []string // 감지된 모호한 표현
	MissingElements  []string // 누락된 요소
}

// GuideAnalyzer analyzes prompt compliance with Claude/ChatGPT guides
type GuideAnalyzer struct {
	prompt *parser.Prompt
}

// NewGuideAnalyzer creates a new guide analyzer
func NewGuideAnalyzer(prompt *parser.Prompt) *GuideAnalyzer {
	return &GuideAnalyzer{
		prompt: prompt,
	}
}

// Analyze analyzes guide compliance
func (g *GuideAnalyzer) Analyze() GuideComplianceResult {
	result := GuideComplianceResult{
		Suggestions:      []string{},
		AmbiguousPhrases: []string{},
		MissingElements:  []string{},
	}

	text := strings.ToLower(g.prompt.Raw)

	// 1. 명확성 분석
	result.ClarityScore = g.analyzeClarity(text, &result)

	// 2. 컨텍스트 분석
	result.ContextScore = g.analyzeContext(text)

	// 3. 예시 분석
	result.ExampleScore = g.analyzeExamples(text)

	// 4. 구조화 분석
	result.StructureScore = g.analyzeStructure()

	// 5. 제약 명시 분석
	result.ConstraintScore = g.analyzeConstraints(text)

	// 개선 제안 생성
	result.generateSuggestions()

	return result
}

// analyzeClarity checks for clarity (detects ambiguous phrases)
func (g *GuideAnalyzer) analyzeClarity(text string, result *GuideComplianceResult) float64 {
	score := 100.0

	// 모호한 표현 패턴
	ambiguousPatterns := map[string]float64{
		"적절히": -10, "필요시": -10, "가능하면": -10, "적당히": -10,
		"보통": -8, "일반적으로": -8, "대략": -10,
		"maybe": -10, "perhaps": -10, "possibly": -10,
		"might": -8, "could": -8, "sort of": -8,
		"아마도": -10, "혹시": -10, "아마": -10,
	}

	// 구체적인 동사 패턴 (가산)
	clearVerbs := []string{
		"생성", "수정", "삭제", "추가", "변경", "확인", "검증",
		"create", "modify", "delete", "add", "change", "verify",
		"update", "remove", "check", "analyze", "implement",
	}

	// 모호한 표현 감지
	var ambiguousContexts []string
	for phrase, penalty := range ambiguousPatterns {
		if strings.Contains(text, phrase) {
			results := g.findPhraseContext(text, phrase)
			ambiguousContexts = append(ambiguousContexts, results...)
			score += penalty
		}
	}
	result.AmbiguousPhrases = ambiguousContexts

	// 구체적인 동사 사용 확인
	verbCount := 0
	for _, verb := range clearVerbs {
		if strings.Contains(text, verb) {
			verbCount++
		}
	}
	if verbCount > 0 {
		score += float64(verbCount) * 2.0 // 동사당 +2점 (최대 +20점)
	}

	// 점수 범위 제한
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

// analyzeContext checks for context provision
func (g *GuideAnalyzer) analyzeContext(text string) float64 {
	score := 0.0

	// 컨텍스트 관련 키워드
	contextIndicators := []string{
		"프로젝트", "프로젝트", "프로젝트 경로", "파일 경로",
		"프레임워크", "라이브러리", "버전", "환경",
		"project", "framework", "library", "version",
		"context", "background", "environment",
		"기존", "이전", "관련", "참고",
	}

	count := 0
	for _, indicator := range contextIndicators {
		if strings.Contains(text, indicator) {
			count++
		}
	}

	// 컨텍스트 제공 점수: 키워드당 10점 (최대 100점)
	score = float64(count) * 10.0
	if score > 100 {
		score = 100
	}

	return score
}

// analyzeExamples checks for example provision
func (g *GuideAnalyzer) analyzeExamples(text string) float64 {
	score := 0.0

	// 예시 관련 키워드
	exampleIndicators := []string{
		"예시", "예", "샘플", "sample", "example",
		"입력:", "출력:", "input:", "output:",
		"다음과 같이", "형식:", "format:",
	}

	count := 0
	for _, indicator := range exampleIndicators {
		if strings.Contains(text, indicator) {
			count++
		}
	}

	// 예시 제공 점수: 키워드당 15점 (최대 100점)
	score = float64(count) * 15.0
	if score > 100 {
		score = 100
	}

	return score
}

// analyzeStructure checks for structured format
func (g *GuideAnalyzer) analyzeStructure() float64 {
	score := 0.0

	// 구조화 지표
	if g.prompt.Role != "" {
		score += 25.0 // ROLE 섹션
	}
	if len(g.prompt.Inputs) > 0 {
		score += 25.0 // INPUTS 섹션
	}
	if len(g.prompt.Invariants) > 0 {
		score += 25.0 // INVARIANTS 섹션
	}
	if len(g.prompt.OutputFormat) > 0 {
		score += 25.0 // OUTPUT FORMAT 섹션
	}

	// 단계별 지시 확인 (번호나 순서 표시)
	stepIndicators := []string{"1.", "2.", "3.", "step", "단계"}
	text := strings.ToLower(g.prompt.Raw)
	stepCount := 0
	for _, indicator := range stepIndicators {
		if strings.Contains(text, indicator) {
			stepCount++
		}
	}
	if stepCount > 0 {
		score += float64(stepCount) * 5.0 // 단계당 +5점 (최대 +20점)
	}

	if score > 100 {
		score = 100
	}

	return score
}

// analyzeConstraints checks for constraint specification
func (g *GuideAnalyzer) analyzeConstraints(text string) float64 {
	score := 0.0

	// 제약 관련 키워드
	constraintIndicators := []string{
		"제약", "제한", "제외", "포함하지 않음",
		"constraint", "limit", "exclude", "not include",
		"사용하지 않음", "사용 금지", "필수", "필수 아님",
		"don't use", "must not", "required", "not required",
	}

	count := 0
	for _, indicator := range constraintIndicators {
		if strings.Contains(text, indicator) {
			count++
		}
	}

	// 제약 명시 점수: 키워드당 12점 (최대 100점)
	score = float64(count) * 12.0
	if score > 100 {
		score = 100
	}

	return score
}

// findPhraseContext finds the context around a phrase
func (g *GuideAnalyzer) findPhraseContext(text, phrase string) []string {
	var results []string
	lowerText := strings.ToLower(text)
	lowerPhrase := strings.ToLower(phrase)

	index := strings.Index(lowerText, lowerPhrase)
	if index >= 0 {
		start := index - 20
		if start < 0 {
			start = 0
		}
		end := index + len(phrase) + 20
		if end > len(text) {
			end = len(text)
		}
		context := text[start:end]
		results = append(results, context)
	}

	return results
}

// generateSuggestions generates specific improvement suggestions
func (r *GuideComplianceResult) generateSuggestions() {
	// 명확성 개선 제안
	if r.ClarityScore < 70 {
		r.Suggestions = append(r.Suggestions, "💡 **명확성 개선**: 모호한 표현을 구체적인 동사로 변경하세요 (예: '생성', '수정', '삭제')")
		if len(r.AmbiguousPhrases) > 0 {
			r.MissingElements = append(r.MissingElements, fmt.Sprintf("모호한 표현 %d개 발견", len(r.AmbiguousPhrases)))
		}
	}

	// 컨텍스트 추가 제안
	if r.ContextScore < 50 {
		r.Suggestions = append(r.Suggestions, "💡 **컨텍스트 추가**: 프로젝트 정보, 파일 경로, 프레임워크 등의 배경 정보를 포함하세요")
		r.MissingElements = append(r.MissingElements, "프로젝트 컨텍스트 정보")
	}

	// 예시 추가 제안
	if r.ExampleScore < 50 {
		r.Suggestions = append(r.Suggestions, "💡 **예시 제공**: 입력/출력 예시나 원하는 형식의 샘플을 포함하세요 (예: '입력: ...', '출력: ...')")
		r.MissingElements = append(r.MissingElements, "입력/출력 예시")
	}

	// 구조화 개선 제안
	if r.StructureScore < 60 {
		r.Suggestions = append(r.Suggestions, "💡 **구조화 개선**: 작업을 단계로 나누거나 구조화된 섹션(ROLE, INPUTS, INVARIANTS)을 사용하세요")
		if r.StructureScore < 40 {
			r.MissingElements = append(r.MissingElements, "구조화된 섹션")
		}
	}

	// 제약 명시 제안
	if r.ConstraintScore < 30 {
		r.Suggestions = append(r.Suggestions, "💡 **제약 조건 명시**: 특정 라이브러리, 프레임워크 제약이나 제외 사항을 명확히 표시하세요")
		r.MissingElements = append(r.MissingElements, "제약 조건 명시")
	}
}
