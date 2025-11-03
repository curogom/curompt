package cli

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"

	"github.com/curogom/curo-prompt/internal/evaluator"
	"github.com/curogom/curo-prompt/internal/model"
	"github.com/curogom/curo-prompt/internal/reporter"
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
			provider, _ := cmd.Flags().GetString("provider") //nolint:errcheck
			filePath, _ := cmd.Flags().GetString("file")     //nolint:errcheck
			outputPath, _ := cmd.Flags().GetString("output") //nolint:errcheck

			// 프롬프트 읽기
			var content []byte
			var err error

			if filePath != "" {
				content, err = os.ReadFile(filePath)
				if err != nil {
					return fmt.Errorf("프롬프트 파일을 읽지 못했습니다: %w", err)
				}
			} else {
				// 파일 미지정: 표준 입력 필요 여부 확인
				stdinFile := os.Stdin
				fd := stdinFile.Fd()
				if isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd) {
					return fmt.Errorf("입력이 필요합니다: --file <경로> 옵션을 사용하거나 표준 입력으로 프롬프트를 제공하세요")
				}

				// 표준 입력에서 읽기
				content, err = io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("표준 입력에서 프롬프트를 읽지 못했습니다: %w", err)
				}
				if len(bytes.TrimSpace(content)) == 0 {
					return fmt.Errorf("빈 입력입니다: --file <경로> 옵션을 사용하거나 표준 입력으로 프롬프트 내용을 제공하세요")
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

			stopSpinner := startSpinner(cmd.ErrOrStderr(), "프롬프트 평가 중입니다")
			defer stopSpinner("")

			result, err := eval.Evaluate(ctx, collectedPrompt)
			if err != nil {
				stopSpinner("프롬프트 평가 실패")
				return fmt.Errorf("평가 실행 중 오류가 발생했습니다: %w", err)
			}

			// 리포트 생성
			markdownReporter := reporter.NewMarkdownReporter()
			report, err := markdownReporter.Generate(result)
			if err != nil {
				stopSpinner("리포트 생성 실패")
				return fmt.Errorf("리포트 생성 중 오류가 발생했습니다: %w", err)
			}

			if outputPath != "" {
				// 파일로 저장
				if err := os.WriteFile(outputPath, []byte(report), 0644); err != nil {
					stopSpinner("리포트 저장 실패")
					return fmt.Errorf("리포트를 파일로 저장하지 못했습니다: %w", err)
				}
				stopSpinner("평가가 완료되었습니다")
				cmd.Printf("리포트가 저장되었습니다: %s\n", outputPath)
			} else {
				stopSpinner("평가가 완료되었습니다")
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
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}

func startSpinner(w io.Writer, message string) func(string) {
	writer, ok := w.(*os.File)
	isTTY := false
	if ok {
		fd := writer.Fd()
		isTTY = isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
	}

	var once sync.Once

	if !isTTY {
		if message != "" {
			fmt.Fprintf(w, "%s...\n", message)
		}
		return func(final string) {
			once.Do(func() {
				if final != "" {
					fmt.Fprintf(w, "%s\n", final)
				}
			})
		}
	}

	spinnerChars := []rune{'|', '/', '-', '\\'}
	updates := make(chan string)
	done := make(chan struct{})

	go func() {
		ticker := time.NewTicker(120 * time.Millisecond)
		defer ticker.Stop()

		idx := 0
		fmt.Fprintf(w, "%s %c", message, spinnerChars[idx])
		idx = (idx + 1) % len(spinnerChars)

		for {
			select {
			case final := <-updates:
				padding := strings.Repeat(" ", utf8.RuneCountInString(message)+2)
				fmt.Fprintf(w, "\r%s\r", padding)
				if final != "" {
					fmt.Fprintf(w, "%s\n", final)
				}
				close(done)
				return
			case <-ticker.C:
				fmt.Fprintf(w, "\r%s %c", message, spinnerChars[idx])
				idx = (idx + 1) % len(spinnerChars)
			}
		}
	}()

	return func(final string) {
		once.Do(func() {
			updates <- final
			<-done
		})
	}
}
