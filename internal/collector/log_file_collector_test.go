package collector

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/curogom/curompt/internal/repository"
)

func TestLogFileCollector_Collect_ClaudeHistory(t *testing.T) {
	// Create temporary log file
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "history.jsonl")

	// Write test data
	logContent := `{"display":"Test prompt 1","pastedContents":{},"timestamp":1759236766528,"project":"/test/project"}
{"display":"Test prompt 2","pastedContents":{},"timestamp":1759236766530,"project":"/test/project"}
`
	err := os.WriteFile(logFile, []byte(logContent), 0644)
	require.NoError(t, err)

	// Create repository
	tmpDB := filepath.Join(tmpDir, "test.db")
	repo, err := repository.NewSQLiteRepository(tmpDB)
	require.NoError(t, err)
	defer repo.Close()

	// Create collector
	collector := NewLogFileCollector(repo, "claude", logFile)

	// Collect
	ctx := context.Background()
	prompts, err := collector.Collect(ctx)
	require.NoError(t, err)

	// Verify
	assert.Len(t, prompts, 2)
	assert.Equal(t, "claude", prompts[0].Tool)
	assert.Equal(t, "Test prompt 1", prompts[0].RawPrompt)
	assert.Equal(t, "Test prompt 2", prompts[1].RawPrompt)
	assert.NotEmpty(t, prompts[0].Metadata["project"])
	assert.Equal(t, "history.jsonl", prompts[0].Metadata["source"])
}

func TestLogFileCollector_Collect_CodexHistory(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "history.jsonl")

	logContent := `{"session_id":"test-session","ts":1758027713,"text":"Test codex prompt 1"}
{"session_id":"test-session","ts":1758027714,"text":"Test codex prompt 2"}
`
	err := os.WriteFile(logFile, []byte(logContent), 0644)
	require.NoError(t, err)

	tmpDB := filepath.Join(tmpDir, "test.db")
	repo, err := repository.NewSQLiteRepository(tmpDB)
	require.NoError(t, err)
	defer repo.Close()

	collector := NewLogFileCollector(repo, "codex", logFile)
	ctx := context.Background()

	prompts, err := collector.Collect(ctx)
	require.NoError(t, err)

	assert.Len(t, prompts, 2)
	assert.Equal(t, "codex", prompts[0].Tool)
	assert.Equal(t, "Test codex prompt 1", prompts[0].RawPrompt)
	assert.Equal(t, "Test codex prompt 2", prompts[1].RawPrompt)
	assert.NotEmpty(t, prompts[0].Metadata["session_id"])
}

func TestLogFileCollector_Collect_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "empty.jsonl")
	err := os.WriteFile(logFile, []byte(""), 0644)
	require.NoError(t, err)

	tmpDB := filepath.Join(tmpDir, "test.db")
	repo, err := repository.NewSQLiteRepository(tmpDB)
	require.NoError(t, err)
	defer repo.Close()

	collector := NewLogFileCollector(repo, "claude", logFile)
	ctx := context.Background()

	prompts, err := collector.Collect(ctx)
	require.NoError(t, err)
	assert.Empty(t, prompts)
}

func TestLogFileCollector_Collect_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "invalid.jsonl")
	err := os.WriteFile(logFile, []byte("not valid json\n{\"display\":\"valid\"}\n"), 0644)
	require.NoError(t, err)

	tmpDB := filepath.Join(tmpDir, "test.db")
	repo, err := repository.NewSQLiteRepository(tmpDB)
	require.NoError(t, err)
	defer repo.Close()

	collector := NewLogFileCollector(repo, "claude", logFile)
	ctx := context.Background()

	prompts, err := collector.Collect(ctx)
	require.NoError(t, err)
	// Should skip invalid lines and collect valid ones
	assert.Len(t, prompts, 1)
	assert.Equal(t, "valid", prompts[0].RawPrompt)
}
