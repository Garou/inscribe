package molecules

import (
	"inscribe/internal/domain"
	"inscribe/internal/tui/components/atoms"

	"github.com/charmbracelet/huh"
)

// K8sContextSelect creates a select field populated with available Kubernetes contexts.
func K8sContextSelect(client domain.KubeClient, value *string) *huh.Select[string] {
	contexts, err := client.ListContexts()
	if err != nil {
		contexts = []string{}
	}
	options := make([]huh.Option[string], len(contexts))
	for i, ctx := range contexts {
		options[i] = huh.NewOption(ctx, ctx)
	}
	return atoms.StyledSelect("Kubernetes Context", options, value)
}

// K8sNamespaceSelect creates a select field populated with namespaces for a given context.
func K8sNamespaceSelect(client domain.KubeClient, context string, value *string) *huh.Select[string] {
	namespaces, err := client.ListNamespaces(context)
	if err != nil {
		namespaces = []string{}
	}
	options := make([]huh.Option[string], len(namespaces))
	for i, ns := range namespaces {
		options[i] = huh.NewOption(ns, ns)
	}
	return atoms.StyledSelect("Namespace", options, value)
}

// K8sCNPGClusterSelect creates a select field populated with CNPG clusters.
func K8sCNPGClusterSelect(client domain.KubeClient, context, namespace string, value *string) *huh.Select[string] {
	clusters, err := client.ListCNPGClusters(context, namespace)
	if err != nil {
		clusters = []string{}
	}
	options := make([]huh.Option[string], len(clusters))
	for i, c := range clusters {
		options[i] = huh.NewOption(c, c)
	}
	return atoms.StyledSelect("CNPG Cluster", options, value)
}
