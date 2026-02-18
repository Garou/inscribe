package engine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseHeader(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		wantNil  bool
		wantType string
	}{
		{
			name:     "template header",
			line:     `{{/* inscribe: type="template" name="cnpg-cluster" command="cluster cnpg" description="CNPG PostgreSQL Cluster" */}}`,
			wantNil:  false,
			wantType: "template",
		},
		{
			name:     "sub-template header",
			line:     `{{/* inscribe: type="sub-template" group="cnpg-resource-templates" description="Production - 4Gi/2CPU" */}}`,
			wantNil:  false,
			wantType: "sub-template",
		},
		{
			name:     "list header",
			line:     `{{/* inscribe: type="list" name="backup-methods" */}}`,
			wantNil:  false,
			wantType: "list",
		},
		{
			name:    "no header",
			line:    `apiVersion: v1`,
			wantNil: true,
		},
		{
			name:    "regular go template comment",
			line:    `{{/* just a comment */}}`,
			wantNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseHeader(tt.line)
			if tt.wantNil {
				if result != nil {
					t.Errorf("expected nil, got %v", result)
				}
				return
			}
			if result == nil {
				t.Fatal("expected non-nil result")
			}
			if result["type"] != tt.wantType {
				t.Errorf("type = %q, want %q", result["type"], tt.wantType)
			}
		})
	}
}

func TestNewRegistry(t *testing.T) {
	dir := setupTestTemplates(t)

	reg, err := NewRegistry(dir)
	if err != nil {
		t.Fatalf("NewRegistry() error: %v", err)
	}

	// Check templates
	templates := reg.ListTemplates()
	if len(templates) != 1 {
		t.Errorf("expected 1 template, got %d", len(templates))
	}

	tmpl, err := reg.GetTemplate("test-template")
	if err != nil {
		t.Fatalf("GetTemplate() error: %v", err)
	}
	if tmpl.Name != "test-template" {
		t.Errorf("template name = %q, want %q", tmpl.Name, "test-template")
	}
	if tmpl.Command != "test cmd" {
		t.Errorf("template command = %q, want %q", tmpl.Command, "test cmd")
	}

	// Check sub-templates
	subs, err := reg.GetSubTemplates("test-group")
	if err != nil {
		t.Fatalf("GetSubTemplates() error: %v", err)
	}
	if len(subs) != 2 {
		t.Errorf("expected 2 sub-templates, got %d", len(subs))
	}

	// Check static lists
	list, err := reg.GetStaticList("test-list")
	if err != nil {
		t.Fatalf("GetStaticList() error: %v", err)
	}
	if len(list.Items) != 2 {
		t.Errorf("expected 2 list items, got %d", len(list.Items))
	}
	if list.Items[0] != "item-one" || list.Items[1] != "item-two" {
		t.Errorf("list items = %v, want [item-one, item-two]", list.Items)
	}
}

func TestRegistryNotFound(t *testing.T) {
	dir := setupTestTemplates(t)
	reg, err := NewRegistry(dir)
	if err != nil {
		t.Fatalf("NewRegistry() error: %v", err)
	}

	_, err = reg.GetTemplate("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent template")
	}

	_, err = reg.GetSubTemplates("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent sub-template group")
	}

	_, err = reg.GetStaticList("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent static list")
	}
}

func setupTestTemplates(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "main.yaml"),
		`{{/* inscribe: type="template" name="test-template" command="test cmd" description="A test" */}}
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ input "name" "dns-name" }}
  namespace: {{ autoList "namespace" }}
data:
  resources:
{{ templateGroup "test-group" | indent 4 }}
`)

	writeFile(t, filepath.Join(dir, "sub1.yaml"),
		`{{/* inscribe: type="sub-template" group="test-group" description="Option A" */}}
keyA: valueA
`)

	writeFile(t, filepath.Join(dir, "sub2.yaml"),
		`{{/* inscribe: type="sub-template" group="test-group" description="Option B" */}}
keyB: valueB
`)

	writeFile(t, filepath.Join(dir, "list.yaml"),
		`{{/* inscribe: type="list" name="test-list" */}}
- item-one
- item-two
`)

	// Non-template file (should be ignored)
	writeFile(t, filepath.Join(dir, "readme.yaml"),
		`# This is not a template
just: data
`)

	return dir
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writing test file %q: %v", path, err)
	}
}
