package parser

import (
	"strings"
)

// claudeTokenizer implements Tokenizer for Claude API
// Note: 실제 구현은 claude-tokenizer 라이브러리를 사용하거나 근사치 계산
type claudeTokenizer struct{}

// NewClaudeTokenizer creates a new Claude tokenizer
func NewClaudeTokenizer() Tokenizer {
	return &claudeTokenizer{}
}

// CountTokens estimates token count for Claude
// Claude는 일반적으로 단어 기반 + 특수 문자를 고려
// 근사치: 평균 4자 = 1 토큰 (영어 기준)
func (t *claudeTokenizer) CountTokens(text string) (int, error) {
	if len(text) == 0 {
		return 0, nil
	}

	// 간단한 근사치: 공백으로 분리된 단어 수 + 특수 문자 고려
	words := strings.Fields(text)
	baseTokens := len(words)

	// 한글 및 멀티바이트 문자 고려 (대략 1.5배)
	multibyteChars := 0
	for _, r := range text {
		if r > 127 {
			multibyteChars++
		}
	}

	// 한글 문자는 평균적으로 더 많은 토큰 사용
	estimatedTokens := baseTokens + (multibyteChars / 2)

	// 최소 1 토큰 보장
	if estimatedTokens < 1 {
		return 1, nil
	}

	return estimatedTokens, nil
}
