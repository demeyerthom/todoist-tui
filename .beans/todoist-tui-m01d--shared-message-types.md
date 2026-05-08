---
# todoist-tui-m01d
title: Shared message types
status: todo
type: task
created_at: 2026-05-07T19:34:43Z
updated_at: 2026-05-07T19:34:43Z
parent: todoist-tui-0rtm
---

## Description

Create `internal/ui/messages.go` with all cross-panel message types. These messages are used by sidebar, tasklist, and detail packages to communicate state changes through the root model.

## Requirements

- `ProjectSelectedMsg{ID string}` — emitted when a project is selected in the sidebar
- `FilterSelectedMsg{ID string, Query string}` — emitted when a filter is selected in the sidebar
- `LabelSelectedMsg{Name string}` — emitted when a label is selected in the sidebar
- `TaskSelectedMsg{ID string}` — emitted when a task row is selected in the task list
- `SyncCompleteMsg{}` — emitted when any sync (full or incremental) completes successfully
- `SyncTickMsg{}` — emitted every 30 seconds to trigger periodic sync
- `ToggleCompletedMsg{}` — emitted when user toggles completed task visibility (`H` key)

## Acceptance Criteria

- `internal/ui/messages.go` exists with all message types as plain Go structs
- `go build ./...` succeeds
- No dependencies on any other package (messages are pure types)
