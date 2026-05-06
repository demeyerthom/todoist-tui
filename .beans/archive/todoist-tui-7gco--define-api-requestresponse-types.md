---
# todoist-tui-7gco
title: Define API request/response types
status: completed
type: task
priority: normal
created_at: 2026-05-03T15:09:37Z
updated_at: 2026-05-06T09:55:36Z
parent: todoist-tui-efa6
blocked_by:
    - todoist-tui-4yfu
---

## Description

Define SyncRequest, SyncResponse, Command, and related types in internal/sync/types.go matching the Todoist Sync API v9 format.

## Requirements

- SyncRequest: sync_token string, resource_types []string, commands []Command
- SyncResponse: sync_token, items []model.Task, projects []model.Project, sections []model.Section, labels []model.Label, filters []model.Filter, temp_id_mapping map[string]string
- Command: type string, args map[string]any, uuid string, temp_id string
- ResourceTypes constant slice: 'items', 'projects', 'sections', 'labels', 'filters'
- SyncEndpoint constant: 'https://api.todoist.com/sync/v9/sync'
- All structs have 'json' tags matching API snake_case

## Acceptance Criteria

- Types defined in internal/sync/types.go
- ResourceTypes includes all 5 resource types
- 'go build ./...' succeeds

## Summary of Changes

Created internal/sync/types.go with:
- SyncRequest, SyncResponse, Command structs with JSON tags matching Todoist Sync API v9
- ResourceTypes var (items, projects, sections, labels, filters)
- SyncEndpoint constant (https://api.todoist.com/sync/v9/sync)
- ErrAuthFailed and ErrSyncFailed sentinel errors
- All types reference model package types correctly
- go build and go vet pass
