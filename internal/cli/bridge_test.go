package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunBridgeDirectRender(t *testing.T) {
	// Set up a simple template that doesn't need k8s or TUI
	dir := t.TempDir()
	outDir := t.TempDir()

	writeFile(t, filepath.Join(dir, "simple.yaml"),
		`{{/* inscribe: type="template" name="simple" command="test" description="Test" */}}
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ input "name" "dns-name" }}
data:
  key: {{ input "value" "string" }}
`)

	err := RunBridge(BridgeConfig{
		TemplateName: "simple",
		TemplateDir:  dir,
		OutputDir:    outDir,
		FlagValues: map[string]string{
			"name":  "my-config",
			"value": "hello",
		},
		Filename: "output.yaml",
	})
	if err != nil {
		t.Fatalf("RunBridge() error: %v", err)
	}

	// Verify output
	data, err := os.ReadFile(filepath.Join(outDir, "output.yaml"))
	if err != nil {
		t.Fatalf("reading output: %v", err)
	}

	content := string(data)
	if !contains(content, "name: my-config") {
		t.Errorf("output should contain 'name: my-config', got:\n%s", content)
	}
	if !contains(content, "key: hello") {
		t.Errorf("output should contain 'key: hello', got:\n%s", content)
	}
}

func TestRunBridgeValidationError(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "simple.yaml"),
		`{{/* inscribe: type="template" name="simple" command="test" description="Test" */}}
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ input "name" "dns-name" }}
`)

	err := RunBridge(BridgeConfig{
		TemplateName: "simple",
		TemplateDir:  dir,
		OutputDir:    t.TempDir(),
		FlagValues: map[string]string{
			"name": "INVALID_DNS",
		},
		Filename: "output.yaml",
	})
	if err == nil {
		t.Error("expected validation error for invalid dns name")
	}
}

func TestRunBridgeTemplateNotFound(t *testing.T) {
	dir := t.TempDir()

	err := RunBridge(BridgeConfig{
		TemplateName: "nonexistent",
		TemplateDir:  dir,
		OutputDir:    t.TempDir(),
		FlagValues:   map[string]string{},
		Filename:     "output.yaml",
	})
	if err == nil {
		t.Error("expected error for nonexistent template")
	}
}

func TestRunBridgeWithSubTemplate(t *testing.T) {
	dir := t.TempDir()
	outDir := t.TempDir()

	writeFile(t, filepath.Join(dir, "main.yaml"),
		`{{/* inscribe: type="template" name="with-sub" command="test" description="Test" */}}
spec:
  resources:
{{ templateGroup "resources" | indent 4 }}
`)

	writeFile(t, filepath.Join(dir, "res-prod.yaml"),
		`{{/* inscribe: type="sub-template" group="resources" description="Production" */}}
memory: "4Gi"
cpu: "2"
`)

	err := RunBridge(BridgeConfig{
		TemplateName: "with-sub",
		TemplateDir:  dir,
		OutputDir:    outDir,
		FlagValues: map[string]string{
			"resources": "Production",
		},
		Filename: "output.yaml",
	})
	if err != nil {
		t.Fatalf("RunBridge() error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(outDir, "output.yaml"))
	if err != nil {
		t.Fatalf("reading output: %v", err)
	}

	content := string(data)
	if !contains(content, "memory: \"4Gi\"") {
		t.Errorf("output should contain sub-template content, got:\n%s", content)
	}
}

func TestRunBridgeWithStaticList(t *testing.T) {
	dir := t.TempDir()
	outDir := t.TempDir()

	writeFile(t, filepath.Join(dir, "main.yaml"),
		`{{/* inscribe: type="template" name="with-list" command="test" description="Test" */}}
method: {{ staticList "methods" }}
`)

	writeFile(t, filepath.Join(dir, "methods.yaml"),
		`{{/* inscribe: type="list" name="methods" */}}
- barmanObjectStore
- volumeSnapshot
`)

	err := RunBridge(BridgeConfig{
		TemplateName: "with-list",
		TemplateDir:  dir,
		OutputDir:    outDir,
		FlagValues: map[string]string{
			"methods": "barmanObjectStore",
		},
		Filename: "output.yaml",
	})
	if err != nil {
		t.Fatalf("RunBridge() error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(outDir, "output.yaml"))
	if err != nil {
		t.Fatalf("reading output: %v", err)
	}

	content := string(data)
	if !contains(content, "method: barmanObjectStore") {
		t.Errorf("output should contain list value, got:\n%s", content)
	}
}

func TestRunBridgeInvalidListValue(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "main.yaml"),
		`{{/* inscribe: type="template" name="with-list" command="test" description="Test" */}}
method: {{ staticList "methods" }}
`)

	writeFile(t, filepath.Join(dir, "methods.yaml"),
		`{{/* inscribe: type="list" name="methods" */}}
- barmanObjectStore
- volumeSnapshot
`)

	err := RunBridge(BridgeConfig{
		TemplateName: "with-list",
		TemplateDir:  dir,
		OutputDir:    t.TempDir(),
		FlagValues: map[string]string{
			"methods": "invalidMethod",
		},
		Filename: "output.yaml",
	})
	if err == nil {
		t.Error("expected error for invalid list value")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writing test file %q: %v", path, err)
	}
}
