package atoms

import "github.com/charmbracelet/huh"

// StyledSelect creates a consistently styled select field.
func StyledSelect(title string, options []huh.Option[string], value *string) *huh.Select[string] {
	return huh.NewSelect[string]().
		Title(title).
		Options(options...).
		Value(value)
}
