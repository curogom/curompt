package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/curogom/curo-prompt/internal/evaluator"
	"github.com/curogom/curo-prompt/internal/model"
	"github.com/curogom/curo-prompt/internal/reporter"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// newEvalCmd creates the eval command
func newEvalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "eval [flags]",
		Short: "프롬프트 평가 및 점수화",
		Long: `프롬프트를 분석하고 점수를 계산합니다.

입력은 파일 경로 또는 표준 입력으로 받을 수 있습니다.

예시:
  curo-prompt eval --file prompts/dev_contract_v2.md
  cat prompts/dev_contract_v2.md | curo-prompt eval
  curo-prompt eval --provider claude --file prompt.md`,
		RunE: func(cmd *cobra.Command, args []string) error {
			provider, _ := cmd.Flags().GetString("provider")
			filePath, _ := cmd.Flags().GetString("file")
			outputPath, _ := cmd.Flags().GetString("output")

			// 프롬프트 읽기
			var content []byte
			var err error

			if filePath != "" {
				content, err = os.ReadFile(filePath)
				if err != nil {
					return fmt.Errorf("failed to read file: %w", err)
				}
			} else {
				// 표준 입력에서 읽기
				content, err = io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("failed to read from stdin: %w", err)
				}
			}

			// 수집된 프롬프트 생성 (임시)
			collectedPrompt := &model.CollectedPrompt{
				ID:         uuid.New().String(),
				Tool:       "manual",
				RawPrompt:  string(content),
				Timestamp:  time.Now().Unix(),
				Command:    fmt.Sprintf("eval --file %s", filePath),
				WorkingDir: getWorkingDir(),
			}

			// 평가 수행
			eval := evaluator.NewEvaluator(provider)
			ctx := context.Background()

			result, err := eval.Evaluate(ctx, collectedPrompt)
			if err != nil {
				return fmt.Errorf("failed to evaluate: %w", err)
			}

			// 리포트 생성
			markdownReporter := reporter.NewMarkdownReporter()
			report, err := markdownReporter.Generate(result)
			if err != nil {
				return fmt.Errorf("failed to generate report: %w", err)
			}

			// 출력
			if outputPath != "" {
				// 파일로 저장
				if err := os.WriteFile(outputPath, []byte(report), 0644); err != nil {
					return fmt.Errorf("failed to write report: %w", err)
				}
				cmd.Printf("리포트가 저장되었습니다: %s\n", outputPath)
			} else {
				// 터미널 출력
				cmd.Print(report)
			}

			return nil
		},
	}

	cmd.Flags().StringP("file", "f", "", "평가할 프롬프트 파일 경로")
	cmd.Flags().StringP("provider", "p", "claude", "LLM Provider (claude, openai)")
	cmd.Flags().Bool("dynamic", false, "동적 평가 수행 (LLM 호출)")
	cmd.Flags().String("output", "", "리포트 출력 파일 경로 (미지정 시 터미널 출력)")

	return cmd
}

func getWorkingDir() string {
	wd, _ := os.Getwd()
	return wd
}
