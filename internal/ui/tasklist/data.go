package tasklist

import (
	"fmt"
	"slices"
	"sort"

	"github.com/demeyerthom/todoist-tui/internal/model"
)

// loadTasks reads all tasks and sections from the store, applies the current
// filter, groups tasks by section, and sorts them. The result is stored in
// m.groups. The bubble-table rows are not rebuilt — that is handled separately.
func (m *Model) loadTasks() error {
	tasks, err := m.store.ListTasks()
	if err != nil {
		return fmt.Errorf("tasklist: load tasks: %w", err)
	}

	sections, err := m.store.ListSections()
	if err != nil {
		return fmt.Errorf("tasklist: load sections: %w", err)
	}

	// Todoist filters are unsupported in M1; clear groups and return.
	if m.currentFilter.IsFilter {
		m.groups = nil
		return nil
	}

	// Build section lookup, excluding deleted and archived sections.
	// When filtering by project, also exclude sections from other projects.
	sectionByID := make(map[string]model.Section, len(sections))
	for _, sec := range sections {
		if sec.IsDeleted || sec.IsArchived {
			continue
		}
		if m.currentFilter.ProjectID != "" && sec.ProjectID != m.currentFilter.ProjectID {
			continue
		}
		sectionByID[sec.ID] = sec
	}

	// Filter tasks.
	filtered := tasks[:0]
	for _, t := range tasks {
		if t.IsDeleted {
			continue
		}
		if !m.showCompleted && t.Checked {
			continue
		}
		if m.currentFilter.ProjectID != "" && t.ProjectID != m.currentFilter.ProjectID {
			continue
		}
		if m.currentFilter.LabelName != "" && !slices.Contains(t.Labels, m.currentFilter.LabelName) {
			continue
		}
		filtered = append(filtered, t)
	}

	// Group tasks by section ID. Tasks whose SectionID does not match any
	// known section are grouped under the empty string (unsectioned).
	groupMap := make(map[string][]model.Task)
	for _, t := range filtered {
		key := t.SectionID
		if _, ok := sectionByID[key]; !ok {
			key = ""
		}
		groupMap[key] = append(groupMap[key], t)
	}

	// Sort tasks within each group by priority (P1 first) then due date
	// (earliest first, no due date last).
	for _, ts := range groupMap {
		sort.Slice(ts, func(i, j int) bool {
			a, b := ts[i], ts[j]
			if a.Priority != b.Priority {
				return a.Priority < b.Priority
			}
			return dueDateLess(a.Due, b.Due)
		})
	}

	// Build SectionGroup slice. Only sections that have at least one task
	// are included. The unsectioned group (key "") is added last.
	var groups []SectionGroup
	for secID, sec := range sectionByID {
		ts, ok := groupMap[secID]
		if !ok {
			continue
		}
		secCopy := sec
		groups = append(groups, SectionGroup{
			Section: &secCopy,
			Tasks:   ts,
		})
	}

	// Sort section groups by SectionOrder ascending.
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Section.SectionOrder < groups[j].Section.SectionOrder
	})

	// Unsectioned group always goes last.
	if ts, ok := groupMap[""]; ok {
		groups = append(groups, SectionGroup{
			Section: nil,
			Tasks:   ts,
		})
	}

	m.groups = groups
	return nil
}

// dueDateLess returns true if a should sort before b. Tasks with an earlier
// due date come first; tasks with no due date (nil) sort last.
func dueDateLess(a, b *model.DueDate) bool {
	if a == nil && b == nil {
		return false
	}
	if a == nil {
		return false
	}
	if b == nil {
		return true
	}
	dateA := a.Date
	if a.Datetime != "" {
		dateA = a.Datetime
	}
	dateB := b.Date
	if b.Datetime != "" {
		dateB = b.Datetime
	}
	return dateA < dateB
}
