package organisms

import (
	"inscribe/internal/domain"
	"inscribe/internal/tui/components/molecules"

	"github.com/charmbracelet/huh"
)

// FieldGroup creates a form group from field definitions, skipping autoDetect fields
// (which are handled separately by context/namespace selection).
func FieldGroup(defs []domain.FieldDefinition, values map[string]*string, registry domain.TemplateRegistry) *huh.Group {
	var fields []huh.Field
	for _, def := range defs {
		val, ok := values[def.Name]
		if !ok {
			s := ""
			values[def.Name] = &s
			val = values[def.Name]
		}

		switch def.Type {
		case domain.FieldManual:
			fields = append(fields, molecules.ManualField(def, val))
		case domain.FieldTemplateGroup:
			fields = append(fields, molecules.TemplatePicker(registry, def.Source, val))
		case domain.FieldList:
			fields = append(fields, molecules.ListPicker(registry, def.Source, val))
		case domain.FieldAutoDetect:
			// Handled by ContextSelectGroup / NamespaceSelectGroup
			continue
		}
	}

	return huh.NewGroup(fields...).Title("Template Fields")
}
