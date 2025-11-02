package gemini

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

// GeminiProvider implements Provider for Google Gemini API
type GeminiProvider struct {
	apiKey     string
	model      string
	httpClient *http.Client
	tokenizer  parser.Tokenizer
}

// NewGeminiProvider creates a new Gemini provider
func NewGeminiProvider(apiKey, model string) provider.Provider {
	if model == "" {
		model = "gemini-pro" // 기본 모델
	}

	return &GeminiProvider{
		apiKey: apiKey,
		model:  model,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		tokenizer: parser.NewOpenAITokenizer(), // Gemini는 OpenAI와 유사한 토큰 계산
	}
}

// Name returns the provider name
func (p *GeminiProvider) Name() string {
	return "gemini"
}

// CalculateTokens calculates token count (Gemini uses similar tokenization to OpenAI)
func (p *GeminiProvider) CalculateTokens(text string) (int, error) {
	return p.tokenizer.CountTokens(text)
}

// Evaluate sends a prompt to Gemini API and returns the response
func (p *GeminiProvider) Evaluate(ctx context.Context, prompt string) (*provider.Response, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", p.model, p.apiKey)

	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{
						"text": prompt,
					},
				},
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
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
			FinishReason string `json:"finishReason"`
		} `json:"candidates"`
		UsageMetadata struct {
			PromptTokenCount     int `json:"promptTokenCount"`
			CandidatesTokenCount int `json:"candidatesTokenCount"`
		} `json:"usageMetadata"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(apiResponse.Candidates) == 0 {
		return nil, fmt.Errorf("no candidates in response")
	}

	candidate := apiResponse.Candidates[0]
	content := ""
	for _, part := range candidate.Content.Parts {
		content += part.Text
	}

	tokenCount := apiResponse.UsageMetadata.PromptTokenCount + apiResponse.UsageMetadata.CandidatesTokenCount

	return &provider.Response{
		Content:      content,
		TokenCount:   tokenCount,
		Model:        p.model,
		FinishReason: candidate.FinishReason,
	}, nil
}
