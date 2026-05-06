package ui

import (
	"context"

	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/demeyerthom/todoist-tui/internal/config"
	"github.com/demeyerthom/todoist-tui/internal/store"
	"github.com/demeyerthom/todoist-tui/internal/sync"
	"github.com/demeyerthom/todoist-tui/internal/ui/keymap"
)

// Panel identifies the currently focused panel in the 3-panel layout.
type Panel int

const (
	PanelSidebar Panel = iota
	PanelMain
	PanelDetail
)

const panelCount = 3

// SyncDoneMsg is emitted when the startup full sync completes successfully.
type SyncDoneMsg struct{}

// SyncErrMsg is emitted when the startup full sync fails.
type SyncErrMsg struct{ Err error }

// Model is the root Bubbletea application model.
// It holds application-wide state and delegates to sub-models for each panel.
type Model struct {
	cfg         *config.Config
	store       *store.Store
	syncClient  *sync.Client
	activePanel Panel
	mode        keymap.Mode
	width       int
	height      int
	err         error
}

// NewModel creates a new root application model with the given dependencies.
// The model starts in Sidebar panel, Normal mode, with zero width and height.
func NewModel(cfg *config.Config, store *store.Store, syncClient *sync.Client) Model {
	return Model{
		cfg:         cfg,
		store:       store,
		syncClient:  syncClient,
		activePanel: PanelSidebar,
		mode:        keymap.ModeNormal,
		width:       0,
		height:      0,
	}
}

// Init implements tea.Model. It kicks off a full sync against the Todoist
// Sync API asynchronously. The result is delivered as SyncDoneMsg or SyncErrMsg.
func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		err := m.syncClient.FullSync(context.Background(), m.store)
		if err != nil {
			return SyncErrMsg{Err: err}
		}
		return SyncDoneMsg{}
	}
}

// Update implements tea.Model. It handles incoming messages and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.activePanel = Panel((int(m.activePanel) + 1) % panelCount)
		case "shift+tab":
			m.activePanel = Panel((int(m.activePanel) - 1 + panelCount) % panelCount)
		}

	case SyncDoneMsg:
		// Full sync completed; store is already populated.

	case SyncErrMsg:
		m.err = msg.Err
	}
	return m, nil
}

// View implements tea.Model. It renders the 3-panel layout:
// sidebar (~20%), main (~50%), detail (~30%) joined horizontally.
// Any error present in m.err is displayed as a banner at the top.
func (m Model) View() string {
	if m.width <= 0 || m.height <= 0 {
		if m.err != nil {
			return "Error: " + m.err.Error()
		}
		return "Initializing..."
	}

	// Render error banner if an error is set
	content := m.panelsView()
	if m.err != nil {
		errBanner := lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.cfg.Theme.Error)).
			Width(m.width).
			Render("Error: " + m.err.Error())
		content = lipgloss.JoinVertical(lipgloss.Top, errBanner, content)
	}
	return content
}

// panelsView renders the 3-panel layout without the error banner.
func (m Model) panelsView() string {
	// Calculate panel widths (~20/50/30 split)
	sidebarWidth := max(1, m.width*20/100)
	mainWidth := max(1, m.width*50/100)
	detailWidth := max(1, m.width-sidebarWidth-mainWidth)

	// Panel styles with rounded borders and focus-aware border color
	sidebarStyle := lipgloss.NewStyle().
		Width(sidebarWidth).
		Height(m.height).
		Border(lipgloss.RoundedBorder()).
		BorderRight(true).
		BorderForeground(m.borderColor(PanelSidebar))

	mainStyle := lipgloss.NewStyle().
		Width(mainWidth).
		Height(m.height).
		Border(lipgloss.RoundedBorder()).
		BorderRight(true).
		BorderForeground(m.borderColor(PanelMain))

	detailStyle := lipgloss.NewStyle().
		Width(detailWidth).
		Height(m.height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.borderColor(PanelDetail))

	// Inner dimensions accounting for borders
	sidebarInnerW := max(1, sidebarWidth-2)
	mainInnerW := max(1, mainWidth-2)
	detailInnerW := max(1, detailWidth-2)
	innerH := max(1, m.height-2)

	// Render each panel with styled placeholder content
	sidebar := sidebarStyle.Render(m.sidebarView(sidebarInnerW, innerH))
	main := mainStyle.Render(m.mainView(mainInnerW, innerH))
	detail := detailStyle.Render(m.detailView(detailInnerW, innerH))

	// Join panels horizontally, aligned at top
	return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, main, detail)
}

// borderColor returns the terminal color for a panel's border based on whether
// it is the currently active (focused) panel.
func (m Model) borderColor(p Panel) lipgloss.TerminalColor {
	if m.activePanel == p {
		return lipgloss.Color(m.cfg.Theme.ActiveBorder)
	}
	return lipgloss.Color(m.cfg.Theme.InactiveBorder)
}

// sidebarView renders the sidebar panel's placeholder content.
// Headings use the Header theme color for visual distinction.
func (m Model) sidebarView(width, height int) string {
	heading := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Foreground(lipgloss.Color(m.cfg.Theme.Header))

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		heading.Render("Projects"),
		heading.Render("Filters"),
		heading.Render("Labels"),
	)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

// mainView renders the main panel's placeholder content.
// MutedText is used to indicate the panel is awaiting input.
func (m Model) mainView(width, height int) string {
	text := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Foreground(lipgloss.Color(m.cfg.Theme.MutedText)).
		Render("Select a project or filter")

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, text)
}

// detailView renders the detail panel's placeholder content.
// MutedText is used to indicate no task is currently selected.
func (m Model) detailView(width, height int) string {
	text := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Foreground(lipgloss.Color(m.cfg.Theme.MutedText)).
		Render("No task selected")

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, text)
}
