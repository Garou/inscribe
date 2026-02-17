package domain

// TemplateMeta describes a main template file.
type TemplateMeta struct {
	Type        string // "template"
	Name        string // e.g., "cnpg-cluster"
	Command     string // e.g., "cluster cnpg"
	Description string
	FilePath    string
}

// SubTemplateMeta describes a sub-template fragment.
type SubTemplateMeta struct {
	Group       string // e.g., "cnpg-resource-templates"
	Description string // e.g., "Production - 4Gi/2CPU"
	Content     string // Raw YAML content (without header comment)
	FilePath    string
}

// StaticListMeta describes a static list of predefined values.
type StaticListMeta struct {
	Name     string
	Items    []string
	FilePath string
}
