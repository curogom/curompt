# Provider Support

> 🇰🇷 **Korean**: [한국어 버전](./PROVIDERS.ko.md)

## Supported Provider List

### ✅ Implemented

1. **Claude (Anthropic)**
   - API: Anthropic API
   - Token calculation: Claude tokenizer
   - Status: ✅ M1 MVP Complete

2. **OpenAI**
   - API: OpenAI API
   - Token calculation: tiktoken
   - Status: ✅ M1 MVP Complete

3. **Gemini (Google)**
   - API: Google Gemini API
   - Token calculation: Gemini tokenizer
   - Status: ✅ M1 MVP Complete

4. **Cursor IDE/CLI**
   - Integration: Cursor CLI + MCP protocol
   - Features: Integration support with Cursor IDE
   - Status: ✅ M1 MVP Complete

### ⚠️ Lower Priority (After M2)

5. **Codex CLI**
   - Integration: AGENTS.md/config.toml recognition
   - Status: Review in M2

6. **Ollama (Local)**
   - Integration: Local Ollama server
   - Status: Review in M2

7. **AWS Bedrock / Google Vertex**
   - Integration: Each platform API
   - Status: Review after M2

## Adding a Provider

The Strategy pattern allows easy addition of new providers:

```go
type Provider interface {
    Evaluate(ctx context.Context, prompt string) (*Response, error)
    CalculateTokens(text string) (int, error)
    Name() string
}
```

### Gemini API Integration

- API endpoint: `https://generativelanguage.googleapis.com/v1beta/models/{model}:generateContent`
- Authentication: API key based
- Supported models: `gemini-pro`, `gemini-pro-vision`, etc.

### Cursor CLI Integration

- Cursor CLI accessible via `cursor` command
- MCP (Model Context Protocol) support
- Session capture and analysis possible

## Configuration Examples

```yaml
provider:
  name: gemini
  api_key: ${GEMINI_API_KEY}  # Use environment variable
  model: gemini-pro
  timeout_sec: 30

# Or use Cursor
provider:
  name: cursor
  cli_path: /usr/local/bin/cursor
  mcp_enabled: true
```

---

**Note**: API keys for each provider should be stored in environment variables or secure configuration files.
