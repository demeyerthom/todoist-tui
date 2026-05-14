---
# todoist-tui-58dy
title: Periodic sync tick
status: completed
type: task
priority: normal
created_at: 2026-05-07T19:38:06Z
updated_at: 2026-05-09T15:03:07Z
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

## Summary of Changes

Added syncing bool field to Model struct to prevent concurrent syncs. Added tickSync() command that fires SyncTickMsg every 30 seconds via tea.Tick. Added doIncrementalSync() command that runs IncrementalSync and returns SyncCompleteMsg or SyncErrMsg. Updated SyncDoneMsg handler to set syncing=false and schedule first tickSync. Added SyncTickMsg handler that skips if syncing=true, otherwise sets syncing=true and runs doIncrementalSync. Updated SyncCompleteMsg handler to set syncing=false, forward to all sub-panels, and schedule next tickSync. Updated SyncErrMsg handler to set syncing=false, show error only for initial sync or auth failures, and reschedule tickSync if initial sync was successful. go build, go vet, and go test all pass.
