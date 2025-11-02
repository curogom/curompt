package redactor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedactor_MasksAPIKeys(t *testing.T) {
	// Red: 실패하는 테스트 작성
	redactor := NewRedactor()
	input := "API key: sk-1234567890abcdef"

	result := redactor.Redact(input)

	assert.NotContains(t, result, "sk-1234567890abcdef")
	assert.Contains(t, result, "[REDACTED]")
}

func TestRedactor_MasksBearerTokens(t *testing.T) {
	redactor := NewRedactor()
	input := "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"

	result := redactor.Redact(input)

	assert.NotContains(t, result, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9")
	assert.Contains(t, result, "[REDACTED]")
}

func TestRedactor_MasksAPIKeyWithColon(t *testing.T) {
	redactor := NewRedactor()
	input := "api_key: abc123def456"

	result := redactor.Redact(input)

	assert.NotContains(t, result, "abc123def456")
	assert.Contains(t, result, "[REDACTED]")
}

func TestRedactor_RedactsEnvReferences(t *testing.T) {
	redactor := NewRedactor()
	input := "Use $API_KEY from .env file"

	result := redactor.Redact(input)

	assert.Contains(t, result, "[REDACTED]")
	assert.NotContains(t, result, "$API_KEY")
}

func TestRedactor_PreservesNonSensitiveText(t *testing.T) {
	redactor := NewRedactor()
	input := "This is a normal prompt without secrets"

	result := redactor.Redact(input)

	assert.Equal(t, input, result)
}

func TestRedactor_MultipleSecrets(t *testing.T) {
	redactor := NewRedactor()
	input := "API key: sk-test123, Bearer token: abc123, .env file has secrets"

	result := redactor.Redact(input)

	assert.NotContains(t, result, "sk-test123")
	assert.NotContains(t, result, "abc123")
	assert.Contains(t, result, "[REDACTED]")
}
