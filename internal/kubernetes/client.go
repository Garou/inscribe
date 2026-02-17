package kubernetes

import (
	"context"
	"fmt"

	"inscribe/internal/domain"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var cnpgClusterGVR = schema.GroupVersionResource{
	Group:    "postgresql.cnpg.io",
	Version:  "v1",
	Resource: "clusters",
}

// Client implements domain.KubeClient using real Kubernetes connections.
type Client struct {
	kubeconfig string
}

var _ domain.KubeClient = (*Client)(nil)

// NewClient creates a new Kubernetes client using the given kubeconfig path.
// If kubeconfig is empty, the default loading rules are used.
func NewClient(kubeconfig string) *Client {
	return &Client{kubeconfig: kubeconfig}
}

func (c *Client) loadingRules() *clientcmd.ClientConfigLoadingRules {
	if c.kubeconfig != "" {
		return &clientcmd.ClientConfigLoadingRules{ExplicitPath: c.kubeconfig}
	}
	return clientcmd.NewDefaultClientConfigLoadingRules()
}

func (c *Client) ListContexts() ([]string, error) {
	config, err := c.loadingRules().Load()
	if err != nil {
		return nil, fmt.Errorf("loading kubeconfig: %w", err)
	}

	var contexts []string
	for name := range config.Contexts {
		contexts = append(contexts, name)
	}
	return contexts, nil
}

func (c *Client) ListNamespaces(ctx string) ([]string, error) {
	clientset, err := c.clientsetForContext(ctx)
	if err != nil {
		return nil, err
	}

	nsList, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing namespaces: %w", err)
	}

	var namespaces []string
	for _, ns := range nsList.Items {
		namespaces = append(namespaces, ns.Name)
	}
	return namespaces, nil
}

func (c *Client) ListCNPGClusters(ctx string, namespace string) ([]string, error) {
	dynClient, err := c.dynamicClientForContext(ctx)
	if err != nil {
		return nil, err
	}

	var resource dynamic.ResourceInterface
	if namespace == "" {
		resource = dynClient.Resource(cnpgClusterGVR)
	} else {
		resource = dynClient.Resource(cnpgClusterGVR).Namespace(namespace)
	}

	list, err := resource.List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing CNPG clusters: %w", err)
	}

	var clusters []string
	for _, item := range list.Items {
		clusters = append(clusters, item.GetName())
	}
	return clusters, nil
}

func (c *Client) clientsetForContext(ctx string) (*k8s.Clientset, error) {
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		c.loadingRules(),
		&clientcmd.ConfigOverrides{CurrentContext: ctx},
	).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("building client config for context %q: %w", ctx, err)
	}
	clientset, err := k8s.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("creating clientset for context %q: %w", ctx, err)
	}
	return clientset, nil
}

func (c *Client) dynamicClientForContext(ctx string) (dynamic.Interface, error) {
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		c.loadingRules(),
		&clientcmd.ConfigOverrides{CurrentContext: ctx},
	).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("building client config for context %q: %w", ctx, err)
	}
	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("creating dynamic client for context %q: %w", ctx, err)
	}
	return dynClient, nil
}
