# Getting Started Guide

> 🇰🇷 **Korean**: [한국어 버전](./GETTING_STARTED.ko.md)

This is a hands-on guide for first-time users of `curompt`.

## Prerequisites

- Go ≥ 1.23
- macOS or Linux
- Project build complete

## 1. Build and Installation

```bash
# Navigate to project directory
cd /path/to/curompt

# Build
make build

# Check binary
./bin/curompt --help
```

Expected output:
```
curompt is a CLI tool for analyzing, evaluating, and optimizing LLM prompts.

Usage:
  curompt [command]

Available Commands (surface):
  scan        Scan, parallel-evaluate, and summarize prompts (default entry)
  suggest     Generate prompt improvement suggestions
  collect     Collect prompts from logs (Claude/Codex)

> **Scan note:** `curompt scan` works on prompts already stored in the local DB (via `curompt collect`/`curompt eval`).  
> If the requested path has no history, the CLI offers to auto-collect from Claude Code or Codex logs (Cursor support arrives in v1.1).
```

## 2. Create Test Prompt Files

First, let's create prompt files for testing:

```bash
# Create test directory
mkdir -p test-prompts

# Create sample prompt file
cat > test-prompts/sample.md << 'EOF'
# ROLE
Senior Software Engineer

# INPUTS
- task: string (feature description to develop)
- context: string (project context)
- requirements: array (requirement list)

# INVARIANTS
- Code must be testable
- All public functions must be documented
- Error handling must be explicit

# OUTPUT FORMAT
Return in JSON format with the following structure:
{
  "code": "string",
  "tests": "string",
  "documentation": "string"
}
EOF
```

Or a more complex example:

```bash
cat > test-prompts/complex.md << 'EOF'
# ROLE
Full-stack Developer

# INPUTS
- feature_description: string
- tech_stack: array
- deadline: date

# INVARIANTS
- Implement both frontend and backend
- Follow RESTful API design principles
- Responsive design required

# OUTPUT FORMAT
Implementation plan in Markdown format
EOF
```

## 3. Basic Testing: scan-first

### 3.1 Evaluate from File

```bash
# Rich summary to console (<=100 lines): stats, distribution, Top-N, coaching
./bin/curompt scan --path test-prompts
```

### 3.2 Evaluate from Standard Input

```bash
# Only the worst 5 prompts
./bin/curompt scan --path test-prompts --top 5
```

Or direct input:

```bash
echo "# ROLE\nEngineer\n\n# INPUTS\n- task: string" | ./bin/curompt eval
```

### 3.3 Save Report to File

```bash
# Save a single merged report file
./bin/curompt scan --path test-prompts --output reports --single-output all_reports.md
```

### 3.4 Change Provider (for token calculation)

```bash
# Provider for token/cost metadata
./bin/curompt scan --path test-prompts --provider claude   # default
./bin/curompt scan --path test-prompts --provider openai
```

## 4. Batch Scanning: scan Command

### 4.1 Scan Current Directory

```bash
# Scan test-prompts directory
./bin/curompt scan --path test-prompts
```

Expected output (summary):
```
요약: 총 N, 평균 XX.X, 중앙값 XX.X, 표준편차 X.X, 최저/최고
분포: 0–39:A, 40–59:B, 60–79:C, 80–100:D | IQR(25%:p25, 75%:p75)
개선 우선순위 Top-10:
 1) 40.0  /path/to/worst.md
 ...
잘 작성된 프롬프트 Top-5:
 1) 90.0  /path/to/best.md
개선 가이드:
 - ROLE을 명확히 ... (중략)
```

### 4.2 Custom Output Directory

```bash
./bin/curompt scan --path test-prompts --output my-reports
```

### 4.3 Specify File Patterns

```bash
# Scan only .md files
./bin/curompt scan --path test-prompts --patterns "*.md"

# Multiple patterns
./bin/curompt scan --path test-prompts --patterns "*.md" --patterns "*.txt"
```

## 5. Improvement Suggestions: suggest Command

### 5.1 Basic Suggestions

```bash
./bin/curompt suggest --file test-prompts/sample.md
```

Expected output:
```
📋 Improvement Suggestions:
==================================================

1. ✅ Prompt is already well-structured!

==================================================

Current Score: 85.3 / 100
```

### 5.2 Testing Prompts Needing Improvement

Create a sample file that needs improvement:

```bash
cat > test-prompts/needs-improvement.md << 'EOF'
# ROLE
Engineer
EOF
```

Check suggestions for this file:

```bash
./bin/curompt suggest --file test-prompts/needs-improvement.md
```

Expected output:
```
📋 Improvement Suggestions:
==================================================

1. 🟡 INPUTS section recommended - Specifying inputs makes prompts clearer
2. 🟡 INVARIANTS section recommended - Specify rules and constraints
3. 🟡 OUTPUT FORMAT section recommended - Specifying output format improves result quality
4. 🟠 Overall score has room for improvement - Apply suggestions to achieve higher score

==================================================

Current Score: 45.2 / 100
```

### 5.3 Suggestions from Standard Input

```bash
cat test-prompts/sample.md | ./bin/curompt suggest
```

## 6. Real Workflow Testing

### 6.1 Complete Workflow

```bash
# 1. Prepare prompt file
cat > my-prompt.md << 'EOF'
# ROLE
Developer

# INPUTS
- feature: string

# OUTPUT FORMAT
Code in Go
EOF

# 2. Evaluate and check score
./bin/curompt eval --file my-prompt.md

# 3. Check improvement suggestions
./bin/curompt suggest --file my-prompt.md

# 4. Save report
./bin/curompt eval --file my-prompt.md --output my-report.md

# 5. Check report
cat my-report.md
```

### 6.2 Batch Evaluation of Multiple Prompts

```bash
# Prepare multiple files in prompts directory
mkdir -p prompts
# ... create multiple prompt files ...

# Batch scan and analyze
./bin/curompt scan --path prompts --output reports --single-output all_reports.md

# Check reports
ls -la reports/
```

## 7. Understanding Results

### 7.1 Score Interpretation

- **90-100 points**: Excellent prompt
- **80-89 points**: Good prompt
- **70-79 points**: Room for improvement
- **Below 60**: Major improvement needed

### 7.2 Understanding Metrics

**Structure**
- Section presence
- Duplicate rule removal
- Required section inclusion

**Conciseness**
- Token density
- Unnecessary description removal

**Risk**
- Ambiguous expression detection
- Sensitive information exposure potential

### 7.3 Report File Contents

Report files (`.md`) include:
- Prompt metadata
- Overall score and per-metric scores
- Static analysis results
- Token information
- Improvement suggestions

## 8. Troubleshooting

### 8.1 Binary Not Found

```bash
# Check build
make build

# Verify binary exists
ls -la ./bin/curompt
```

### 8.2 Cannot Read File

```bash
# Check file path (use absolute path)
./bin/curompt eval --file /full/path/to/prompt.md

# Check file permissions
ls -la test-prompts/sample.md
```

### 8.3 Score is 0 or Unexpected Results

- Check prompt format (Markdown section structure)
- Check file encoding (UTF-8)
- Check file content (output with `cat`)

### 8.4 Report Files Not Generated

```bash
# Check output directory permissions
mkdir -p reports
chmod 755 reports

# Try again
./bin/curompt scan --path test-prompts --output reports
```

## 9. Advanced Usage

### 9.1 Using Pipelines

```bash
# Generate prompt → Evaluate → Suggest
echo "# ROLE\nEngineer" | \
  ./bin/curompt eval | \
  tee eval-output.txt | \
  ./bin/curompt suggest
```

### 9.2 Batch Processing Script

```bash
#!/bin/bash
# evaluate-all.sh

for file in prompts/*.md; do
    echo "Evaluating: $file"
    ./bin/curompt eval --file "$file" --output "reports/$(basename $file .md)_report.md"
done
```

## 10. Next Steps

Once testing is successfully completed:

1. **Test with real prompts**: Use prompts from your own projects
2. **Analyze reports**: Review generated reports and apply improvements
3. **Explore documentation**: See [Architecture](./ARCHITECTURE.md), [Config Guide](./CONFIG.md)
4. **Contribute**: Welcome to report issues or suggest improvements

---

**If you encounter issues**: Report via GitHub Issues or check the documentation.
