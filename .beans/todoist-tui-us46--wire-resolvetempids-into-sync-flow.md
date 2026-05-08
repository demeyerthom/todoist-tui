---
# todoist-tui-us46
title: Wire resolveTempIDs into sync flow
status: todo
type: task
created_at: 2026-05-07T19:35:04Z
updated_at: 2026-05-07T19:35:04Z
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
