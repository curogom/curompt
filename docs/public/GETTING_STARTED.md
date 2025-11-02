# Getting Started Guide

> 🇰🇷 **Korean**: [한국어 버전](./GETTING_STARTED.ko.md)

This is a hands-on guide for first-time users of `curo-prompt`.

## Prerequisites

- Go ≥ 1.23
- macOS or Linux
- Project build complete

## 1. Build and Installation

```bash
# Navigate to project directory
cd /path/to/curo-prompt

# Build
make build

# Check binary
./bin/curo-prompt --help
```

Expected output:
```
curo-prompt is a CLI tool for analyzing, evaluating, and optimizing LLM prompts.

Usage:
  curo-prompt [command]

Available Commands:
  eval        Evaluate and score prompts
  scan        Scan and analyze prompt files in repository
  suggest     Generate prompt improvement suggestions
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

## 3. Basic Testing: eval Command

### 3.1 Evaluate from File

```bash
./bin/curo-prompt eval --file test-prompts/sample.md
```

Expected output:
- Prompt analysis results
- Overall score (0-100)
- Per-metric scores
- Token count
- Improvement suggestions

### 3.2 Evaluate from Standard Input

```bash
cat test-prompts/sample.md | ./bin/curo-prompt eval
```

Or direct input:

```bash
echo "# ROLE\nEngineer\n\n# INPUTS\n- task: string" | ./bin/curo-prompt eval
```

### 3.3 Save Report to File

```bash
./bin/curo-prompt eval --file test-prompts/sample.md --output reports/sample_report.md
```

Report file will be saved to `reports/sample_report.md`.

### 3.4 Change Provider (for token calculation)

```bash
# Use Claude (default)
./bin/curo-prompt eval --file test-prompts/sample.md --provider claude

# Use OpenAI
./bin/curo-prompt eval --file test-prompts/sample.md --provider openai
```

## 4. Batch Scanning: scan Command

### 4.1 Scan Current Directory

```bash
# Scan test-prompts directory
./bin/curo-prompt scan --repo test-prompts
```

Expected output:
```
Found prompt files: 2

[1/2] Analyzing: test-prompts/sample.md
  ✅ Complete - Score: 85.3/100, Report: reports/sample.md_report.md

[2/2] Analyzing: test-prompts/complex.md
  ✅ Complete - Score: 72.1/100, Report: reports/complex.md_report.md

Total 2 files analyzed
Report location: reports
```

### 4.2 Custom Output Directory

```bash
./bin/curo-prompt scan --repo test-prompts --output my-reports
```

### 4.3 Specify File Patterns

```bash
# Scan only .md files
./bin/curo-prompt scan --repo test-prompts --patterns "*.md"

# Multiple patterns
./bin/curo-prompt scan --repo test-prompts --patterns "*.md" --patterns "*.txt"
```

## 5. Improvement Suggestions: suggest Command

### 5.1 Basic Suggestions

```bash
./bin/curo-prompt suggest --file test-prompts/sample.md
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
./bin/curo-prompt suggest --file test-prompts/needs-improvement.md
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
cat test-prompts/sample.md | ./bin/curo-prompt suggest
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
./bin/curo-prompt eval --file my-prompt.md

# 3. Check improvement suggestions
./bin/curo-prompt suggest --file my-prompt.md

# 4. Save report
./bin/curo-prompt eval --file my-prompt.md --output my-report.md

# 5. Check report
cat my-report.md
```

### 6.2 Batch Evaluation of Multiple Prompts

```bash
# Prepare multiple files in prompts directory
mkdir -p prompts
# ... create multiple prompt files ...

# Batch scan and analyze
./bin/curo-prompt scan --repo prompts --output reports

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
ls -la ./bin/curo-prompt
```

### 8.2 Cannot Read File

```bash
# Check file path (use absolute path)
./bin/curo-prompt eval --file /full/path/to/prompt.md

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
./bin/curo-prompt scan --repo test-prompts --output reports
```

## 9. Advanced Usage

### 9.1 Using Pipelines

```bash
# Generate prompt → Evaluate → Suggest
echo "# ROLE\nEngineer" | \
  ./bin/curo-prompt eval | \
  tee eval-output.txt | \
  ./bin/curo-prompt suggest
```

### 9.2 Batch Processing Script

```bash
#!/bin/bash
# evaluate-all.sh

for file in prompts/*.md; do
    echo "Evaluating: $file"
    ./bin/curo-prompt eval --file "$file" --output "reports/$(basename $file .md)_report.md"
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
