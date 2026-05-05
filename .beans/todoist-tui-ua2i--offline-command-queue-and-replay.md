---
# todoist-tui-ua2i
title: Offline command queue and replay
status: todo
type: feature
created_at: 2026-05-03T14:58:35Z
updated_at: 2026-05-03T14:58:35Z
parent: todoist-tui-ta89
---

Queue write commands when offline (internal/queue). Commands are stored locally and replayed in order when connectivity is restored. The command queue persists across app restarts via bbolt. On reconnect, drain the queue by sending all pending commands to the Sync API.
