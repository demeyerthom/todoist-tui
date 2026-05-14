package sidebar

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/demeyerthom/todoist-tui/internal/ui/msg"
)

// Update implements tea.Model. It handles keyboard input for cursor navigation,
// section/project expand/collapse, and item selection. It also responds to
// SyncCompleteMsg by reloading the item list and WindowSizeMsg by updating
// panel dimensions.
func (m Model) Update(teaMsg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := teaMsg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

	case msg.SyncCompleteMsg:
		if err := (&m).loadItems(); err != nil {
			return m, nil
		}
		(&m).clampCursor()
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

// handleKeyMsg processes a Bubbletea key event, dispatching to the
// appropriate handler based on the configured key bindings.
func (m Model) handleKeyMsg(keyMsg tea.KeyMsg) (tea.Model, tea.Cmd) {
	keyStr := keyMsg.String()

	// Look up configured key bindings from the Normal mode keymap.
	// Fall back to sensible defaults for actions not in the keymap.
	keyUp := m.cfg.Keymap.Normal["up"]
	keyDown := m.cfg.Keymap.Normal["down"]
	keyEnter := m.cfg.Keymap.Normal["enter"]
	keyGoTop := m.cfg.Keymap.Normal["go_top"]
	keyGoBottom := m.cfg.Keymap.Normal["go_bottom"]

	keyExpand := m.cfg.Keymap.Normal["expand"]
	if keyExpand == "" {
		keyExpand = "l"
	}
	keyCollapse := m.cfg.Keymap.Normal["collapse"]
	if keyCollapse == "" {
		keyCollapse = "h"
	}

	switch {
	case keyMatches(keyStr, keyDown):
		if m.cursor < len(m.items)-1 {
			m.cursor++
		}

	case keyMatches(keyStr, keyUp):
		if m.cursor > 0 {
			m.cursor--
		}

	case keyMatches(keyStr, keyGoTop):
		m.cursor = 0

	case keyMatches(keyStr, keyGoBottom):
		if len(m.items) > 0 {
			m.cursor = len(m.items) - 1
		}

	case keyMatches(keyStr, keyEnter):
		return m.handleEnter()

	case keyMatches(keyStr, keyExpand):
		return m.handleExpand()

	case keyMatches(keyStr, keyCollapse):
		return m.handleCollapse()
	}

	return m, nil
}

// handleEnter processes the Enter key based on the item at the current cursor
// position. Section headers toggle expansion. Expandable projects toggle
// subtree visibility. Leaf projects, filters, and labels emit selection
// messages for the root model to consume.
func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	if len(m.items) == 0 || m.cursor >= len(m.items) {
		return m, nil
	}

	item := m.items[m.cursor]
	switch item.Kind {
	case "section-header":
		key := sectionKey(item.Name)
		if key != "" {
			m.expandedSections[key] = !m.expandedSections[key]
			if err := (&m).loadItems(); err != nil {
				return m, nil
			}
			(&m).clampCursor()
		}

	case "project":
		if item.Expandable {
			m.expandedProjects[item.ID] = !m.expandedProjects[item.ID]
			if err := (&m).loadItems(); err != nil {
				return m, nil
			}
			(&m).clampCursor()
		} else {
			return m, emitProjectSelected(item.ID)
		}

	case "filter":
		query := ""
		if f, err := m.store.GetFilter(item.ID); err == nil && f != nil {
			query = f.Query
		}
		return m, emitFilterSelected(item.ID, query)

	case "label":
		return m, emitLabelSelected(item.Name)
	}

	return m, nil
}

// handleExpand expands the item at the current cursor position if it is
// expandable and currently collapsed. After expansion, items are reloaded
// and the cursor is clamped to the new bounds.
func (m Model) handleExpand() (tea.Model, tea.Cmd) {
	if len(m.items) == 0 || m.cursor >= len(m.items) {
		return m, nil
	}

	item := m.items[m.cursor]
	switch item.Kind {
	case "section-header":
		key := sectionKey(item.Name)
		if key != "" && !item.Expanded {
			m.expandedSections[key] = true
			if err := (&m).loadItems(); err != nil {
				return m, nil
			}
			(&m).clampCursor()
		}

	case "project":
		if item.Expandable && !item.Expanded {
			m.expandedProjects[item.ID] = true
			if err := (&m).loadItems(); err != nil {
				return m, nil
			}
			(&m).clampCursor()
		}
	}

	return m, nil
}

// handleCollapse collapses the item at the current cursor position if it is
// expandable and currently expanded. After collapsing, items are reloaded
// and the cursor is clamped to the new bounds.
func (m Model) handleCollapse() (tea.Model, tea.Cmd) {
	if len(m.items) == 0 || m.cursor >= len(m.items) {
		return m, nil
	}

	item := m.items[m.cursor]
	switch item.Kind {
	case "section-header":
		key := sectionKey(item.Name)
		if key != "" && item.Expanded {
			m.expandedSections[key] = false
			if err := (&m).loadItems(); err != nil {
				return m, nil
			}
			(&m).clampCursor()
		}

	case "project":
		if item.Expandable && item.Expanded {
			m.expandedProjects[item.ID] = false
			if err := (&m).loadItems(); err != nil {
				return m, nil
			}
			(&m).clampCursor()
		}
	}

	return m, nil
}

// sectionKey extracts the section map key from a section header name.
// Section header names have the format "Projects (N)", "Filters (N)",
// or "Labels (N)". Returns "" if the name does not match a known section.
func sectionKey(name string) string {
	switch {
	case strings.HasPrefix(name, "Projects"):
		return "projects"
	case strings.HasPrefix(name, "Filters"):
		return "filters"
	case strings.HasPrefix(name, "Labels"):
		return "labels"
	default:
		return ""
	}
}

// clampCursor ensures the cursor position is within the valid range for the
// current items slice. If the items list is empty, cursor is set to 0.
func (m *Model) clampCursor() {
	if len(m.items) == 0 {
		m.cursor = 0
		return
	}
	if m.cursor >= len(m.items) {
		m.cursor = len(m.items) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

// keyMatches compares a Bubbletea key event string against a configured key
// binding. For single-character bindings, comparison is case-sensitive
// (e.g., "g" ≠ "G"). For multi-character bindings (e.g., "Enter", "Esc"),
// comparison is case-insensitive because Bubbletea reports key names in
// lowercase while the default config uses title case.
func keyMatches(eventKey, binding string) bool {
	if eventKey == binding {
		return true
	}
	if len(binding) <= 1 {
		return false
	}
	return strings.EqualFold(eventKey, binding)
}

// emitProjectSelected returns a command that emits a [msg.ProjectSelectedMsg]
// with the given project ID.
func emitProjectSelected(id string) tea.Cmd {
	return func() tea.Msg {
		return msg.ProjectSelectedMsg{ID: id}
	}
}

// emitFilterSelected returns a command that emits a [msg.FilterSelectedMsg]
// with the given filter ID and query string.
func emitFilterSelected(id, query string) tea.Cmd {
	return func() tea.Msg {
		return msg.FilterSelectedMsg{ID: id, Query: query}
	}
}

// emitLabelSelected returns a command that emits a [msg.LabelSelectedMsg]
// with the given label name.
func emitLabelSelected(name string) tea.Cmd {
	return func() tea.Msg {
		return msg.LabelSelectedMsg{Name: name}
	}
}
