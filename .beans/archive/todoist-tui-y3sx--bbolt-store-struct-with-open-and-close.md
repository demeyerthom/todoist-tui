---
# todoist-tui-y3sx
title: bbolt Store struct with Open and Close
status: completed
type: task
priority: normal
created_at: 2026-05-05T20:15:49Z
updated_at: 2026-05-05T20:16:19Z
---

Implement Store type with Open/Close using bbolt, storing the database at the XDG data directory. StoreConfig with Path and Mode, Store holds *bbolt.DB, New() opens DB, DBPath() returns XDG data path, Close() shuts down, parent dirs created if missing, ErrStoreNotOpen sentinel error.

## Summary of Changes

- Created `internal/store/store.go` with:
  - `ErrStoreNotOpen` sentinel error
  - `DefaultMode` constant (0o600)
  - `StoreConfig` struct with Path and Mode fields
  - `Store` struct wrapping `*bbolt.DB`
  - `DBPath()` function returning XDG data dir path (`/todoist-tui/todoist.db` or `~/.local/share/todoist-tui/todoist.db`)
  - `New(cfg StoreConfig) (*Store, error)` constructor that opens/creates the bbolt DB, creating parent dirs if missing
  - `Close() error` method that shuts down the DB and returns `ErrStoreNotOpen` if already closed
- Updated `internal/store/doc.go` to remove blank import (bbolt is now imported directly in store.go)
- Both `go build ./...` and `go vet ./...` pass cleanly.
