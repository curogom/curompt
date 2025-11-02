package redactor

import (
	"regexp"
)

// patternRedactor implements Redactor using regex patterns
type patternRedactor struct {
	patterns []*regexp.Regexp
}

// NewRedactor creates a new redactor with default patterns
func NewRedactor() Redactor {
	patterns := []string{
		// API keys (sk-xxx)
		`(?i)sk-[a-z0-9]+`,
		// API key with colon/equals (api_key: xxx or api_key=xxx)
		`(?i)api[_-]?key\s*[:=]\s*[a-z0-9]+`,
		// Bearer tokens (Bearer token: xxx or Bearer xxx)
		`(?i)bearer\s+(?:token\s*[:=]\s*)?[a-z0-9\-_\.]+`,
		// .env references
		`(?i)\.env`,
		// Environment variables ($VAR_NAME)
		`(?i)\$[A-Z_][A-Z0-9_]*`,
	}

	compiledPatterns := make([]*regexp.Regexp, 0, len(patterns))
	for _, pattern := range patterns {
		compiled, err := regexp.Compile(pattern)
		if err == nil {
			compiledPatterns = append(compiledPatterns, compiled)
		}
	}

	return &patternRedactor{
		patterns: compiledPatterns,
	}
}

// Redact masks sensitive information in text
func (r *patternRedactor) Redact(text string) string {
	result := text

	// 모든 패턴을 순차적으로 적용
	for _, pattern := range r.patterns {
		// Find all matches in current result
		matches := pattern.FindAllStringIndex(result, -1)
		if len(matches) == 0 {
			continue
		}

		// Replace from end to start to preserve indices
		for i := len(matches) - 1; i >= 0; i-- {
			match := matches[i]

			// Replace with [REDACTED]
			if match[0] < len(result) && match[1] <= len(result) {
				result = result[:match[0]] + "[REDACTED]" + result[match[1]:]
			}
		}
	}

	return result
}
