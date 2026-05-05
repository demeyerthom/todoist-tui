---
# todoist-tui-w5h6
title: Sync client tests
status: todo
type: task
created_at: 2026-05-03T15:09:59Z
updated_at: 2026-05-03T15:09:59Z
parent: todoist-tui-efa6
blocked_by:
    - todoist-tui-7tn2
    - todoist-tui-mwib
---

## Description

Tests for the sync client using httptest mock servers.

## Requirements

- Full sync test: mock server returns full response with all entities, verify store populated
- Incremental sync test: mock server returns delta, verify store updated
- Auth failure test: mock returns 401, verify ErrAuthFailed
- Network error test: connection refused, verify ErrSyncFailed
- Temp ID mapping test: mock returns temp_id_mapping, verify IDs replaced in store
- Empty sync token triggers full sync
- Use httptest.NewServer for mock
- Use t.TempDir() for test store

## Acceptance Criteria

- 'go test ./internal/sync/...' passes
- Covers FullSync, IncrementalSync, auth failure, network error, temp_id_mapping
