package reporter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/curogom/curompt/internal/analyzer"
	"github.com/curogom/curompt/internal/evaluator"
	"github.com/curogom/curompt/internal/model"
	"github.com/curogom/curompt/internal/scorer"
)

func TestMarkdownReporter_Generate(t *testing.T) {
	// Red: 실패하는 테스트 작성
	reporter := NewMarkdownReporter()

	result := &evaluator.EvaluationResult{
		Prompt: &model.CollectedPrompt{
			ID:         "test-1",
			Tool:       "codex",
			RawPrompt:  "# ROLE\nEngineer",
			Timestamp:  1234567890,
			Command:    "codex exec test",
			WorkingDir: "/tmp",
		},
		Analysis: &analyzer.AnalysisResult{
			HasRole:         true,
			HasInputs:       false,
			HasInvariants:   false,
			HasOutputFormat: false,
			SectionCount:    1,
			DuplicateRules:  []string{},
		},
		Score: &scorer.ScoreResult{
			OverallScore: 75.5,
			Metrics: scorer.Metrics{
				Structure:            80.0,
				Conciseness:          70.0,
				InstructionFollowing: 80.0,
				Risk:                 70.0,
			},
		},
		TokenCount:    100,
		TokenProvider: "claude",
	}

	report, err := reporter.Generate(result)

	require.NoError(t, err)
	assert.NotEmpty(t, report)
	assert.Contains(t, report, "프롬프트 평가 리포트")
	assert.Contains(t, report, "75.5")
	assert.Contains(t, report, "codex")
}

func TestMarkdownReporter_WithMissingSections(t *testing.T) {
	reporter := NewMarkdownReporter()

	result := &evaluator.EvaluationResult{
		Prompt: &model.CollectedPrompt{
			ID:        "test-1",
			Tool:      "codex",
			RawPrompt: "# ROLE\nEngineer",
		},
		Analysis: &analyzer.AnalysisResult{
			HasRole:         true,
			HasInputs:       false,
			HasInvariants:   false,
			HasOutputFormat: false,
			SectionCount:    1,
		},
		Score: &scorer.ScoreResult{
			OverallScore: 50.0,
			Metrics: scorer.Metrics{
				Structure: 50.0,
			},
		},
		TokenCount: 50,
	}

	report, err := reporter.Generate(result)
	require.NoError(t, err)

	// MissingSections가 없어도 리포트는 생성되어야 함
	assert.NotEmpty(t, report)
}

func TestMarkdownReporter_ShowsConfidenceWhenEnabled(t *testing.T) {
	t.Setenv("CUROMPT_SHOW_STRUCTURE_CONFIDENCE", "true")

	reporter := NewMarkdownReporter()

	result := &evaluator.EvaluationResult{
		Prompt: &model.CollectedPrompt{
			ID:   "test-2",
			Tool: "codex",
		},
		Analysis: &analyzer.AnalysisResult{
			HasRole:               true,
			HasInputs:             true,
			HasInvariants:         false,
			HasOutputFormat:       true,
			RoleConfidence:        0.8,
			InputsConfidence:      0.7,
			InvariantsConfidence:  0.3,
			OutputFormatConfidence: 0.9,
		},
		Score: &scorer.ScoreResult{
			OverallScore: 80.0,
			Metrics: scorer.Metrics{
				Structure:            75.0,
				Conciseness:          70.0,
				InstructionFollowing: 90.0,
				Risk:                 65.0,
			},
		},
		TokenCount: 120,
	}

	report, err := reporter.Generate(result)
	require.NoError(t, err)

	assert.Contains(t, report, "ROLE 확신도")
	assert.Contains(t, report, "INPUTS 확신도")
	assert.Contains(t, report, "INVARIANTS 확신도")
	assert.Contains(t, report, "OUTPUT FORMAT 확신도")
}
