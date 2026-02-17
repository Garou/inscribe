package domain

// KubeClient provides access to Kubernetes cluster information.
type KubeClient interface {
	ListContexts() ([]string, error)
	ListNamespaces(context string) ([]string, error)
	ListCNPGClusters(context string, namespace string) ([]string, error)
}

// TemplateRegistry indexes and retrieves templates and sub-templates.
type TemplateRegistry interface {
	GetTemplate(name string) (*TemplateMeta, error)
	GetSubTemplates(group string) ([]SubTemplateMeta, error)
	GetStaticList(name string) (*StaticListMeta, error)
	ListTemplates() []TemplateMeta
	ListTemplatesByCommandPrefix(prefix string) []TemplateMeta
}

// ManifestWriter writes rendered manifest content to files.
type ManifestWriter interface {
	Write(content string, outputDir string, filename string) (string, error)
}
