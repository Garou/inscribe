package molecules

import (
	"inscribe/internal/domain"
	"inscribe/internal/tui/components/atoms"

	"github.com/charmbracelet/huh"
)

// ManualField creates a text input with domain validation for a manual field.
func ManualField(def domain.FieldDefinition, value *string) *huh.Input {
	input := atoms.StyledInput(def.Name, "Enter "+def.Name, value)

	if def.ValidationType != "" {
		vt := def.ValidationType
		input = input.Validate(func(s string) error {
			_, err := domain.ParseValue(vt, s)
			return err
		})
	}

	return input
}
