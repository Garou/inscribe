package atoms

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Brand colors
	Primary   = lipgloss.Color("#7C3AED") // Purple
	Secondary = lipgloss.Color("#06B6D4") // Cyan
	Muted     = lipgloss.Color("#6B7280") // Gray
	Success   = lipgloss.Color("#10B981") // Green
	Danger    = lipgloss.Color("#EF4444") // Red

	// Text styles
	TitleStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true)

	DescriptionStyle = lipgloss.NewStyle().
				Foreground(Muted)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(Success)
)

// Theme returns the huh theme used throughout the application.
func Theme() *huh.Theme {
	return huh.ThemeCatppuccin()
}
