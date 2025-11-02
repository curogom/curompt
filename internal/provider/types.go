package provider

import "context"

// Response contains the LLM provider response
type Response struct {
	Content      string
	TokenCount   int
	Model        string
	FinishReason string
}

// Provider interface for LLM providers
type Provider interface {
	// Evaluate evaluates a prompt and returns a response
	Evaluate(ctx context.Context, prompt string) (*Response, error)

	// CalculateTokens calculates token count for text
	CalculateTokens(text string) (int, error)

	// Name returns the provider name
	Name() string
}
