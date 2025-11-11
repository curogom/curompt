package repository

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/curogom/curompt/internal/model"
	"github.com/curogom/curompt/internal/parser"
)

func TestPromptRepository_SaveAndFind(t *testing.T) {
	// Red: 실패하는 테스트 작성
	// 임시 DB 파일 사용
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	repo, err := NewSQLiteRepository(dbPath)
	require.NoError(t, err)
	defer func() {
		_ = repo.Close()
	}()

	ctx := context.Background()

	prompt := &model.CollectedPrompt{
		ID:        "test-1",
		Tool:      "codex",
		RawPrompt: "# ROLE\nEngineer",
		Prompt: &parser.Prompt{
			Role: "Engineer",
			Raw:  "# ROLE\nEngineer",
		},
		Timestamp:  1234567890,
		Command:    "codex exec test",
		WorkingDir: "/tmp",
		Metadata:   map[string]string{"key": "value"},
	}

	err = repo.Save(ctx, prompt)
	require.NoError(t, err)

	found, err := repo.FindByID(ctx, "test-1")
	require.NoError(t, err)
	assert.Equal(t, prompt.ID, found.ID)
	assert.Equal(t, prompt.Tool, found.Tool)
	assert.Equal(t, prompt.RawPrompt, found.RawPrompt)
}

func TestPromptRepository_FindByTool(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	repo, err := NewSQLiteRepository(dbPath)
	require.NoError(t, err)
	defer func() {
		_ = repo.Close()
	}()

	ctx := context.Background()

	// 두 개의 프롬프트 저장
	prompt1 := &model.CollectedPrompt{
		ID:        "test-1",
		Tool:      "codex",
		RawPrompt: "Prompt 1",
		Timestamp: 1234567890,
	}
	prompt2 := &model.CollectedPrompt{
		ID:        "test-2",
		Tool:      "cursor",
		RawPrompt: "Prompt 2",
		Timestamp: 1234567891,
	}

	require.NoError(t, repo.Save(ctx, prompt1))
	require.NoError(t, repo.Save(ctx, prompt2))

	// codex로 수집된 프롬프트만 조회
	codexPrompts, err := repo.FindByTool(ctx, "codex")
	require.NoError(t, err)
	assert.Equal(t, 1, len(codexPrompts))
	assert.Equal(t, "test-1", codexPrompts[0].ID)
}

func TestPromptRepository_FindRecent(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	repo, err := NewSQLiteRepository(dbPath)
	require.NoError(t, err)
	defer func() {
		_ = repo.Close()
	}()

	ctx := context.Background()

	// 여러 프롬프트 저장
	for i := 0; i < 5; i++ {
		prompt := &model.CollectedPrompt{
			ID:        fmt.Sprintf("test-%d", i),
			Tool:      "codex",
			RawPrompt: fmt.Sprintf("Prompt %d", i),
			Timestamp: int64(1234567890 + i),
		}
		require.NoError(t, repo.Save(ctx, prompt))
	}

	// 최근 3개 조회
	recent, err := repo.FindRecent(ctx, 3)
	require.NoError(t, err)
	assert.Equal(t, 3, len(recent))
	// 최신순으로 정렬되어야 함
	assert.Equal(t, "test-4", recent[0].ID)
}
