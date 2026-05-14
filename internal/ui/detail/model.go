package detail

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/demeyerthom/todoist-tui/internal/config"
	"github.com/demeyerthom/todoist-tui/internal/store"
	"github.com/demeyerthom/todoist-tui/internal/ui/theme"
)

// Model is the detail panel Bubbletea model. It displays the full details of
// the currently selected task.
type Model struct {
	store  *store.Store
	cfg    *config.Config
	styles *theme.Styles

	width  int
	height int

	// taskID is the ID of the currently selected task. Empty when no task is selected.
	taskID string
}

// NewModel creates a detail panel Model with the given dependencies.
func NewModel(cfg *config.Config, store *store.Store, styles *theme.Styles) Model {
	return Model{
		cfg:    cfg,
		store:  store,
		styles: styles,
		width:  0,
		height: 0,
		taskID: "",
	}
}

// Init implements tea.Model. No initial command is needed.
func (m Model) Init() tea.Cmd {
	return nil
}
