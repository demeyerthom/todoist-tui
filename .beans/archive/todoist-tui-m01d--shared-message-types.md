---
# todoist-tui-m01d
title: Shared message types
status: completed
type: task
priority: normal
created_at: 2026-05-07T19:34:43Z
updated_at: 2026-05-08T16:32:57Z
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

## Summary of Changes

Created `internal/ui/messages.go` with 7 cross-panel message types as plain Go structs:

- **ProjectSelectedMsg** — sidebar project selection (ID string)
- **FilterSelectedMsg** — sidebar filter selection (ID, Query)
- **LabelSelectedMsg** — sidebar label selection (Name string)
- **TaskSelectedMsg** — task list row selection (ID string)
- **SyncCompleteMsg** — emitted after any sync completes (distinct from SyncDoneMsg for startup sync only)
- **SyncTickMsg** — 30-second periodic sync trigger
- **ToggleCompletedMsg** — completed task visibility toggle

All messages are pure types with no methods and no imports beyond the `ui` package declaration. `go build ./...` succeeds.

## Summary of Changes\n\nCreated  with 7 cross-panel message types: ProjectSelectedMsg, FilterSelectedMsg, LabelSelectedMsg, TaskSelectedMsg, SyncCompleteMsg, SyncTickMsg, ToggleCompletedMsg. All are pure Go structs with no external dependencies. SyncCompleteMsg is distinct from the existing SyncDoneMsg (startup-only) to cover all sync events including periodic background syncs.
