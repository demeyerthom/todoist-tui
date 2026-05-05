---
# todoist-tui-p2ij
title: bbolt Store struct with Open and Close
status: todo
type: task
created_at: 2026-05-03T15:09:06Z
updated_at: 2026-05-03T15:09:06Z
parent: todoist-tui-g83v
blocked_by:
    - todoist-tui-umza
---

## Description

Implement the Store type with Open/Close using bbolt, storing the database at the XDG data directory.

## Requirements

- StoreConfig struct with Path string (defaults to XDG data dir) and Mode os.FileMode
- Store struct holds *bbolt.DB internally
- New(cfg StoreConfig) (*Store, error) opens/creates the bbolt DB
- DBPath() returns ~/.local/share/todoist-tui/todoist.db (using os.UserCacheDir or XDG_DATA_HOME fallback)
- Close() error calls db.Close()
- Parent directories created if missing
- Sentinel error: ErrStoreNotOpen

## Acceptance Criteria

- Store opens bbolt DB at XDG data path
- Close() cleanly shuts down
- Parent directories created if missing
- 'go build ./...' succeeds
