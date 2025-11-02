package parser

import (
	"strings"
)

// openAITokenizer implements Tokenizer for OpenAI API
// Note: 실제 구현은 tiktoken 라이브러리를 사용하거나 근사치 계산
type openAITokenizer struct{}

// NewOpenAITokenizer creates a new OpenAI tokenizer
func NewOpenAITokenizer() Tokenizer {
	return &openAITokenizer{}
}

// CountTokens estimates token count for OpenAI
// OpenAI는 일반적으로 4자 = 1 토큰 (영어 기준)
func (t *openAITokenizer) CountTokens(text string) (int, error) {
	if len(text) == 0 {
		return 0, nil
	}

	// 간단한 근사치: 문자 수 / 4
	// 실제로는 BPE 인코딩을 사용하므로 더 복잡함
	bytes := len([]byte(text))
	estimatedTokens := bytes / 4

	// 한글 및 멀티바이트 문자는 더 많은 토큰 사용
	multibyteChars := 0
	for _, r := range text {
		if r > 127 {
			multibyteChars++
		}
	}

	// 멀티바이트 문자는 평균 2-3 토큰
	estimatedTokens += multibyteChars

	// 단어 경계도 고려
	words := len(strings.Fields(text))
	if words > 0 {
		// 단어 수의 10%를 추가 (단어 경계 토큰)
		estimatedTokens += words / 10
	}

	// 최소 1 토큰 보장
	if estimatedTokens < 1 {
		return 1, nil
	}

	return estimatedTokens, nil
}
