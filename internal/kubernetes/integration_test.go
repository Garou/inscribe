//go:build integration

package kubernetes

import (
	"testing"
)

func TestRealClientListContexts(t *testing.T) {
	client := NewClient("")
	contexts, err := client.ListContexts()
	if err != nil {
		t.Fatalf("ListContexts() error: %v", err)
	}
	t.Logf("Contexts: %v", contexts)
	if len(contexts) == 0 {
		t.Error("expected at least one context")
	}
}

func TestRealClientListNamespaces(t *testing.T) {
	client := NewClient("")
	namespaces, err := client.ListNamespaces("minikube")
	if err != nil {
		t.Fatalf("ListNamespaces() error: %v", err)
	}
	t.Logf("Namespaces: %v", namespaces)
	if len(namespaces) == 0 {
		t.Error("expected at least one namespace")
	}
}

func TestRealClientListCNPGClusters(t *testing.T) {
	client := NewClient("")
	clusters, err := client.ListCNPGClusters("minikube", "default")
	if err != nil {
		t.Fatalf("ListCNPGClusters() error: %v", err)
	}
	t.Logf("CNPG Clusters (default): %v", clusters)

	allClusters, err := client.ListCNPGClusters("minikube", "")
	if err != nil {
		t.Fatalf("ListCNPGClusters() all ns error: %v", err)
	}
	t.Logf("CNPG Clusters (all): %v", allClusters)
}
