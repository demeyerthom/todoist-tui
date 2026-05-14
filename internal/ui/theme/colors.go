package theme

import "github.com/charmbracelet/lipgloss"

// TodoistColor maps a Todoist project color name to a lipgloss hex color.
// If the name is not recognized, fallback is returned.
func TodoistColor(name string, fallback lipgloss.TerminalColor) lipgloss.TerminalColor {
	switch name {
	case "berry_red":
		return lipgloss.Color("#ED3351")
	case "red":
		return lipgloss.Color("#F74C3C")
	case "orange":
		return lipgloss.Color("#F8931F")
	case "lime_green":
		return lipgloss.Color("#A4D825")
	case "green":
		return lipgloss.Color("#24B34B")
	case "teal":
		return lipgloss.Color("#14A5A0")
	case "blue":
		return lipgloss.Color("#4285F4")
	case "sky_blue":
		return lipgloss.Color("#40B4E5")
	case "purple":
		return lipgloss.Color("#9C6ADE")
	case "pink":
		return lipgloss.Color("#E860A0")
	case "grey":
		return lipgloss.Color("#808080")
	case "magenta":
		return lipgloss.Color("#E64EB2")
	case "peach":
		return lipgloss.Color("#FBA56B")
	case "yellow":
		return lipgloss.Color("#F9D74A")
	case "ivory":
		return lipgloss.Color("#FFFFC2")
	default:
		return fallback
	}
}
