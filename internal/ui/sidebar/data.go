package sidebar

import (
	"fmt"
	"sort"

	"github.com/demeyerthom/todoist-tui/internal/model"
)

// loadItems reads projects, filters, and labels from the store and builds
// the flat sidebar item list. Sections and subtrees respect the current
// expansion state from expandedSections and expandedProjects.
func (m *Model) loadItems() error {
	// Load and filter projects.
	allProjects, err := m.store.ListProjects()
	if err != nil {
		return fmt.Errorf("sidebar: list projects: %w", err)
	}
	projects := make([]model.Project, 0, len(allProjects))
	for _, p := range allProjects {
		if p.IsDeleted || p.IsArchived {
			continue
		}
		projects = append(projects, p)
	}

	// Load and filter filters.
	allFilters, err := m.store.ListFilters()
	if err != nil {
		return fmt.Errorf("sidebar: list filters: %w", err)
	}
	filters := make([]model.Filter, 0, len(allFilters))
	for _, f := range allFilters {
		if !f.IsDeleted {
			filters = append(filters, f)
		}
	}

	// Load and filter labels.
	allLabels, err := m.store.ListLabels()
	if err != nil {
		return fmt.Errorf("sidebar: list labels: %w", err)
	}
	labels := make([]model.Label, 0, len(allLabels))
	for _, l := range allLabels {
		if !l.IsDeleted {
			labels = append(labels, l)
		}
	}

	m.items = nil

	// Projects section.
	m.items = append(m.items, SidebarItem{
		Kind:       "section-header",
		Name:       fmt.Sprintf("Projects (%d)", len(projects)),
		Expandable: true,
		Expanded:   m.expandedSections["projects"],
	})
	if m.expandedSections["projects"] {
		m.items = append(m.items, buildProjectTree(projects, m.expandedProjects)...)
	}

	// Filters section.
	m.items = append(m.items, SidebarItem{
		Kind:       "section-header",
		Name:       fmt.Sprintf("Filters (%d)", len(filters)),
		Expandable: true,
		Expanded:   m.expandedSections["filters"],
	})
	if m.expandedSections["filters"] {
		sort.Slice(filters, func(i, j int) bool { return filters[i].Name < filters[j].Name })
		for _, f := range filters {
			m.items = append(m.items, SidebarItem{
				Kind:  "filter",
				ID:    f.ID,
				Name:  f.Name,
				Color: f.Color,
			})
		}
	}

	// Labels section.
	m.items = append(m.items, SidebarItem{
		Kind:       "section-header",
		Name:       fmt.Sprintf("Labels (%d)", len(labels)),
		Expandable: true,
		Expanded:   m.expandedSections["labels"],
	})
	if m.expandedSections["labels"] {
		sort.Slice(labels, func(i, j int) bool { return labels[i].Name < labels[j].Name })
		for _, l := range labels {
			m.items = append(m.items, SidebarItem{
				Kind:  "label",
				ID:    l.ID,
				Name:  l.Name,
				Color: l.Color,
			})
		}
	}

	return nil
}

// buildProjectTree converts a flat project list into a flat, ordered list of
// SidebarItems with correct indentation, respecting the expansion state of
// parent projects. Inbox projects are pinned to the top. Orphaned projects
// (whose parent is not in the set) are treated as root-level.
func buildProjectTree(projects []model.Project, expandedProjects map[string]bool) []SidebarItem {
	// Build project lookup and classify into roots and children.
	projectMap := make(map[string]model.Project, len(projects))
	for _, p := range projects {
		projectMap[p.ID] = p
	}

	children := make(map[string][]model.Project)
	var roots []model.Project

	for _, p := range projects {
		if p.ParentID == "" {
			roots = append(roots, p)
		} else if _, ok := projectMap[p.ParentID]; ok {
			children[p.ParentID] = append(children[p.ParentID], p)
		} else {
			// Orphan: parent not in project set; treat as root.
			roots = append(roots, p)
		}
	}

	// Sort each sibling group alphabetically.
	for id := range children {
		sort.Slice(children[id], func(i, j int) bool {
			return children[id][i].Name < children[id][j].Name
		})
	}

	// Sort roots: Inbox first, then alphabetical by name.
	sort.Slice(roots, func(i, j int) bool {
		if roots[i].IsInbox != roots[j].IsInbox {
			return roots[i].IsInbox
		}
		return roots[i].Name < roots[j].Name
	})

	var items []SidebarItem
	for _, p := range roots {
		items = append(items, addProjectSubtree(p, 0, children, expandedProjects)...)
	}
	return items
}

// addProjectSubtree recursively adds a project and its visible children to the
// flat item list. Children are only shown when the parent is expanded.
func addProjectSubtree(p model.Project, indent int, children map[string][]model.Project, expandedProjects map[string]bool) []SidebarItem {
	hasChildren := len(children[p.ID]) > 0
	expanded := expandedProjects[p.ID]

	item := SidebarItem{
		Kind:       "project",
		ID:         p.ID,
		Name:       p.Name,
		Color:      p.Color,
		Indent:     indent,
		IsInbox:    p.IsInbox,
		Expandable: hasChildren,
		Expanded:   expanded,
	}
	items := []SidebarItem{item}

	if !hasChildren || !expanded {
		return items
	}

	for _, child := range children[p.ID] {
		items = append(items, addProjectSubtree(child, indent+1, children, expandedProjects)...)
	}
	return items
}
