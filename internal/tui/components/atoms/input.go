package atoms

import "github.com/charmbracelet/huh"

// StyledInput creates a consistently styled text input field.
func StyledInput(title, placeholder string, value *string) *huh.Input {
	return huh.NewInput().
		Title(title).
		Placeholder(placeholder).
		Value(value)
}
