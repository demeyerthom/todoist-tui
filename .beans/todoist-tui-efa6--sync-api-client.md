---
# todoist-tui-efa6
title: Sync API client
status: todo
type: feature
priority: normal
created_at: 2026-05-03T14:58:18Z
updated_at: 2026-05-03T15:11:30Z
parent: todoist-tui-j3br
---

Todoist Sync API v9 client: perform full sync on first launch (sync_token=*), store response data in bbolt, persist returned sync_token for incremental syncs. Resolve temp_id_mapping from sync responses to replace placeholder IDs with real server IDs. Auth via personal API token only. Use google/uuid for command idempotency UUIDs.



## Task Dependency Graph

1. `todoist-tui-7gco` Define API request/response types ← F1-dir-structure
2. `todoist-tui-c4lv` HTTP client with auth ← 7gco
3. `todoist-tui-hgg5` Full sync implementation ← c4lv + F3-store-CRUD
4. `todoist-tui-7tn2` Incremental sync ← hgg5
5. `todoist-tui-mwib` Temp ID mapping resolution ← hgg5
6. `todoist-tui-w5h6` Sync client tests ← 7tn2 + mwib

## Cross-feature dependencies
- Depends on F1 (project init) and F3 (store/models)
- F5 (TUI) init sync depends on this feature
