# 프롬프트 수집 테스트 가이드

> 🇰🇷 **Korean**: [한국어 버전](./TESTING_COLLECTION.ko.md)

## Overview

This guide helps you test whether `curo-prompt` successfully collects prompts when you use external tools like Claude Code, Codex CLI, or Cursor CLI.

## Important Notes

### Claude Code CLI

- **Claude Code**: CLI tool - can be wrapped via `wrap` command
- **Codex CLI / Cursor CLI**: Command-line tools - can be wrapped via `wrap` command

## Testing Methods

### Method 1: File Scanning (Currently Working) ✅

The simplest and most reliable method:

```bash
# 1. Create a test prompt file
cat > test-prompt.md << 'EOF'
# ROLE
Senior Software Engineer

# INPUTS
- task: string
- context: string

# INVARIANTS
- Code must be testable
- Error handling required

# OUTPUT FORMAT
JSON format
EOF

# 2. Scan and collect
curo-prompt scan --repo .

# 3. Verify collection
curo-prompt list
```

**Expected Result:**
- Prompt saved in database
- Report generated in `reports/`
- `list` command shows the prompt

### Method 2: CLI Wrapper (For CLI Tools) ⚠️

For command-line tools like Codex CLI or Cursor CLI:

```bash
# Wrap a CLI command
curo-prompt wrap codex exec "TASK: Add authentication feature"

# Verify collection
curo-prompt list --tool codex
```

**Limitations:**
- Only works for CLI tools, not IDEs
- Requires tool-specific prompt extraction logic
- Currently basic implementation

### Method 3: Claude Code CLI Wrapper ✅

For Claude Code CLI:

```bash
# Wrap Claude Code CLI command
curo-prompt wrap claude-code exec "TASK: Add authentication feature"

# Or with other Claude Code commands
curo-prompt wrap claude-code chat "Implement login"

# Verify collection
curo-prompt list --tool claude-code
```

**How it works:**
1. Executes the wrapped command
2. Extracts prompt from command arguments
3. Saves to database
4. Original command output is preserved

## Step-by-Step Test Procedure

### Test 1: Basic Collection Test

```bash
# 1. Clear existing data (optional)
rm -rf ~/.curo-prompt/db.sqlite

# 2. Create test prompt
echo "# ROLE\nEngineer" > test.md

# 3. Scan and collect
curo-prompt scan --repo .

# 4. Check if saved
curo-prompt list

# Expected: Should show 1 prompt
```

### Test 2: Multiple Prompts Collection

```bash
# 1. Create multiple prompts
mkdir -p test-prompts
echo "# ROLE\nEngineer1" > test-prompts/p1.md
echo "# ROLE\nEngineer2" > test-prompts/p2.md

# 2. Scan all
curo-prompt scan --repo test-prompts

# 3. Verify count
curo-prompt list --limit 10

# Expected: Should show 2 prompts
```

### Test 3: Collection After Real Usage

```bash
# 1. Use Claude Code IDE normally
# (create prompts, make requests, etc.)

# 2. After using Claude Code, try to collect:
# Option A: If log parsing implemented
curo-prompt collect --from claude-code

# Option B: If you saved prompts as files
curo-prompt scan --repo ./my-prompts

# 3. Check collection
curo-prompt list
```

## Verification Checklist

After running collection:

- [ ] Database file exists: `~/.curo-prompt/db.sqlite`
- [ ] `list` command shows collected prompts
- [ ] Prompts have correct tool identifier
- [ ] Timestamps are correct
- [ ] Reports are generated (for scan command)

## Troubleshooting

### Issue: No prompts collected

**Check:**
1. Database exists: `ls -la ~/.curo-prompt/db.sqlite`
2. Check tool identifier: `curo-prompt list --tool scan`
3. Verify file paths in scan command

**Solution:**
```bash
# Re-scan with verbose output
curo-prompt scan --repo . --patterns "*.md"
curo-prompt list
```

### Issue: Wrapper doesn't capture prompts

**For CLI Tools:**
- Verify command format
- Check if tool-specific parser needed
- Try with quoted prompts: `curo-prompt wrap codex exec "PROMPT: ..."`

### Issue: Claude Code prompts not collected

**For CLI Tools:**
- Use `wrap` command to wrap Claude Code CLI
- Check command format: `curo-prompt wrap claude-code [your-command]`
- Verify prompt extraction logic matches your command format

**If wrap doesn't work:**
1. Check if Claude Code CLI is in PATH: `which claude-code`
2. Verify command format matches expected patterns
3. Try saving prompts manually to files and use `scan` command

## Next Steps for Implementation

To improve Claude Code CLI support:

1. **Enhanced Prompt Extraction**
   - Analyze Claude Code CLI command formats
   - Improve `extractPrompt` function for better pattern matching
   - Support different Claude Code subcommands

2. **Better Error Handling**
   - Validate CLI tool installation
   - Provide clearer error messages
   - Handle edge cases in prompt extraction

3. **Alternative: Manual File Export**
   - If wrap doesn't capture properly
   - Export prompts manually to files
   - Use `scan` command

---

**For M1 MVP**: File scanning (`scan` command) is fully functional and sufficient for testing prompt collection and evaluation.

