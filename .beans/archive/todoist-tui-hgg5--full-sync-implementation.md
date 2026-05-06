---
# todoist-tui-hgg5
title: Full sync implementation
status: completed
type: task
priority: normal
created_at: 2026-05-03T15:09:45Z
updated_at: 2026-05-06T10:00:20Z
parent: todoist-tui-efa6
blocked_by:
    - todoist-tui-c4lv
    - todoist-tui-xbbs
---

## Description

Implement FullSync that sends sync_token=* with all resource types, writes results to the store, and persists the sync token.

## Requirements

- FullSync(ctx context.Context, store *store.Store) error
- Sends SyncRequest with sync_token='*', all ResourceTypes
- Receives SyncResponse, writes each entity collection to bbolt store:
  - projects → store.Put for each project
  - items → store.Put for each task
  - sections → store.Put for each section
  - labels → store.Put for each label
  - filters → store.Put for each filter
- Persists returned sync_token via store.SetSyncToken
- Returns ErrAuthFailed on 401, ErrSyncFailed on non-200

## Acceptance Criteria

- FullSync stores all entities in bbolt
- Sync token is persisted after successful sync
- Auth errors return ErrAuthFailed
- Network errors return ErrSyncFailed

## Summary of Changes

Created internal/sync/sync.go with:
- FullSync method on Client that sends sync_token='*' with all ResourceTypes
- Iterates all entity collections (Tasks, Projects, Sections, Labels, Filters) and stores them via typed Store methods
- Persists sync_token via store.SetSyncToken after all entities are written
- Propagates ErrAuthFailed/ErrSyncFailed from DoSync
- Store errors wrapped with entity type and ID context
- go build and go vet pass
