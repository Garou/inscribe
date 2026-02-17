package output

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriterWrite(t *testing.T) {
	w := NewWriter()
	dir := t.TempDir()

	content := "apiVersion: v1\nkind: ConfigMap\n"
	path, err := w.Write(content, dir, "test.yaml")
	if err != nil {
		t.Fatalf("Write() error: %v", err)
	}

	expected := filepath.Join(dir, "test.yaml")
	if path != expected {
		t.Errorf("path = %q, want %q", path, expected)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading written file: %v", err)
	}
	if string(data) != content {
		t.Errorf("file content = %q, want %q", string(data), content)
	}
}

func TestWriterCreatesDirectory(t *testing.T) {
	w := NewWriter()
	dir := filepath.Join(t.TempDir(), "sub", "dir")

	_, err := w.Write("test", dir, "out.yaml")
	if err != nil {
		t.Fatalf("Write() error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "out.yaml")); os.IsNotExist(err) {
		t.Error("expected file to exist")
	}
}
