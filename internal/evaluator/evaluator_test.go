package evaluator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/curogom/curompt/internal/model"
	"github.com/curogom/curompt/internal/parser"
)

func TestEvaluator_Evaluate(t *testing.T) {
	// Red: 실패하는 테스트 작성
	eval := NewEvaluator("claude")

	collectedPrompt := &model.CollectedPrompt{
		ID:        "test-1",
		Tool:      "codex",
		RawPrompt: "# ROLE\nEngineer\n\n# INPUTS\n- task: string",
		Timestamp: 1234567890,
		Command:   "codex exec test",
	}

	ctx := context.Background()
	result, err := eval.Evaluate(ctx, collectedPrompt)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Analysis)
	assert.NotNil(t, result.Score)
	assert.Greater(t, result.Score.OverallScore, 0.0)
	assert.LessOrEqual(t, result.Score.OverallScore, 100.0)
}

func TestEvaluator_WithParsedPrompt(t *testing.T) {
	eval := NewEvaluator("claude")

	collectedPrompt := &model.CollectedPrompt{
		ID:        "test-1",
		Tool:      "codex",
		RawPrompt: "# ROLE\nEngineer\n\n# INPUTS\n- task: string",
		Prompt: &parser.Prompt{
			Role:   "Engineer",
			Inputs: []string{"task: string"},
			Raw:    "# ROLE\nEngineer\n\n# INPUTS\n- task: string",
		},
		Timestamp: 1234567890,
	}

	ctx := context.Background()
	result, err := eval.Evaluate(ctx, collectedPrompt)

	require.NoError(t, err)
	assert.Equal(t, "Engineer", result.Prompt.Prompt.Role)
}

func TestEvaluator_CalculatesTokenCount(t *testing.T) {
	eval := NewEvaluator("claude")

	collectedPrompt := &model.CollectedPrompt{
		ID:        "test-1",
		Tool:      "codex",
		RawPrompt: "Hello, world!",
		Timestamp: 1234567890,
	}

	ctx := context.Background()
	result, err := eval.Evaluate(ctx, collectedPrompt)

	require.NoError(t, err)
	assert.Greater(t, result.TokenCount, 0)
}
