package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/curogom/curompt/internal/evaluator"
	"github.com/curogom/curompt/internal/reporter"
	"github.com/curogom/curompt/internal/repository"
)

// newCollectCmd creates the collect command for log file collection
func newCollectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collect [flags]",
		Short: "로그 파일에서 프롬프트 수집",
		Long: `로그 파일이나 히스토리 파일에서 프롬프트를 자동으로 수집합니다.

지원 도구:
  - claude: Claude Code history.jsonl 파일
  - codex: Codex CLI history.jsonl 파일
  - cursor: Cursor IDE 로그 (예정)

예시:
  # Claude Code: 현재 프로젝트만 수집 (CLAUDE.md가 있는 디렉토리)
  cd /path/to/project
  curompt collect --from claude
  
  # Claude Code: 모든 프로젝트의 프롬프트 수집
  curompt collect --from claude --all
  
  # Codex: 현재 디렉토리만 수집 (CLAUDE.md 불필요)
  cd /path/to/project
  curompt collect --from codex
  
  # Codex: 모든 프로젝트의 프롬프트 수집
  curompt collect --from codex --all
  
  # 특정 파일 지정
  curompt collect --from claude --file ~/.claude/history.jsonl
  
  # 수집 후 자동 평가
  curompt collect --from claude --eval`,
		RunE: func(cmd *cobra.Command, args []string) error {
			from, _ := cmd.Flags().GetString("from")         //nolint:errcheck
			filePath, _ := cmd.Flags().GetString("file")     //nolint:errcheck
			collectAll, _ := cmd.Flags().GetBool("all")      //nolint:errcheck
			eval, _ := cmd.Flags().GetBool("eval")           //nolint:errcheck
			output, _ := cmd.Flags().GetString("output")     //nolint:errcheck
			provider, _ := cmd.Flags().GetString("provider") //nolint:errcheck

			if from == "" {
				return fmt.Errorf("--from 옵션은 필수입니다 (claude, codex, cursor)")
			}

			// 프로젝트 루트 확인 (--all이 아닌 경우)
			var projectRoot string
			if !collectAll {
				wd, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("현재 디렉토리 확인 실패: %w", err)
				}

				// Claude Code: CLAUDE.md 파일로 프로젝트 루트 찾기
				// Codex: 현재 디렉토리를 프로젝트 경로로 사용 (session 파일의 cwd와 매칭)
				if from == "claude" || from == "claude-code" {
					projectRoot = findProjectRoot(wd)
					if projectRoot == "" {
						return fmt.Errorf("프로젝트 루트를 찾을 수 없습니다. CLAUDE.md 또는 Claude.md 파일이 있는 디렉토리로 이동하거나 --all 옵션을 사용하세요")
					}
				} else if from == "codex" {
					// Codex의 경우 현재 디렉토리를 프로젝트 경로로 사용
					projectRoot = wd
				} else {
					// 기타 도구는 CLAUDE.md 기반
					projectRoot = findProjectRoot(wd)
					if projectRoot == "" {
						return fmt.Errorf("프로젝트 루트를 찾을 수 없습니다. CLAUDE.md 또는 Claude.md 파일이 있는 디렉토리로 이동하거나 --all 옵션을 사용하세요")
					}
				}
				cmd.Printf("프로젝트 루트: %s\n", projectRoot)
			}

			// 저장소 초기화
			dbPath := filepath.Join(os.Getenv("HOME"), ".curompt", "db.sqlite")
			if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
				return fmt.Errorf("DB 디렉토리 생성 실패: %w", err)
			}

			repo, err := repository.NewSQLiteRepository(dbPath)
			if err != nil {
				return fmt.Errorf("저장소 초기화 실패: %w", err)
			}
			defer func() {
				_ = repo.Close()
			}()

			// 로그 파일 경로 결정
			if filePath == "" {
				filePath = getDefaultLogPath(from)
				if filePath == "" {
					return fmt.Errorf("로그 파일 경로를 찾을 수 없습니다. --file 옵션으로 직접 지정해주세요")
				}
			}

			// Collector 실행
			var projectFilter string
			if !collectAll && projectRoot != "" {
				projectFilter = projectRoot
			}

			ctx := context.Background()
			prompts, savedCount, err := collectFromLog(ctx, repo, from, filePath, projectFilter, cmd)
			if err != nil {
				return fmt.Errorf("수집 실패: %w", err)
			}

			cmd.Printf("저장 완료: %d개\n", savedCount)

			// 평가 옵션이 있으면 각 프롬프트 평가
			if eval {
				evaluator := evaluator.NewEvaluator(provider)
				reporter := reporter.NewMarkdownReporter()
				outputDir := output
				if outputDir == "" {
					outputDir = "reports"
				}
				if err := os.MkdirAll(outputDir, 0755); err != nil {
					return fmt.Errorf("출력 디렉토리 생성 실패: %w", err)
				}

				cmd.Printf("\n평가 수행 중...\n\n")
				for i, prompt := range prompts {
					cmd.Printf("[%d/%d] 평가 중...\n", i+1, len(prompts))

					result, err := evaluator.Evaluate(ctx, prompt)
					if err != nil {
						cmd.Printf("  ⚠️  평가 실패: %v\n", err)
						continue
					}

					report, err := reporter.Generate(result)
					if err != nil {
						cmd.Printf("  ⚠️  리포트 생성 실패: %v\n", err)
						continue
					}

					reportPath := filepath.Join(outputDir, fmt.Sprintf("%s_report.md", prompt.ID[:8]))
					if err := os.WriteFile(reportPath, []byte(report), 0644); err != nil {
						cmd.Printf("  ⚠️  저장 실패: %v\n", err)
						continue
					}

					cmd.Printf("  ✅ 완료 - 점수: %.1f/100, 리포트: %s\n", result.Score.OverallScore, reportPath)
				}
			}

			return nil
		},
	}

	cmd.Flags().String("from", "", "수집할 도구 (claude, codex, cursor)")
	cmd.Flags().StringP("file", "f", "", "로그 파일 경로 (미지정 시 기본 경로 사용)")
	cmd.Flags().Bool("all", false, "모든 프로젝트의 프롬프트 수집 (미지정 시 현재 디렉토리 프로젝트만)")
	cmd.Flags().Bool("eval", false, "수집한 프롬프트들을 평가하여 리포트 생성")
	cmd.Flags().StringP("output", "o", "reports", "리포트 출력 디렉토리 (--eval 사용 시)")
	cmd.Flags().StringP("provider", "p", "claude", "LLM Provider (토큰 계산용, --eval 사용 시)")

	return cmd
}

// getDefaultLogPath returns the default log file path for a tool
