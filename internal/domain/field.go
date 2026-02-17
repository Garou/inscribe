package domain

// FieldType classifies how a template field is populated.
type FieldType int

const (
	FieldManual        FieldType = iota // User-provided with validation
	FieldAutoDetect                     // Pulled from k8s
	FieldTemplateGroup                  // Pick from sub-template group
	FieldList                           // Pick from static predefined list
)

// FieldDefinition is extracted from a template during the first pass.
type FieldDefinition struct {
	Name           string
	Type           FieldType
	ValidationType string // For manual: "dns-name", "integer", "string", etc.
	Source         string // For autoDetect: "namespace", "cnpg-clusters"; for templateGroup/list: group/list name
	Order          int
}

// FieldValue holds a collected value for a field.
type FieldValue struct {
	Definition FieldDefinition
	Value      string
}
