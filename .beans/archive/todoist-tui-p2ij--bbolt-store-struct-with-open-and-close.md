---
# todoist-tui-p2ij
title: bbolt Store struct with Open and Close
status: completed
type: task
priority: normal
created_at: 2026-05-03T15:09:06Z
updated_at: 2026-05-05T20:17:58Z
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

## Summary of Changes\n\nImplemented internal/store package with bbolt Store struct:\n- StoreConfig struct with Path (defaults to XDG data dir) and Mode (defaults to 0o600)\n- Store struct holding *bbolt.DB internally\n- New(cfg StoreConfig) opens/creates bbolt DB with parent directory creation\n- DBPath() returns XDG_DATA_HOME/todoist-tui/todoist.db with ~/.local/share fallback and /tmp last resort\n- Close() calls db.Close() and nils the reference; returns ErrStoreNotOpen on double-close\n- ErrStoreNotOpen sentinel error\n- doc.go updated to remove blank import (bbolt now directly imported)
