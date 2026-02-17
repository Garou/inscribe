package molecules

import (
	"inscribe/internal/domain"
	"inscribe/internal/tui/components/atoms"

	"github.com/charmbracelet/huh"
)

// TemplatePicker creates a select field with sub-template options, showing descriptions.
func TemplatePicker(registry domain.TemplateRegistry, group string, value *string) *huh.Select[string] {
	subs, err := registry.GetSubTemplates(group)
	if err != nil {
		return atoms.StyledSelect(group, nil, value)
	}

	options := make([]huh.Option[string], len(subs))
	for i, sub := range subs {
		options[i] = huh.NewOption(sub.Description, sub.Content)
	}
	return atoms.StyledSelect(group, options, value)
}
