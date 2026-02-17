package engine

import (
	"sync"
	"testing"
	"text/template"

	"inscribe/internal/domain"

	"bytes"
)

func TestExtractorFuncMap(t *testing.T) {
	var fields []domain.FieldDefinition
	var mu sync.Mutex
	fm := NewExtractorFuncMap(&fields, &mu)

	tmplStr := `name={{ manual "name" "dns-name" }} ns={{ autoDetect "namespace" }} group={{ templateGroup "resources" }} pick={{ list "items" }}`
	tmpl, err := template.New("test").Funcs(fm).Parse(tmplStr)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, nil); err != nil {
		t.Fatalf("execute error: %v", err)
	}

	if len(fields) != 4 {
		t.Fatalf("expected 4 fields, got %d", len(fields))
	}

	// Check field types and order
	expected := []struct {
		name  string
		ftype domain.FieldType
		order int
	}{
		{"name", domain.FieldManual, 0},
		{"namespace", domain.FieldAutoDetect, 1},
		{"resources", domain.FieldTemplateGroup, 2},
		{"items", domain.FieldList, 3},
	}

	for i, e := range expected {
		if fields[i].Name != e.name {
			t.Errorf("field[%d].Name = %q, want %q", i, fields[i].Name, e.name)
		}
		if fields[i].Type != e.ftype {
			t.Errorf("field[%d].Type = %d, want %d", i, fields[i].Type, e.ftype)
		}
		if fields[i].Order != e.order {
			t.Errorf("field[%d].Order = %d, want %d", i, fields[i].Order, e.order)
		}
	}
}

func TestRendererFuncMap(t *testing.T) {
	values := map[string]string{
		"name":      "mydb",
		"namespace": "default",
		"resources": "key: value",
		"items":     "picked-item",
	}
	fm := NewRendererFuncMap(values)

	tmplStr := `name={{ manual "name" "dns-name" }} ns={{ autoDetect "namespace" }} group={{ templateGroup "resources" }} pick={{ list "items" }}`
	tmpl, err := template.New("test").Funcs(fm).Parse(tmplStr)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, nil); err != nil {
		t.Fatalf("execute error: %v", err)
	}

	got := buf.String()
	want := "name=mydb ns=default group=key: value pick=picked-item"
	if got != want {
		t.Errorf("rendered = %q, want %q", got, want)
	}
}

func TestIndentString(t *testing.T) {
	tests := []struct {
		name    string
		spaces  int
		content string
		want    string
	}{
		{
			name:    "single line",
			spaces:  4,
			content: "hello",
			want:    "    hello",
		},
		{
			name:    "multi line",
			spaces:  4,
			content: "line1\nline2\nline3",
			want:    "    line1\n    line2\n    line3",
		},
		{
			name:    "empty lines preserved",
			spaces:  2,
			content: "line1\n\nline3",
			want:    "  line1\n\n  line3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := indentString(tt.spaces, tt.content)
			if got != tt.want {
				t.Errorf("indentString(%d, %q) = %q, want %q", tt.spaces, tt.content, got, tt.want)
			}
		})
	}
}
