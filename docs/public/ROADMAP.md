# Roadmap

> 🇰🇷 **Korean**: [한국어 버전](./ROADMAP.ko.md)

> Detailed development plans and duration estimates are available in internal project documents.

## Development Strategy

Incremental development starting from core features. Each milestone is an independently deployable unit.

## M1: MVP (Core Features) - Approximately 5-8 weeks

**Goal**: Complete basic prompt analysis and evaluation tool

### Core Features
- ✅ Basic CLI commands: `scan`, `eval`, `suggest`
- ✅ Static prompt analysis
  - Markdown section structure analysis (ROLE, INPUTS, OUTPUT FORMAT, etc.)
  - Duplicate rule detection
  - Token calculation (Claude, OpenAI)
- ✅ Scoring engine (0-100 overall score)
  - structure, conciseness, instruction_following, risk metrics
- ✅ Markdown report generation
- ✅ Redaction (secret masking)
- ✅ Provider adapter (Claude first, OpenAI added)
- ✅ Basic refactoring suggestions (duplicate removal, token reduction)
- ✅ SQLite storage

**Deliverables**: Usable CLI tool, basic documentation

---

## M2: Dynamic Evaluation - Approximately 3-4 weeks (After M1)

**Goal**: Prompt performance evaluation through actual LLM calls

### Additional Features
- ✅ Multi-sample generation (temperature grid support)
- ✅ Self-consistency validation (key/format match rate across samples)
- ✅ JSON Schema validation
- ✅ Cost and latency tracking
  - Token usage
  - p50/p95 latency
  - Cost calculation
- ✅ JSON report format

**Considerations**: API costs incurred, implementation recommended after necessity review

---

## M3: Automatic Refactoring - Approximately 3-5 weeks

**Goal**: Automatic application of prompt optimization suggestions

### Additional Features
- ✅ Deep prompt structure analysis
- ✅ Few-shot example summarization algorithm
  - Example importance evaluation
  - Minimum example selection
- ✅ Cache optimization guide
- ✅ Auto-apply engine
  - Backup and rollback functionality
  - Before/after comparison reports

**Considerations**: Algorithm research needed, implement basic suggestions first

---

## M4: Team Premium - Approximately 3-5 weeks (Optional)

**Goal**: Advanced features for team usage

### Additional Features
- ✅ Aggregation dashboard
  - Team-wide prompt statistics
  - Trend analysis
- ✅ Policy enforcement
  - Prompt policy validation
  - Auto-rejection rules
- ✅ Internal proxy deployment
  - Centralized monitoring
  - Usage limits

---

## Development Duration Summary

| Milestone | Experienced Developer | Intermediate Developer | Part-time |
|-----------|---------------------|----------------------|-----------|
| **M1 MVP** | 4.7 weeks | 7.4 weeks | 15 weeks |
| **M2 Dynamic Evaluation** | 2.6 weeks | 4.2 weeks | 8.4 weeks |
| **M3 Refactoring** | 3 weeks | 5 weeks | 10 weeks |
| **M4 Team Premium** | 3 weeks | 5 weeks | 10 weeks |
| **Total (M1-M3)** | **10.3 weeks** | **16.6 weeks** | **33.4 weeks** |

## Priorities

### Required (M1)
1. Static analysis (structure, tokens, rules)
2. Scoring engine
3. Basic reports (Markdown)
4. Redaction (security)
5. Provider support (Claude first)

### High Priority
1. Refactoring suggestions (basic)
2. SQLite storage

### Optional
1. Dynamic evaluation (M2) - Consider costs and complexity
2. Automatic refactoring (M3) - Algorithm research needed
3. Team features (M4) - After confirming demand

## Risk Management

- **Dynamic evaluation complexity**: Recommended to re-evaluate necessity after M1 completion
- **Provider API changes**: Abstracted via adapter pattern
- **Algorithm research**: Basic implementation first, refinement later

## Next Steps

1. **Week 1-2**: Project structure and basic CLI implementation
2. **Core features first**: Static analysis → Scoring → Reports
3. **One provider only**: Claude first (API stability)
4. **Dynamic evaluation on hold**: Re-evaluate necessity after M1 completion

---

**Note**: For detailed technology stack information, see [TECH_STACK.md](./TECH_STACK.md).
