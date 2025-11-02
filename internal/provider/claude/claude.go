package claude

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/curogom/curo-prompt/internal/parser"
	"github.com/curogom/curo-prompt/internal/provider"
)

// ClaudeProvider implements Provider for Anthropic Claude API
type ClaudeProvider struct {
	apiKey     string
	model      string
	httpClient *http.Client
	tokenizer  parser.Tokenizer
}

// NewClaudeProvider creates a new Claude provider
func NewClaudeProvider(apiKey, model string) provider.Provider {
	if model == "" {
		model = "claude-3-5-sonnet-20241022" // 기본 모델
	}

	return &ClaudeProvider{
		apiKey: apiKey,
		model:  model,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		tokenizer: parser.NewClaudeTokenizer(),
	}
}

// Name returns the provider name
func (p *ClaudeProvider) Name() string {
	return "claude"
}

// CalculateTokens calculates token count using Claude tokenizer
func (p *ClaudeProvider) CalculateTokens(text string) (int, error) {
	return p.tokenizer.CountTokens(text)
}

// Evaluate sends a prompt to Claude API and returns the response
func (p *ClaudeProvider) Evaluate(ctx context.Context, prompt string) (*provider.Response, error) {
	url := "https://api.anthropic.com/v1/messages"

	requestBody := map[string]interface{}{
		"model":      p.model,
		"max_tokens": 4096,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("API error: %d - failed to read body: %w", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var apiResponse struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Model      string `json:"model"`
		StopReason string `json:"stop_reason"`
		Usage      struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Content 추출
	content := ""
	for _, block := range apiResponse.Content {
		if block.Type == "text" {
			content += block.Text
		}
	}

	return &provider.Response{
		Content:      content,
		TokenCount:   apiResponse.Usage.InputTokens + apiResponse.Usage.OutputTokens,
		Model:        apiResponse.Model,
		FinishReason: apiResponse.StopReason,
	}, nil
}
