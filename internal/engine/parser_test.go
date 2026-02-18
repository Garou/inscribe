package engine

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"inscribe/internal/domain"
)

func TestParserExtractFields(t *testing.T) {
	dir := setupTestTemplates(t)
	reg, err := NewRegistry(dir)
	if err != nil {
		t.Fatalf("NewRegistry() error: %v", err)
	}

	parser := NewParser(reg)
	fields, err := parser.ExtractFields("test-template")
	if err != nil {
		t.Fatalf("ExtractFields() error: %v", err)
	}

	if len(fields) != 3 {
		t.Fatalf("expected 3 fields, got %d: %+v", len(fields), fields)
	}

	// Check extracted fields
	if fields[0].Name != "name" || fields[0].Type != domain.FieldInput {
		t.Errorf("field[0] = %+v, want name/FieldInput", fields[0])
	}
	if fields[0].ValidationType != "dns-name" {
		t.Errorf("field[0].ValidationType = %q, want %q", fields[0].ValidationType, "dns-name")
	}
	if fields[1].Name != "namespace" || fields[1].Type != domain.FieldAutoList {
		t.Errorf("field[1] = %+v, want namespace/FieldAutoList", fields[1])
	}
	if fields[2].Name != "test-group" || fields[2].Type != domain.FieldTemplateGroup {
		t.Errorf("field[2] = %+v, want test-group/FieldTemplateGroup", fields[2])
	}
}

func TestParserRender(t *testing.T) {
	dir := setupTestTemplates(t)
	reg, err := NewRegistry(dir)
	if err != nil {
		t.Fatalf("NewRegistry() error: %v", err)
	}

	parser := NewParser(reg)
	values := map[string]string{
		"name":       "mydb",
		"namespace":  "production",
		"test-group": "keyA: valueA",
	}

	result, err := parser.Render("test-template", values)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	if !strings.Contains(result, "name: mydb") {
		t.Errorf("result should contain 'name: mydb', got:\n%s", result)
	}
	if !strings.Contains(result, "namespace: production") {
		t.Errorf("result should contain 'namespace: production', got:\n%s", result)
	}
	if !strings.Contains(result, "keyA: valueA") {
		t.Errorf("result should contain 'keyA: valueA', got:\n%s", result)
	}
}

func TestParserRoundTrip(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "simple.yaml"),
		`{{/* inscribe: type="template" name="simple" command="simple" description="Simple test" */}}
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ input "name" "dns-name" }}
data:
  key: {{ input "value" "string" }}
`)

	reg, err := NewRegistry(dir)
	if err != nil {
		t.Fatalf("NewRegistry() error: %v", err)
	}

	parser := NewParser(reg)

	// Pass 1: extract
	fields, err := parser.ExtractFields("simple")
	if err != nil {
		t.Fatalf("ExtractFields() error: %v", err)
	}
	if len(fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(fields))
	}

	// Pass 2: render
	values := map[string]string{
		"name":  "my-config",
		"value": "hello-world",
	}
	result, err := parser.Render("simple", values)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}

	expected := `apiVersion: v1
kind: ConfigMap
metadata:
  name: my-config
data:
  key: hello-world
`
	if result != expected {
		t.Errorf("rendered output mismatch:\ngot:\n%s\nwant:\n%s", result, expected)
	}
}

func TestParserTemplateNotFound(t *testing.T) {
	dir := t.TempDir()
	reg, err := NewRegistry(dir)
	if err != nil {
		t.Fatalf("NewRegistry() error: %v", err)
	}

	parser := NewParser(reg)

	_, err = parser.ExtractFields("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent template")
	}

	_, err = parser.Render("nonexistent", nil)
	if err == nil {
		t.Error("expected error for nonexistent template")
	}
}

func TestStripHeader(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "with header",
			content: "{{/* inscribe: type=\"template\" name=\"test\" */}}\napiVersion: v1\n",
			want:    "apiVersion: v1\n",
		},
		{
			name:    "without header",
			content: "apiVersion: v1\nkind: Pod\n",
			want:    "apiVersion: v1\nkind: Pod\n",
		},
		{
			name:    "single line no header",
			content: "apiVersion: v1",
			want:    "apiVersion: v1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write to temp file and test stripHeader
			tmpFile := filepath.Join(t.TempDir(), "test.yaml")
			if err := os.WriteFile(tmpFile, []byte(tt.content), 0644); err != nil {
				t.Fatalf("writing temp file: %v", err)
			}
			got := stripHeader(tt.content)
			if got != tt.want {
				t.Errorf("stripHeader() = %q, want %q", got, tt.want)
			}
		})
	}
}
