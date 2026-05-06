---
# todoist-tui-qzhx
title: Implement FullSync method on Client
status: completed
type: task
priority: normal
created_at: 2026-05-06T09:58:24Z
updated_at: 2026-05-06T09:58:52Z
---

Add FullSync method to sync.Client that sends sync_token=* with all resource types, writes results to the store, and persists the sync token.

## Summary of Changes

- Added  with  method on 
-  sends  via 
- Writes each entity collection to bbolt using typed store accessors (, , , , )
- Persists returned  via 
- Returns  on 401,  on non-200 (from )
- Returns store errors wrapped with context (entity type + ID)
-  and  pass
