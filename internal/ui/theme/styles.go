package theme

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/demeyerthom/todoist-tui/internal/config"
)

// Styles provides Lipgloss style builders that read color values from a
// theme configuration. Build instances with [NewStyles].
type Styles struct {
	theme *config.ThemeConfig
}

// NewStyles creates a Styles backed by the theme configuration in cfg.
func NewStyles(cfg *config.Config) *Styles {
	return &Styles{theme: &cfg.Theme}
}

// colorFromTheme converts a theme config color string to a lipgloss terminal
// color. Both named ANSI colors and hex values (#RRGGBB) are supported.
func colorFromTheme(s string) lipgloss.TerminalColor {
	return lipgloss.Color(s)
}

// Header returns a bold style using the Header theme color as foreground.
func (s *Styles) Header() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(colorFromTheme(s.theme.Header))
}

// MutedText returns a style using the MutedText theme color as foreground.
func (s *Styles) MutedText() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(colorFromTheme(s.theme.MutedText))
}

// SectionSep returns a style for horizontal separator lines, using the
// MutedText theme color as foreground.
func (s *Styles) SectionSep() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(colorFromTheme(s.theme.MutedText))
}

// PriorityDot returns a bold style using the TaskPriorityHigh theme color
// as foreground.
func (s *Styles) PriorityDot() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(colorFromTheme(s.theme.TaskPriorityHigh))
}

// ActiveItem returns a bold style with the SelectedRow theme color as
// background and the Header theme color as foreground.
func (s *Styles) ActiveItem() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Background(colorFromTheme(s.theme.SelectedRow)).
		Foreground(colorFromTheme(s.theme.Header))
}

// InactiveItem returns a style using the MutedText theme color as foreground.
func (s *Styles) InactiveItem() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(colorFromTheme(s.theme.MutedText))
}

// ColorDot returns a bold style that renders a colored dot character (●).
// The color argument should be a hex string or ANSI color name.
func (s *Styles) ColorDot(color string) lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(color))
}
