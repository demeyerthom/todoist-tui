---
# todoist-tui-p3n5
title: Sync status line
status: completed
type: task
priority: normal
created_at: 2026-05-07T19:38:19Z
updated_at: 2026-05-14T15:00:23Z
parent: todoist-tui-fykz
blocked_by:
    - todoist-tui-58dy
---

## Description

Render persistent sync status in the command bar area of `app.go`'s `View()`.

## Requirements

### `Model` struct changes
- Add `syncStatus string` field
- Add `lastSyncTime time.Time` field

### Status values
- `""` — no sync yet / before first sync
- `"syncing"` — sync in progress
- `"synced"` — last sync completed successfully
- `"failed"` — last sync failed

### Status transitions
- On `SyncDoneMsg`: `syncStatus = "synced"`, `lastSyncTime = time.Now()`
- On `SyncErrMsg`: `syncStatus = "failed"`
- On `SyncTickMsg`: `syncStatus = "syncing"`
- On periodic sync completion: `syncStatus = "synced"`, `lastSyncTime = time.Now()`
- On periodic sync error: `syncStatus = "failed"`

### Rendering in `View()`
- Right-aligned in the bottom bar area (always visible, not just in command mode)
- Status text:
  - Syncing: "Syncing..."
  - Synced: "Synced HH:MM" using `lastSyncTime.Format("15:04")"
  - Failed: "Sync failed" styled with `cfg.Theme.Error`
- When in command mode: command bar on left, sync status on right
- When not in command mode: sync status alone on the bottom line

## Acceptance Criteria

- Sync status visible at all times in bottom bar
- Status text correct for each state
- Synced time formatted as HH:MM
- Failed status styled with error color
- `go build ./...` succeeds
