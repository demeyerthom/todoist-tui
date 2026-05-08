---
# todoist-tui-58dy
title: Periodic sync tick
status: todo
type: task
created_at: 2026-05-07T19:38:06Z
updated_at: 2026-05-07T19:38:06Z
parent: todoist-tui-fykz
blocked_by:
    - todoist-tui-lexk
---

## Description

Add a `SyncTickMsg` and a `tickSync()` command to `internal/ui/app.go` that fires every 30 seconds and triggers an incremental sync.

## Requirements

### Message and command
- `SyncTickMsg{}` message type (in `internal/ui/messages.go`)
- `tickSync() tea.Cmd` returns `tea.Tick(30*time.Second, func() tea.Msg { return SyncTickMsg{} })`

### `Model` struct changes
- Add `syncing bool` field — prevents concurrent syncs

### `Update()` changes
- After initial `SyncDoneMsg`:
  - Set `synced = true`, `syncing = false`
  - Return `tickSync()` to schedule first periodic tick
- On `SyncTickMsg`:
  - If `syncing == true`: return `tickSync()` (skip, sync already in progress)
  - Set `syncing = true`
  - Run `IncrementalSync` as a `tea.Cmd` (same pattern as initial `FullSync`)
- On incremental sync completion:
  - Set `syncing = false`
  - Send `SyncCompleteMsg` to all sub-panels
  - Return `tickSync()` to schedule next tick
- On incremental sync error:
  - Set `syncing = false`
  - Set `m.err` only for auth failures (401)
  - Return `tickSync()` to schedule next tick

### Sync state tracking
- `syncing` prevents concurrent syncs
- `synced` tracks whether initial sync completed
- Periodic sync only starts after initial sync succeeds

## Acceptance Criteria

- Incremental sync fires every 30 seconds after initial sync
- Concurrent syncs prevented via `syncing` flag
- `SyncCompleteMsg` emitted after each successful periodic sync
- `go build ./...` succeeds
