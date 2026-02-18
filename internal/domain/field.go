package domain

// FieldType classifies how a template field is populated.
type FieldType int

const (
	FieldInput         FieldType = iota // User-provided with validation
	FieldAutoList                       // Pulled from k8s
	FieldTemplateGroup                  // Pick from sub-template group
	FieldStaticList                     // Pick from static predefined list
)

// FieldDefinition is extracted from a template during the first pass.
type FieldDefinition struct {
	Name           string
	Type           FieldType
	ValidationType string // For input: "dns-name", "integer", "string", etc.
	Source         string // For autoList: "namespace", "cnpg-clusters"; for templateGroup/staticList: group/list name
	Order          int
}

// FieldValue holds a collected value for a field.
type FieldValue struct {
	Definition FieldDefinition
	Value      string
}
