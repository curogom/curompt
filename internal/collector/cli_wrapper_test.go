package collector

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/curogom/curompt/internal/repository"
)

func TestCLIWrapperCollector_WrapCommand(t *testing.T) {
	// Red: 실패하는 테스트 작성
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	repo, err := repository.NewSQLiteRepository(dbPath)
	require.NoError(t, err)
	defer repo.Close()

	collector := NewCLIWrapperCollector(repo, "codex")

	ctx := context.Background()

	// echo 명령으로 테스트 (실제 도구 대신)
	err = collector.WrapCommand(ctx, "echo", []string{"Hello"})

	// 명령은 성공해야 함
	assert.NoError(t, err)
}

func TestCLIWrapperCollector_ExtractPrompt(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	repo, err := repository.NewSQLiteRepository(dbPath)
	require.NoError(t, err)
	defer repo.Close()

	collector := NewCLIWrapperCollector(repo, "codex")

	// 프롬프트 추출 테스트
	prompt := collector.extractPrompt("codex", []string{"exec", "TASK: Add feature"})
	assert.Contains(t, prompt, "TASK")

	// 큰따옴표로 감싸진 경우
	prompt2 := collector.extractPrompt("cursor", []string{"ask", `"How to implement X?"`})
	assert.Equal(t, "How to implement X?", prompt2)
}

func TestCLIWrapperCollector_SavesPrompt(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	repo, err := repository.NewSQLiteRepository(dbPath)
	require.NoError(t, err)
	defer repo.Close()

	collector := NewCLIWrapperCollector(repo, "codex")

	ctx := context.Background()

	// echo로 테스트하되 프롬프트가 포함된 명령 시뮬레이션
	// 실제로는 codex exec "TASK: ..." 같은 명령이지만, 테스트용으로 echo 사용
	err = collector.WrapCommand(ctx, "echo", []string{`"# ROLE\nEngineer"`})
	require.NoError(t, err)

	// 저장소에서 최근 프롬프트 확인
	recent, err := repo.FindRecent(ctx, 1)
	require.NoError(t, err)
	assert.Greater(t, len(recent), 0)
}

func TestCLIWrapperCollector_Name(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	repo, err := repository.NewSQLiteRepository(dbPath)
	require.NoError(t, err)
	defer repo.Close()

	collector := NewCLIWrapperCollector(repo, "cursor")

	name := collector.Name()
	assert.Equal(t, "cli-wrapper-cursor", name)
}
