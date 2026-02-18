package engine

import (
	"fmt"
	"strings"
	"sync"
	"text/template"

	"inscribe/internal/domain"
)

// NewExtractorFuncMap returns a FuncMap for pass 1 (field extraction).
// Each custom function appends a FieldDefinition to the collector and returns a placeholder.
func NewExtractorFuncMap(collector *[]domain.FieldDefinition, mu *sync.Mutex) template.FuncMap {
	order := 0
	return template.FuncMap{
		"input": func(name, validationType string) string {
			mu.Lock()
			defer mu.Unlock()
			*collector = append(*collector, domain.FieldDefinition{
				Name:           name,
				Type:           domain.FieldInput,
				ValidationType: validationType,
				Order:          order,
			})
			order++
			return fmt.Sprintf("__PLACEHOLDER_%s__", name)
		},
		"autoList": func(source string) string {
			mu.Lock()
			defer mu.Unlock()
			*collector = append(*collector, domain.FieldDefinition{
				Name:   source,
				Type:   domain.FieldAutoList,
				Source: source,
				Order:  order,
			})
			order++
			return fmt.Sprintf("__PLACEHOLDER_%s__", source)
		},
		"templateGroup": func(group string) string {
			mu.Lock()
			defer mu.Unlock()
			*collector = append(*collector, domain.FieldDefinition{
				Name:   group,
				Type:   domain.FieldTemplateGroup,
				Source: group,
				Order:  order,
			})
			order++
			return fmt.Sprintf("__PLACEHOLDER_%s__", group)
		},
		"staticList": func(listName string) string {
			mu.Lock()
			defer mu.Unlock()
			*collector = append(*collector, domain.FieldDefinition{
				Name:   listName,
				Type:   domain.FieldStaticList,
				Source: listName,
				Order:  order,
			})
			order++
			return fmt.Sprintf("__PLACEHOLDER_%s__", listName)
		},
		"indent": func(spaces int, content string) string {
			return indentString(spaces, content)
		},
	}
}

// NewRendererFuncMap returns a FuncMap for pass 2 (rendering with collected values).
func NewRendererFuncMap(values map[string]string) template.FuncMap {
	return template.FuncMap{
		"input": func(name, validationType string) string {
			if v, ok := values[name]; ok {
				return v
			}
			return ""
		},
		"autoList": func(source string) string {
			if v, ok := values[source]; ok {
				return v
			}
			return ""
		},
		"templateGroup": func(group string) string {
			if v, ok := values[group]; ok {
				return v
			}
			return ""
		},
		"staticList": func(listName string) string {
			if v, ok := values[listName]; ok {
				return v
			}
			return ""
		},
		"indent": func(spaces int, content string) string {
			return indentString(spaces, content)
		},
	}
}

// indentString indents each line of content by the given number of spaces.
func indentString(spaces int, content string) string {
	pad := strings.Repeat(" ", spaces)
	lines := strings.Split(content, "\n")
	for i := range lines {
		if lines[i] != "" {
			lines[i] = pad + lines[i]
		}
	}
	return strings.Join(lines, "\n")
}
