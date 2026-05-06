package ui

import (
	"context"
	"fmt"

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
	commandBuf  string
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
		commandBuf:  "",
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
		switch m.mode {
		case keymap.ModeNormal:
			switch msg.String() {
			case ":":
				m.mode = keymap.ModeCommand
				m.commandBuf = ""
			case "ctrl+c":
				return m, tea.Quit
			case "tab":
				m.activePanel = Panel((int(m.activePanel) + 1) % panelCount)
			case "shift+tab":
				m.activePanel = Panel((int(m.activePanel) - 1 + panelCount) % panelCount)
			}

		case keymap.ModeInsert:
			if msg.String() == "ctrl+c" {
				return m, tea.Quit
			}

		case keymap.ModeCommand:
			switch msg.String() {
			case "enter":
				switch m.commandBuf {
				case "q", "quit":
					return m, tea.Quit
				default:
					m.err = fmt.Errorf("unknown command: %s", m.commandBuf)
					m.mode = keymap.ModeNormal
					m.commandBuf = ""
				}
			case "esc":
				m.commandBuf = ""
				m.mode = keymap.ModeNormal
			case "backspace":
				if len(m.commandBuf) > 0 {
					m.commandBuf = m.commandBuf[:len(m.commandBuf)-1]
				} else {
					m.commandBuf = ""
					m.mode = keymap.ModeNormal
				}
			case "ctrl+c":
				m.commandBuf = ""
				m.mode = keymap.ModeNormal
			default:
				if len(msg.Runes) > 0 {
					m.commandBuf += string(msg.Runes)
				}
			}
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
// When in command mode, a command bar is displayed at the bottom.
func (m Model) View() string {
	if m.width <= 0 || m.height <= 0 {
		if m.err != nil {
			return "Error: " + m.err.Error()
		}
		return "Initializing..."
	}

	// Compute effective panel height, accounting for error banner and command bar
	panelHeight := m.height
	if m.err != nil {
		panelHeight--
	}
	if m.mode == keymap.ModeCommand {
		panelHeight--
	}

	content := m.panelsView(panelHeight)
	if m.err != nil {
		errBanner := lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.cfg.Theme.Error)).
			Width(m.width).
			Render("Error: " + m.err.Error())
		content = lipgloss.JoinVertical(lipgloss.Top, errBanner, content)
	}
	if m.mode == keymap.ModeCommand {
		cmdBar := lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.cfg.Theme.CommandBar)).
			Width(m.width).
			Render(":" + m.commandBuf)
		content = lipgloss.JoinVertical(lipgloss.Top, content, cmdBar)
	}
	return content
}

// panelsView renders the 3-panel layout without the error banner or command bar.
func (m Model) panelsView(panelHeight int) string {
	// Calculate panel widths (~20/50/30 split)
	sidebarWidth := max(1, m.width*20/100)
	mainWidth := max(1, m.width*50/100)
	detailWidth := max(1, m.width-sidebarWidth-mainWidth)

	// Panel styles with rounded borders and focus-aware border color
	sidebarStyle := lipgloss.NewStyle().
		Width(sidebarWidth).
		Height(panelHeight).
		Border(lipgloss.RoundedBorder()).
		BorderRight(true).
		BorderForeground(m.borderColor(PanelSidebar))

	mainStyle := lipgloss.NewStyle().
		Width(mainWidth).
		Height(panelHeight).
		Border(lipgloss.RoundedBorder()).
		BorderRight(true).
		BorderForeground(m.borderColor(PanelMain))

	detailStyle := lipgloss.NewStyle().
		Width(detailWidth).
		Height(panelHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.borderColor(PanelDetail))

	// Inner dimensions accounting for borders
	sidebarInnerW := max(1, sidebarWidth-2)
	mainInnerW := max(1, mainWidth-2)
	detailInnerW := max(1, detailWidth-2)
	innerH := max(1, panelHeight-2)

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

// Cleanup releases resources held by the model, including closing the store.
// Call this after tea.Program.Run() returns to ensure the bbolt database is
// cleanly closed.
func (m Model) Cleanup() error {
	return m.store.Close()
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
