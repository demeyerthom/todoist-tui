---
# todoist-tui-us46
title: Wire resolveTempIDs into sync flow
status: completed
type: task
priority: normal
created_at: 2026-05-07T19:35:04Z
updated_at: 2026-05-09T14:52:25Z
parent: todoist-tui-fykz
---

## Description

Call `resolveTempIDs()` from `FullSync()` and `IncrementalSync()` in `internal/sync/sync.go` after storing all entities and before persisting the sync token.

## Requirements

- In `FullSync()`: after all entity storage loops, before `SetSyncToken()`, call `resolveTempIDs(s, resp.TempIDMapping)`
- In `IncrementalSync()`: same placement — after entity loops, before `SetSyncToken()`
- Skip if `resp.TempIDMapping` is nil or empty (no-op for most syncs)
- Error from `resolveTempIDs()` should abort the sync and return the error

## Acceptance Criteria

- `go build ./...` succeeds
- `go test ./internal/sync/...` passes (existing tests still pass)
- Add test case: sync with non-empty `TempIDMapping` triggers `resolveTempIDs`

## Summary of Changes

Added resolveTempIDs calls in internal/sync/sync.go: in FullSync() after entity storage loops and before SetSyncToken(), and in IncrementalSync() at the same placement. Both calls are guarded by len(resp.TempIDMapping) > 0 to skip when no temp IDs exist. Errors from resolveTempIDs are wrapped and returned, aborting the sync. Added TestFullSync_ResolvesTempIDs test case verifying tmp-task-1 is resolved to real-task-1, temp entry is cleaned up, and sync token is persisted. All 9 sync tests pass.
