---
# todoist-tui-fykz
title: Periodic background sync
status: todo
type: feature
created_at: 2026-05-03T14:58:27Z
updated_at: 2026-05-03T14:58:27Z
parent: todoist-tui-96st
---

Incremental sync every 30 seconds in the background. Uses stored sync_token from bbolt for incremental syncs. Sync runs as a Bubbletea background command to avoid blocking the UI.
