package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/spf13/cobra"

	"github.com/curogom/curompt/internal/evaluator"
	"github.com/curogom/curompt/internal/model"
	"github.com/curogom/curompt/internal/reporter"
	"github.com/curogom/curompt/internal/repository"
)

// newScanCmd creates the scan command
func newScanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scan [flags]",
		Short: "수집된 프롬프트를 분석하고 리포트를 생성",
		Long: `저장소(DB)에 저장된 프롬프트 기록을 불러와 정적 분석/점수화를 수행합니다.

예시:
  curompt scan                # 전체 히스토리 분석
  curompt scan --path ~/dev   # 특정 경로에서 실행된 프롬프트만 분석
  curompt scan --path repo --output reports/`,
		RunE: func(cmd *cobra.Command, args []string) error {
			pathFilter, _ := cmd.Flags().GetString("path")         //nolint:errcheck
			legacyRepo, _ := cmd.Flags().GetString("repo")         //nolint:errcheck
			outputDir, _ := cmd.Flags().GetString("output")        //nolint:errcheck
			singleOut, _ := cmd.Flags().GetString("single-output") //nolint:errcheck
			doClean, _ := cmd.Flags().GetBool("clean")             //nolint:errcheck
			provider, _ := cmd.Flags().GetString("provider")       //nolint:errcheck
			concurrency, _ := cmd.Flags().GetInt("concurrency")    //nolint:errcheck
			fullOutput, _ := cmd.Flags().GetBool("full")           //nolint:errcheck
			maxLines, _ := cmd.Flags().GetInt("max-lines")         //nolint:errcheck
			topN, _ := cmd.Flags().GetInt("top")                   //nolint:errcheck
			summaryMode, _ := cmd.Flags().GetString("summary")     //nolint:errcheck
			showAdvice, _ := cmd.Flags().GetBool("advice")         //nolint:errcheck

			if pathFilter == "" {
				pathFilter = legacyRepo
			}
			pathFilter = strings.TrimSpace(pathFilter)

			var targetPath string
			if pathFilter != "" {
				abs, err := filepath.Abs(pathFilter)
				if err != nil {
					return fmt.Errorf("경로 해석 실패: %w", err)
				}
				targetPath = filepath.Clean(abs)
			}

			// 출력 디렉터리는 사용자가 지정한 경우에만 사용
			useFileOutput := outputDir != ""

			// 출력 디렉토리 준비 및 선택적 정리
			if useFileOutput {
				if err := os.MkdirAll(outputDir, 0755); err != nil {
					return fmt.Errorf("failed to create output directory: %w", err)
				}
				if doClean {
					entries, _ := os.ReadDir(outputDir) //nolint:errcheck
					for _, e := range entries {
						if !e.IsDir() && strings.HasSuffix(strings.ToLower(e.Name()), ".md") {
							_ = os.Remove(filepath.Join(outputDir, e.Name())) //nolint:errcheck
						}
					}
				}
			}

			// 저장소 초기화 (필수)
			dbPath := filepath.Join(os.Getenv("HOME"), ".curompt", "db.sqlite")
			if _, err := os.Stat(dbPath); os.IsNotExist(err) {
				return fmt.Errorf("저장된 프롬프트가 없습니다. 먼저 collect 명령으로 히스토리를 수집하세요")
			}

			if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
				return fmt.Errorf("저장소 경로 생성 실패: %w", err)
			}

	repo, err := repository.NewSQLiteRepository(dbPath)
	if err != nil {
		return fmt.Errorf("저장소 초기화 실패: %w", err)
	}
	defer func() {
		_ = repo.Close()
	}()

			ctx := context.Background()

			allPrompts, err := repo.FindAll(ctx)
			if err != nil {
				return fmt.Errorf("프롬프트 로드 실패: %w", err)
			}

			filteredPrompts := filterPromptsByPath(allPrompts, targetPath)
			if len(filteredPrompts) == 0 {
				added, err := maybeAutoCollect(ctx, repo, targetPath, cmd)
				if err != nil {
					return err
				}
				if added == 0 {
					cmd.Printf("분석할 프롬프트가 없습니다. collect 또는 eval 명령으로 데이터를 저장한 뒤 다시 실행하세요.\n")
					return nil
				}
				allPrompts, err = repo.FindAll(ctx)
				if err != nil {
					return fmt.Errorf("프롬프트 재조회 실패: %w", err)
				}
				filteredPrompts = filterPromptsByPath(allPrompts, targetPath)
				if len(filteredPrompts) == 0 {
					cmd.Printf("새로 수집한 프롬프트 중 조건과 일치하는 항목이 없습니다.\n")
					return nil
				}
			}

			cmd.Printf("선택된 프롬프트: %d개\n\n", len(filteredPrompts))

			// 병렬 평가 설정
			if concurrency <= 0 {
				concurrency = runtime.NumCPU()
				if concurrency < 1 {
					concurrency = 1
				}
			}

			type evalResult struct {
				prompt      *model.CollectedPrompt
				displayPath string
				report      string
				score       float64
				err         error
			}

			eval := evaluator.NewEvaluator(provider)
			markdownReporter := reporter.NewMarkdownReporter()

			jobs := make(chan *model.CollectedPrompt)
			var wg sync.WaitGroup
			var mu sync.Mutex
			var results []evalResult

			worker := func() {
				defer wg.Done()
				for prompt := range jobs {
					result, err := eval.Evaluate(ctx, prompt)
					if err != nil {
						mu.Lock()
						results = append(results, evalResult{prompt: prompt, displayPath: resolvePromptDisplayPath(prompt), err: fmt.Errorf("평가 실패: %w", err)})
						mu.Unlock()
						continue
					}
					report, err := markdownReporter.Generate(result)
					if err != nil {
						mu.Lock()
						results = append(results, evalResult{prompt: prompt, displayPath: resolvePromptDisplayPath(prompt), err: fmt.Errorf("리포트 생성 실패: %w", err)})
						mu.Unlock()
						continue
					}
					mu.Lock()
					results = append(results, evalResult{
						prompt:      prompt,
						displayPath: resolvePromptDisplayPath(prompt),
						report:      report,
						score:       result.Score.OverallScore,
					})
					mu.Unlock()
				}
			}

			for i := 0; i < concurrency; i++ {
				wg.Add(1)
				go worker()
			}
			for _, p := range filteredPrompts {
				jobs <- p
			}
			close(jobs)
			wg.Wait()

			var mergedReports []string
			successCount := 0
			linesPrinted := 0
			for _, r := range results {
				if r.err != nil {
					cmd.Printf("  ⚠️  %s: %v\n", r.displayPath, r.err)
					continue
				}

				if useFileOutput && singleOut == "" {
					filenameSeed := filepath.Base(r.displayPath)
					if filenameSeed == "" {
						filenameSeed = r.prompt.ID[:8]
					}
					reportFileName := fmt.Sprintf("%s_report.md", sanitizeFileName(filenameSeed))
					reportPath := filepath.Join(outputDir, reportFileName)
					if err := os.WriteFile(reportPath, []byte(r.report), 0644); err != nil {
						cmd.Printf("  ⚠️  리포트 저장 실패: %v\n", err)
						continue
					}
					cmd.Printf("  ✅ 완료 - 점수: %.1f/100, 리포트: %s\n", r.score, reportPath)
				} else if useFileOutput && singleOut != "" {
					header := fmt.Sprintf("# %s\n\n", r.displayPath)
					mergedReports = append(mergedReports, header+r.report+"\n\n")
					cmd.Printf("  ✅ 완료 - 점수: %.1f/100\n", r.score)
				} else {
					// 콘솔 출력: 기본은 요약(집계 및 하위 점수 Top-N), --full 시 전체 리포트
					if fullOutput {
						// 전체 출력 강제
						sep := strings.Repeat("-", 80)
						cmd.Printf("%s\n프롬프트: %s\n%s\n%s\n\n", sep, r.displayPath, sep, r.report)
						cmd.Printf("  ✅ 완료 - 점수: %.1f/100 (stdout 출력)\n", r.score)
					}
				}
				successCount++
			}

			// 단일 파일 병합 저장
			if useFileOutput && singleOut != "" && len(mergedReports) > 0 {
				merged := strings.Join(mergedReports, "")
				outPath := filepath.Join(outputDir, singleOut)
				if err := os.WriteFile(outPath, []byte(merged), 0644); err != nil {
					cmd.Printf("  ⚠️  단일 리포트 저장 실패: %v\n", err)
				} else {
					cmd.Printf("단일 리포트 저장: %s\n", outPath)
				}
			}

			// 콘솔 요약 출력 처리 (useFileOutput이 아닌 경우)
			if !useFileOutput && !fullOutput {
				if topN <= 0 {
					topN = 10
				}
				// 통계: 평균/최고/최저/중앙값/표준편차/퍼센타일
				var sum, min, max float64
				count := 0
				first := true
				var scores []float64
				for _, r := range results {
					if r.err != nil {
						continue
					}
					s := r.score
					sum += s
					if first || s < min {
						min = s
					}
					if first || s > max {
						max = s
					}
					first = false
					count++
					scores = append(scores, s)
				}
				avg := 0.0
				if count > 0 {
					avg = sum / float64(count)
				}

				// 보조 통계 계산
				calcMedian := func(vals []float64) float64 {
					if len(vals) == 0 {
						return 0
					}
					v := append([]float64(nil), vals...)
					sort.Float64s(v)
					m := len(v) / 2
					if len(v)%2 == 0 {
						return (v[m-1] + v[m]) / 2
					}
					return v[m]
				}
				calcStdDev := func(vals []float64, mean float64) float64 {
					if len(vals) == 0 {
						return 0
					}
					var ss float64
					for _, x := range vals {
						d := x - mean
						ss += d * d
					}
					return ss / float64(len(vals))
				}
				percentile := func(vals []float64, p float64) float64 {
					if len(vals) == 0 {
						return 0
					}
					v := append([]float64(nil), vals...)
					sort.Float64s(v)
					if p <= 0 {
						return v[0]
					}
					if p >= 1 {
						return v[len(v)-1]
					}
					idx := int(p*float64(len(v)-1) + 0.5)
					if idx < 0 {
						idx = 0
					}
					if idx >= len(v) {
						idx = len(v) - 1
					}
					return v[idx]
				}

				median := calcMedian(scores)
				stddev := calcStdDev(scores, avg)
				p25 := percentile(scores, 0.25)
				p75 := percentile(scores, 0.75)

				// 하위/상위 점수 Top-N 선택
				good := make([]evalResult, 0, len(results))
				for _, r := range results {
					if r.err == nil {
						good = append(good, r)
					}
				}
				sort.Slice(good, func(i, j int) bool { return good[i].score < good[j].score })
				if topN > len(good) {
					topN = len(good)
				}

				// 분포 요약: 구간별 개수
				var b0_39, b40_59, b60_79, b80_100 int
				for _, s := range scores {
					switch {
					case s < 40:
						b0_39++
					case s < 60:
						b40_59++
					case s < 80:
						b60_79++
					default:
						b80_100++
					}
				}

				// 요약 모드 처리
				if summaryMode == "short" {
					cmd.Printf("요약: 총 %d, 평균 %.1f, 중앙값 %.1f, 최저 %.1f, 최고 %.1f\n", count, avg, median, min, max)
				} else {
					cmd.Printf("요약: 총 %d, 평균 %.1f, 중앙값 %.1f, 표준편차 %.1f, 최저 %.1f, 최고 %.1f\n", count, avg, median, stddev, min, max)
					cmd.Printf("분포: 0–39:%d, 40–59:%d, 60–79:%d, 80–100:%d  |  IQR(25%%:%.1f, 75%%:%.1f)\n", b0_39, b40_59, b60_79, b80_100, p25, p75)
					if topN > 0 {
						cmd.Printf("개선 우선순위 Top-%d:\n", topN)
						limit := topN
						if maxLines > 0 && limit > maxLines {
							limit = maxLines
						}
						for i := 0; i < limit; i++ {
							cmd.Printf("%2d) %.1f  %s\n", i+1, good[i].score, good[i].displayPath)
							linesPrinted++
							if maxLines > 0 && linesPrinted >= maxLines {
								cmd.Printf("... (요약 출력 제한 %d라인에 도달)\n", maxLines)
								break
							}
						}
						// 상위 점수 Top-5 칭찬 목록
						if linesPrinted < maxLines {
							bestN := 5
							if bestN > len(good) {
								bestN = len(good)
							}
							if bestN > 0 {
								cmd.Printf("\n잘 작성된 프롬프트 Top-%d:\n", bestN)
								// good은 오름차순이므로 뒤에서부터 선택
								for i := 0; i < bestN; i++ {
									idx := len(good) - 1 - i
									cmd.Printf("%2d) %.1f  %s\n", i+1, good[idx].score, good[idx].displayPath)
									linesPrinted++
									if maxLines > 0 && linesPrinted >= maxLines {
										cmd.Printf("... (요약 출력 제한 %d라인에 도달)\n", maxLines)
										break
									}
								}
							}
						}
					}

					// 코칭형 조언 블록 (정적 가이드)
					if showAdvice && (maxLines == 0 || linesPrinted+8 <= maxLines) {
						cmd.Printf("\n개선 가이드:\n")
						cmd.Printf("- ROLE을 명확히: 한 문장으로 기대 역할/톤을 정의하세요.\n")
						cmd.Printf("- INPUTS를 구조화: 입력 항목을 불릿/키:타입 형태로 명시하세요.\n")
						cmd.Printf("- INVARIANTS 분리: 절대 규칙/금지어를 별도 섹션에 모으세요.\n")
						cmd.Printf("- OUTPUT FORMAT 스키마: JSON 스키마/키 순서/예제를 함께 제시하세요.\n")
						cmd.Printf("- 간결성 유지: 중복 규칙 제거, 길면 요약/표로 정리하세요.\n")
						cmd.Printf("- 리스크 관리: 토큰/비용 고려, 민감정보는 마스킹 후 예시 제공하세요.\n")
					}
				}
			}

			cmd.Printf("\n총 %d개 프롬프트 분석 완료\n", successCount)
			if useFileOutput {
				cmd.Printf("리포트 저장 위치: %s\n", outputDir)
			}

			return nil
		},
	}

	cmd.Flags().String("path", "", "특정 경로에서 생성된 프롬프트만 분석")
	cmd.Flags().String("repo", "", "(deprecated) --path로 대체되었습니다")
	if err := cmd.Flags().MarkHidden("repo"); err != nil {
		cmd.Printf("Warning: failed to hide --repo flag: %v\n", err)
	}
	cmd.Flags().StringP("output", "o", "", "리포트 출력 디렉토리 (미지정 시 콘솔 출력)")
	cmd.Flags().String("single-output", "", "여러 리포트를 하나의 파일로 병합해 저장 (— --output 필요)")
	cmd.Flags().Bool("clean", false, "--output 디렉토리의 기존 .md 리포트를 정리하고 새로 생성")
	cmd.Flags().StringP("provider", "p", "claude", "LLM Provider (토큰 계산용)")
	cmd.Flags().Int("concurrency", 0, "동시 평가 작업 수 (기본값: CPU 코어 수)")
	cmd.Flags().Bool("full", false, "콘솔에 전체 리포트를 출력 (기본은 요약)")
	cmd.Flags().Int("max-lines", 100, "콘솔 요약 최대 출력 라인 수 (기본: 100)")
	cmd.Flags().Int("top", 10, "요약 출력 시 하위 점수 Top-N 표시 (기본: 10)")
	cmd.Flags().String("summary", "rich", "요약 모드 (rich|short)")
	cmd.Flags().Bool("advice", true, "요약에 코칭형 개선 가이드 출력")

	return cmd
}

// sanitizeFileName sanitizes file name for report output
func sanitizeFileName(name string) string {
	name = strings.TrimSuffix(name, filepath.Ext(name))
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "\\", "_")
	return name
}

func filterPromptsByPath(prompts []*model.CollectedPrompt, target string) []*model.CollectedPrompt {
	if target == "" {
		return prompts
	}
	var filtered []*model.CollectedPrompt
	for _, p := range prompts {
		projectPath := resolvePromptProjectPath(p)
		if projectPath == "" {
			continue
		}
		if isPathWithin(normalizeAbsPath(projectPath), target) {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

func resolvePromptProjectPath(p *model.CollectedPrompt) string {
	if p.Metadata != nil {
		if project := p.Metadata["project"]; project != "" {
			return project
		}
		if filePath := p.Metadata["file_path"]; filePath != "" {
			return filepath.Dir(filePath)
		}
	}
	if p.WorkingDir != "" {
		return p.WorkingDir
	}
	return ""
}

func resolvePromptDisplayPath(p *model.CollectedPrompt) string {
	if p.Metadata != nil {
		if filePath := p.Metadata["file_path"]; filePath != "" {
			return filePath
		}
		if project := p.Metadata["project"]; project != "" {
			return project
		}
	}
	if p.WorkingDir != "" {
		return p.WorkingDir
	}
	return p.ID
}

func normalizeAbsPath(path string) string {
	if path == "" {
		return ""
	}
	if !filepath.IsAbs(path) {
		abs, err := filepath.Abs(path)
		if err == nil {
			path = abs
		}
	}
	return filepath.Clean(path)
}

func isPathWithin(child, parent string) bool {
	if parent == "" {
		return true
	}
	parent = normalizeAbsPath(parent)
	child = normalizeAbsPath(child)
	if child == "" || parent == "" {
		return false
	}
	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}
	return rel == "." || (!strings.HasPrefix(rel, ".."))
}

func maybeAutoCollect(ctx context.Context, repo repository.PromptRepository, targetPath string, cmd *cobra.Command) (int, error) {
	if !isInteractive() {
		cmd.Printf("경로에 해당하는 프롬프트가 없습니다. 먼저 'curompt collect --from claude|codex'를 실행하세요.\n")
		return 0, nil
	}

	message := "해당 경로에서 저장된 프롬프트가 없습니다. 지금 히스토리를 수집할까요?"
	proceed, err := promptYesNo(message, false)
	if err != nil {
		return 0, err
	}
	if !proceed {
		return 0, nil
	}

	source, err := promptSourceSelection([]autoCollectOption{
		{Label: "Claude Code history.jsonl", Value: "claude"},
		{Label: "Codex CLI history.jsonl", Value: "codex"},
	})
	if err != nil {
		return 0, err
	}
	if source == "" {
		cmd.Printf("수집이 취소되었습니다.\n")
		return 0, nil
	}

	logPath := getDefaultLogPath(source)
	if logPath == "" || !fileExists(logPath) {
		logPath, err = promptLine("로그 파일 경로를 입력하세요", logPath)
		if err != nil {
			return 0, err
		}
	}
	if logPath == "" {
		return 0, fmt.Errorf("로그 파일 경로가 필요합니다")
	}

	projectFilter := defaultProjectPath(source, targetPath)
	if source == "claude" && projectFilter == "" {
		projectFilter, err = promptLine("CLAUDE.md가 있는 프로젝트 경로를 입력하세요", "")
		if err != nil {
			return 0, err
		}
	}
	projectFilter = normalizeAbsPath(projectFilter)

	_, saved, err := collectFromLog(ctx, repo, source, logPath, projectFilter, cmd)
	if err != nil {
		return 0, err
	}
	return saved, nil
}
