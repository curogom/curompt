package evaluator

import (
	"context"

	"github.com/curogom/curompt/internal/analyzer"
	"github.com/curogom/curompt/internal/model"
	"github.com/curogom/curompt/internal/parser"
	"github.com/curogom/curompt/internal/scorer"
)

// EvaluationResult contains the complete evaluation result
type EvaluationResult struct {
	Prompt          *model.CollectedPrompt
	Analysis        *analyzer.AnalysisResult
	Score           *scorer.ScoreResult
	TokenCount      int
	TokenProvider   string                          // claude, openai 등
	GuideCompliance *analyzer.GuideComplianceResult // 가이드 준수 분석 결과
}

// Evaluator evaluates collected prompts
type Evaluator struct {
	parser        parser.Parser
	analyzer      analyzer.Analyzer
	scorer        scorer.Scorer
	tokenizer     parser.Tokenizer
	tokenProvider string
}

// NewEvaluator creates a new evaluator
func NewEvaluator(tokenProvider string) *Evaluator {
	var tokenizer parser.Tokenizer
	switch tokenProvider {
	case "openai":
		tokenizer = parser.NewOpenAITokenizer()
	default:
		tokenizer = parser.NewClaudeTokenizer()
	}

	return &Evaluator{
		parser:        parser.NewParser(),
		analyzer:      analyzer.NewAnalyzer(),
		scorer:        scorer.NewScorer(),
		tokenizer:     tokenizer,
		tokenProvider: tokenProvider,
	}
}

// Evaluate evaluates a collected prompt
func (e *Evaluator) Evaluate(ctx context.Context, collectedPrompt *model.CollectedPrompt) (*EvaluationResult, error) {
	// 프롬프트 파싱 (이미 파싱되어 있으면 재사용, 없으면 파싱)
	var parsedPrompt *parser.Prompt
	var err error

	if collectedPrompt.Prompt != nil {
		parsedPrompt = collectedPrompt.Prompt
	} else {
		parsedPrompt, err = e.parser.Parse(collectedPrompt.RawPrompt)
		if err != nil {
			return nil, err
		}
	}

	// 정적 분석
	analysis := e.analyzer.Analyze(parsedPrompt)

	// 가이드 준수 분석
	guideAnalyzer := analyzer.NewGuideAnalyzer(parsedPrompt)
	guideCompliance := guideAnalyzer.Analyze()

	// 토큰 계산
	tokenCount, err := e.tokenizer.CountTokens(collectedPrompt.RawPrompt)
	if err != nil {
		return nil, err
	}

	// 메트릭 계산
	metrics := scorer.Metrics{
		Structure:            e.calculateStructureScore(parsedPrompt, analysis),
		Conciseness:          e.calculateConcisenessScore(collectedPrompt.RawPrompt, tokenCount),
		InstructionFollowing: 80.0, // 기본값 (M2에서 실제 구현)
		SelfConsistency:      0.0,  // 동적 평가 필요 (M2)
		LatencyCost:          0.0,  // 동적 평가 필요 (M2)
		Risk:                 e.calculateRiskScore(parsedPrompt),
	}

	// 점수화
	weights := scorer.DefaultWeights()
	scoreResult := e.scorer.Score(metrics, weights)

	return &EvaluationResult{
		Prompt:          collectedPrompt,
		Analysis:        analysis,
		Score:           &scoreResult,
		TokenCount:      tokenCount,
		TokenProvider:   e.tokenProvider,
		GuideCompliance: &guideCompliance,
	}, nil
}

// calculateStructureScore calculates structure metric
func (e *Evaluator) calculateStructureScore(prompt *parser.Prompt, analysis *analyzer.AnalysisResult) float64 {
	calc := scorer.NewStructureMetricCalculator(prompt, analysis)
	score, err := calc.Calculate()
	if err != nil {
		// Return default score on error
		return 50.0
	}
	return score
}

// calculateConcisenessScore calculates conciseness metric
func (e *Evaluator) calculateConcisenessScore(rawPrompt string, tokenCount int) float64 {
	// 간단한 근사치로 계산 (실제로는 ConcisenessMetricCalculator 사용)
	prompt := &parser.Prompt{Raw: rawPrompt}
	calc := scorer.NewConcisenessMetricCalculator(prompt, e.tokenizer)
	score, err := calc.Calculate()
	if err != nil {
		// Return default score on error
		return 50.0
	}
	return score
}

// calculateRiskScore calculates risk metric
func (e *Evaluator) calculateRiskScore(prompt *parser.Prompt) float64 {
	calc := scorer.NewRiskMetricCalculator(prompt)
	score, err := calc.Calculate()
	if err != nil {
		// Return default score on error
		return 50.0
	}
	return score
}
