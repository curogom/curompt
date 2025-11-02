package redactor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedactor_MasksAPIKeys(t *testing.T) {
	// Red: 실패하는 테스트 작성
	redactor := NewRedactor()
	// NOTE: This is a dummy API key for testing only, not a real secret
	input := "API key: sk-test-dummy-key-1234567890abcdef"

	result := redactor.Redact(input)

	assert.NotContains(t, result, "sk-test-dummy-key-1234567890abcdef")
	assert.Contains(t, result, "[REDACTED]")
}

func TestRedactor_MasksBearerTokens(t *testing.T) {
	redactor := NewRedactor()
	// NOTE: This is a dummy JWT token for testing only, not a real secret
	// Base64 encoded "dummy-header.dummy-payload.dummy-signature"
	input := "Authorization: Bearer dummy-header.dummy-payload.dummy-signature"

	result := redactor.Redact(input)

	assert.NotContains(t, result, "dummy-header.dummy-payload.dummy-signature")
	assert.Contains(t, result, "[REDACTED]")
}

func TestRedactor_MasksAPIKeyWithColon(t *testing.T) {
	redactor := NewRedactor()
	// NOTE: This is a dummy API key for testing only, not a real secret
	input := "api_key: test-dummy-key-abc123def456"

	result := redactor.Redact(input)

	assert.NotContains(t, result, "test-dummy-key-abc123def456")
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
	// NOTE: All tokens/keys here are dummy values for testing only, not real secrets
	input := "API key: sk-test-dummy-123, Bearer token: test-dummy-token-abc123, .env file has secrets"

	result := redactor.Redact(input)

	assert.NotContains(t, result, "sk-test-dummy-123")
	assert.NotContains(t, result, "test-dummy-token-abc123")
	assert.Contains(t, result, "[REDACTED]")
}
