package tasklist

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/demeyerthom/todoist-tui/internal/ui/msg"
)

// Update implements tea.Model. It handles keyboard navigation, filter
// selection messages, sync completion, completed-task toggling, and
// window resize events. Key bindings are read from the Normal mode
// keymap with vim-style defaults.
func (m Model) Update(teaMsg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := teaMsg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

	case msg.ProjectSelectedMsg:
		m.currentFilter = Filter{ProjectID: msg.ID}
		_ = (&m).loadTasks()
		m.selectedID = firstTaskID(m.groups)
		return m, emitTaskSelected(m.selectedID)

	case msg.LabelSelectedMsg:
		m.currentFilter = Filter{LabelName: msg.Name}
		_ = (&m).loadTasks()
		m.selectedID = firstTaskID(m.groups)
		return m, emitTaskSelected(m.selectedID)

	case msg.FilterSelectedMsg:
		m.currentFilter = Filter{FilterQuery: msg.Query, IsFilter: true}
		_ = (&m).loadTasks()
		m.selectedID = ""
		return m, emitTaskSelected(m.selectedID)

	case msg.SyncCompleteMsg:
		old := m.selectedID
		_ = (&m).loadTasks()
		if !taskIDExists(m.groups, m.selectedID) {
			m.selectedID = ""
		}
		if m.selectedID != old {
			return m, emitTaskSelected(m.selectedID)
		}
		return m, nil

	case msg.ToggleCompletedMsg:
		old := m.selectedID
		m.showCompleted = !m.showCompleted
		_ = (&m).loadTasks()
		if !taskIDExists(m.groups, m.selectedID) {
			m.selectedID = ""
		}
		if m.selectedID != old {
			return m, emitTaskSelected(m.selectedID)
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

// handleKeyMsg processes a Bubbletea key event, dispatching to the
// appropriate handler based on the configured Normal mode key bindings.
func (m Model) handleKeyMsg(keyMsg tea.KeyMsg) (tea.Model, tea.Cmd) {
	keyStr := keyMsg.String()

	keyUp := m.cfg.Keymap.Normal["up"]
	if keyUp == "" {
		keyUp = "k"
	}
	keyDown := m.cfg.Keymap.Normal["down"]
	if keyDown == "" {
		keyDown = "j"
	}
	keyGoTop := m.cfg.Keymap.Normal["go_top"]
	if keyGoTop == "" {
		keyGoTop = "g"
	}
	keyGoBottom := m.cfg.Keymap.Normal["go_bottom"]
	if keyGoBottom == "" {
		keyGoBottom = "G"
	}
	keyToggle := m.cfg.Keymap.Normal["toggle_completed"]
	if keyToggle == "" {
		keyToggle = "H"
	}

	switch {
	case keyMatches(keyStr, keyDown):
		return m.navigate(1)

	case keyMatches(keyStr, keyUp):
		return m.navigate(-1)

	case keyMatches(keyStr, keyGoTop):
		return m.goToBoundary(0, 1)

	case keyMatches(keyStr, keyGoBottom):
		return m.goToBoundary(len(buildRows(m.groups, m.styles, m.cfg))-1, -1)

	case keyMatches(keyStr, keyToggle):
		old := m.selectedID
		m.showCompleted = !m.showCompleted
		_ = (&m).loadTasks()
		if !taskIDExists(m.groups, m.selectedID) {
			m.selectedID = ""
		}
		if m.selectedID != old {
			return m, emitTaskSelected(m.selectedID)
		}
		return m, nil
	}

	return m, nil
}

// navigate moves the selected row by direction rows (1 for down, -1 for up),
// skipping section separator rows. The selection does not wrap around.
func (m Model) navigate(direction int) (tea.Model, tea.Cmd) {
	rows := buildRows(m.groups, m.styles, m.cfg)
	currentIdx := findHighlightedRow(rows, m.selectedID)

	for i := currentIdx + direction; i >= 0 && i < len(rows); i += direction {
		isSection, _ := rows[i].Data["is_section"].(bool)
		if isSection {
			continue
		}
		taskID, _ := rows[i].Data["task_id"].(string)
		if taskID != "" && taskID != m.selectedID {
			m.selectedID = taskID
			return m, emitTaskSelected(m.selectedID)
		}
	}

	return m, nil
}

// goToBoundary moves the selection to the first non-section task row starting
// at start and stepping by step (1 for forward, -1 for backward).
func (m Model) goToBoundary(start, step int) (tea.Model, tea.Cmd) {
	if start < 0 {
		return m, nil
	}

	rows := buildRows(m.groups, m.styles, m.cfg)
	for i := start; i >= 0 && i < len(rows); i += step {
		isSection, _ := rows[i].Data["is_section"].(bool)
		if isSection {
			continue
		}
		taskID, _ := rows[i].Data["task_id"].(string)
		if taskID != "" {
			if taskID != m.selectedID {
				m.selectedID = taskID
				return m, emitTaskSelected(m.selectedID)
			}
			return m, nil
		}
	}

	return m, nil
}

// keyMatches compares a Bubbletea key event string against a configured key
// binding. For single-character bindings, comparison is case-sensitive
// (e.g., "g" ≠ "G"). For multi-character bindings (e.g., "Enter", "Esc"),
// comparison is case-insensitive.
func keyMatches(eventKey, binding string) bool {
	if eventKey == binding {
		return true
	}
	if len(binding) <= 1 {
		return false
	}
	return strings.EqualFold(eventKey, binding)
}

// emitTaskSelected returns a command that emits a [msg.TaskSelectedMsg]
// with the given task ID.
func emitTaskSelected(id string) tea.Cmd {
	return func() tea.Msg {
		return msg.TaskSelectedMsg{ID: id}
	}
}

// firstTaskID returns the ID of the first task in the first group that has
// tasks. Returns "" if no tasks exist in any group.
func firstTaskID(groups []SectionGroup) string {
	for _, g := range groups {
		if len(g.Tasks) > 0 {
			return g.Tasks[0].ID
		}
	}
	return ""
}

// taskIDExists returns true if any task in any group has the given ID.
func taskIDExists(groups []SectionGroup, id string) bool {
	if id == "" {
		return false
	}
	for _, g := range groups {
		for _, t := range g.Tasks {
			if t.ID == id {
				return true
			}
		}
	}
	return false
}
