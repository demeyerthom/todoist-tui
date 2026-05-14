package sidebar

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/demeyerthom/todoist-tui/internal/config"
	"github.com/demeyerthom/todoist-tui/internal/store"
	"github.com/demeyerthom/todoist-tui/internal/ui/theme"
)

// SidebarItem represents a single renderable row in the sidebar panel.
// Kind determines how the item is styled and how it responds to user input.
type SidebarItem struct {
	// Kind identifies the item type: "project", "filter", "label",
	// "section-header", or "subtree-header".
	Kind string
	// ID is the project, filter, or label ID. Empty for header items.
	ID string
	// Name is the display text for this item.
	Name string
	// Color is the Todoist color name, used for projects and labels.
	Color string
	// Indent is the indentation level. 0 for top-level items, 1+ for nested.
	Indent int
	// IsInbox is true only for the Inbox project.
	IsInbox bool
	// Expandable is true for section headers and subtree headers.
	Expandable bool
	// Expanded is the current expansion state of this item.
	Expanded bool
}

// Model is the sidebar panel Bubbletea model. It manages the list of
// projects, filters, and labels, tracks cursor position, and handles
// section/project expansion and collapse.
type Model struct {
	store  *store.Store
	cfg    *config.Config
	styles *theme.Styles

	width  int
	height int

	// cursor is the current position in the flat items slice.
	cursor int

	// expandedSections tracks which sections are expanded: "projects",
	// "filters", "labels".
	expandedSections map[string]bool
	// expandedProjects tracks subtree expansion state by project ID.
	expandedProjects map[string]bool

	// items is the flat, ordered list of renderable sidebar rows.
	// Rebuilt when data is loaded or when expand/collapse state changes.
	items []SidebarItem
}

// NewModel creates a new sidebar model with the given dependencies.
// All maps are initialized, cursor is 0, and panel dimensions are 0.
func NewModel(cfg *config.Config, store *store.Store, styles *theme.Styles) Model {
	return Model{
		cfg:              cfg,
		store:            store,
		styles:           styles,
		cursor:           0,
		width:            0,
		height:           0,
		expandedSections: make(map[string]bool),
		expandedProjects: make(map[string]bool),
	}
}

// Init implements tea.Model. No initial command is needed.
func (m Model) Init() tea.Cmd {
	return nil
}


