package organisms

import (
	"inscribe/internal/domain"
	"inscribe/internal/tui/components/molecules"

	"github.com/charmbracelet/huh"
)

// ContextSelectGroup creates a form group for selecting Kubernetes context and namespace.
func ContextSelectGroup(client domain.KubeClient, context, namespace *string) *huh.Group {
	return huh.NewGroup(
		molecules.K8sContextSelect(client, context),
	).Title("Kubernetes Connection")
}

// NamespaceSelectGroup creates a form group for selecting a namespace after context is chosen.
func NamespaceSelectGroup(client domain.KubeClient, context string, namespace *string) *huh.Group {
	return huh.NewGroup(
		molecules.K8sNamespaceSelect(client, context, namespace),
	).Title("Namespace Selection")
}
