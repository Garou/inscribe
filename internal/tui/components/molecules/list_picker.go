package molecules

import (
	"inscribe/internal/domain"
	"inscribe/internal/tui/components/atoms"

	"github.com/charmbracelet/huh"
)

// ListPicker creates a select field with items from a static list.
func ListPicker(registry domain.TemplateRegistry, listName string, value *string) *huh.Select[string] {
	list, err := registry.GetStaticList(listName)
	if err != nil {
		return atoms.StyledSelect(listName, nil, value)
	}

	options := make([]huh.Option[string], len(list.Items))
	for i, item := range list.Items {
		options[i] = huh.NewOption(item, item)
	}
	return atoms.StyledSelect(listName, options, value)
}
