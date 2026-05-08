---
# todoist-tui-fykz
title: Periodic background sync
status: todo
type: feature
priority: normal
created_at: 2026-05-03T14:58:27Z
updated_at: 2026-05-07T19:38:35Z
parent: todoist-tui-96st
blocked_by:
    - todoist-tui-us18
---

Incremental sync every 30 seconds in the background. Uses stored sync_token from bbolt for incremental syncs. Sync runs as a Bubbletea background command to avoid blocking the UI.

## Design decisions

- **Interval**: 30 seconds, fired via Bubbletea `tea.Tick`.
- **Status line**: Persistent status in the command bar area — "Syncing..." during sync, "Synced HH:MM" on success, "Sync failed" on error. Left-aligned for commands, right-aligned for sync status.
- **Error handling**: Network/HTTP errors shown in status line only. Auth failures (401) also displayed in the error banner at top.
- **Startup**: Initial app launch shows "Loading..." in all panels during full sync.
- **`resolveTempIDs`**: Wire up now — call from `FullSync()` and `IncrementalSync()` even though M1 is read-only (no temp IDs created locally). Gets real-world testing of the resolution logic.
- **Panel refresh**: `SyncCompleteMsg` propagated through `Update()` — each panel re-reads from store independently.
