package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClaudeTokenizer_CountTokens(t *testing.T) {
	// Red: 실패하는 테스트 작성
	tokenizer := NewClaudeTokenizer()

	result, err := tokenizer.CountTokens("Hello, world!")
	require.NoError(t, err)
	assert.Greater(t, result, 0)
}

func TestClaudeTokenizer_CountTokensForLongText(t *testing.T) {
	tokenizer := NewClaudeTokenizer()
	text := "# ROLE\nSenior engineer.\n\n# INPUTS\n- task: string"

	result, err := tokenizer.CountTokens(text)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, result, 5) // 최소 5개 토큰 이상
}

func TestOpenAITokenizer_CountTokens(t *testing.T) {
	tokenizer := NewOpenAITokenizer()

	result, err := tokenizer.CountTokens("Hello, world!")
	require.NoError(t, err)
	assert.Greater(t, result, 0)
}

func TestTokenizers_SimilarTextShouldHaveSimilarCounts(t *testing.T) {
	text := "This is a test prompt with some content."
	claudeTokenizer := NewClaudeTokenizer()
	openAITokenizer := NewOpenAITokenizer()

	claudeCount, err := claudeTokenizer.CountTokens(text)
	require.NoError(t, err)

	openAICount, err := openAITokenizer.CountTokens(text)
	require.NoError(t, err)

	// 두 토큰 카운트가 크게 다르지 않아야 함 (최대 50% 차이)
	ratio := float64(claudeCount) / float64(openAICount)
	assert.Greater(t, ratio, 0.5)
	assert.Less(t, ratio, 2.0)
}
