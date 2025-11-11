package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/curogom/curompt/internal/collector"
	"github.com/curogom/curompt/internal/model"
	"github.com/curogom/curompt/internal/repository"
)

func collectFromLog(ctx context.Context, repo repository.PromptRepository, tool, logPath, projectFilter string, output outputWriter) ([]*model.CollectedPrompt, int, error) {
	if logPath == "" {
		return nil, 0, fmt.Errorf("로그 경로가 필요합니다")
	}
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		return nil, 0, fmt.Errorf("로그 파일을 찾을 수 없습니다: %s", logPath)
	}

	logCollector := collector.NewLogFileCollector(repo, tool, logPath)
	if projectFilter != "" {
		logCollector.SetProjectFilter(projectFilter)
	}

	if output != nil {
		output.Printf("로그 파일에서 프롬프트 수집 중: %s\n", logPath)
		output.Printf("도구: %s\n", tool)
		if projectFilter != "" {
			output.Printf("필터: %s\n", projectFilter)
		} else {
			output.Printf("필터: 모든 프로젝트\n")
		}
		output.Println()
	}

	prompts, err := logCollector.Collect(ctx)
	if err != nil {
		return nil, 0, err
	}

	saved := 0
	for i, prompt := range prompts {
		if err := repo.Save(ctx, prompt); err != nil {
			if output != nil {
				output.Printf("  [%d/%d] ⚠️  저장 실패: %v\n", i+1, len(prompts), err)
			}
			continue
		}
		saved++
	}

	if output != nil {
		output.Printf("수집된 프롬프트: %d개 (저장 %d개)\n\n", len(prompts), saved)
	}

	return prompts, saved, nil
}

type outputWriter interface {
	Println(...any)
	Printf(string, ...any)
}

type autoCollectOption struct {
	Label string
	Value string
}

func getDefaultLogPath(tool string) string {
	home := os.Getenv("HOME")
	switch tool {
	case "claude", "claude-code":
		return filepath.Join(home, ".claude", "history.jsonl")
	case "codex":
		return filepath.Join(home, ".codex", "history.jsonl")
	default:
		return ""
	}
}

// findProjectRoot searches upward for CLAUDE.md or Claude.md. Used mainly for Claude Code.
func findProjectRoot(startDir string) string {
	current := startDir
	for {
		claudeMd := filepath.Join(current, "CLAUDE.md")
		claudeMdLower := filepath.Join(current, "Claude.md")

		if _, err := os.Stat(claudeMd); err == nil {
			return current
		}
		if _, err := os.Stat(claudeMdLower); err == nil {
			return current
		}

		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return ""
}

// promptLine asks the user for input with an optional default.
func promptLine(question, defaultValue string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", question, defaultValue)
	} else {
		fmt.Printf("%s: ", question)
	}
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue, nil
	}
	return filepath.Clean(input), nil
}

func promptYesNo(question string, defaultYes bool) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	def := "n"
	if defaultYes {
		def = "y"
	}
	fmt.Printf("%s [y/n, 기본 %s]: ", question, def)
	answer, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	answer = strings.TrimSpace(strings.ToLower(answer))
	if answer == "" {
		return defaultYes, nil
	}
	return answer == "y" || answer == "yes", nil
}

func isInteractive() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func defaultProjectPath(source, targetPath string) string {
	if targetPath != "" {
		return targetPath
	}
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	switch source {
	case "claude", "claude-code":
		if root := findProjectRoot(wd); root != "" {
			return root
		}
		return ""
	case "codex":
		return wd
	default:
		return ""
	}
}

func promptSourceSelection(options []autoCollectOption) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("수집할 도구를 선택하세요:")
	for idx, opt := range options {
		fmt.Printf("  %d) %s\n", idx+1, opt.Label)
	}
	fmt.Print("번호 선택 (Enter 취소): ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return "", nil
	}
	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > len(options) {
		return "", fmt.Errorf("잘못된 선택입니다")
	}
	return options[choice-1].Value, nil
}
