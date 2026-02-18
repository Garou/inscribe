package atoms

import (
	"github.com/charmbracelet/huh"
)

// Theme returns the huh theme used throughout the application.
func Theme() *huh.Theme {
	return huh.ThemeCatppuccin()
}
