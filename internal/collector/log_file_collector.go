package collector

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"github.com/curogom/curompt/internal/model"
	"github.com/curogom/curompt/internal/parser"
	"github.com/curogom/curompt/internal/repository"
)

// LogFileCollector collects prompts from log files
type LogFileCollector struct {
	repository    repository.PromptRepository
	parser        parser.Parser
	tool          string // claude, codex, cursor
	logPath       string // Path to log file
	projectFilter string // Optional: filter by project path (empty = all projects)
}

// NewLogFileCollector creates a new log file collector
func NewLogFileCollector(repo repository.PromptRepository, tool string, logPath string) *LogFileCollector {
	return &LogFileCollector{
		repository:    repo,
		parser:        parser.NewParser(),
		tool:          tool,
		logPath:       logPath,
		projectFilter: "",
	}
}

// SetProjectFilter sets the project path filter (only collect prompts from this project)
func (c *LogFileCollector) SetProjectFilter(projectPath string) {
	c.projectFilter = projectPath
}

// Collect collects prompts from the log file
func (c *LogFileCollector) Collect(ctx context.Context) ([]*model.CollectedPrompt, error) {
	if c.logPath == "" {
		return nil, fmt.Errorf("log path is required")
	}

	file, err := os.Open(c.logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	var collectedPrompts []*model.CollectedPrompt

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Tool-specific parsing
		var promptText string
		var timestamp int64
		var metadata map[string]string

		switch c.tool {
		case "claude", "claude-code":
			p, ts, m, err := c.parseClaudeHistoryLine(line)
			if err != nil {
				// Skip invalid lines but continue
				continue
			}
			// 프로젝트 필터 적용
			if c.projectFilter != "" {
				projectPath := m["project"]
				if projectPath != c.projectFilter {
					// 이 프로젝트가 아니면 스킵
					continue
				}
			}
			promptText = p
			timestamp = ts
			metadata = m
		case "codex":
			p, ts, m, err := c.parseCodexHistoryLine(line)
			if err != nil {
				continue
			}
			// 프로젝트 필터 적용
			if c.projectFilter != "" {
				projectPath := m["project"]
				if projectPath != c.projectFilter {
					// 이 프로젝트가 아니면 스킵
					continue
				}
			}
			promptText = p
			timestamp = ts
			metadata = m
		case "cursor":
			p, ts, m, err := c.parseCursorHistoryLine(line)
			if err != nil {
				continue
			}
			promptText = p
			timestamp = ts
			metadata = m
		default:
			return nil, fmt.Errorf("unsupported tool for log collection: %s", c.tool)
		}

		if promptText == "" {
			continue
		}

		// Parse prompt
		parsedPrompt, err := c.parser.Parse(promptText)
		if err != nil {
			// Skip invalid prompts
			continue
		}

		// Create collected prompt
		collected := &model.CollectedPrompt{
			ID:         uuid.New().String(),
			Tool:       c.tool,
			Prompt:     parsedPrompt,
			RawPrompt:  promptText,
			Timestamp:  timestamp,
			Command:    fmt.Sprintf("collect --from %s --file %s", c.tool, c.logPath),
			WorkingDir: filepath.Dir(c.logPath),
			Metadata:   metadata,
		}

		collectedPrompts = append(collectedPrompts, collected)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read log file: %w", err)
	}

	return collectedPrompts, nil
}

// parseClaudeHistoryLine parses a line from Claude Code history.jsonl
func (c *LogFileCollector) parseClaudeHistoryLine(line string) (prompt string, timestamp int64, metadata map[string]string, err error) {
	var data struct {
		Display        string                 `json:"display"`
		PastedContents map[string]interface{} `json:"pastedContents"`
		Timestamp      int64                  `json:"timestamp"`
		Project        string                 `json:"project"`
	}

	if err := json.Unmarshal([]byte(line), &data); err != nil {
		return "", 0, nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Extract prompt from display field
	prompt = data.Display

	// Convert timestamp from milliseconds to seconds
	timestamp = data.Timestamp / 1000

	// Build metadata
	metadata = make(map[string]string)
	if data.Project != "" {
		metadata["project"] = data.Project
	}
	if len(data.PastedContents) > 0 {
		// Convert pastedContents to string representation
		if pastedBytes, err := json.Marshal(data.PastedContents); err == nil {
			metadata["pastedContents"] = string(pastedBytes)
		}
	}
	metadata["source"] = "history.jsonl"

	return prompt, timestamp, metadata, nil
}

// parseCodexHistoryLine parses a line from Codex history.jsonl
func (c *LogFileCollector) parseCodexHistoryLine(line string) (prompt string, timestamp int64, metadata map[string]string, err error) {
	var data struct {
		SessionID string `json:"session_id"`
		Timestamp int64  `json:"ts"`
		Text      string `json:"text"`
	}

	if err := json.Unmarshal([]byte(line), &data); err != nil {
		return "", 0, nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	prompt = data.Text
	timestamp = data.Timestamp

	metadata = make(map[string]string)
	if data.SessionID != "" {
		metadata["session_id"] = data.SessionID
		// session_id를 사용하여 session 파일에서 프로젝트 경로 찾기
		if projectPath := c.findCodexProjectPath(data.SessionID); projectPath != "" {
			metadata["project"] = projectPath
		}
	}
	metadata["source"] = "history.jsonl"

	return prompt, timestamp, metadata, nil
}

// findCodexProjectPath finds the project path from Codex session file
func (c *LogFileCollector) findCodexProjectPath(sessionID string) string {
	home := os.Getenv("HOME")
	sessionsDir := filepath.Join(home, ".codex", "sessions")

	// session 파일 찾기 (YYYY/MM/DD/rollout-TIMESTAMP-SESSION_ID.jsonl 형식)
	pattern := filepath.Join(sessionsDir, "**", "*"+sessionID+"*.jsonl")
	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) == 0 {
		// 패턴 매칭 실패 시 직접 검색
		matches = c.findSessionFile(sessionsDir, sessionID)
	}

	if len(matches) == 0 {
		return ""
	}

	// session 파일의 첫 줄(session_meta)에서 cwd 추출
	return c.extractCwdFromSessionFile(matches[0])
}

// findSessionFile recursively searches for session file
func (c *LogFileCollector) findSessionFile(dir string, sessionID string) []string {
	var matches []string
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error { //nolint:errcheck
		if err != nil {
			return nil
		}
		if !info.IsDir() && strings.Contains(info.Name(), sessionID) && strings.HasSuffix(info.Name(), ".jsonl") {
			matches = append(matches, path)
		}
		return nil
	})
	return matches
}

// extractCwdFromSessionFile extracts cwd from session file's session_meta entry
func (c *LogFileCollector) extractCwdFromSessionFile(sessionFile string) string {
	file, err := os.Open(sessionFile)
	if err != nil {
		return ""
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var entry struct {
			Type    string `json:"type"`
			Payload struct {
				Cwd string `json:"cwd"`
			} `json:"payload"`
		}

		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			continue
		}

		if entry.Type == "session_meta" && entry.Payload.Cwd != "" {
			return entry.Payload.Cwd
		}
	}

	return ""
}

// parseCursorHistoryLine parses a line from Cursor history/log files
// Note: Cursor format may vary, this is a basic implementation
func (c *LogFileCollector) parseCursorHistoryLine(line string) (prompt string, timestamp int64, metadata map[string]string, err error) {
	// Try JSON format first
	var data struct {
		Message   string `json:"message"`
		Prompt    string `json:"prompt"`
		Timestamp int64  `json:"timestamp"`
		SessionID string `json:"session_id"`
	}

	if err := json.Unmarshal([]byte(line), &data); err == nil {
		// JSON format detected
		if data.Prompt != "" {
			prompt = data.Prompt
		} else if data.Message != "" {
			prompt = data.Message
		}
		timestamp = data.Timestamp
		metadata = make(map[string]string)
		if data.SessionID != "" {
			metadata["session_id"] = data.SessionID
		}
		metadata["source"] = "cursor_log"
		return prompt, timestamp, metadata, nil
	}

	// If not JSON, try plain text format with timestamp
	// This is a fallback - may need adjustment based on actual Cursor log format
	return "", 0, nil, fmt.Errorf("unsupported cursor log format")
}

// Name returns the collector name
func (c *LogFileCollector) Name() string {
	return fmt.Sprintf("log-file-%s", c.tool)
}
