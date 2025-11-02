package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/curogom/curo-prompt/internal/evaluator"
	"github.com/curogom/curo-prompt/internal/model"
	"github.com/curogom/curo-prompt/internal/parser"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// newSuggestCmd creates the suggest command
func newSuggestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "suggest [flags]",
		Short: "프롬프트 개선 제안",
		Long: `프롬프트를 분석하고 개선 제안을 제공합니다.

예시:
  curo-prompt suggest --file prompts/dev_contract_v2.md
  cat prompt.md | curo-prompt suggest`,
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath, _ := cmd.Flags().GetString("file")     //nolint:errcheck
			provider, _ := cmd.Flags().GetString("provider") //nolint:errcheck
			apply, _ := cmd.Flags().GetBool("apply")         //nolint:errcheck

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

			// 프롬프트 파싱
			p := parser.NewParser()
			parsedPrompt, err := p.Parse(string(content))
			if err != nil {
				return fmt.Errorf("failed to parse prompt: %w", err)
			}

			// 평가 수행
			collectedPrompt := &model.CollectedPrompt{
				ID:         uuid.New().String(),
				Tool:       "suggest",
				RawPrompt:  string(content),
				Prompt:     parsedPrompt,
				Timestamp:  time.Now().Unix(),
				Command:    fmt.Sprintf("suggest --file %s", filePath),
				WorkingDir: getWorkingDir(),
			}

			eval := evaluator.NewEvaluator(provider)
			ctx := context.Background()

			result, err := eval.Evaluate(ctx, collectedPrompt)
			if err != nil {
				return fmt.Errorf("failed to evaluate: %w", err)
			}

			// 개선 제안 생성
			suggestions := generateSuggestions(result)

			// 출력
			cmd.Println("\n📋 개선 제안:")
			cmd.Println(strings.Repeat("=", 50))
			for i, suggestion := range suggestions {
				cmd.Printf("\n%d. %s\n", i+1, suggestion)
			}
			cmd.Println("\n" + strings.Repeat("=", 50))

			// 현재 점수 표시
			cmd.Printf("\n현재 점수: %.1f / 100\n", result.Score.OverallScore)

			// 자동 적용 옵션 (미구현)
			if apply {
				cmd.Println("\n⚠️  자동 적용 기능은 아직 구현되지 않았습니다.")
				cmd.Println("제안 사항을 수동으로 적용해주세요.")
			}

			return nil
		},
	}

	cmd.Flags().StringP("file", "f", "", "개선 제안할 프롬프트 파일 경로")
	cmd.Flags().StringP("provider", "p", "claude", "LLM Provider (토큰 계산용)")
	cmd.Flags().Bool("apply", false, "제안 사항 자동 적용 (M3에서 구현 예정)")

	return cmd
}

// generateSuggestions generates improvement suggestions based on evaluation result
func generateSuggestions(result *evaluator.EvaluationResult) []string {
	var suggestions []string

	analysis := result.Analysis

	// 누락된 섹션
	if !analysis.HasRole {
		suggestions = append(suggestions, "🔴 ROLE 섹션 추가 필요 - 프롬프트의 역할을 명확히 정의하세요")
	}
	if !analysis.HasInputs {
		suggestions = append(suggestions, "🟡 INPUTS 섹션 추가 권장 - 입력값을 명시하면 프롬프트가 더 명확해집니다")
	}
	if !analysis.HasInvariants {
		suggestions = append(suggestions, "🟡 INVARIANTS 섹션 추가 권장 - 규칙과 제약사항을 명시하세요")
	}
	if !analysis.HasOutputFormat {
		suggestions = append(suggestions, "🟡 OUTPUT FORMAT 섹션 추가 권장 - 출력 형식을 명시하면 결과 품질이 향상됩니다")
	}

	// 중복 규칙
	if len(analysis.DuplicateRules) > 0 {
		suggestions = append(suggestions, fmt.Sprintf("🟠 중복 규칙 %d개 발견 - 중복된 규칙을 제거하여 프롬프트를 간결하게 만드세요", len(analysis.DuplicateRules)))
	}

	// 토큰 수 최적화
	if result.TokenCount > 1000 {
		suggestions = append(suggestions, fmt.Sprintf("🟠 토큰 수가 많습니다 (%d 토큰) - 불필요한 내용을 제거하여 토큰을 절감하세요", result.TokenCount))
	}

	// 점수 기반 제안
	if result.Score.OverallScore < 60 {
		suggestions = append(suggestions, "🔴 종합 점수가 낮습니다 - 위 제안 사항들을 적용하여 프롬프트를 개선하세요")
	} else if result.Score.OverallScore < 80 {
		suggestions = append(suggestions, "🟡 종합 점수 개선 여지가 있습니다 - 제안 사항을 적용하여 더 높은 점수를 목표로 하세요")
	}

	// 메트릭별 제안
	if result.Score.Metrics.Structure < 70 {
		suggestions = append(suggestions, "🟠 구조 점수가 낮습니다 - 필수 섹션을 추가하고 중복을 제거하세요")
	}
	if result.Score.Metrics.Conciseness < 70 {
		suggestions = append(suggestions, "🟠 간결성 점수가 낮습니다 - 불필요한 설명을 제거하세요")
	}
	if result.Score.Metrics.Risk < 70 {
		suggestions = append(suggestions, "🟠 위험 점수가 낮습니다 - 모호한 표현을 구체적으로 수정하세요")
	}

	// 제안이 없으면
	if len(suggestions) == 0 {
		suggestions = append(suggestions, "✅ 프롬프트가 이미 잘 구성되어 있습니다!")
	}

	return suggestions
}
