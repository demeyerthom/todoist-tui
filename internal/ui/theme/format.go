package theme

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/demeyerthom/todoist-tui/internal/model"
)

// FormatDueDate returns a human-readable relative date string and a
// theme-derived terminal color for the given DueDate.
//
// Rules:
//   - nil due → ("", nil)
//   - Same day as now → "Today" in TaskDueToday color
//   - Next day → "Tomorrow" in TaskDueToday color
//   - 2–7 days ahead → "In N days" in NormalText color
//   - Further future → "Jan 2" formatted date in NormalText color
//   - Past → "Overdue Nd" in TaskOverdue color
func (s *Styles) FormatDueDate(due *model.DueDate, now time.Time) (string, lipgloss.TerminalColor) {
	if due == nil || due.Date == "" {
		return "", nil
	}

	dueDate, err := time.Parse("2006-01-02", due.Date)
	if err != nil {
		return "", nil
	}

	loc := now.Location()
	nowDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	dueDay := time.Date(dueDate.Year(), dueDate.Month(), dueDate.Day(), 0, 0, 0, 0, loc)
	days := int(dueDay.Sub(nowDay).Hours() / 24)

	switch {
	case days == 0:
		return "Today", colorFromTheme(s.theme.TaskDueToday)
	case days == 1:
		return "Tomorrow", colorFromTheme(s.theme.TaskDueToday)
	case days > 1 && days <= 7:
		return fmt.Sprintf("In %d days", days), colorFromTheme(s.theme.NormalText)
	case days > 7:
		return dueDay.Format("Jan 2"), colorFromTheme(s.theme.NormalText)
	default:
		return fmt.Sprintf("Overdue %dd", -days), colorFromTheme(s.theme.TaskOverdue)
	}
}
