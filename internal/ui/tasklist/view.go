package tasklist

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

// View renders the task list panel using a bubble-table.
//
// Returns an empty string when width or height are ≤ 0 (no space to render).
//
// When there are no tasks to display, a contextual empty-state message is
// centered based on the active filter:
//   - Project filter: "No tasks in project"
//   - Label filter:   "No tasks with label"
//   - Todoist filter: "No tasks match filter"
//   - Default:        "No tasks"
//
// When tasks are present, a bubble-table is built from the current groups
// using [buildRows] and [defaultColumns]. The row matching [Model.selectedID]
// is highlighted. Section separator rows are never highlighted.
func (m Model) View() string {
	if m.width <= 0 || m.height <= 0 {
		return ""
	}

	if len(m.groups) == 0 || allGroupsEmpty(m.groups) {
		return m.renderEmptyState()
	}

	return m.renderTable()
}

// renderEmptyState returns a centered, muted message string appropriate for
// the current filter context.
func (m Model) renderEmptyState() string {
	msg := emptyStateMessage(m.currentFilter)
	style := m.styles.MutedText().Align(lipgloss.Center, lipgloss.Center)
	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		style.Render(msg),
	)
}

// renderTable constructs a bubble-table from current groups, highlights the
// row matching m.selectedID, and returns the rendered string.
func (m Model) renderTable() string {
	rows := buildRows(m.groups, m.styles, m.cfg)
	columns := defaultColumns(m.width)

	t := table.New(columns).
		WithRows(rows).
		WithTargetWidth(m.width).
		WithBaseStyle(lipgloss.NewStyle()).
		HighlightStyle(lipgloss.NewStyle().
			Background(lipgloss.Color(m.cfg.Theme.SelectedRow))).
		Focused(true).
		WithHeaderVisibility(false).
		WithNoPagination()

	highlightIdx := findHighlightedRow(rows, m.selectedID)
	t = t.WithHighlightedRow(highlightIdx)

	return t.View()
}

// findHighlightedRow returns the index of the row whose "task_id" metadata
// matches selectedID. Section separator rows are skipped. Returns 0 when no
// match is found.
func findHighlightedRow(rows []table.Row, selectedID string) int {
	if selectedID == "" {
		return 0
	}

	for i, row := range rows {
		isSection, _ := row.Data["is_section"].(bool)
		if isSection {
			continue
		}

		taskID, _ := row.Data["task_id"].(string)
		if taskID == selectedID {
			return i
		}
	}

	return 0
}

// allGroupsEmpty returns true when every group in the slice has zero tasks.
func allGroupsEmpty(groups []SectionGroup) bool {
	for _, g := range groups {
		if len(g.Tasks) > 0 {
			return false
		}
	}
	return true
}

// emptyStateMessage returns a contextual empty-state message based on the
// active filter dimensions.
func emptyStateMessage(f Filter) string {
	switch {
	case f.ProjectID != "":
		return "No tasks in project"
	case f.LabelName != "":
		return "No tasks with label"
	case f.IsFilter:
		return "No tasks match filter"
	default:
		return "No tasks"
	}
}
