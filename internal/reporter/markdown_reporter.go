package reporter

import (
	"fmt"
	"strings"
	"time"

	"github.com/curogom/curompt/internal/evaluator"
)

// markdownReporter implements Reporter for Markdown format
type markdownReporter struct{}

// NewMarkdownReporter creates a new Markdown reporter
func NewMarkdownReporter() Reporter {
	return &markdownReporter{}
}

// Name returns the reporter name
func (r *markdownReporter) Name() string {
	return "markdown"
}

// Generate generates a Markdown report
func (r *markdownReporter) Generate(result *evaluator.EvaluationResult) (string, error) {
	var sb strings.Builder

	// 헤더
	sb.WriteString("# 프롬프트 평가 리포트\n\n")

	// 메타데이터
	sb.WriteString("## 메타데이터\n\n")
	sb.WriteString(fmt.Sprintf("- **도구**: %s\n", result.Prompt.Tool))
	sb.WriteString(fmt.Sprintf("- **수집 시간**: %s\n", time.Unix(result.Prompt.Timestamp, 0).Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("- **명령어**: `%s`\n", result.Prompt.Command))
	if result.Prompt.WorkingDir != "" {
		sb.WriteString(fmt.Sprintf("- **작업 디렉토리**: %s\n", result.Prompt.WorkingDir))
	}
	sb.WriteString("\n")

	// 종합 점수
	sb.WriteString("## 종합 점수\n\n")
	sb.WriteString(fmt.Sprintf("### %.1f / 100\n\n", result.Score.OverallScore))

	// 메트릭 상세
	sb.WriteString("## 메트릭 상세\n\n")
	sb.WriteString("| 메트릭 | 점수 | 가중치 | 기여도 |\n")
	sb.WriteString("|--------|------|--------|--------|\n")
	sb.WriteString(fmt.Sprintf("| Structure | %.1f | 20%% | %.1f |\n",
		result.Score.Metrics.Structure, result.Score.Metrics.Structure*0.2))
	sb.WriteString(fmt.Sprintf("| Conciseness | %.1f | 15%% | %.1f |\n",
		result.Score.Metrics.Conciseness, result.Score.Metrics.Conciseness*0.15))
	sb.WriteString(fmt.Sprintf("| Instruction Following | %.1f | 30%% | %.1f |\n",
		result.Score.Metrics.InstructionFollowing, result.Score.Metrics.InstructionFollowing*0.3))
	sb.WriteString(fmt.Sprintf("| Risk | %.1f | 10%% | %.1f |\n",
		result.Score.Metrics.Risk, result.Score.Metrics.Risk*0.1))
	sb.WriteString("\n")

	// 분석 결과
	sb.WriteString("## 정적 분석 결과\n\n")
	sb.WriteString(fmt.Sprintf("- **ROLE 섹션**: %s\n", r.boolToStatus(result.Analysis.HasRole)))
	sb.WriteString(fmt.Sprintf("- **INPUTS 섹션**: %s\n", r.boolToStatus(result.Analysis.HasInputs)))
	sb.WriteString(fmt.Sprintf("- **INVARIANTS 섹션**: %s\n", r.boolToStatus(result.Analysis.HasInvariants)))
	sb.WriteString(fmt.Sprintf("- **OUTPUT FORMAT 섹션**: %s\n", r.boolToStatus(result.Analysis.HasOutputFormat)))
	sb.WriteString(fmt.Sprintf("- **섹션 수**: %d\n", result.Analysis.SectionCount))

	if len(result.Analysis.DuplicateRules) > 0 {
		sb.WriteString(fmt.Sprintf("- **중복 규칙**: %d개 발견\n", len(result.Analysis.DuplicateRules)))
		for _, dup := range result.Analysis.DuplicateRules {
			sb.WriteString(fmt.Sprintf("  - %s\n", dup))
		}
	}
	sb.WriteString("\n")

	// 토큰 정보
	sb.WriteString("## 토큰 정보\n\n")
	sb.WriteString(fmt.Sprintf("- **토큰 수**: %d\n", result.TokenCount))
	sb.WriteString(fmt.Sprintf("- **토큰 계산기**: %s\n", result.TokenProvider))
	sb.WriteString("\n")

	// 가이드 준수 분석 (있는 경우)
	if result.GuideCompliance != nil {
		sb.WriteString("## 프롬프트 가이드 준수 분석\n\n")
		sb.WriteString("| 항목 | 점수 |\n")
		sb.WriteString("|------|------|\n")
		sb.WriteString(fmt.Sprintf("| 명확성 (Clarity) | %.1f/100 |\n", result.GuideCompliance.ClarityScore))
		sb.WriteString(fmt.Sprintf("| 컨텍스트 제공 (Context) | %.1f/100 |\n", result.GuideCompliance.ContextScore))
		sb.WriteString(fmt.Sprintf("| 예시 제공 (Examples) | %.1f/100 |\n", result.GuideCompliance.ExampleScore))
		sb.WriteString(fmt.Sprintf("| 구조화 (Structure) | %.1f/100 |\n", result.GuideCompliance.StructureScore))
		sb.WriteString(fmt.Sprintf("| 제약 명시 (Constraints) | %.1f/100 |\n", result.GuideCompliance.ConstraintScore))
		sb.WriteString("\n")

		// 모호한 표현 감지 시 표시
		if len(result.GuideCompliance.AmbiguousPhrases) > 0 {
			sb.WriteString("### 감지된 모호한 표현\n")
			for i, phrase := range result.GuideCompliance.AmbiguousPhrases {
				if i < 3 { // 최대 3개만 표시
					sb.WriteString(fmt.Sprintf("- `%s`\n", phrase))
				}
			}
			sb.WriteString("\n")
		}
	}

	// 개선 제안
	sb.WriteString("## 개선 제안\n\n")

	// 가이드 기반 개선 제안 (우선 표시)
	if result.GuideCompliance != nil && len(result.GuideCompliance.Suggestions) > 0 {
		for _, suggestion := range result.GuideCompliance.Suggestions {
			sb.WriteString(fmt.Sprintf("%s\n", suggestion))
		}
		sb.WriteString("\n")
	}

	// 누락된 섹션 체크
	missingSections := []string{}
	if !result.Analysis.HasInputs {
		missingSections = append(missingSections, "INPUTS")
	}
	if !result.Analysis.HasInvariants {
		missingSections = append(missingSections, "INVARIANTS")
	}
	if !result.Analysis.HasOutputFormat {
		missingSections = append(missingSections, "OUTPUT FORMAT")
	}

	if len(missingSections) > 0 {
		sb.WriteString("### 누락된 섹션\n")
		for _, section := range missingSections {
			sb.WriteString(fmt.Sprintf("- `%s` 섹션을 추가하세요\n", section))
		}
		sb.WriteString("\n")
	}

	if len(result.Analysis.DuplicateRules) > 0 {
		sb.WriteString("### 중복 제거\n")
		sb.WriteString("- 중복된 규칙을 제거하여 프롬프트를 간결하게 만드세요\n")
		sb.WriteString("\n")
	}

	if result.Score.Metrics.Conciseness < 70 {
		sb.WriteString("### 토큰 절감\n")
		sb.WriteString(fmt.Sprintf("- 현재 토큰 수가 많습니다 (%d 토큰). 불필요한 내용을 제거하세요.\n", result.TokenCount))
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

// boolToStatus converts boolean to status string
func (r *markdownReporter) boolToStatus(b bool) string {
	if b {
		return "✅ 있음"
	}
	return "❌ 없음"
}
