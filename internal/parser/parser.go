package parser

import (
	"bufio"
	"strings"
)

// markdownParser implements Parser interface using simple string parsing
type markdownParser struct{}

// NewParser creates a new Markdown-based prompt parser
func NewParser() Parser {
	return &markdownParser{}
}

// Parse parses a prompt from Markdown content
func (p *markdownParser) Parse(content string) (*Prompt, error) {
	prompt := &Prompt{
		Raw: content,
	}

	scanner := bufio.NewScanner(strings.NewReader(content))
	var currentSection string
	var sectionContent []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "# ") {
			// 이전 섹션 저장
			p.saveSection(prompt, currentSection, sectionContent)

			// 새 섹션 시작
			currentSection = strings.TrimSpace(strings.TrimPrefix(line, "# "))
			sectionContent = []string{}
			continue
		}

		if currentSection != "" && line != "" {
			sectionContent = append(sectionContent, line)
		}
	}

	// 마지막 섹션 저장
	p.saveSection(prompt, currentSection, sectionContent)

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return prompt, nil
}

// saveSection saves parsed section content to the appropriate field
func (p *markdownParser) saveSection(prompt *Prompt, section string, content []string) {
	if len(content) == 0 {
		return
	}

	switch strings.ToUpper(section) {
	case "ROLE":
		// ROLE은 일반적으로 한 줄 또는 여러 줄을 하나로 합침
		prompt.Role = strings.Join(content, " ")
	case "INPUTS":
		prompt.Inputs = p.cleanListItems(content)
	case "INVARIANTS":
		prompt.Invariants = p.cleanListItems(content)
	case "OUTPUT FORMAT", "OUTPUTFORMAT":
		prompt.OutputFormat = content
	}
}

// cleanListItems removes "- " prefix from list items
func (p *markdownParser) cleanListItems(items []string) []string {
	cleaned := make([]string, 0, len(items))
	for _, item := range items {
		cleaned = append(cleaned, strings.TrimPrefix(item, "- "))
	}
	return cleaned
}
