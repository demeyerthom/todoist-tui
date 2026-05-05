---
# todoist-tui-vc7f
title: Optimistic UI updates with temp IDs
status: todo
type: feature
created_at: 2026-05-03T14:58:33Z
updated_at: 2026-05-03T14:58:33Z
parent: todoist-tui-ta89
---

Optimistic UI updates using tmp- prefixed temp IDs. When a write command is issued, immediately reflect the change in the UI with a temp ID. Sync API responses include temp_id_mapping that maps temp IDs to real server IDs — resolve these to update the local store.
