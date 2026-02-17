package kubernetes

import (
	"testing"
)

func TestMockClientListContexts(t *testing.T) {
	mock := NewMockClient()
	contexts, err := mock.ListContexts()
	if err != nil {
		t.Fatalf("ListContexts() error: %v", err)
	}
	if len(contexts) != 3 {
		t.Errorf("expected 3 contexts, got %d", len(contexts))
	}
}

func TestMockClientListNamespaces(t *testing.T) {
	mock := NewMockClient()

	ns, err := mock.ListNamespaces("minikube")
	if err != nil {
		t.Fatalf("ListNamespaces() error: %v", err)
	}
	if len(ns) != 3 {
		t.Errorf("expected 3 namespaces for minikube, got %d", len(ns))
	}

	_, err = mock.ListNamespaces("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent context")
	}
}

func TestMockClientListCNPGClusters(t *testing.T) {
	mock := NewMockClient()

	// Specific namespace
	clusters, err := mock.ListCNPGClusters("production", "app-prod")
	if err != nil {
		t.Fatalf("ListCNPGClusters() error: %v", err)
	}
	if len(clusters) != 2 {
		t.Errorf("expected 2 clusters, got %d", len(clusters))
	}

	// All namespaces
	clusters, err = mock.ListCNPGClusters("production", "")
	if err != nil {
		t.Fatalf("ListCNPGClusters() all namespaces error: %v", err)
	}
	if len(clusters) != 2 {
		t.Errorf("expected 2 clusters across all namespaces, got %d", len(clusters))
	}

	// No clusters in context
	clusters, err = mock.ListCNPGClusters("minikube", "")
	if err != nil {
		t.Fatalf("ListCNPGClusters() empty error: %v", err)
	}
	if len(clusters) != 0 {
		t.Errorf("expected 0 clusters for minikube, got %d", len(clusters))
	}
}
