---
# todoist-tui-hwqz
title: Root model integration and keymap
status: todo
type: task
created_at: 2026-05-07T19:37:30Z
updated_at: 2026-05-07T19:37:30Z
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
