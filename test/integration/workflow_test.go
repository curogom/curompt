package integration

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/curogom/curo-prompt/internal/evaluator"
	"github.com/curogom/curo-prompt/internal/model"
	"github.com/curogom/curo-prompt/internal/repository"
)

func TestWorkflow_CollectAnalyzeScoreReport(t *testing.T) {
	// Red: 실패하는 테스트 작성
	// 전체 워크플로우 테스트: 수집 → 저장 → 분석 → 점수화 → 리포트

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// 1. 저장소 초기화
	repo, err := repository.NewSQLiteRepository(dbPath)
	require.NoError(t, err)
	defer repo.Close()

	ctx := context.Background()

	// 2. 프롬프트 수집 (시뮬레이션)
	collectedPrompt := &model.CollectedPrompt{
		ID:        uuid.New().String(),
		Tool:      "test",
		RawPrompt: "# ROLE\nEngineer\n\n# INPUTS\n- task: string",
		Timestamp: time.Now().Unix(),
		Command:   "test command",
	}

	// 3. 저장
	err = repo.Save(ctx, collectedPrompt)
	require.NoError(t, err)

	// 4. 저장된 프롬프트 조회
	found, err := repo.FindByID(ctx, collectedPrompt.ID)
	require.NoError(t, err)
	assert.Equal(t, collectedPrompt.RawPrompt, found.RawPrompt)

	// 5. 평가 수행
	eval := evaluator.NewEvaluator("claude")
	result, err := eval.Evaluate(ctx, found)
	require.NoError(t, err)

	// 6. 결과 검증
	assert.NotNil(t, result)
	assert.Greater(t, result.Score.OverallScore, 0.0)
	assert.LessOrEqual(t, result.Score.OverallScore, 100.0)
	assert.NotNil(t, result.Analysis)
	assert.Greater(t, result.TokenCount, 0)
}

func TestWorkflow_MultiplePrompts(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	repo, err := repository.NewSQLiteRepository(dbPath)
	require.NoError(t, err)
	defer repo.Close()

	ctx := context.Background()
	eval := evaluator.NewEvaluator("claude")

	// 여러 프롬프트 저장 및 평가
	prompts := []*model.CollectedPrompt{
		{
			ID:        uuid.New().String(),
			Tool:      "codex",
			RawPrompt: "# ROLE\nEngineer",
			Timestamp: time.Now().Unix(),
		},
		{
			ID:        uuid.New().String(),
			Tool:      "cursor",
			RawPrompt: "# ROLE\nDesigner\n\n# INPUTS\n- task: string",
			Timestamp: time.Now().Unix() + 1,
		},
	}

	for _, prompt := range prompts {
		err := repo.Save(ctx, prompt)
		require.NoError(t, err)

		result, err := eval.Evaluate(ctx, prompt)
		require.NoError(t, err)
		assert.Greater(t, result.Score.OverallScore, 0.0)
	}

	// 최근 프롬프트 조회
	recent, err := repo.FindRecent(ctx, 2)
	require.NoError(t, err)
	assert.Equal(t, 2, len(recent))
}

func TestWorkflow_ScanEvaluateReport(t *testing.T) {
	// scan 명령 워크플로우 테스트
	tmpDir := t.TempDir()

	// 테스트용 프롬프트 파일 생성
	testPromptFile := filepath.Join(tmpDir, "test_prompt.md")
	testPrompt := "# ROLE\nEngineer\n\n# INPUTS\n- task: string"
	err := os.WriteFile(testPromptFile, []byte(testPrompt), 0644)
	require.NoError(t, err)

	// 저장소
	dbPath := filepath.Join(tmpDir, "test.db")
	repo, err := repository.NewSQLiteRepository(dbPath)
	require.NoError(t, err)
	defer repo.Close()

	ctx := context.Background()

	// 파일 읽기
	content, err := os.ReadFile(testPromptFile)
	require.NoError(t, err)

	// 수집된 프롬프트 생성
	collectedPrompt := &model.CollectedPrompt{
		ID:        uuid.New().String(),
		Tool:      "scan",
		RawPrompt: string(content),
		Timestamp: time.Now().Unix(),
		Command:   "scan --repo " + tmpDir,
	}

	// 저장 및 평가
	err = repo.Save(ctx, collectedPrompt)
	require.NoError(t, err)

	eval := evaluator.NewEvaluator("claude")
	result, err := eval.Evaluate(ctx, collectedPrompt)
	require.NoError(t, err)

	assert.NotNil(t, result)
	assert.NotNil(t, result.Score)
}
