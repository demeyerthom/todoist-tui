---
# todoist-tui-ghej
title: Task list Update and selection
status: todo
type: task
created_at: 2026-05-07T19:36:57Z
updated_at: 2026-05-07T19:36:57Z
parent: todoist-tui-rcpd
blocked_by:
    - todoist-tui-zpjy
---

## Description

Create `internal/ui/tasklist/update.go` with the `Update()` method handling key events, selection changes, and filter changes.

## Requirements

### `Update(msg tea.Msg) (tea.Model, tea.Cmd)` method on `Model`

#### Key handling (normal mode, using keymap)
- `j`/`down`: move selection down one row (skip separator rows)
- `k`/`up`: move selection up one row (skip separator rows)
- `H`/`toggle_completed`: toggle `showCompleted`, call `loadTasks()`, preserve selection if possible
- `g`/`go_top`: move selection to first task row
- `G`/`go_bottom`: move selection to last task row

#### Message handling
- `ProjectSelectedMsg{ID}`: set `currentFilter = Filter{ProjectID: id}`, call `loadTasks()`, reset selection
- `LabelSelectedMsg{Name}`: set `currentFilter = Filter{LabelName: name}`, call `loadTasks()`, reset selection
- `FilterSelectedMsg{ID, Query}`: set `currentFilter = Filter{FilterQuery: query, IsFilter: true}`, clear groups (unsupported in M1)
- `SyncCompleteMsg`: call `loadTasks()`, preserve `selectedID` if task still exists
- `tea.WindowSizeMsg`: update width/height, rebuild table

#### Selection tracking
- `selectedID` tracks the currently selected task ID
- On selection change, emits `TaskSelectedMsg{ID: selectedID}`
- Selection resets when filter changes
- Selection preserved when `SyncCompleteMsg` arrives if task still exists

## Acceptance Criteria

- Selection moves correctly with j/k/g/G, skipping separator rows
- H toggles completed visibility and reloads
- Filter messages trigger correct data reload
- TaskSelectedMsg emitted on selection change
- `go build ./...` succeeds
