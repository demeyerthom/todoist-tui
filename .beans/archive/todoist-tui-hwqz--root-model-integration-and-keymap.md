---
# todoist-tui-hwqz
title: Root model integration and keymap
status: completed
type: task
priority: normal
created_at: 2026-05-07T19:37:30Z
updated_at: 2026-05-09T14:37:58Z
parent: todoist-tui-us18
blocked_by:
    - todoist-tui-546a
    - todoist-tui-ghej
    - todoist-tui-n9zx
---

## Description

Update `internal/ui/keymap/keymap.go` and `internal/ui/app.go` to support new navigation keys and delegate to sub-panel models instead of inline placeholder views.

## Requirements

### Keymap updates (`internal/ui/keymap/keymap.go`)
- Add to `DefaultKeyMap().Normal`:
  - `"collapse": "h"`
  - `"expand": "l"`
  - `"toggle_completed": "H"`

### Root model updates (`internal/ui/app.go`)

#### `Model` struct changes
- Add fields: `sidebar sidebar.Model`, `tasklist tasklist.Model`, `detail detail.Model`
- Add field: `synced bool` (initially false, set true on `SyncDoneMsg`)
- Remove inline placeholder rendering methods: `sidebarView()`, `mainView()`, `detailView()`

#### `NewModel()` updates
- Constructs `sidebar.NewModel(cfg, store, styles)`
- Constructs `tasklist.NewModel(cfg, store, styles)`
- Constructs `detail.NewModel(cfg, store, styles)`
- Passes `theme.NewStyles(cfg)` to all sub-panels

#### `Update()` changes
- Root-level keys handled first: `:`, `ctrl+c`, `tab`, `shift+tab`, `1`, `2`, `3`
- Panel switching keys (`1`, `2`, `3`) added alongside `tab`/`shift+tab`
- After root keys, delegate to focused sub-panel's `Update()`:
  ```go
  switch m.activePanel {
  case PanelSidebar:
      newSidebar, cmd := m.sidebar.Update(msg)
      m.sidebar = newSidebar.(sidebar.Model)
      return m, cmd
  case PanelMain:
      // same for tasklist
  case PanelDetail:
      // same for detail
  }
  ```
- Cross-panel messages routed by root:
  - `ProjectSelectedMsg`, `LabelSelectedMsg`, `FilterSelectedMsg` → sent to tasklist
  - `TaskSelectedMsg` → sent to detail
  - `SyncCompleteMsg` → sent to all three sub-panels
  - `SyncDoneMsg` → sets `synced = true`, sends init commands to sub-panels

#### `View()` changes
- Calls `m.sidebar.View()`, `m.tasklist.View()`, `m.detail.View()` instead of inline methods
- Each sub-panel view wrapped in its border style

## Acceptance Criteria

- Root model delegates to sub-panels correctly
- Cross-panel messages routed to correct panels
- Panel switching with Tab/Shift+Tab and 1/2/3 works
- `go build ./...` succeeds

## Summary of Changes

Updated internal/ui/keymap/keymap.go: added collapse (h), expand (l), toggle_completed (H) keybindings to Normal mode.

Rewrote internal/ui/app.go: added sidebar, tasklist, detail sub-model fields and synced bool to Model struct; NewModel creates shared theme.Styles and initializes all sub-panels; Update routes WindowSizeMsg to all panels with computed inner dimensions; Update handles cross-panel messages (ProjectSelectedMsg/LabelSelectedMsg/FilterSelectedMsg → tasklist, TaskSelectedMsg → detail, SyncCompleteMsg/SyncDoneMsg → all panels); Update delegates unrecognized keys to focused panel via delegateToPanel helper; panel switching with Tab/Shift+Tab and 1/2/3; View calls sub-panel View() methods instead of placeholders; removed sidebarView/mainView/detailView methods. go build and go vet pass.
