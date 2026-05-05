---
# todoist-tui-7tn2
title: Incremental sync implementation
status: todo
type: task
created_at: 2026-05-03T15:09:50Z
updated_at: 2026-05-03T15:09:50Z
parent: todoist-tui-efa6
blocked_by:
    - todoist-tui-hgg5
---

## Description

Implement IncrementalSync that reads the stored sync token and sends an incremental sync request.

## Requirements

- IncrementalSync(ctx context.Context, store *store.Store) error
- Reads stored sync_token via store.GetSyncToken
- If token is empty, delegates to FullSync
- Otherwise sends SyncRequest with stored token + all ResourceTypes
- Merges results into store (same as FullSync: Put for each entity)
- Updates sync_token on success
- Deleted items (is_deleted flag) are removed from store

## Acceptance Criteria

- Empty sync token triggers FullSync
- Incremental sync updates store with only changed entities
- Deleted entities are removed
- Sync token updated after success
