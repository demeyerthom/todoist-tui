---
# todoist-tui-cn66
title: SyncToken helpers
status: todo
type: task
created_at: 2026-05-03T15:09:23Z
updated_at: 2026-05-03T15:09:23Z
parent: todoist-tui-g83v
blocked_by:
    - todoist-tui-p2ij
---

## Description

Implement GetSyncToken() and SetSyncToken() using the sync_meta bucket.

## Requirements

- GetSyncToken() (string, error) — reads key 'sync_token' from sync_meta bucket; returns empty string if not found (no error)
- SetSyncToken(token string) error — writes key 'sync_token' to sync_meta bucket
- Also: GetLastSyncTime() / SetLastSyncTime() for timestamp of last successful sync
- Uses the generic Put/Get helpers from the CRUD task

## Acceptance Criteria

- GetSyncToken returns '' when no token stored (first launch triggers full sync)
- SetSyncToken persists, GetSyncToken retrieves
- Timestamp round-trip works
