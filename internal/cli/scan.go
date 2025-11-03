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
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/curogom/curo-prompt/internal/evaluator"
	"github.com/curogom/curo-prompt/internal/model"
	"github.com/curogom/curo-prompt/internal/reporter"
	"github.com/curogom/curo-prompt/internal/repository"
)

// newScanCmd creates the scan command
func newScanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scan [flags]",
		Short: "레포지토리 내 프롬프트 파일 스캔 및 분석",
		Long: `레포지토리 내 프롬프트 파일을 찾아 분석하고 리포트를 생성합니다.

예시:
  curo-prompt scan --repo .
  curo-prompt scan --repo ./prompts --output reports/`,
		RunE: func(cmd *cobra.Command, args []string) error {
			repoPath, _ := cmd.Flags().GetString("repo")           //nolint:errcheck
			outputDir, _ := cmd.Flags().GetString("output")        //nolint:errcheck
			singleOut, _ := cmd.Flags().GetString("single-output") //nolint:errcheck
			doClean, _ := cmd.Flags().GetBool("clean")             //nolint:errcheck
			patterns, _ := cmd.Flags().GetStringSlice("patterns")  //nolint:errcheck
			provider, _ := cmd.Flags().GetString("provider")       //nolint:errcheck
			concurrency, _ := cmd.Flags().GetInt("concurrency")    //nolint:errcheck
			fullOutput, _ := cmd.Flags().GetBool("full")           //nolint:errcheck
			maxLines, _ := cmd.Flags().GetInt("max-lines")         //nolint:errcheck
			topN, _ := cmd.Flags().GetInt("top")                   //nolint:errcheck
			summaryMode, _ := cmd.Flags().GetString("summary")     //nolint:errcheck
			showAdvice, _ := cmd.Flags().GetBool("advice")         //nolint:errcheck

			if repoPath == "" {
				repoPath = "."
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

			// 프롬프트 파일 찾기 (배치 수집)
			files, err := findPromptFiles(repoPath, patterns)
			if err != nil {
				return fmt.Errorf("failed to find prompt files: %w", err)
			}

			if len(files) == 0 {
				cmd.Printf("프롬프트 파일을 찾을 수 없습니다. (경로: %s, 패턴: %v)\n", repoPath, patterns)
				return nil
			}

			cmd.Printf("발견된 프롬프트 파일: %d개\n\n", len(files))

			// 저장소 초기화 (선택적)
			dbPath := filepath.Join(os.Getenv("HOME"), ".curo-prompt", "db.sqlite")
			if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
				cmd.Printf("Warning: failed to create db directory: %v\n", err)
				dbPath = "" // 저장소 없이 진행
			}

			var repo repository.PromptRepository
			if dbPath != "" {
				repo, err = repository.NewSQLiteRepository(dbPath)
				if err != nil {
					cmd.Printf("Warning: failed to initialize repository: %v\n", err)
				} else {
					defer repo.Close()
				}
			}

			// 병렬 평가 설정
			if concurrency <= 0 {
				concurrency = runtime.NumCPU()
				if concurrency < 1 {
					concurrency = 1
				}
			}

			type evalResult struct {
				filePath string
				absPath  string
				report   string
				score    float64
				err      error
			}

			eval := evaluator.NewEvaluator(provider)
			ctx := context.Background()
			markdownReporter := reporter.NewMarkdownReporter()

			// Serialize DB writes to avoid SQLITE_BUSY under parallel evaluation
			var saveWg sync.WaitGroup
			saveJobs := make(chan *model.CollectedPrompt, 128)
			if repo != nil {
				saveWg.Add(1)
				go func() {
					defer saveWg.Done()
					for cp := range saveJobs {
						// retry with simple backoff
						var lastErr error
						for attempt := 0; attempt < 5; attempt++ {
							if err := repo.Save(ctx, cp); err != nil {
								lastErr = err
								time.Sleep(time.Duration(50*(1<<attempt)) * time.Millisecond)
								continue
							}
							lastErr = nil
							break
						}
						if lastErr != nil {
							// non-fatal: print a warning
							cmd.Printf("  ⚠️  저장 실패(최대 재시도 후): %v\n", lastErr)
						}
					}
				}()
			}

			jobs := make(chan string)
			var wg sync.WaitGroup
			var mu sync.Mutex
			var results []evalResult

			worker := func() {
				defer wg.Done()
				for file := range jobs {
					content, err := os.ReadFile(file)
					if err != nil {
						mu.Lock()
						results = append(results, evalResult{filePath: file, err: fmt.Errorf("파일 읽기 실패: %w", err)})
						mu.Unlock()
						continue
					}
					absPath, err := filepath.Abs(file)
					if err != nil {
						absPath = file
					}
					collectedPrompt := &model.CollectedPrompt{
						ID:         uuid.New().String(),
						Tool:       "scan",
						RawPrompt:  string(content),
						Timestamp:  time.Now().Unix(),
						Command:    fmt.Sprintf("scan --repo %s --file %s", repoPath, file),
						WorkingDir: repoPath,
						Metadata: map[string]string{
							"file_path": absPath,
						},
					}
					result, err := eval.Evaluate(ctx, collectedPrompt)
					if err != nil {
						mu.Lock()
						results = append(results, evalResult{filePath: file, absPath: absPath, err: fmt.Errorf("평가 실패: %w", err)})
						mu.Unlock()
						continue
					}
					report, err := markdownReporter.Generate(result)
					if err != nil {
						mu.Lock()
						results = append(results, evalResult{filePath: file, absPath: absPath, err: fmt.Errorf("리포트 생성 실패: %w", err)})
						mu.Unlock()
						continue
					}
					if repo != nil {
						// enqueue for serialized save
						saveJobs <- collectedPrompt
					}
					mu.Lock()
					results = append(results, evalResult{filePath: file, absPath: absPath, report: report, score: result.Score.OverallScore})
					mu.Unlock()
				}
			}

			for i := 0; i < concurrency; i++ {
				wg.Add(1)
				go worker()
			}
			for _, f := range files {
				jobs <- f
			}
			close(jobs)
			wg.Wait()

			var mergedReports []string
			successCount := 0
			linesPrinted := 0
			for _, r := range results {
				if r.err != nil {
					cmd.Printf("  ⚠️  %s: %v\n", r.filePath, r.err)
					continue
				}

				if useFileOutput && singleOut == "" {
					reportFileName := fmt.Sprintf("%s_report.md", sanitizeFileName(filepath.Base(r.filePath)))
					reportPath := filepath.Join(outputDir, reportFileName)
					if err := os.WriteFile(reportPath, []byte(r.report), 0644); err != nil {
						cmd.Printf("  ⚠️  리포트 저장 실패: %v\n", err)
						continue
					}

					// finalize DB writes (once after processing all results)
					if repo != nil {
						close(saveJobs)
						saveWg.Wait()
					}
					cmd.Printf("  ✅ 완료 - 점수: %.1f/100, 리포트: %s\n", r.score, reportPath)
				} else if useFileOutput && singleOut != "" {
					header := fmt.Sprintf("# %s\n\n", r.absPath)
					mergedReports = append(mergedReports, header+r.report+"\n\n")
					cmd.Printf("  ✅ 완료 - 점수: %.1f/100\n", r.score)
				} else {
					// 콘솔 출력: 기본은 요약(집계 및 하위 점수 Top-N), --full 시 전체 리포트
					if fullOutput {
						// 전체 출력 강제
						sep := strings.Repeat("-", 80)
						cmd.Printf("%s\n파일: %s\n%s\n%s\n\n", sep, r.absPath, sep, r.report)
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
							cmd.Printf("%2d) %.1f  %s\n", i+1, good[i].score, good[i].absPath)
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
									cmd.Printf("%2d) %.1f  %s\n", i+1, good[idx].score, good[idx].absPath)
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

			cmd.Printf("\n총 %d개 파일 분석 완료\n", successCount)
			if useFileOutput {
				cmd.Printf("리포트 저장 위치: %s\n", outputDir)
			}

			return nil
		},
	}

	cmd.Flags().StringP("repo", "r", ".", "스캔할 레포지토리 경로")
	cmd.Flags().StringP("output", "o", "", "리포트 출력 디렉토리 (미지정 시 콘솔 출력)")
	cmd.Flags().String("single-output", "", "여러 리포트를 하나의 파일로 병합해 저장 (— --output 필요)")
	cmd.Flags().Bool("clean", false, "--output 디렉토리의 기존 .md 리포트를 정리하고 새로 생성")
	cmd.Flags().StringSlice("patterns", []string{"*.md", "*.txt"}, "프롬프트 파일 패턴")
	cmd.Flags().StringP("provider", "p", "claude", "LLM Provider (토큰 계산용)")
	cmd.Flags().Int("concurrency", 0, "동시 평가 작업 수 (기본값: CPU 코어 수)")
	cmd.Flags().Bool("full", false, "콘솔에 전체 리포트를 출력 (기본은 요약)")
	cmd.Flags().Int("max-lines", 100, "콘솔 요약 최대 출력 라인 수 (기본: 100)")
	cmd.Flags().Int("top", 10, "요약 출력 시 하위 점수 Top-N 표시 (기본: 10)")
	cmd.Flags().String("summary", "rich", "요약 모드 (rich|short)")
	cmd.Flags().Bool("advice", true, "요약에 코칭형 개선 가이드 출력")

	return cmd
}

// findPromptFiles finds prompt files matching patterns in the directory
func findPromptFiles(rootPath string, patterns []string) ([]string, error) {
	var files []string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// 권한 문제 등 접근 불가 경로는 건너뜀
			if os.IsPermission(err) {
				return nil
			}
			// 기타 오류도 스캔을 멈추지 않고 계속 진행
			return nil
		}

		if info.IsDir() {
			// 특정 디렉토리는 제외
			name := info.Name()
			if strings.HasPrefix(name, ".") && name != "." {
				return filepath.SkipDir
			}
			switch strings.ToLower(name) {
			case ".git", ".hg", ".svn", ".trash", "node_modules", "vendor", "dist", "build", ".next", ".cache", "library":
				return filepath.SkipDir
			}
			return nil
		}

		// 패턴 매칭
		for _, pattern := range patterns {
			matched, _ := filepath.Match(pattern, info.Name()) //nolint:errcheck
			if matched {
				files = append(files, path)
				break
			}
		}

		return nil
	})

	return files, err
}

// sanitizeFileName sanitizes file name for report output
func sanitizeFileName(name string) string {
	name = strings.TrimSuffix(name, filepath.Ext(name))
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "\\", "_")
	return name
}
