package claude

import (
	"context"
	"os"
	"testing"

	"github.com/curogom/curo-prompt/internal/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClaudeProvider_Evaluate(t *testing.T) {
	// Red: 실패하는 테스트 작성
	// API 키가 없으면 스킵하거나 에러 반환
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("ANTHROPIC_API_KEY not set, skipping integration test")
	}

	p := NewClaudeProvider(apiKey, "claude-3-5-sonnet-20241022")

	ctx := context.Background()
	response, err := p.Evaluate(ctx, "Hello, world!")

	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Content)
	assert.Equal(t, "claude-3-5-sonnet-20241022", response.Model)
}

func TestClaudeProvider_CalculateTokens(t *testing.T) {
	p := NewClaudeProvider("test-key", "claude-3-5-sonnet-20241022")

	count, err := p.CalculateTokens("Hello, world!")

	require.NoError(t, err)
	assert.Greater(t, count, 0)
}

func TestClaudeProvider_Name(t *testing.T) {
	p := NewClaudeProvider("test-key", "claude-3-5-sonnet-20241022")

	name := p.Name()
	assert.Equal(t, "claude", name)
}

func TestClaudeProvider_ImplementsProviderInterface(t *testing.T) {
	var _ provider.Provider = (*ClaudeProvider)(nil)
}
