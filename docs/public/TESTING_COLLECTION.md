# 프롬프트 수집 테스트 가이드

> 🇰🇷 **Korean**: [한국어 버전](./TESTING_COLLECTION.ko.md)

## Overview

This guide helps you test whether `curompt` successfully collects prompts when you use external tools like Claude Code, Codex CLI, or Cursor CLI.

> **Scan note:** `curompt scan` analyzes prompts already saved in the local DB.  
> If no history exists for a path, the CLI can auto-collect from Claude Code or Codex logs (Cursor support lands in v1.1).

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
curompt scan --path .

# 3. Verify collection
curompt list
```

**Expected Result:**
- Prompt saved in database
- Report generated in `reports/`
- `list` command shows the prompt

### Method 2: CLI Wrapper (For CLI Tools) ⚠️

For command-line tools like Codex CLI or Cursor CLI:

```bash
# Wrap a CLI command
curompt wrap codex exec "TASK: Add authentication feature"

# Verify collection
curompt list --tool codex
```

**Limitations:**
- Only works for CLI tools, not IDEs
- Requires tool-specific prompt extraction logic
- Currently basic implementation

### Method 3: Claude Code CLI Wrapper ✅

For Claude Code CLI:

```bash
# Wrap Claude Code CLI command
curompt wrap claude-code exec "TASK: Add authentication feature"

# Or with other Claude Code commands
curompt wrap claude-code chat "Implement login"

# Verify collection
curompt list --tool claude-code
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
rm -rf ~/.curompt/db.sqlite

# 2. Create test prompt
echo "# ROLE\nEngineer" > test.md

# 3. Scan and collect
curompt scan --path .

# 4. Check if saved
curompt list

# Expected: Should show 1 prompt
```

### Test 2: Multiple Prompts Collection

```bash
# 1. Create multiple prompts
mkdir -p test-prompts
echo "# ROLE\nEngineer1" > test-prompts/p1.md
echo "# ROLE\nEngineer2" > test-prompts/p2.md

# 2. Scan all
curompt scan --path test-prompts

# 3. Verify count
curompt list --limit 10

# Expected: Should show 2 prompts
```

### Test 3: Collection After Real Usage

```bash
# 1. Use Claude Code IDE normally
# (create prompts, make requests, etc.)

# 2. After using Claude Code, try to collect:
# Option A: If log parsing implemented
curompt collect --from claude-code

# Option B: If you saved prompts as files
curompt scan --path ./my-prompts

# 3. Check collection
curompt list
```

## Verification Checklist

After running collection:

- [ ] Database file exists: `~/.curompt/db.sqlite`
- [ ] `list` command shows collected prompts
- [ ] Prompts have correct tool identifier
- [ ] Timestamps are correct
- [ ] Reports are generated (for scan command)

## Troubleshooting

### Issue: No prompts collected

**Check:**
1. Database exists: `ls -la ~/.curompt/db.sqlite`
2. Check tool identifier: `curompt list --tool scan`
3. Verify file paths in scan command

**Solution:**
```bash
# Re-scan with verbose output
curompt scan --path . --patterns "*.md"
curompt list
```

### Issue: Wrapper doesn't capture prompts

**For CLI Tools:**
- Verify command format
- Check if tool-specific parser needed
- Try with quoted prompts: `curompt wrap codex exec "PROMPT: ..."`

### Issue: Claude Code prompts not collected

**For CLI Tools:**
- Use `wrap` command to wrap Claude Code CLI
- Check command format: `curompt wrap claude-code [your-command]`
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
