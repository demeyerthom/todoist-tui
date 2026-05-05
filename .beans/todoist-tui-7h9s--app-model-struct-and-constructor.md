---
# todoist-tui-7h9s
title: App model struct and constructor
status: todo
type: task
created_at: 2026-05-03T15:10:09Z
updated_at: 2026-05-03T15:10:09Z
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
