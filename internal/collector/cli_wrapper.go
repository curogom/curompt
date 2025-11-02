package collector

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/curogom/curo-prompt/internal/model"
	"github.com/curogom/curo-prompt/internal/parser"
	"github.com/curogom/curo-prompt/internal/repository"
)

// CLIWrapperCollector wraps CLI commands and collects prompts
type CLIWrapperCollector struct {
	repository repository.PromptRepository
	parser     parser.Parser
	tool       string // codex, cursor 등
}

// NewCLIWrapperCollector creates a new CLI wrapper collector
func NewCLIWrapperCollector(repo repository.PromptRepository, tool string) *CLIWrapperCollector {
	return &CLIWrapperCollector{
		repository: repo,
		parser:     parser.NewParser(),
		tool:       tool,
	}
}

// WrapCommand wraps a CLI command and collects prompts
func (c *CLIWrapperCollector) WrapCommand(ctx context.Context, command string, args []string) error {
	// 원래 명령 실행
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 현재 디렉토리 설정
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}
	cmd.Dir = wd

	// 명령 실행
	err = cmd.Run()

	// 프롬프트 추출 (도구별로 다른 방법)
	// TODO: 실제 구현은 각 도구별로 다름
	promptText := c.extractPrompt(command, args)

	if promptText != "" {
		// 프롬프트 파싱
		parsedPrompt, parseErr := c.parser.Parse(promptText)
		if parseErr != nil {
			// 파싱 실패해도 원본 저장
			parsedPrompt = nil
		}

		// 수집된 프롬프트 생성
		collected := &model.CollectedPrompt{
			ID:         uuid.New().String(),
			Tool:       c.tool,
			Prompt:     parsedPrompt,
			RawPrompt:  promptText,
			Timestamp:  time.Now().Unix(),
			Command:    fmt.Sprintf("%s %s", command, strings.Join(args, " ")),
			WorkingDir: wd,
			Metadata: map[string]string{
				"exit_code": fmt.Sprintf("%d", cmd.ProcessState.ExitCode()),
			},
		}

		// 저장
		saveErr := c.repository.Save(ctx, collected)
		if saveErr != nil {
			// 저장 실패는 로깅만 하고 명령 실행은 계속
			fmt.Fprintf(os.Stderr, "Warning: failed to save prompt: %v\n", saveErr)
		}
	}

	return err
}

// extractPrompt extracts prompt from command arguments
// 도구별로 다른 추출 방법 필요
func (c *CLIWrapperCollector) extractPrompt(command string, args []string) string {
	// 간단한 구현: args에서 프롬프트 찾기
	// 실제로는 도구별로 다른 파싱 필요
	commandStr := strings.Join(args, " ")

	// codex나 cursor 명령에서 프롬프트 추출
	// 예: codex exec "TASK: ..." -> "TASK: ..."
	for i, arg := range args {
		if arg == "exec" && i+1 < len(args) {
			return args[i+1]
		}
		// 큰따옴표나 작은따옴표로 감싸진 텍스트 추출
		if strings.HasPrefix(arg, `"`) && strings.HasSuffix(arg, `"`) {
			return strings.Trim(arg, `"`)
		}
		if strings.HasPrefix(arg, `'`) && strings.HasSuffix(arg, `'`) {
			return strings.Trim(arg, `'`)
		}
	}

	// 전체 args를 프롬프트로 간주 (fallback)
	return commandStr
}

// Collect collects prompts from wrapped commands (현재는 지원 안 함)
func (c *CLIWrapperCollector) Collect(ctx context.Context) ([]*model.CollectedPrompt, error) {
	// CLI 래퍼는 WrapCommand를 사용하므로 Collect는 사용하지 않음
	return nil, fmt.Errorf("CLIWrapperCollector does not support Collect(), use WrapCommand() instead")
}

// Name returns the collector name
func (c *CLIWrapperCollector) Name() string {
	return fmt.Sprintf("cli-wrapper-%s", c.tool)
}
