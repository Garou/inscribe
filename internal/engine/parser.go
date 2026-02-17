package engine

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"
	"text/template"

	"inscribe/internal/domain"
)

// Parser handles the two-pass template processing.
type Parser struct {
	registry domain.TemplateRegistry
}

// NewParser creates a new template parser with the given registry.
func NewParser(registry domain.TemplateRegistry) *Parser {
	return &Parser{registry: registry}
}

// ExtractFields performs pass 1: parses the template and extracts all FieldDefinitions.
func (p *Parser) ExtractFields(templateName string) ([]domain.FieldDefinition, error) {
	meta, err := p.registry.GetTemplate(templateName)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(meta.FilePath)
	if err != nil {
		return nil, fmt.Errorf("reading template %q: %w", meta.FilePath, err)
	}

	// Strip the header line
	templateContent := stripHeader(string(content))

	var fields []domain.FieldDefinition
	var mu sync.Mutex
	funcMap := NewExtractorFuncMap(&fields, &mu)

	tmpl, err := template.New(templateName).Funcs(funcMap).Parse(templateContent)
	if err != nil {
		return nil, fmt.Errorf("parsing template %q: %w", templateName, err)
	}

	// Execute with nil data - we just want the side effects (field collection)
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, nil); err != nil {
		return nil, fmt.Errorf("executing extraction pass for %q: %w", templateName, err)
	}

	return fields, nil
}

// Render performs pass 2: renders the template with the given values.
func (p *Parser) Render(templateName string, values map[string]string) (string, error) {
	meta, err := p.registry.GetTemplate(templateName)
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(meta.FilePath)
	if err != nil {
		return "", fmt.Errorf("reading template %q: %w", meta.FilePath, err)
	}

	templateContent := stripHeader(string(content))

	funcMap := NewRendererFuncMap(values)

	tmpl, err := template.New(templateName).Funcs(funcMap).Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("parsing template %q: %w", templateName, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, nil); err != nil {
		return "", fmt.Errorf("rendering template %q: %w", templateName, err)
	}

	return buf.String(), nil
}

// stripHeader removes the first line if it contains an inscribe header.
func stripHeader(content string) string {
	lines := strings.SplitN(content, "\n", 2)
	if len(lines) < 2 {
		return content
	}
	if headerRegexp.MatchString(lines[0]) {
		return lines[1]
	}
	return content
}
