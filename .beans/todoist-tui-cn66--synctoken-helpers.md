---
# todoist-tui-cn66
title: SyncToken helpers
status: completed
type: task
priority: normal
created_at: 2026-05-03T15:09:23Z
updated_at: 2026-05-05T20:21:55Z
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

## Progress\n\n- [x] Implemented GetSyncToken() — returns empty string with nil error when no token stored\n- [x] Implemented SetSyncToken() — persists sync_token key in sync_meta bucket\n- [x] Implemented GetLastSyncTime() — returns empty string with nil error when not found\n- [x] Implemented SetLastSyncTime() — stores time.Time as RFC3339 string\n- [x] Uses bbolt transactions directly (not generic Put/Get helpers, which don't exist yet)\n- [x] go build ./... passes\n- [x] go vet ./... passes

## Summary of Changes\n\nAdded internal/store/sync.go:\n- GetSyncToken() returns sync_token from sync_meta bucket, empty string if not found\n- SetSyncToken(token) writes sync_token to sync_meta bucket\n- GetLastSyncTime() returns last_sync_time as RFC3339 string, empty string if not found\n- SetLastSyncTime(t time.Time) writes RFC3339-formatted timestamp\n- Uses direct bbolt transactions (not generic CRUD helpers) to handle not-found gracefully\n- Package-level constants for bucket and key names
