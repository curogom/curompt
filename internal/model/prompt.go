package model

import "github.com/curogom/curompt/internal/parser"

// CollectedPrompt represents a collected prompt with metadata
type CollectedPrompt struct {
	ID         string            // 고유 ID
	Tool       string            // 도구 이름 (codex, cursor, claude-code 등)
	Prompt     *parser.Prompt    // 파싱된 프롬프트
	RawPrompt  string            // 원본 프롬프트 텍스트
	Timestamp  int64             // 수집 시간 (Unix timestamp)
	Command    string            // 실행된 명령어
	WorkingDir string            // 작업 디렉토리
	Metadata   map[string]string // 추가 메타데이터
}
