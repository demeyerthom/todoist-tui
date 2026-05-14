package tasklist

import (
	"github.com/evertras/bubble-table/table"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/demeyerthom/todoist-tui/internal/config"
	"github.com/demeyerthom/todoist-tui/internal/model"
	"github.com/demeyerthom/todoist-tui/internal/store"
	"github.com/demeyerthom/todoist-tui/internal/ui/theme"
)

// Filter represents the current filtering context applied to the task list.
// Only one filter dimension is active at a time; others are zero-value.
type Filter struct {
	ProjectID   string // filter by project ID (empty = not filtering by project)
	LabelName   string // filter by label name (empty = not filtering by label)
	FilterQuery string // raw Todoist filter query (empty = not a filter)
	IsFilter    bool   // true if this is a Todoist filter (unsupported in M1)
}

// SectionGroup represents a section header and its associated tasks.
// When Section is nil the group represents unsectioned tasks.
type SectionGroup struct {
	Section *model.Section // nil for unsectioned tasks
	Tasks   []model.Task   // tasks in this group
}

// Model is the main task list panel Bubbletea model.
// It renders tasks grouped by section with compact/expanded display modes.
type Model struct {
	store         *store.Store
	cfg           *config.Config
	styles        *theme.Styles
	width         int
	height        int
	selectedID    string
	showCompleted bool
	currentFilter Filter
	groups        []SectionGroup
	table         table.Model
}

// NewModel creates a task list Model with default state and no loaded data.
// The bubble-table instance is initialized with empty columns; column
// configuration is applied in a later initialisation step.
func NewModel(cfg *config.Config, store *store.Store, styles *theme.Styles) Model {
	return Model{
		store:         store,
		cfg:           cfg,
		styles:        styles,
		width:         0,
		height:        0,
		selectedID:    "",
		showCompleted: false,
		currentFilter: Filter{},
		groups:        nil,
		table:         table.New([]table.Column{}),
	}
}

// Init satisfies the tea.Model interface. Returns nil because the task list
// does not run any initial commands; data is loaded in response to messages.
func (m Model) Init() tea.Cmd {
	return nil
}

