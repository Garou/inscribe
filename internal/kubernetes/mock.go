package kubernetes

import (
	"fmt"

	"inscribe/internal/domain"
)

// MockClient implements domain.KubeClient for testing.
type MockClient struct {
	Contexts     []string
	Namespaces   map[string][]string            // context → namespaces
	CNPGClusters map[string]map[string][]string  // context → namespace → clusters
}

var _ domain.KubeClient = (*MockClient)(nil)

// NewMockClient creates a MockClient with sensible defaults.
func NewMockClient() *MockClient {
	return &MockClient{
		Contexts: []string{"minikube", "production", "staging"},
		Namespaces: map[string][]string{
			"minikube":   {"default", "kube-system", "cnpg-system"},
			"production": {"default", "kube-system", "app-prod", "cnpg-system"},
			"staging":    {"default", "kube-system", "app-staging", "cnpg-system"},
		},
		CNPGClusters: map[string]map[string][]string{
			"production": {
				"app-prod": {"main-db", "analytics-db"},
			},
			"staging": {
				"app-staging": {"staging-db"},
			},
		},
	}
}

func (m *MockClient) ListContexts() ([]string, error) {
	return m.Contexts, nil
}

func (m *MockClient) ListNamespaces(context string) ([]string, error) {
	ns, ok := m.Namespaces[context]
	if !ok {
		return nil, fmt.Errorf("context %q not found", context)
	}
	return ns, nil
}

func (m *MockClient) ListCNPGClusters(context string, namespace string) ([]string, error) {
	contextClusters, ok := m.CNPGClusters[context]
	if !ok {
		return []string{}, nil
	}

	if namespace == "" {
		var all []string
		for _, clusters := range contextClusters {
			all = append(all, clusters...)
		}
		return all, nil
	}

	clusters, ok := contextClusters[namespace]
	if !ok {
		return []string{}, nil
	}
	return clusters, nil
}
