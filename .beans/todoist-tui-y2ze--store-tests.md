---
# todoist-tui-y2ze
title: Store tests
status: completed
type: task
priority: normal
created_at: 2026-05-03T15:09:28Z
updated_at: 2026-05-05T20:24:10Z
parent: todoist-tui-g83v
blocked_by:
    - todoist-tui-xbbs
    - todoist-tui-cn66
---

## Description

Table-driven tests for store operations: CRUD, sync token persistence, bucket creation.

## Requirements

- Test Open creates DB file and all buckets
- Test Put/Get round-trip for each model type (Task, Project, etc.)
- Test List returns all items in a bucket
- Test Delete removes items
- Test GetSyncToken/SetSyncToken round-trip
- Test GetSyncToken returns empty string when not set
- Use t.TempDir() for test DB isolation
- Test Close is idempotent

## Acceptance Criteria

- 'go test ./internal/store/...' passes
- Covers Open, Put, Get, Delete, List, SyncToken, Close

## Summary of Changes\n\nCreated internal/store/store_test.go with table-driven tests covering:\n- TestOpen: creates DB file and all buckets, creates parent directories\n- TestClose: idempotent, returns ErrStoreNotOpen on uninitialized\n- TestPutGet: round-trip for all 6 model types, error on nonexistent key\n- TestList: returns all items, empty bucket returns empty slice\n- TestDelete: removes item, nonexistent key is no-op\n- TestSyncToken: empty string when not set, round-trip, overwrites\n- TestLastSyncTime: empty string when not set, round-trip with RFC3339\n\nAll 18 subtests pass. Uses t.TempDir() for isolation.
