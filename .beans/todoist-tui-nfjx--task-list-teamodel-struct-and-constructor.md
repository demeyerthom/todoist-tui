---
# todoist-tui-nfjx
title: Task list tea.Model struct and constructor
status: todo
type: task
created_at: 2026-05-07T19:36:09Z
updated_at: 2026-05-07T19:36:09Z
parent: todoist-tui-rcpd
blocked_by:
    - todoist-tui-m01d
---

## Description

Create `internal/ui/tasklist/model.go` with the `Model` struct and message types for the task list panel.

## Requirements

### `Model` struct fields
- `store *store.Store` — reads tasks, sections
- `cfg *config.Config` — theme and keymap access
- `styles *theme.Styles` — lipgloss style definitions
- `width int`, `height int` — panel dimensions
- `selectedID string` — currently selected task ID
- `showCompleted bool` — whether completed tasks are visible (default: false)
- `currentFilter Filter` — current filter context
- `groups []SectionGroup` — loaded task groups
- `table table.Model` — bubble-table instance

### `Filter` struct
- `ProjectID string` — filter by project ID (empty = not filtering by project)
- `LabelName string` — filter by label name (empty = not filtering by label)
- `FilterQuery string` — raw Todoist filter query (empty = not a filter)
- `IsFilter bool` — true if this is a Todoist filter (unsupported in M1)

### `SectionGroup` struct
- `Section *model.Section` — nil for unsectioned tasks
- `Tasks []model.Task` — tasks in this group

### Constructor
- `NewModel(cfg *config.Config, store *store.Store, styles *theme.Styles) Model`
- `Init() tea.Cmd` returns nil

## Acceptance Criteria

- `internal/ui/tasklist/model.go` exists and compiles
- `NewModel()` returns properly initialized model
- `go build ./...` succeeds
