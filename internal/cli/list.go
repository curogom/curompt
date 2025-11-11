package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/spf13/cobra"

	"github.com/curogom/curompt/internal/evaluator"
	"github.com/curogom/curompt/internal/model"
	"github.com/curogom/curompt/internal/reporter"
	"github.com/curogom/curompt/internal/repository"
)

// newListCmd creates the list command for viewing stored prompts
func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [flags]",
		Short: "저장된 프롬프트 목록 조회",
		Long: `저장소에 저장된 프롬프트 목록을 조회합니다.

예시:
  curompt list                    # 최근 프롬프트 10개
  curompt list --limit 20         # 최근 프롬프트 20개
  curompt list --tool codex       # codex로 수집된 프롬프트만
  curompt list --tool cursor      # cursor로 수집된 프롬프트만
  curompt list --eval             # 목록 조회 후 각각 평가 수행`,
		RunE: func(cmd *cobra.Command, args []string) error {
			limit, _ := cmd.Flags().GetInt("limit")          //nolint:errcheck
			tool, _ := cmd.Flags().GetString("tool")         //nolint:errcheck
			eval, _ := cmd.Flags().GetBool("eval")           //nolint:errcheck
			output, _ := cmd.Flags().GetString("output")     //nolint:errcheck
			provider, _ := cmd.Flags().GetString("provider") //nolint:errcheck

			if limit <= 0 {
				limit = 10
			}

			// 저장소 초기화
			dbPath := filepath.Join(os.Getenv("HOME"), ".curompt", "db.sqlite")
			if _, err := os.Stat(dbPath); os.IsNotExist(err) {
				return fmt.Errorf("저장소가 초기화되지 않았습니다. 먼저 프롬프트를 수집하거나 스캔해주세요")
			}

			repo, err := repository.NewSQLiteRepository(dbPath)
			if err != nil {
				return fmt.Errorf("저장소 초기화 실패: %w", err)
			}
			defer func() {
				_ = repo.Close()
			}()

			ctx := context.Background()

			var prompts []*PromptListItem
			if tool != "" {
				// 도구별 조회
				all, err := repo.FindByTool(ctx, tool)
				if err != nil {
					return fmt.Errorf("조회 실패: %w", err)
				}
				// 최신순 정렬 및 제한
				sort.Slice(all, func(i, j int) bool {
					return all[i].Timestamp > all[j].Timestamp
				})
				if len(all) > limit {
					all = all[:limit]
				}
				prompts = convertToListItem(all)
			} else {
				// 최근 프롬프트 조회
				recent, err := repo.FindRecent(ctx, limit)
				if err != nil {
					return fmt.Errorf("조회 실패: %w", err)
				}
				prompts = convertToListItem(recent)
			}

			if len(prompts) == 0 {
				cmd.Printf("저장된 프롬프트가 없습니다.\n")
				cmd.Printf("\n프롬프트를 수집하려면:\n")
				cmd.Printf("  - collect 명령으로 Claude/Codex 로그를 수집하세요.\n")
				cmd.Printf("  - eval 명령으로 개별 프롬프트를 평가해 저장하세요.\n")
				return nil
			}

			// 목록 출력
			cmd.Printf("저장된 프롬프트: %d개\n\n", len(prompts))
			for i, item := range prompts {
				cmd.Printf("[%d] ID: %s\n", i+1, item.ID[:8])
				cmd.Printf("    도구: %s\n", item.Tool)
				cmd.Printf("    시간: %s\n", time.Unix(item.Timestamp, 0).Format("2006-01-02 15:04:05"))
				if item.Command != "" {
					cmd.Printf("    명령: %s\n", item.Command)
				}
				if item.WorkingDir != "" {
					cmd.Printf("    경로: %s\n", item.WorkingDir)
				}
				cmd.Printf("    프롬프트: %s\n", truncate(item.RawPrompt, 80))
				cmd.Println()
			}

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
				for i, item := range prompts {
					cmd.Printf("[%d/%d] 평가 중: %s\n", i+1, len(prompts), item.ID[:8])

					// 프롬프트 로드 (전체 데이터)
					full, err := repo.FindByID(ctx, item.ID)
					if err != nil {
						cmd.Printf("  ⚠️  로드 실패: %v\n", err)
						continue
					}

					// 평가
					result, err := evaluator.Evaluate(ctx, full)
					if err != nil {
						cmd.Printf("  ⚠️  평가 실패: %v\n", err)
						continue
					}

					// 리포트 생성
					report, err := reporter.Generate(result)
					if err != nil {
						cmd.Printf("  ⚠️  리포트 생성 실패: %v\n", err)
						continue
					}

					// 리포트 저장
					reportPath := filepath.Join(outputDir, fmt.Sprintf("%s_report.md", item.ID[:8]))
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

	cmd.Flags().IntP("limit", "l", 10, "조회할 프롬프트 개수")
	cmd.Flags().StringP("tool", "t", "", "특정 도구로 수집된 프롬프트만 조회 (codex, cursor, scan 등)")
	cmd.Flags().Bool("eval", false, "조회한 프롬프트들을 평가하여 리포트 생성")
	cmd.Flags().StringP("output", "o", "reports", "리포트 출력 디렉토리 (--eval 사용 시)")
	cmd.Flags().StringP("provider", "p", "claude", "LLM Provider (토큰 계산용, --eval 사용 시)")

	return cmd
}

// PromptListItem is a simplified version for list display
type PromptListItem struct {
	ID         string
	Tool       string
	Timestamp  int64
	Command    string
	WorkingDir string
	RawPrompt  string
}

// convertToListItem converts CollectedPrompt to PromptListItem
func convertToListItem(prompts []*model.CollectedPrompt) []*PromptListItem {
	items := make([]*PromptListItem, len(prompts))
	for i, p := range prompts {
		items[i] = &PromptListItem{
			ID:         p.ID,
			Tool:       p.Tool,
			Timestamp:  p.Timestamp,
			Command:    p.Command,
			WorkingDir: p.WorkingDir,
			RawPrompt:  p.RawPrompt,
		}
	}
	return items
}

// truncate truncates string to max length
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
