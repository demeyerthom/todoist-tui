package detail

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/demeyerthom/todoist-tui/internal/model"
	"github.com/demeyerthom/todoist-tui/internal/ui/msg"
)

// Update handles Bubbletea messages for the detail panel.
func (m Model) Update(teaMsg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := teaMsg.(type) {
	case msg.TaskSelectedMsg:
		m.taskID = msg.ID
	case msg.SyncCompleteMsg:
		// Task data may have changed; the next View() will re-read from store.
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

// View renders the detail panel. When no task is selected, it shows a placeholder.
// When a task is selected, it renders key-value pairs for the task's fields.
func (m Model) View() string {
	if m.taskID == "" {
		return m.styles.MutedText().
			Width(m.width).
			Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).
			Render("Select a task to view details")
	}

	task, err := m.store.GetTask(m.taskID)
	if err != nil || task == nil {
		return m.styles.MutedText().
			Width(m.width).
			Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).
			Render("Task not found")
	}

	var b strings.Builder
	b.WriteString(m.styles.Header().Render(task.Content))
	b.WriteString("\n")

	m.writeField(&b, "Description", m.descriptionValue(task))
	m.writeField(&b, "Project", m.projectName(task))
	m.writeField(&b, "Section", m.sectionName(task))
	m.writeField(&b, "Priority", m.priorityValue(task))
	m.writeField(&b, "Due", m.dueValue(task))
	m.writeField(&b, "Labels", m.labelsValue(task))
	m.writeField(&b, "Subtasks", m.subtaskCount(task))

	return b.String()
}

// writeField writes a "Key: Value" line. The key is rendered with MutedText style.
func (m Model) writeField(b *strings.Builder, key, value string) {
	keyStyle := m.styles.MutedText().Width(12)
	b.WriteString(keyStyle.Render(key + ":"))
	b.WriteString(" ")
	b.WriteString(value)
	b.WriteString("\n")
}

// descriptionValue returns the task description or a muted "none" placeholder.
func (m Model) descriptionValue(task *model.Task) string {
	if task.Description == "" {
		return m.styles.MutedText().Render("none")
	}
	return task.Description
}

// projectName returns the project name for the task, or "Inbox" if not found.
func (m Model) projectName(task *model.Task) string {
	if task.ProjectID == "" {
		return "Inbox"
	}
	p, err := m.store.GetProject(task.ProjectID)
	if err != nil || p == nil {
		return "Inbox"
	}
	return p.Name
}

// sectionName returns the section name for the task, or a muted "none" if empty.
func (m Model) sectionName(task *model.Task) string {
	if task.SectionID == "" {
		return m.styles.MutedText().Render("none")
	}
	s, err := m.store.GetSection(task.SectionID)
	if err != nil || s == nil {
		return m.styles.MutedText().Render("none")
	}
	return s.Name
}

// priorityValue returns the priority label with theme coloring.
func (m Model) priorityValue(task *model.Task) string {
	var color string
	var label string
	switch task.Priority {
	case 1:
		label = "P1"
		color = m.cfg.Theme.TaskPriorityHigh
	case 2:
		label = "P2"
		color = m.cfg.Theme.TaskPriorityMedium
	case 3:
		label = "P3"
		color = m.cfg.Theme.TaskPriorityLow
	default:
		label = "P4"
		color = m.cfg.Theme.NormalText
	}
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(color)).Render(label)
}

// dueValue returns the formatted due date string, or a muted "none" if no due date.
func (m Model) dueValue(task *model.Task) string {
	text, color := m.styles.FormatDueDate(task.Due, time.Now())
	if text == "" {
		return m.styles.MutedText().Render("none")
	}
	return lipgloss.NewStyle().Foreground(color).Render(text)
}

// labelsValue returns comma-separated label names, or a muted "none" if empty.
func (m Model) labelsValue(task *model.Task) string {
	if len(task.Labels) == 0 {
		return m.styles.MutedText().Render("none")
	}
	var names []string
	for _, id := range task.Labels {
		l, err := m.store.GetLabel(id)
		if err != nil || l == nil {
			names = append(names, id)
			continue
		}
		names = append(names, l.Name)
	}
	return strings.Join(names, ", ")
}

// subtaskCount returns the count of subtasks for the current task.
func (m Model) subtaskCount(task *model.Task) string {
	tasks, err := m.store.ListTasks()
	if err != nil {
		return "0"
	}
	count := 0
	for _, t := range tasks {
		if t.ParentID == task.ID && !t.IsDeleted {
			count++
		}
	}
	return fmt.Sprintf("%d", count)
}
