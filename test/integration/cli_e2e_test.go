package integration

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/curogom/curo-prompt/internal/evaluator"
	"github.com/curogom/curo-prompt/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCLI_EvalCommand tests the eval command end-to-end
func TestCLI_EvalCommand(t *testing.T) {
	// Skip if binary not found
	binaryPath := findBinary(t)
	if binaryPath == "" {
		t.Skip("curo-prompt binary not found, skipping E2E test")
	}

	tmpDir := t.TempDir()

	// Create test prompt file
	testPrompt := `# ROLE
Senior Engineer

# INPUTS
- task: string
- context: string

# INVARIANTS
- Must be testable
- Must be documented

# OUTPUT FORMAT
JSON`
	testFile := filepath.Join(tmpDir, "test_prompt.md")
	err := os.WriteFile(testFile, []byte(testPrompt), 0644)
	require.NoError(t, err)

	// Run eval command
	cmd := exec.Command(binaryPath, "eval", "--file", testFile, "--provider", "claude")
	output, err := cmd.CombinedOutput()

	// Command should succeed
	assert.NoError(t, err, "eval command failed: %s", string(output))

	// Output should contain score
	outputStr := string(output)
	assert.Contains(t, outputStr, "점수", "output should contain score")
	assert.Contains(t, outputStr, "ROLE", "output should contain parsed sections")
}

// TestCLI_ScanCommand tests the scan command end-to-end
func TestCLI_ScanCommand(t *testing.T) {
	binaryPath := findBinary(t)
	if binaryPath == "" {
		t.Skip("curo-prompt binary not found, skipping E2E test")
	}

	tmpDir := t.TempDir()

	// Create test prompts directory
	promptsDir := filepath.Join(tmpDir, "prompts")
	err := os.MkdirAll(promptsDir, 0755)
	require.NoError(t, err)

	// Create test prompt files
	testPrompt1 := `# ROLE
Engineer`
	testPrompt2 := `# ROLE
Designer

# INPUTS
- task: string`

	err = os.WriteFile(filepath.Join(promptsDir, "prompt1.md"), []byte(testPrompt1), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(promptsDir, "prompt2.md"), []byte(testPrompt2), 0644)
	require.NoError(t, err)

	// Create output directory
	outputDir := filepath.Join(tmpDir, "reports")

	// Run scan command
	cmd := exec.Command(binaryPath, "scan", "--repo", promptsDir, "--output", outputDir, "--provider", "claude")
	output, err := cmd.CombinedOutput()

	// Command should succeed
	assert.NoError(t, err, "scan command failed: %s", string(output))

	// Reports should be generated
	outputStr := string(output)
	assert.Contains(t, outputStr, "분석 중", "output should contain analysis progress")

	// Check if report files exist
	files, err := os.ReadDir(outputDir)
	require.NoError(t, err)
	assert.Greater(t, len(files), 0, "report files should be generated")
}

// TestCLI_SuggestCommand tests the suggest command end-to-end
func TestCLI_SuggestCommand(t *testing.T) {
	binaryPath := findBinary(t)
	if binaryPath == "" {
		t.Skip("curo-prompt binary not found, skipping E2E test")
	}

	tmpDir := t.TempDir()

	// Create test prompt file (missing sections)
	testPrompt := `# ROLE
Engineer`
	testFile := filepath.Join(tmpDir, "test_prompt.md")
	err := os.WriteFile(testFile, []byte(testPrompt), 0644)
	require.NoError(t, err)

	// Run suggest command
	cmd := exec.Command(binaryPath, "suggest", "--file", testFile, "--provider", "claude")
	output, err := cmd.CombinedOutput()

	// Command should succeed
	assert.NoError(t, err, "suggest command failed: %s", string(output))

	// Output should contain suggestions
	outputStr := string(output)
	assert.Contains(t, outputStr, "개선 제안", "output should contain suggestions")
	assert.Contains(t, outputStr, "점수", "output should contain score")
}

// TestCLI_StdinInput tests reading from stdin
func TestCLI_StdinInput(t *testing.T) {
	binaryPath := findBinary(t)
	if binaryPath == "" {
		t.Skip("curo-prompt binary not found, skipping E2E test")
	}

	// Create test prompt
	testPrompt := `# ROLE
Engineer

# INPUTS
- task: string`

	// Run eval command with stdin
	cmd := exec.Command(binaryPath, "eval", "--provider", "claude")
	cmd.Stdin = strings.NewReader(testPrompt)
	output, err := cmd.CombinedOutput()

	// Command should succeed
	assert.NoError(t, err, "eval with stdin failed: %s", string(output))

	// Output should contain score
	outputStr := string(output)
	assert.Contains(t, outputStr, "점수", "output should contain score")
}

// TestWorkflow_EvaluatorIntegration tests the complete evaluator workflow
func TestWorkflow_EvaluatorIntegration(t *testing.T) {
	ctx := context.Background()

	// Create evaluator
	eval := evaluator.NewEvaluator("claude")

	// Create test prompt
	collectedPrompt := &model.CollectedPrompt{
		ID:        "test-id",
		Tool:      "test",
		RawPrompt: "# ROLE\nEngineer\n\n# INPUTS\n- task: string\n- context: string\n\n# INVARIANTS\n- Must be testable\n\n# OUTPUT FORMAT\nJSON",
	}

	// Evaluate
	result, err := eval.Evaluate(ctx, collectedPrompt)
	require.NoError(t, err)

	// Verify result structure
	assert.NotNil(t, result)
	assert.NotNil(t, result.Score)
	assert.Greater(t, result.Score.OverallScore, 0.0)
	assert.LessOrEqual(t, result.Score.OverallScore, 100.0)
	assert.NotNil(t, result.Analysis)
	assert.Greater(t, result.TokenCount, 0)

	// Verify analysis results
	assert.True(t, result.Analysis.HasRole)
	assert.True(t, result.Analysis.HasInputs)
	assert.True(t, result.Analysis.HasInvariants)
	assert.True(t, result.Analysis.HasOutputFormat)
}

// findBinary finds the curo-prompt binary
func findBinary(t *testing.T) string {
	// Try common locations
	locations := []string{
		"./bin/curo-prompt",
		"../bin/curo-prompt",
		"../../bin/curo-prompt",
	}

	for _, loc := range locations {
		if _, err := os.Stat(loc); err == nil {
			abs, err := filepath.Abs(loc)
			if err == nil {
				return abs
			}
		}
	}

	// Try to find in PATH
	path, err := exec.LookPath("curo-prompt")
	if err == nil {
		return path
	}

	return ""
}

