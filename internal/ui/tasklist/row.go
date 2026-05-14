package tasklist

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"

	"github.com/demeyerthom/todoist-tui/internal/config"
	"github.com/demeyerthom/todoist-tui/internal/model"
	"github.com/demeyerthom/todoist-tui/internal/ui/theme"
)

// taskToRow converts a model.Task into a bubble-table Row with columns
// for priority, content, due date, and labels. Hidden metadata columns
// ("task_id" and "is_section") are included for selection tracking.
func taskToRow(task model.Task, styles *theme.Styles, cfg *config.Config) table.Row {
	// Priority dot.
	priorityColor := taskPriorityColor(task, cfg)
	priorityCell := table.NewStyledCell("●", lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(priorityColor)))

	// Due date.
	var dueCell any = ""
	dueStr, dueColor := styles.FormatDueDate(task.Due, time.Now())
	if dueColor != nil {
		dueCell = table.NewStyledCell(dueStr, lipgloss.NewStyle().Foreground(dueColor))
	}

	return table.NewRow(table.RowData{
		"priority":   priorityCell,
		"content":    task.Content,
		"due":        dueCell,
		"labels":     strings.Join(task.Labels, ", "),
		"task_id":    task.ID,
		"is_section": false,
	})
}

// taskPriorityColor returns the theme color string for a task's priority level.
// Priority 1 = high (red), 2 = medium (yellow), 3 = low (blue),
// 4 or unexpected = normal text color.
func taskPriorityColor(task model.Task, cfg *config.Config) string {
	switch task.Priority {
	case 1:
		return cfg.Theme.TaskPriorityHigh
	case 2:
		return cfg.Theme.TaskPriorityMedium
	case 3:
		return cfg.Theme.TaskPriorityLow
	default:
		return cfg.Theme.NormalText
	}
}

// sectionRow creates a section separator row for the bubble-table.
// The section name is placed in the content column as a styled separator
// line ("─── SectionName ───") and the entire row is styled with
// [theme.Styles.SectionSep]. Hidden metadata marks it as a section row.
func sectionRow(name string, styles *theme.Styles) table.Row {
	return table.NewRow(table.RowData{
		"priority":   "",
		"content":    fmt.Sprintf("─── %s ───", name),
		"due":        "",
		"labels":     "",
		"task_id":    "",
		"is_section": true,
	}).WithStyle(styles.SectionSep())
}

// buildRows converts a list of SectionGroups into a flat slice of
// bubble-table rows. Each section group produces a section separator
// row followed by rows for each task in the group.
func buildRows(groups []SectionGroup, styles *theme.Styles, cfg *config.Config) []table.Row {
	var rows []table.Row
	for _, g := range groups {
		name := "No Section"
		if g.Section != nil {
			name = g.Section.Name
		}
		rows = append(rows, sectionRow(name, styles))

		for _, t := range g.Tasks {
			rows = append(rows, taskToRow(t, styles, cfg))
		}
	}
	return rows
}

// defaultColumns returns the standard column definitions for the task list
// table. The priority column is narrow and fixed; content flexes to fill
// available space; due and labels columns have fixed widths.
func defaultColumns(width int) []table.Column {
	return []table.Column{
		table.NewColumn("priority", "", 3),
		table.NewFlexColumn("content", "Task", 1),
		table.NewColumn("due", "Due", 12),
		table.NewColumn("labels", "Labels", 20),
	}
}
