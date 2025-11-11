package parser

import (
	"bufio"
	"errors"
	"math"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

const (
	sectionRole         = "ROLE"
	sectionInputs       = "INPUTS"
	sectionInvariants   = "INVARIANTS"
	sectionOutputFormat = "OUTPUT_FORMAT"

	defaultFuzzySimilarity   = 0.8
	roleImplicitConfidence   = 0.6
	aliasBaseConfidence      = 0.85
	fuzzyBaseConfidence      = 0.9
	contentConfidenceBoost   = 0.2
	headerConfidence         = 1.0
	defaultConfigRelativeDir = "configs"
	defaultConfigFile        = "sections.yaml"
)

// markdownParser implements Parser interface using enriched parsing and confidence scoring
type markdownParser struct {
	detector *sectionDetector
}

// NewParser creates a new Markdown-based prompt parser
func NewParser() Parser {
	d := newSectionDetector()
	d.loadOverrides(filepath.Join(defaultConfigRelativeDir, defaultConfigFile))
	return &markdownParser{
		detector: d,
	}
}

// Parse parses a prompt from Markdown content
func (p *markdownParser) Parse(content string) (*Prompt, error) {
	prompt := &Prompt{
		Raw: content,
	}

	scanner := bufio.NewScanner(strings.NewReader(content))
	var currentSection string
	var currentConfidence float64
	sectionContent := []string{}

	flushSection := func() {
		if currentSection == "" {
			return
		}
		p.saveSection(prompt, currentSection, sectionContent, currentConfidence)
		currentSection = ""
		currentConfidence = 0.0
		sectionContent = []string{}
	}

	for scanner.Scan() {
		rawLine := scanner.Text()
		line := strings.TrimSpace(rawLine)

		if line == "" {
			continue
		}

		// Markdown heading detection
		if sec, conf, rest := p.detector.detectHeading(rawLine); sec != "" {
			flushSection()
			currentSection = sec
			currentConfidence = conf
			if rest != "" {
				sectionContent = append(sectionContent, strings.TrimSpace(rest))
			}
			continue
		}

		// Alias line detection (e.g., "입력:", "Output Format:")
		if sec, conf, rest := p.detector.detectAliasLine(line); sec != "" {
			flushSection()
			currentSection = sec
			currentConfidence = conf
			if rest != "" {
				sectionContent = append(sectionContent, strings.TrimSpace(rest))
			}
			continue
		}

		// Implicit role sentence (e.g., "당신은 ...", "You are ...")
		if currentSection == "" {
			if sec, conf, text := p.detector.detectImplicitLine(line); sec != "" {
				flushSection()
				currentSection = sec
				currentConfidence = conf
				if text != "" {
					sectionContent = append(sectionContent, strings.TrimSpace(text))
				}
				continue
			}
		}

		if currentSection != "" {
			sectionContent = append(sectionContent, line)
		}
	}

	flushSection()

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return prompt, nil
}

// saveSection saves parsed section content and updates confidence
func (p *markdownParser) saveSection(prompt *Prompt, section string, content []string, baseConfidence float64) {
	if section == "" {
		return
	}

	canonical := canonicalSection(section)
	cleaned := trimNonEmpty(content)
	contentBoost := p.detector.contentConfidence(canonical, cleaned)
	confidence := clampConfidence(baseConfidence + contentBoost)

	switch canonical {
	case sectionRole:
		if len(cleaned) > 0 {
			prompt.Role = strings.Join(cleaned, " ")
		}
		prompt.RoleConfidence = maxFloat(prompt.RoleConfidence, confidence)
	case sectionInputs:
		items := p.cleanListItems(cleaned)
		if len(items) > 0 {
			prompt.Inputs = items
		}
		prompt.InputsConfidence = maxFloat(prompt.InputsConfidence, confidence)
	case sectionInvariants:
		items := p.cleanListItems(cleaned)
		if len(items) > 0 {
			prompt.Invariants = items
		}
		prompt.InvariantsConfidence = maxFloat(prompt.InvariantsConfidence, confidence)
	case sectionOutputFormat:
		if len(cleaned) > 0 {
			prompt.OutputFormat = cleaned
		}
		prompt.OutputFormatConfidence = maxFloat(prompt.OutputFormatConfidence, confidence)
	}
}

// cleanListItems removes bullet prefixes from list items
func (p *markdownParser) cleanListItems(items []string) []string {
	cleaned := make([]string, 0, len(items))
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		switch {
		case strings.HasPrefix(trimmed, "- "):
			trimmed = strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))
		case strings.HasPrefix(trimmed, "* "):
			trimmed = strings.TrimSpace(strings.TrimPrefix(trimmed, "* "))
		case strings.HasPrefix(trimmed, "• "):
			trimmed = strings.TrimSpace(strings.TrimPrefix(trimmed, "• "))
		}
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	return cleaned
}

type sectionDetector struct {
	aliases              map[string][]string
	fuzzySimilarity      float64
	implicitRolePrefixes []string
}

type sectionConfig struct {
	Sections   map[string][]string `yaml:"sections"`
	Thresholds struct {
		Fuzzy float64 `yaml:"fuzzy"`
	} `yaml:"thresholds"`
}

func newSectionDetector() *sectionDetector {
	return &sectionDetector{
		aliases: map[string][]string{
			sectionRole: {
				"ROLE", "역할", "당신은", "you are", "as a", "담당", "position",
			},
			sectionInputs: {
				"INPUTS", "입력", "inputs", "parameters", "요구 입력", "필요 입력",
			},
			sectionInvariants: {
				"INVARIANTS", "불변", "제약", "constraints", "limitations", "규칙", "금지",
			},
			sectionOutputFormat: {
				"OUTPUT FORMAT", "OUTPUTFORMAT", "출력 형식", "응답 형식", "format", "schema", "response", "output",
			},
		},
		fuzzySimilarity: defaultFuzzySimilarity,
		implicitRolePrefixes: []string{
			"당신은", "너는", "you are", "you're", "as a", "역할:", "role:",
		},
	}
}

func (d *sectionDetector) loadOverrides(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) || errors.Is(err, filepath.ErrBadPattern) {
			return
		}
		return
	}

	var cfg sectionConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return
	}

	for key, aliases := range cfg.Sections {
		canonical := canonicalSection(key)
		if _, ok := d.aliases[canonical]; !ok {
			d.aliases[canonical] = []string{}
		}
		d.aliases[canonical] = appendUnique(d.aliases[canonical], aliases...)
	}

	if cfg.Thresholds.Fuzzy > 0 && cfg.Thresholds.Fuzzy < 1 {
		d.fuzzySimilarity = cfg.Thresholds.Fuzzy
	}
}

func (d *sectionDetector) detectHeading(line string) (string, float64, string) {
	trimmed := strings.TrimSpace(line)
	if len(trimmed) == 0 || trimmed[0] != '#' {
		return "", 0, ""
	}

	// Remove leading '#' characters
	title := strings.TrimLeft(trimmed, "#")
	title = strings.TrimSpace(title)
	if title == "" {
		return "", 0, ""
	}

	section, confidence := d.matchAlias(title)
	if section != "" {
		if confidence < headerConfidence {
			confidence = headerConfidence
		}
		return section, confidence, ""
	}

	if idx := strings.Index(title, ":"); idx > 0 {
		label := strings.TrimSpace(title[:idx])
		rest := strings.TrimSpace(title[idx+1:])
		if section, conf := d.matchAlias(label); section != "" {
			if conf < aliasBaseConfidence {
				conf = aliasBaseConfidence
			}
			return section, conf, rest
		}
	}

	if section, conf := d.fuzzyMatch(title); section != "" {
		return section, conf, ""
	}

	return "", 0, ""
}

func (d *sectionDetector) detectAliasLine(line string) (string, float64, string) {
	if line == "" {
		return "", 0, ""
	}

	if idx := strings.Index(line, ":"); idx > 0 {
		label := strings.TrimSpace(line[:idx])
		rest := strings.TrimSpace(line[idx+1:])
		if section, conf := d.matchAlias(label); section != "" {
			if conf < aliasBaseConfidence {
				conf = aliasBaseConfidence
			}
			return section, conf, rest
		}
	}

	if strings.HasSuffix(line, "섹션") {
		if section, conf := d.matchAlias(strings.TrimSuffix(line, "섹션")); section != "" {
			return section, conf, ""
		}
	}

	return "", 0, ""
}

func (d *sectionDetector) detectImplicitLine(line string) (string, float64, string) {
	lower := strings.ToLower(strings.TrimSpace(line))
	for _, prefix := range d.implicitRolePrefixes {
		if strings.HasPrefix(lower, strings.ToLower(prefix)) {
			return sectionRole, roleImplicitConfidence, line
		}
	}
	return "", 0, ""
}

func (d *sectionDetector) contentConfidence(section string, content []string) float64 {
	if len(content) == 0 {
		return 0.0
	}

	switch section {
	case sectionRole, sectionInputs, sectionInvariants, sectionOutputFormat:
		return contentConfidenceBoost
	default:
		return 0.0
	}
}

func (d *sectionDetector) matchAlias(label string) (string, float64) {
	trimmed := strings.TrimSpace(label)
	if trimmed == "" {
		return "", 0
	}

	normalized := normalizeKey(trimmed)

	for section, aliases := range d.aliases {
		if normalized == normalizeKey(section) {
			return section, headerConfidence
		}
		for _, alias := range aliases {
			if strings.EqualFold(trimmed, alias) || normalized == normalizeKey(alias) {
				return section, aliasBaseConfidence
			}
		}
	}

	return "", 0
}

func (d *sectionDetector) fuzzyMatch(label string) (string, float64) {
	normalized := normalizeKey(label)
	if normalized == "" {
		return "", 0
	}

	bestSection := ""
	bestScore := 0.0

	for section, aliases := range d.aliases {
		candidates := append(append([]string{}, aliases...), section)
		for _, candidate := range candidates {
			candidateNormalized := normalizeKey(candidate)
			score := similarity(normalized, candidateNormalized)
			if isSingleTransposition(normalized, candidateNormalized) {
				score = 0.95
			}
			if score >= d.fuzzySimilarity && score > bestScore {
				bestScore = score
				bestSection = section
			}
		}
	}

	if bestSection == "" {
		return "", 0
	}

	confidence := fuzzyBaseConfidence
	if bestScore < 1.0 {
		confidence = fuzzyBaseConfidence - (1.0-bestScore)*0.1
	}
	return bestSection, clampConfidence(confidence)
}

func canonicalSection(section string) string {
	switch normalizeKey(section) {
	case "ROLE":
		return sectionRole
	case "INPUTS":
		return sectionInputs
	case "INVARIANTS":
		return sectionInvariants
	case "OUTPUTFORMAT":
		return sectionOutputFormat
	default:
		return strings.ToUpper(strings.ReplaceAll(section, " ", "_"))
	}
}

func normalizeKey(value string) string {
	var builder strings.Builder
	for _, r := range strings.ToUpper(value) {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func similarity(a, b string) float64 {
	if len(a) == 0 && len(b) == 0 {
		return 1.0
	}
	maxLen := float64(maxInt(len(a), len(b)))
	if maxLen == 0 {
		return 0
	}
	distance := float64(levenshteinDistance(a, b))
	return 1.0 - (distance / maxLen)
}

func levenshteinDistance(a, b string) int {
	la := len(a)
	lb := len(b)

	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}

	matrix := make([][]int, la+1)
	for i := range matrix {
		matrix[i] = make([]int, lb+1)
		matrix[i][0] = i
	}
	for j := 0; j <= lb; j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= la; i++ {
		for j := 1; j <= lb; j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			matrix[i][j] = minInt(
				matrix[i-1][j]+1,
				matrix[i][j-1]+1,
				matrix[i-1][j-1]+cost,
			)
		}
	}

	return matrix[la][lb]
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func clampConfidence(v float64) float64 {
	switch {
	case v < 0:
		return 0
	case v > 1:
		return 1
	default:
		return v
	}
}

func trimNonEmpty(items []string) []string {
	trimmed := make([]string, 0, len(items))
	for _, item := range items {
		if t := strings.TrimSpace(item); t != "" {
			trimmed = append(trimmed, t)
		}
	}
	return trimmed
}

func appendUnique(dst []string, src ...string) []string {
	existing := make(map[string]struct{}, len(dst))
	for _, item := range dst {
		existing[strings.ToLower(strings.TrimSpace(item))] = struct{}{}
	}
	for _, item := range src {
		key := strings.ToLower(strings.TrimSpace(item))
		if key == "" {
			continue
		}
		if _, ok := existing[key]; ok {
			continue
		}
		dst = append(dst, item)
		existing[key] = struct{}{}
	}
	return dst
}

func minInt(values ...int) int {
	min := math.MaxInt32
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return min
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func isSingleTransposition(a, b string) bool {
	if len(a) != len(b) || len(a) < 2 {
		return false
	}

	mismatches := make([]int, 0, 2)
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			mismatches = append(mismatches, i)
			if len(mismatches) > 2 {
				return false
			}
		}
	}

	if len(mismatches) != 2 {
		return false
	}

	i, j := mismatches[0], mismatches[1]
	if j != i+1 {
		return false
	}

	return a[i] == b[j] && a[j] == b[i]
}
