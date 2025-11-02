package parser

// Tokenizer calculates tokens for different LLM providers
type Tokenizer interface {
	CountTokens(text string) (int, error)
}
