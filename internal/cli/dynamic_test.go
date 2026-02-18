package cli

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildDynamicCommandsTreeStructure(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "cluster.yaml"),
		`{{/* inscribe: type="template" name="cnpg-cluster" command="cluster cnpg" description="CNPG PostgreSQL Cluster" */}}
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ input "name" "dns-name" }}
`)

	cmds := BuildDynamicCommands(dir)
	if len(cmds) != 1 {
		t.Fatalf("expected 1 parent command, got %d", len(cmds))
	}

	parent := cmds[0]
	if parent.Use != "cluster" {
		t.Errorf("expected parent Use='cluster', got %q", parent.Use)
	}

	if !parent.HasSubCommands() {
		t.Fatal("expected parent to have subcommands")
	}

	leaf, _, err := parent.Find([]string{"cnpg"})
	if err != nil {
		t.Fatalf("finding leaf command: %v", err)
	}
	if leaf.Use != "cnpg" {
		t.Errorf("expected leaf Use='cnpg', got %q", leaf.Use)
	}
}

func TestBuildDynamicCommandsDynamicFlags(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "cluster.yaml"),
		`{{/* inscribe: type="template" name="cnpg-cluster" command="cluster cnpg" description="CNPG Cluster" */}}
apiVersion: v1
kind: Cluster
metadata:
  name: {{ input "name" "dns-name" }}
  namespace: {{ autoList "namespace" }}
spec:
  instances: {{ input "instances" "integer" }}
`)

	cmds := BuildDynamicCommands(dir)
	if len(cmds) != 1 {
		t.Fatalf("expected 1 parent, got %d", len(cmds))
	}

	leaf, _, _ := cmds[0].Find([]string{"cnpg"})

	// Check dynamic field flags
	for _, flagName := range []string{"name", "namespace", "instances"} {
		f := leaf.Flags().Lookup(flagName)
		if f == nil {
			t.Errorf("expected flag --%s to be registered", flagName)
		}
	}

	// Check standard flags
	for _, flagName := range []string{"context", "filename"} {
		f := leaf.Flags().Lookup(flagName)
		if f == nil {
			t.Errorf("expected standard flag --%s to be registered", flagName)
		}
	}
}

func TestBuildDynamicCommandsSharedParent(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "cluster.yaml"),
		`{{/* inscribe: type="template" name="cnpg-cluster" command="cluster cnpg" description="CNPG Cluster" */}}
name: {{ input "name" "dns-name" }}
`)
	writeFile(t, filepath.Join(dir, "backup.yaml"),
		`{{/* inscribe: type="template" name="cnpg-backup" command="backup cnpg" description="CNPG Backup" */}}
name: {{ input "name" "dns-name" }}
`)

	cmds := BuildDynamicCommands(dir)
	if len(cmds) != 2 {
		t.Fatalf("expected 2 parent commands, got %d", len(cmds))
	}

	names := make(map[string]bool)
	for _, cmd := range cmds {
		names[cmd.Use] = true
	}
	if !names["cluster"] || !names["backup"] {
		t.Errorf("expected 'cluster' and 'backup' parents, got %v", names)
	}
}

func TestBuildDynamicCommandsMultipleLeafsSameParent(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "cnpg-cluster.yaml"),
		`{{/* inscribe: type="template" name="cnpg-cluster" command="cluster cnpg" description="CNPG Cluster" */}}
name: {{ input "name" "dns-name" }}
`)
	writeFile(t, filepath.Join(dir, "other-cluster.yaml"),
		`{{/* inscribe: type="template" name="other-cluster" command="cluster other" description="Other Cluster" */}}
name: {{ input "name" "dns-name" }}
`)

	cmds := BuildDynamicCommands(dir)
	if len(cmds) != 1 {
		t.Fatalf("expected 1 parent command, got %d", len(cmds))
	}

	parent := cmds[0]
	if parent.Use != "cluster" {
		t.Errorf("expected parent Use='cluster', got %q", parent.Use)
	}

	subs := parent.Commands()
	if len(subs) != 2 {
		t.Fatalf("expected 2 leaf commands under cluster, got %d", len(subs))
	}

	leafNames := make(map[string]bool)
	for _, sub := range subs {
		leafNames[sub.Use] = true
	}
	if !leafNames["cnpg"] || !leafNames["other"] {
		t.Errorf("expected 'cnpg' and 'other' leaves, got %v", leafNames)
	}
}

func TestBuildDynamicCommandsInvalidDir(t *testing.T) {
	cmds := BuildDynamicCommands("/nonexistent/path/that/does/not/exist")
	if cmds != nil {
		t.Errorf("expected nil for invalid dir, got %d commands", len(cmds))
	}
}

func TestBuildDynamicCommandsEmptyDir(t *testing.T) {
	dir := t.TempDir()
	cmds := BuildDynamicCommands(dir)
	if cmds != nil {
		t.Errorf("expected nil for empty dir, got %d commands", len(cmds))
	}
}

func TestFlagDescriptionTemplateGroup(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "cluster.yaml"),
		`{{/* inscribe: type="template" name="cnpg-cluster" command="cluster cnpg" description="CNPG Cluster" */}}
{{ templateGroup "resources" | indent 4 }}
`)
	writeFile(t, filepath.Join(dir, "res-prod.yaml"),
		`{{/* inscribe: type="sub-template" group="resources" description="Production - 4Gi/2CPU" */}}
memory: "4Gi"
`)
	writeFile(t, filepath.Join(dir, "res-qa.yaml"),
		`{{/* inscribe: type="sub-template" group="resources" description="QA - 2Gi/1CPU" */}}
memory: "2Gi"
`)

	cmds := BuildDynamicCommands(dir)
	if len(cmds) != 1 {
		t.Fatalf("expected 1 parent, got %d", len(cmds))
	}

	leaf, _, _ := cmds[0].Find([]string{"cnpg"})
	f := leaf.Flags().Lookup("resources")
	if f == nil {
		t.Fatal("expected --resources flag")
	}

	if !strings.Contains(f.Usage, "Production - 4Gi/2CPU") {
		t.Errorf("expected flag description to contain sub-template option, got %q", f.Usage)
	}
	if !strings.Contains(f.Usage, "QA - 2Gi/1CPU") {
		t.Errorf("expected flag description to contain sub-template option, got %q", f.Usage)
	}
}

func TestFlagDescriptionList(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "backup.yaml"),
		`{{/* inscribe: type="template" name="cnpg-backup" command="backup cnpg" description="CNPG Backup" */}}
method: {{ staticList "methods" }}
`)
	writeFile(t, filepath.Join(dir, "methods.yaml"),
		`{{/* inscribe: type="list" name="methods" */}}
- barmanObjectStore
- volumeSnapshot
`)

	cmds := BuildDynamicCommands(dir)
	if len(cmds) != 1 {
		t.Fatalf("expected 1 parent, got %d", len(cmds))
	}

	leaf, _, _ := cmds[0].Find([]string{"cnpg"})
	f := leaf.Flags().Lookup("methods")
	if f == nil {
		t.Fatal("expected --methods flag")
	}

	if !strings.Contains(f.Usage, "barmanObjectStore") {
		t.Errorf("expected flag description to list items, got %q", f.Usage)
	}
	if !strings.Contains(f.Usage, "volumeSnapshot") {
		t.Errorf("expected flag description to list items, got %q", f.Usage)
	}
}

func TestFlagDescriptionManual(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "tmpl.yaml"),
		`{{/* inscribe: type="template" name="test" command="test cmd" description="Test" */}}
name: {{ input "name" "dns-name" }}
`)

	cmds := BuildDynamicCommands(dir)
	leaf, _, _ := cmds[0].Find([]string{"cmd"})
	f := leaf.Flags().Lookup("name")
	if f == nil {
		t.Fatal("expected --name flag")
	}

	if !strings.Contains(f.Usage, "dns-name") {
		t.Errorf("expected flag description to mention validation type, got %q", f.Usage)
	}
}

func TestFlagDescriptionAutoDetect(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "tmpl.yaml"),
		`{{/* inscribe: type="template" name="test" command="test cmd" description="Test" */}}
ns: {{ autoList "namespace" }}
`)

	cmds := BuildDynamicCommands(dir)
	leaf, _, _ := cmds[0].Find([]string{"cmd"})
	f := leaf.Flags().Lookup("namespace")
	if f == nil {
		t.Fatal("expected --namespace flag")
	}

	if !strings.Contains(f.Usage, "auto-listed") {
		t.Errorf("expected flag description to note auto-listing, got %q", f.Usage)
	}
}
