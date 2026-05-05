---
# todoist-tui-ta89
title: 'M2: Writes'
status: todo
type: milestone
priority: normal
created_at: 2026-05-03T14:58:09Z
updated_at: 2026-05-03T14:59:16Z
blocked_by:
    - todoist-tui-96st
---

Full write support: Quick Add, edit forms, complete/delete tasks and projects, optimistic updates, offline queue.



## Key Decisions
- Temp IDs use tmp- prefix for optimistic updates
- Command queue persists across app restarts via bbolt
- On reconnect: drain queue, send all pending commands to Sync API
- Quick Add leverages server-side natural language parsing
