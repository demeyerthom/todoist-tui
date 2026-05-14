package ui

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/demeyerthom/todoist-tui/internal/config"
	"github.com/demeyerthom/todoist-tui/internal/store"
	"github.com/demeyerthom/todoist-tui/internal/sync"
	"github.com/demeyerthom/todoist-tui/internal/ui/detail"
	"github.com/demeyerthom/todoist-tui/internal/ui/keymap"
	"github.com/demeyerthom/todoist-tui/internal/ui/sidebar"
	"github.com/demeyerthom/todoist-tui/internal/ui/tasklist"
	"github.com/demeyerthom/todoist-tui/internal/ui/theme"
)

// stackWidthThreshold is the terminal width below which the 3-panel layout
// switches from horizontal side-by-side to vertical stacking.
const stackWidthThreshold = 80

// Panel identifies the currently focused panel in the 3-panel layout.
type Panel int

const (
	PanelSidebar Panel = iota
	PanelMain
	PanelDetail
)

const panelCount = 3

// SyncErrMsg is emitted when a sync fails.
type SyncErrMsg struct{ Err error }

// Model is the root Bubbletea application model.
// It holds application-wide state and delegates to sub-models for each panel.
type Model struct {
	cfg         *config.Config
	store       *store.Store
	syncClient  *sync.Client
	styles      *theme.Styles
	activePanel Panel
	mode        keymap.Mode
	commandBuf  string
	width       int
	height      int
	err         error
	syncErr     error
	synced      bool
	syncing     bool
	sidebar     sidebar.Model
	tasklist    tasklist.Model
	detail      detail.Model
}

// NewModel creates a new root application model with the given dependencies.
// The model starts in Sidebar panel, Normal mode, with zero width and height.
func NewModel(cfg *config.Config, store *store.Store, syncClient *sync.Client) Model {
	styles := theme.NewStyles(cfg)
	return Model{
		cfg:         cfg,
		store:       store,
		syncClient:  syncClient,
		styles:      styles,
		activePanel: PanelSidebar,
		mode:        keymap.ModeNormal,
		commandBuf:  "",
		width:       0,
		height:      0,
		synced:      false,
		syncing:     false,
		sidebar:     sidebar.NewModel(cfg, store, styles),
		tasklist:    tasklist.NewModel(cfg, store, styles),
		detail:      detail.NewModel(cfg, store, styles),
	}
}

// Init implements tea.Model. It kicks off an incremental sync (which
// delegates to FullSync when no sync token exists) against the Todoist
// Sync API asynchronously. The result is delivered as SyncCompleteMsg or
// SyncErrMsg.
func (m Model) Init() tea.Cmd {
	return doIncrementalSync(m.syncClient, m.store)
}

// tickSync returns a command that fires a SyncTickMsg after 30 seconds.
// The tick is rescheduled after each sync completes (success or failure) to
// maintain a continuous 30-second background sync cycle.
func tickSync() tea.Cmd {
	return tea.Tick(30*time.Second, func(t time.Time) tea.Msg {
		return SyncTickMsg{}
	})
}

// doIncrementalSync returns a command that performs an incremental sync
// against the Todoist Sync API. On success it emits SyncCompleteMsg; on
// failure it emits SyncErrMsg.
func doIncrementalSync(syncClient *sync.Client, store *store.Store) tea.Cmd {
	return func() tea.Msg {
		err := syncClient.IncrementalSync(context.Background(), store)
		if err != nil {
			return SyncErrMsg{Err: err}
		}
		return SyncCompleteMsg{}
	}
}

// Update implements tea.Model. It handles incoming messages and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Account for error banner, sync error bar, and command bar in panel height
		panelHeight := m.height
		if m.err != nil {
			panelHeight--
		}
		if m.syncErr != nil {
			panelHeight--
		}
		if m.mode == keymap.ModeCommand {
			panelHeight--
		}
		innerH := max(1, panelHeight-2) // subtract top/bottom borders

		var cmds []tea.Cmd

		if m.width < stackWidthThreshold {
			// Stacked mode: all panels get full width
			innerW := max(1, m.width-2)
			newSidebar, cmd := m.sidebar.Update(tea.WindowSizeMsg{
				Width:  innerW,
				Height: innerH,
			})
			m.sidebar = newSidebar.(sidebar.Model)
			cmds = append(cmds, cmd)

			newTasklist, cmd := m.tasklist.Update(tea.WindowSizeMsg{
				Width:  innerW,
				Height: innerH,
			})
			m.tasklist = newTasklist.(tasklist.Model)
			cmds = append(cmds, cmd)

			newDetail, cmd := m.detail.Update(tea.WindowSizeMsg{
				Width:  innerW,
				Height: innerH,
			})
			m.detail = newDetail.(detail.Model)
			cmds = append(cmds, cmd)
		} else {
			// Horizontal mode: proportional widths with minimum enforcement
			sidebarWidth, mainWidth, detailWidth := m.panelWidths()

			newSidebar, cmd := m.sidebar.Update(tea.WindowSizeMsg{
				Width:  max(1, sidebarWidth-2),
				Height: innerH,
			})
			m.sidebar = newSidebar.(sidebar.Model)
			cmds = append(cmds, cmd)

			newTasklist, cmd := m.tasklist.Update(tea.WindowSizeMsg{
				Width:  max(1, mainWidth-2),
				Height: innerH,
			})
			m.tasklist = newTasklist.(tasklist.Model)
			cmds = append(cmds, cmd)

			newDetail, cmd := m.detail.Update(tea.WindowSizeMsg{
				Width:  max(1, detailWidth-2),
				Height: innerH,
			})
			m.detail = newDetail.(detail.Model)
			cmds = append(cmds, cmd)
		}

		return m, tea.Batch(cmds...)

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
			case "1":
				m.activePanel = PanelSidebar
			case "2":
				m.activePanel = PanelMain
			case "3":
				m.activePanel = PanelDetail
			default:
				return m.delegateToPanel(msg)
			}
			return m, nil

		case keymap.ModeInsert:
			if msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
			return m.delegateToPanel(msg)

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

	case ProjectSelectedMsg:
		newTasklist, cmd := m.tasklist.Update(msg)
		m.tasklist = newTasklist.(tasklist.Model)
		return m, cmd

	case LabelSelectedMsg:
		newTasklist, cmd := m.tasklist.Update(msg)
		m.tasklist = newTasklist.(tasklist.Model)
		return m, cmd

	case FilterSelectedMsg:
		newTasklist, cmd := m.tasklist.Update(msg)
		m.tasklist = newTasklist.(tasklist.Model)
		return m, cmd

	case TaskSelectedMsg:
		newDetail, cmd := m.detail.Update(msg)
		m.detail = newDetail.(detail.Model)
		return m, cmd

	case SyncCompleteMsg:
		m.synced = true
		m.syncing = false
		m.syncErr = nil
		var cmds []tea.Cmd
		newSidebar, cmd := m.sidebar.Update(msg)
		m.sidebar = newSidebar.(sidebar.Model)
		cmds = append(cmds, cmd)
		newTasklist, cmd := m.tasklist.Update(msg)
		m.tasklist = newTasklist.(tasklist.Model)
		cmds = append(cmds, cmd)
		newDetail, cmd := m.detail.Update(msg)
		m.detail = newDetail.(detail.Model)
		cmds = append(cmds, cmd)
		cmds = append(cmds, tickSync())
		return m, tea.Batch(cmds...)

	case SyncErrMsg:
		m.syncing = false
		m.syncErr = msg.Err
		if !m.synced || errors.Is(msg.Err, sync.ErrAuthFailed) {
			m.err = msg.Err
		}
		// Schedule next tick only if initial sync was successful
		if m.synced {
			return m, tickSync()
		}

	case SyncTickMsg:
		// Prevent concurrent syncs; if already syncing, just schedule next tick
		if m.syncing {
			return m, tickSync()
		}
		m.syncing = true
		return m, doIncrementalSync(m.syncClient, m.store)
	}

	return m, nil
}

// delegateToPanel sends a message to the currently focused sub-panel.
func (m Model) delegateToPanel(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.activePanel {
	case PanelSidebar:
		newModel, cmd := m.sidebar.Update(msg)
		m.sidebar = newModel.(sidebar.Model)
		return m, cmd
	case PanelMain:
		newModel, cmd := m.tasklist.Update(msg)
		m.tasklist = newModel.(tasklist.Model)
		return m, cmd
	case PanelDetail:
		newModel, cmd := m.detail.Update(msg)
		m.detail = newModel.(detail.Model)
		return m, cmd
	default:
		return m, nil
	}
}

// View implements tea.Model. It renders the 3-panel layout:
// sidebar (~20%), main (~50%), detail (~30%) joined horizontally.
// A fatal error (m.err) is displayed as a banner at the top.
// A transient sync error (m.syncErr) is displayed as a status line at the bottom.
// When in command mode, a command bar is displayed at the bottom.
func (m Model) View() string {
	if m.width <= 0 || m.height <= 0 {
		if m.err != nil {
			return "Error: " + m.err.Error()
		}
		return "Initializing..."
	}

	// Compute effective panel height, accounting for error banner, sync status, and command bar
	panelHeight := m.height
	if m.err != nil {
		panelHeight--
	}
	if m.syncErr != nil {
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
	if m.syncErr != nil {
		syncBar := lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.cfg.Theme.Warning)).
			Width(m.width).
			Render("Sync error: " + m.syncErr.Error())
		content = lipgloss.JoinVertical(lipgloss.Top, content, syncBar)
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
// At width >= stackWidthThreshold, panels are side-by-side (horizontal).
// Below the threshold, panels are stacked vertically with the focused panel
// taking full height and unfocused panels collapsed to header strips.
func (m Model) panelsView(panelHeight int) string {
	if m.width < stackWidthThreshold {
		return m.stackedView(panelHeight)
	}
	return m.horizontalView(panelHeight)
}

// panelWidths returns the outer widths for the three panels in horizontal mode.
// Widths are proportional (sidebar 20%, main 50%, detail 30%) with minimum
// enforcement. If the total width cannot satisfy all minimums, detail is
// reduced first, then main, then sidebar.
func (m Model) panelWidths() (sidebarWidth, mainWidth, detailWidth int) {
	sidebarWidth = m.width * 20 / 100
	mainWidth = m.width * 50 / 100
	detailWidth = m.width - sidebarWidth - mainWidth

	// Enforce minimums by taking from lower-priority panels.
	if detailWidth < 30 {
		mainWidth -= 30 - detailWidth
		detailWidth = 30
	}
	if mainWidth < 40 {
		sidebarWidth -= 40 - mainWidth
		mainWidth = 40
	}
	if sidebarWidth < 20 {
		sidebarWidth = 20
		remaining := m.width - sidebarWidth
		mainWidth = max(40, remaining-30)
		detailWidth = remaining - mainWidth
	}

	// Clamp to positive values.
	sidebarWidth = max(1, sidebarWidth)
	mainWidth = max(1, mainWidth)
	detailWidth = max(1, detailWidth)
	return
}

// horizontalView renders the three panels side-by-side.
// Each sub-panel's View() is wrapped in a border style with focus-aware
// border coloring. Sidebar and main panels have a right border to create
// clean visual joins.
// When not yet synced, "Loading..." is rendered centered in each panel.
func (m Model) horizontalView(panelHeight int) string {
	sidebarWidth, mainWidth, detailWidth := m.panelWidths()

	sidebarStyle := lipgloss.NewStyle().
		Width(max(1, sidebarWidth-2)).
		Height(max(1, panelHeight-2)).
		Border(lipgloss.RoundedBorder()).
		BorderRight(true).
		BorderForeground(m.borderColor(PanelSidebar))

	mainStyle := lipgloss.NewStyle().
		Width(max(1, mainWidth-2)).
		Height(max(1, panelHeight-2)).
		Border(lipgloss.RoundedBorder()).
		BorderRight(true).
		BorderForeground(m.borderColor(PanelMain))

	detailStyle := lipgloss.NewStyle().
		Width(max(1, detailWidth-2)).
		Height(max(1, panelHeight-2)).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.borderColor(PanelDetail))

	if !m.synced {
		sidebarInnerW := max(1, sidebarWidth-2)
		mainInnerW := max(1, mainWidth-2)
		detailInnerW := max(1, detailWidth-2)
		innerH := max(1, panelHeight-2)

		sidebar := sidebarStyle.Render(
			lipgloss.Place(sidebarInnerW, innerH, lipgloss.Center, lipgloss.Center,
				m.styles.MutedText().Render("Loading...")))
		main := mainStyle.Render(
			lipgloss.Place(mainInnerW, innerH, lipgloss.Center, lipgloss.Center,
				m.styles.MutedText().Render("Loading...")))
		detailPanel := detailStyle.Render(
			lipgloss.Place(detailInnerW, innerH, lipgloss.Center, lipgloss.Center,
				m.styles.MutedText().Render("Loading...")))
		return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, main, detailPanel)
	}

	sidebar := sidebarStyle.Render(m.sidebar.View())
	main := mainStyle.Render(m.tasklist.View())
	detailPanel := detailStyle.Render(m.detail.View())

	return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, main, detailPanel)
}

// stackedView renders the three panels stacked vertically.
// The focused panel takes the full available height; unfocused panels
// collapse to 1-line header strips.
// When not yet synced, "Loading..." is rendered centered in the focused
// panel and as header text in unfocused panels instead of sub-panel contents.
func (m Model) stackedView(panelHeight int) string {
	labels := map[Panel]string{
		PanelSidebar: "Projects",
		PanelMain:    "Tasks",
		PanelDetail:  "Details",
	}

	var panels []string
	for _, p := range []Panel{PanelSidebar, PanelMain, PanelDetail} {
		if p == m.activePanel {
			style := lipgloss.NewStyle().
				Width(max(1, m.width-2)).
				Height(panelHeight - 2).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(m.borderColor(p))

			if !m.synced {
				innerW := max(1, m.width-2)
				innerH := max(1, panelHeight-4) // subtract outer height and border
				panels = append(panels, style.Render(
					lipgloss.Place(innerW, innerH, lipgloss.Center, lipgloss.Center,
						m.styles.MutedText().Render("Loading..."))))
			} else {
				var content string
				switch p {
				case PanelSidebar:
					content = m.sidebar.View()
				case PanelMain:
					content = m.tasklist.View()
				case PanelDetail:
					content = m.detail.View()
				}
				panels = append(panels, style.Render(content))
			}
		} else {
			headerStyle := lipgloss.NewStyle().
				Width(max(1, m.width-2)).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(m.borderColor(p))
			if !m.synced {
				panels = append(panels, headerStyle.Render(
					m.styles.MutedText().Render("Loading...")))
			} else {
				panels = append(panels, headerStyle.Render(m.styles.Header().Render(labels[p])))
			}
		}
	}

	return lipgloss.JoinVertical(lipgloss.Top, panels...)
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
