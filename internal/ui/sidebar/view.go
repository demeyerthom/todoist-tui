package sidebar

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/demeyerthom/todoist-tui/internal/config"
	"github.com/demeyerthom/todoist-tui/internal/ui/theme"
)

// View renders the sidebar panel as a string with ANSI styling suitable for
// display in the terminal. It computes a visible window that guarantees the
// cursor item is always on-screen, then renders each visible item with
// appropriate indentation, color dots, and expand/collapse indicators.
//
// Section headers use bold text styled with [theme.Styles.Header].
// Project and label items include a colored dot resolved via
// [theme.TodoistColor]. Filter items use a 2-space indent to align with
// project names after the dot. The active item (at cursor position) is
// styled with [theme.Styles.ActiveItem]; all other items use
// [theme.Styles.InactiveItem].
func (m Model) View() string {
	if m.height <= 0 || len(m.items) == 0 {
		return ""
	}

	start := m.cursor - m.height + 1
	if start < 0 {
		start = 0
	}

	end := start + m.height
	if end > len(m.items) {
		end = len(m.items)
		start = end - m.height
		if start < 0 {
			start = 0
		}
	}

	var lines []string
	for i := start; i < end; i++ {
		content := renderItem(m.items[i], m.styles, m.cfg)
		if i == m.cursor {
			content = m.styles.ActiveItem().Render(content)
		} else {
			content = m.styles.InactiveItem().Render(content)
		}
		lines = append(lines, content)
	}
	return strings.Join(lines, "\n")
}

// renderItem formats a single SidebarItem's text content with item-level
// styling (color dots, indentation, expand/collapse arrows). It does not
// apply row-level active/inactive highlighting — that is handled by the
// caller [Model.View].
func renderItem(item SidebarItem, styles *theme.Styles, cfg *config.Config) string {
	switch item.Kind {
	case "section-header":
		arrow := "▸"
		if item.Expanded {
			arrow = "▾"
		}
		return styles.Header().Render(arrow + " " + item.Name)

	case "project":
		indent := strings.Repeat(" ", item.Indent*2)
		dotColor := theme.TodoistColor(item.Color, lipgloss.Color(cfg.Theme.NormalText))
		dotStyle := lipgloss.NewStyle().Bold(true).Foreground(dotColor)
		dot := dotStyle.Render("●")

		name := item.Name
		if item.Expandable {
			if item.Expanded {
				name = "▾ " + name
			} else {
				name = "▸ " + name
			}
		}
		return indent + dot + " " + name

	case "filter":
		return "  " + item.Name

	case "label":
		dotColor := theme.TodoistColor(item.Color, lipgloss.Color(cfg.Theme.NormalText))
		dotStyle := lipgloss.NewStyle().Bold(true).Foreground(dotColor)
		dot := dotStyle.Render("●")
		return dot + " " + item.Name

	default:
		return item.Name
	}
}
