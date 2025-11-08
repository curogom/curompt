package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewRootCmd creates and returns the root command
func NewRootCmd(version, buildTime, gitCommit string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "curompt",
		Short: "CLI 기반 LLM 프롬프트 분석·평가·최적화 도구",
		Long: `curompt는 CLI 기반 개발자의 LLM 프롬프트를 분석·평가·최적화하는 도구입니다.

주요 기능:
- 프롬프트 정적 분석 (섹션 구조, 중복 규칙, 토큰 계산)
- 동적 평가 (스키마 적합률, 자체 일관성, 지연·비용)
- 점수화 (0-100 종합 점수)
- 리팩터 제안 (토큰 절감, 규칙 분리, few-shot 축약)`,
		Version: fmt.Sprintf("%s (built: %s, commit: %s)", version, buildTime, gitCommit),
		Run: func(cmd *cobra.Command, args []string) {
			// 기본 동작: 도움말 출력
			if err := cmd.Help(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	// 서브커맨드 추가
	rootCmd.AddCommand(newScanCmd())
	rootCmd.AddCommand(newEvalCmd())
	rootCmd.AddCommand(newSuggestCmd())
	rootCmd.AddCommand(newListCmd())
	rootCmd.AddCommand(newWrapCmd())
	rootCmd.AddCommand(newCollectCmd())

	return rootCmd
}
