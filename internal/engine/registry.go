package engine

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"inscribe/internal/domain"
)

// Registry implements domain.TemplateRegistry by scanning a directory for templates.
type Registry struct {
	templates    map[string]*domain.TemplateMeta
	subTemplates map[string][]domain.SubTemplateMeta
	staticLists  map[string]*domain.StaticListMeta
}

var _ domain.TemplateRegistry = (*Registry)(nil)

// headerRegexp matches the inscribe header comment: {{/* inscribe: key="value" ... */}}
var headerRegexp = regexp.MustCompile(`\{\{/\*\s*inscribe:\s*(.+?)\s*\*/\}\}`)

// kvRegexp matches key="value" pairs within the header.
var kvRegexp = regexp.MustCompile(`(\w+)="([^"]*)"`)

// NewRegistry scans the given directory recursively and builds a template registry.
func NewRegistry(dir string) (*Registry, error) {
	r := &Registry{
		templates:    make(map[string]*domain.TemplateMeta),
		subTemplates: make(map[string][]domain.SubTemplateMeta),
		staticLists:  make(map[string]*domain.StaticListMeta),
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}
		return r.processFile(path)
	})
	if err != nil {
		return nil, fmt.Errorf("scanning template directory %q: %w", dir, err)
	}

	return r, nil
}

func (r *Registry) processFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("opening %q: %w", path, err)
	}
	defer func() { _ = f.Close() }()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		return nil // empty file
	}
	firstLine := scanner.Text()

	header := parseHeader(firstLine)
	if header == nil {
		return nil // no inscribe header, skip
	}

	switch header["type"] {
	case "template":
		r.templates[header["name"]] = &domain.TemplateMeta{
			Type:        "template",
			Name:        header["name"],
			Command:     header["command"],
			Description: header["description"],
			FilePath:    path,
		}
	case "sub-template":
		content, err := readContentAfterHeader(scanner)
		if err != nil {
			return fmt.Errorf("reading sub-template content from %q: %w", path, err)
		}
		r.subTemplates[header["group"]] = append(r.subTemplates[header["group"]], domain.SubTemplateMeta{
			Group:       header["group"],
			Description: header["description"],
			Content:     content,
			FilePath:    path,
		})
	case "list":
		items, err := parseListItems(scanner)
		if err != nil {
			return fmt.Errorf("reading list items from %q: %w", path, err)
		}
		r.staticLists[header["name"]] = &domain.StaticListMeta{
			Name:     header["name"],
			Items:    items,
			FilePath: path,
		}
	default:
		return fmt.Errorf("unknown inscribe type %q in %q", header["type"], path)
	}

	return nil
}

// parseHeader extracts key-value pairs from an inscribe header line.
func parseHeader(line string) map[string]string {
	match := headerRegexp.FindStringSubmatch(line)
	if match == nil {
		return nil
	}
	kvs := kvRegexp.FindAllStringSubmatch(match[1], -1)
	if len(kvs) == 0 {
		return nil
	}
	result := make(map[string]string)
	for _, kv := range kvs {
		result[kv[1]] = kv[2]
	}
	return result
}

// readContentAfterHeader reads all remaining content after the header line.
func readContentAfterHeader(scanner *bufio.Scanner) (string, error) {
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("scanning content: %w", err)
	}
	return strings.Join(lines, "\n"), nil
}

// parseListItems reads YAML list items (lines starting with "- ").
func parseListItems(scanner *bufio.Scanner) ([]string, error) {
	var items []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "- ") {
			items = append(items, strings.TrimPrefix(line, "- "))
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning list items: %w", err)
	}
	return items, nil
}

func (r *Registry) GetTemplate(name string) (*domain.TemplateMeta, error) {
	t, ok := r.templates[name]
	if !ok {
		return nil, fmt.Errorf("template %q not found", name)
	}
	return t, nil
}

func (r *Registry) GetSubTemplates(group string) ([]domain.SubTemplateMeta, error) {
	st, ok := r.subTemplates[group]
	if !ok {
		return nil, fmt.Errorf("sub-template group %q not found", group)
	}
	return st, nil
}

func (r *Registry) GetStaticList(name string) (*domain.StaticListMeta, error) {
	sl, ok := r.staticLists[name]
	if !ok {
		return nil, fmt.Errorf("static list %q not found", name)
	}
	return sl, nil
}

func (r *Registry) ListTemplates() []domain.TemplateMeta {
	var result []domain.TemplateMeta
	for _, t := range r.templates {
		result = append(result, *t)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}

func (r *Registry) ListTemplatesByCommandPrefix(prefix string) []domain.TemplateMeta {
	var result []domain.TemplateMeta
	for _, t := range r.templates {
		if strings.HasPrefix(t.Command, prefix) {
			result = append(result, *t)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}
