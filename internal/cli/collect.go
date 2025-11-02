package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/curogom/curo-prompt/internal/collector"
	"github.com/curogom/curo-prompt/internal/evaluator"
	"github.com/curogom/curo-prompt/internal/repository"
	"github.com/curogom/curo-prompt/internal/reporter"
	"github.com/spf13/cobra"
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
  curo-prompt collect --from claude
  
  # Claude Code: 모든 프로젝트의 프롬프트 수집
  curo-prompt collect --from claude --all
  
  # Codex: 현재 디렉토리만 수집 (CLAUDE.md 불필요)
  cd /path/to/project
  curo-prompt collect --from codex
  
  # Codex: 모든 프로젝트의 프롬프트 수집
  curo-prompt collect --from codex --all
  
  # 특정 파일 지정
  curo-prompt collect --from claude --file ~/.claude/history.jsonl
  
  # 수집 후 자동 평가
  curo-prompt collect --from claude --eval`,
		RunE: func(cmd *cobra.Command, args []string) error {
			from, _ := cmd.Flags().GetString("from")
			filePath, _ := cmd.Flags().GetString("file")
			collectAll, _ := cmd.Flags().GetBool("all")
			eval, _ := cmd.Flags().GetBool("eval")
			output, _ := cmd.Flags().GetString("output")
			provider, _ := cmd.Flags().GetString("provider")

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
			dbPath := filepath.Join(os.Getenv("HOME"), ".curo-prompt", "db.sqlite")
			if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
				return fmt.Errorf("DB 디렉토리 생성 실패: %w", err)
			}

			repo, err := repository.NewSQLiteRepository(dbPath)
			if err != nil {
				return fmt.Errorf("저장소 초기화 실패: %w", err)
			}
			defer repo.Close()

			// 로그 파일 경로 결정
			if filePath == "" {
				filePath = getDefaultLogPath(from)
				if filePath == "" {
					return fmt.Errorf("로그 파일 경로를 찾을 수 없습니다. --file 옵션으로 직접 지정해주세요")
				}
			}

			// 파일 존재 확인
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return fmt.Errorf("로그 파일을 찾을 수 없습니다: %s", filePath)
			}

			// Collector 생성
			var projectFilter string
			if !collectAll && projectRoot != "" {
				projectFilter = projectRoot
			}
			logCollector := collector.NewLogFileCollector(repo, from, filePath)
			if projectFilter != "" {
				logCollector.SetProjectFilter(projectFilter)
			}

			// 수집 실행
			ctx := context.Background()
			cmd.Printf("로그 파일에서 프롬프트 수집 중: %s\n", filePath)
			cmd.Printf("도구: %s\n", from)
			if projectFilter != "" {
				cmd.Printf("필터: %s (현재 프로젝트만)\n", projectFilter)
			} else {
				cmd.Printf("필터: 모든 프로젝트\n")
			}
			cmd.Println()

			prompts, err := logCollector.Collect(ctx)
			if err != nil {
				return fmt.Errorf("수집 실패: %w", err)
			}

			if len(prompts) == 0 {
				cmd.Printf("수집된 프롬프트가 없습니다.\n")
				return nil
			}

			cmd.Printf("수집된 프롬프트: %d개\n\n", len(prompts))

			// 저장
			savedCount := 0
			for i, prompt := range prompts {
				if err := repo.Save(ctx, prompt); err != nil {
					cmd.Printf("  [%d/%d] ⚠️  저장 실패: %v\n", i+1, len(prompts), err)
					continue
				}
				savedCount++
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
func getDefaultLogPath(tool string) string {
	home := os.Getenv("HOME")
	switch tool {
	case "claude", "claude-code":
		return filepath.Join(home, ".claude", "history.jsonl")
	case "codex":
		return filepath.Join(home, ".codex", "history.jsonl")
	case "cursor":
		// Cursor는 workspaceStorage에 있지만 형식이 다를 수 있음
		// 기본적으로는 프로젝트별 .cursor 디렉토리 확인
		// TODO: Cursor 히스토리 파일 위치 정확히 파악 필요
		return ""
	default:
		return ""
	}
}

// findProjectRoot finds the project root by looking for CLAUDE.md or Claude.md
// It searches upward from the current directory
func findProjectRoot(startDir string) string {
	current := startDir
	for {
		// Check for CLAUDE.md or Claude.md
		claudeMd := filepath.Join(current, "CLAUDE.md")
		claudeMdLower := filepath.Join(current, "Claude.md")
		
		if _, err := os.Stat(claudeMd); err == nil {
			return current
		}
		if _, err := os.Stat(claudeMdLower); err == nil {
			return current
		}

		// Move to parent directory
		parent := filepath.Dir(current)
		if parent == current {
			// Reached filesystem root
			break
		}
		current = parent
	}
	return ""
}

