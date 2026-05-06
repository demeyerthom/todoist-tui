---
# todoist-tui-7h9s
title: App model struct and constructor
status: completed
type: task
priority: normal
created_at: 2026-05-03T15:10:09Z
updated_at: 2026-05-06T11:18:18Z
parent: todoist-tui-cwq9
blocked_by:
    - todoist-tui-4yfu
---

## Description

Define the root Bubbletea Model struct and constructor in internal/ui/app.go.

## Requirements

- Model struct holds: cfg *config.Config, store *store.Store, sync *sync.Client, activePanel Panel, mode keymap.Mode, width int, height int, err error
- Panel enum: PanelSidebar, PanelMain, PanelDetail
- NewModel(cfg *config.Config, store *store.Store, sync *sync.Client) Model returns initialized model
- Implements bubbletea.Model: Init() tea.Cmd, Update(msg tea.Msg) (tea.Model, tea.Cmd), View() string
- Init() returns nil for now (sync cmd added in separate task)
- Update() returns model unchanged, no-op for now
- View() returns placeholder string for now

## Acceptance Criteria

- Model struct defined with all required fields
- NewModel initializes all fields
- bubbletea.Model interface satisfied
- 'go build ./...' compiles

## Summary of Changes

Created two files:

1. **`internal/ui/keymap/mode.go`** — Mode type as `int` with `ModeNormal`, `ModeInsert`, `ModeCommand` constants and a `String()` method.

2. **`internal/ui/app.go`** — Root Bubbletea Model struct with fields: `cfg`, `store`, `syncClient`, `activePanel`, `mode`, `width`, `height`, `err`. Panel enum (`PanelSidebar`, `PanelMain`, `PanelDetail`). `NewModel` constructor initializes defaults (Sidebar panel, Normal mode, zero dimensions). Implements `tea.Model` interface (`Init`, `Update`, `View`) as no-op stubs.

Both `go build ./...` and `go vet ./...` pass cleanly.

## Summary of Changes

Created internal/ui/keymap/mode.go with Mode type (ModeNormal, ModeInsert, ModeCommand) and String() method.
Created internal/ui/app.go with Panel enum (PanelSidebar, PanelMain, PanelDetail), Model struct with all required fields, NewModel constructor, and bubbletea.Model interface implementation (Init/Update/View as no-ops).
go build and go vet pass.
