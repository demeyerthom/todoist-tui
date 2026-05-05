---
# todoist-tui-y2ze
title: Store tests
status: todo
type: task
created_at: 2026-05-03T15:09:28Z
updated_at: 2026-05-03T15:09:28Z
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
