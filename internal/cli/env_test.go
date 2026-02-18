package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestEnvCmdDefaultDir(t *testing.T) {
	cmd := newEnvCmd()

	// Capture output
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{})

	// Set templateDir to a known value (simulating the global)
	oldDir := templateDir
	templateDir = "templates"
	defer func() { templateDir = oldDir }()

	if err := cmd.Execute(); err != nil {
		t.Fatalf("env command error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "INSCRIBE_TEMPLATE_DIR=") {
		t.Errorf("expected output to contain INSCRIBE_TEMPLATE_DIR=, got: %s", output)
	}
}

func TestEnvCmdWithArg(t *testing.T) {
	cmd := newEnvCmd()

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"/custom/path"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("env command error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "/custom/path") {
		t.Errorf("expected output to contain /custom/path, got: %s", output)
	}
}

func TestEnvCmdTooManyArgs(t *testing.T) {
	cmd := newEnvCmd()

	var buf bytes.Buffer
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"arg1", "arg2"})

	if err := cmd.Execute(); err == nil {
		t.Error("expected error for too many args")
	}
}
