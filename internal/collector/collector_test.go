package collector

import (
	"context"
	"testing"

	"github.com/curogom/curo-prompt/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestCollector_Interface(t *testing.T) {
	// Red: 실패하는 테스트 작성
	// Collector 인터페이스가 올바르게 정의되어 있는지 확인
	var c Collector = &mockCollector{}

	name := c.Name()
	assert.NotEmpty(t, name)
}

func TestCollector_Collect(t *testing.T) {
	var c Collector = &mockCollector{}

	ctx := context.Background()
	prompts, err := c.Collect(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, prompts)
}

// mockCollector is a mock implementation for testing
type mockCollector struct{}

func (m *mockCollector) Collect(ctx context.Context) ([]*model.CollectedPrompt, error) {
	return []*model.CollectedPrompt{
		{
			ID:        "test-1",
			Tool:      "mock",
			RawPrompt: "# ROLE\nTest",
			Timestamp: 1234567890,
			Command:   "mock command",
		},
	}, nil
}

func (m *mockCollector) Name() string {
	return "mock"
}
