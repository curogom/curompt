# Workflow & Prompt Collection

> 🇰🇷 **Korean**: [한국어 버전](./WORKFLOW.ko.md)

## Overview

This document explains how `curo-prompt` collects, stores, and manages prompts.

## Current Workflow

### 1. Initial Setup

When you first run `curo-prompt`, the database is automatically initialized:

```bash
# First run - DB will be created at ~/.curo-prompt/db.sqlite
curo-prompt scan --repo .
```

The database is created at: `~/.curo-prompt/db.sqlite`

### 2. Prompt Collection Methods

#### Method 1: File Scanning (Current - MVP)

Scan existing prompt files in your repository:

```bash
# Scan current directory for prompt files
curo-prompt scan --repo .

# Scan specific directory
curo-prompt scan --repo ./prompts

# Custom file patterns
curo-prompt scan --repo . --patterns "*.md" --patterns "*.txt"
```

**What happens (all in one command):**
1. Finds all prompt files matching patterns
2. **Collects** each prompt (creates CollectedPrompt)
3. **Evaluates** each prompt (calls Evaluate)
4. **Saves** to database (`~/.curo-prompt/db.sqlite`)
5. **Generates** reports in `reports/` directory

**Note**: The `scan` command does everything in one go: collection → evaluation → storage → report generation.

#### Method 2: Direct Evaluation (Doesn't Save to DB)

Evaluate a single prompt file without saving:

```bash
# Evaluate without saving
curo-prompt eval --file prompt.md

# Evaluate and save report
curo-prompt eval --file prompt.md --output reports/
```

**Note**: `eval` command currently **does not save** to database. Only `scan` command saves.

#### Method 3: CLI Wrapper Collection ⚠️

Wrap external CLI tools to capture prompts:

```bash
# Wrap one-shot CLI commands
curo-prompt wrap claude --print "Your prompt"
curo-prompt wrap codex exec "TASK: Add feature"
curo-prompt wrap cursor chat "Implement login"
```

**How it works:**
1. Executes the wrapped CLI command
2. Extracts prompt from command arguments
3. Saves to database automatically
4. Preserves original command output

**Limitations:**
- Only works for **one-shot commands** (non-interactive)
- Does NOT capture prompts from **interactive sessions**
  - Claude Code interactive session (`claude` without `--print`)
  - Multi-turn conversations
  - Session history

**For Claude Code Interactive Sessions**: See [Claude Code Collection Guide](./CLAUDE_CODE_COLLECTION.md)

**Status**: Basic implementation ready for one-shot commands. Interactive session capture requires additional implementation (M2).

### 3. Viewing Stored Prompts

Use the `list` command to view stored prompts:

```bash
# List recent 10 prompts
curo-prompt list

# List more prompts
curo-prompt list --limit 20

# List prompts from specific tool
curo-prompt list --tool scan
curo-prompt list --tool codex

# List and re-evaluate
curo-prompt list --eval
```

**Output example:**
```
저장된 프롬프트: 5개

[1] ID: a1b2c3d4
    도구: scan
    시간: 2025-01-15 14:30:22
    명령: scan --repo ./prompts
    경로: /Users/user/project
    프롬프트: # ROLE\nSenior Engineer\n\n# INPUTS\n- task: string...
```

### 4. Re-evaluating Stored Prompts

Re-evaluate prompts stored in the database:

```bash
# List and evaluate all recent prompts
curo-prompt list --eval --limit 10

# Evaluate specific tool's prompts
curo-prompt list --tool scan --eval --output reports/
```

## Database Schema

Stored information:

- `id`: Unique prompt ID
- `tool`: Collection tool (scan, codex, cursor, etc.)
- `raw_prompt`: Original prompt text
- `role`, `inputs`, `invariants`, `output_format`: Parsed sections
- `timestamp`: Collection time (Unix timestamp)
- `command`: Command that collected it
- `working_dir`: Working directory at collection time
- `metadata`: Additional metadata (JSON)

## Current Limitations

### ❌ Not Implemented Yet

1. **Automatic CLI Wrapper Collection**
   - Collector exists but tool-specific parsing incomplete
   - Need to implement parsing for codex, cursor commands

2. **Eval Command Doesn't Save**
   - `eval` command evaluates but doesn't save to database
   - Use `scan` command if you want to save

3. **Log File Collection**
   - Not yet implemented
   - Planned: Parse `.cursor/`, `.codex/` log files

4. **Session Capture**
   - Not yet implemented
   - Planned: Capture stdin/stdout from CLI tools

### ✅ Currently Working

1. **File Scanning & Storage**
   - `scan` command finds files, evaluates, and saves to database
   
2. **Prompt Listing**
   - `list` command shows stored prompts
   - Filter by tool, limit results

3. **Re-evaluation**
   - `list --eval` re-evaluates stored prompts

## Recommended Workflow

### For New Users

1. **Initial Setup**
   ```bash
   # Scan your existing prompt files
   curo-prompt scan --repo ./prompts
   ```

2. **View Collected Prompts**
   ```bash
   # See what was collected
   curo-prompt list
   ```

3. **Review Reports**
   ```bash
   # Check generated reports
   ls reports/
   cat reports/your_prompt_report.md
   ```

### For Regular Use

1. **After Creating New Prompts**
   ```bash
   # Scan new prompts
   curo-prompt scan --repo ./prompts
   ```

2. **Quick Evaluation (Without Saving)**
   ```bash
   # Just evaluate without saving
   curo-prompt eval --file new_prompt.md
   ```

3. **Re-evaluate Stored Prompts**
   ```bash
   # Re-evaluate all recent prompts
   curo-prompt list --eval
   ```

## Future Improvements (M2+)

- Automatic collection via CLI wrappers
- Log file parsing
- Session capture
- Git hook integration
- Batch operations on stored prompts

---

**Questions?** Check [Getting Started Guide](./GETTING_STARTED.md) or [Architecture](./ARCHITECTURE.md).

