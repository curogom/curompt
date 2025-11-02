package provider

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProvider_Evaluate(t *testing.T) {
	// Red: 실패하는 테스트 작성
	// Mock provider를 만들어야 하지만, 먼저 인터페이스 계약 테스트
	ctx := context.Background()
	prompt := "Test prompt"

	// 인터페이스가 올바르게 정의되어 있는지 확인
	var p Provider = &mockProvider{}

	response, err := p.Evaluate(ctx, prompt)
	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Content)
}

func TestProvider_CalculateTokens(t *testing.T) {
	var p Provider = &mockProvider{}

	count, err := p.CalculateTokens("Hello, world!")
	require.NoError(t, err)
	assert.Greater(t, count, 0)
}

func TestProvider_Name(t *testing.T) {
	var p Provider = &mockProvider{}

	name := p.Name()
	assert.NotEmpty(t, name)
}

// mockProvider is a mock implementation for testing
type mockProvider struct{}

func (m *mockProvider) Evaluate(ctx context.Context, prompt string) (*Response, error) {
	return &Response{
		Content:    "Mock response",
		TokenCount: 10,
		Model:      "mock-model",
	}, nil
}

func (m *mockProvider) CalculateTokens(text string) (int, error) {
	return len(text) / 4, nil // 간단한 근사치
}

func (m *mockProvider) Name() string {
	return "mock"
}
