package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"time"

	"github.com/curogom/curo-prompt/internal/evaluator"
	"github.com/curogom/curo-prompt/internal/model"
	"github.com/curogom/curo-prompt/internal/reporter"
	"github.com/curogom/curo-prompt/internal/repository"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
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
			repoPath, _ := cmd.Flags().GetString("repo")          //nolint:errcheck
			outputDir, _ := cmd.Flags().GetString("output")       //nolint:errcheck
			patterns, _ := cmd.Flags().GetStringSlice("patterns") //nolint:errcheck
			provider, _ := cmd.Flags().GetString("provider")      //nolint:errcheck

			if repoPath == "" {
				repoPath = "."
			}

			if outputDir == "" {
				outputDir = "reports"
			}

			// 출력 디렉토리 생성
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %w", err)
			}

			// 프롬프트 파일 찾기
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

			// 각 파일 분석
			eval := evaluator.NewEvaluator(provider)
			ctx := context.Background()
			markdownReporter := reporter.NewMarkdownReporter()

			successCount := 0
			for i, file := range files {
				cmd.Printf("[%d/%d] 분석 중: %s\n", i+1, len(files), file)

				// 파일 읽기
				content, err := os.ReadFile(file)
				if err != nil {
					cmd.Printf("  ⚠️  파일 읽기 실패: %v\n", err)
					continue
				}

				// 수집된 프롬프트 생성
				absPath, err := filepath.Abs(file)
				if err != nil {
					absPath = file // Fallback to relative path
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

				// 평가 수행
				result, err := eval.Evaluate(ctx, collectedPrompt)
				if err != nil {
					cmd.Printf("  ⚠️  평가 실패: %v\n", err)
					continue
				}

				// 리포트 생성
				report, err := markdownReporter.Generate(result)
				if err != nil {
					cmd.Printf("  ⚠️  리포트 생성 실패: %v\n", err)
					continue
				}

				// 저장소에 저장 (있으면)
				if repo != nil {
					if err := repo.Save(ctx, collectedPrompt); err != nil {
						cmd.Printf("  ⚠️  저장 실패: %v\n", err)
					}
				}

				// 리포트 파일 저장
				reportFileName := fmt.Sprintf("%s_report.md", sanitizeFileName(filepath.Base(file)))
				reportPath := filepath.Join(outputDir, reportFileName)
				if err := os.WriteFile(reportPath, []byte(report), 0644); err != nil {
					cmd.Printf("  ⚠️  리포트 저장 실패: %v\n", err)
					continue
				}

				cmd.Printf("  ✅ 완료 - 점수: %.1f/100, 리포트: %s\n", result.Score.OverallScore, reportPath)
				successCount++
			}

			cmd.Printf("\n총 %d개 파일 분석 완료\n", successCount)
			if outputDir != "" {
				cmd.Printf("리포트 저장 위치: %s\n", outputDir)
			}

			return nil
		},
	}

	cmd.Flags().StringP("repo", "r", ".", "스캔할 레포지토리 경로")
	cmd.Flags().StringP("output", "o", "reports", "리포트 출력 디렉토리")
	cmd.Flags().StringSlice("patterns", []string{"*.md", "*.txt"}, "프롬프트 파일 패턴")
	cmd.Flags().StringP("provider", "p", "claude", "LLM Provider (토큰 계산용)")

	return cmd
}

// findPromptFiles finds prompt files matching patterns in the directory
func findPromptFiles(rootPath string, patterns []string) ([]string, error) {
	var files []string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// 특정 디렉토리는 제외
			if strings.HasPrefix(info.Name(), ".") && info.Name() != "." {
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
