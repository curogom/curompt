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

### 2. Prompt Collection Methods (Scan-first)

#### Method 1: File Scanning (Default)

Scan existing prompt files in your repository:

```bash
# Scan current directory for prompt files
curo-prompt scan --repo .

# Scan specific directory
curo-prompt scan --repo ./prompts

# Custom file patterns
curo-prompt scan --repo . --patterns "*.md" --patterns "*.txt"
```

**What happens by default (all in one):**
1. Finds all prompt files (batch collect)
2. Evaluates them in parallel (concurrency = CPU cores)
3. Shows a rich summary to stdout (<= 100 lines): stats, distribution, Top-N, coaching
4. Optionally saves reports only if `--output` is specified (individual or single-file)

**Notes**:
- Default is console summary only. Use `--output DIR` to write files.
- Single merged report: `--single-output all_reports.md` (requires `--output`).

#### Method 2: Log File Collection ✅ (Optional)

#### Method 3: CLI Wrapper Collection ⚠️ (Optional)

Collect prompts from tool history/log files:

```bash
# Collect from Claude Code (current project only)
cd /path/to/project
curo-prompt collect --from claude

# Collect from all projects
curo-prompt collect --from claude --all

# Collect from Codex (current directory as project)
cd /path/to/project
curo-prompt collect --from codex

# Collect from all Codex prompts
curo-prompt collect --from codex --all

# Collect and auto-evaluate
curo-prompt collect --from claude --eval
```

**How it works:**
1. Reads history/log files from `~/.claude/history.jsonl` or `~/.codex/history.jsonl`
2. Parses prompts from log entries
3. Extracts project information (from session files for Codex)
4. Filters by project (if `--all` not used)
5. Saves to database automatically

**Project Filtering:**

**Claude Code:**
- Requires `CLAUDE.md` or `Claude.md` file in project root
- Without `--all`: Only collects prompts from current project
- With `--all`: Collects from all projects

**Codex:**
- Uses current directory as project path (no `CLAUDE.md` needed)
- Matches against session file's `cwd` field
- Without `--all`: Only collects prompts from current directory
- With `--all`: Collects from all projects

**Examples:**

```bash
# Claude Code: Project-specific collection
cd ~/projects/my-app  # Must have CLAUDE.md
curo-prompt collect --from claude
# → Collects only prompts from ~/projects/my-app

# Codex: Project-specific collection
cd ~/projects/my-app  # No CLAUDE.md needed
curo-prompt collect --from codex
# → Collects only prompts from ~/projects/my-app

# Collect from all projects
curo-prompt collect --from claude --all
curo-prompt collect --from codex --all
# → Collects from all projects in history
```

#### Deprecated: Direct Evaluation (`eval`)

Single-file evaluation is covered by `scan` (pointing `--repo` to a file or using patterns).
If needed, `eval` can be used but is no longer highlighted in docs.

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

3. **Cursor IDE Collection**
   - Cursor log file parsing not yet implemented
   - Planned: Parse Cursor workspace logs

4. **Session Capture**
   - Not yet implemented
   - Planned: Capture stdin/stdout from CLI tools

### ✅ Currently Working

1. **File Scanning & Storage**
   - `scan` command finds files, evaluates, and saves to database
   
2. **Log File Collection**
   - `collect` command parses history files from Claude Code and Codex
   - Project-specific filtering support
   - Automatic project detection from session files (Codex)

3. **Prompt Listing**
   - `list` command shows stored prompts
   - Filter by tool, limit results

4. **Re-evaluation**
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

