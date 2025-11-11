package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/curogom/curompt/internal/collector"
	"github.com/curogom/curompt/internal/repository"
)

// newWrapCmd creates the wrap command for CLI tool collection
func newWrapCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wrap [tool] [command...]",
		Short: "외부 CLI 도구 래핑하여 프롬프트 수집",
		Long: `외부 CLI 도구(Codex, Cursor CLI 등)를 래핑하여 프롬프트를 자동 수집합니다.

예시:
  curompt wrap codex exec "TASK: Add feature"
  curompt wrap cursor chat "Implement login"

Claude Code CLI를 래핑하여 프롬프트를 자동으로 수집합니다.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return fmt.Errorf("사용법: curompt wrap [tool] [command] [args...]")
			}

			tool := args[0]
			command := args[1]
			commandArgs := args[2:]

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

			// Collector 생성
			collector := collector.NewCLIWrapperCollector(repo, tool)

			// 명령 래핑 및 실행
			ctx := context.Background()
			cmd.Printf("도구 래핑: %s\n", tool)
			cmd.Printf("명령 실행: %s %v\n\n", command, commandArgs)

			if err := collector.WrapCommand(ctx, command, commandArgs); err != nil {
				return fmt.Errorf("명령 실행 실패: %w", err)
			}

			cmd.Printf("\n✅ 프롬프트 수집 완료 (도구: %s)\n", tool)
			cmd.Printf("저장소 확인: curompt list --tool %s\n", tool)

			return nil
		},
	}

	return cmd
}
