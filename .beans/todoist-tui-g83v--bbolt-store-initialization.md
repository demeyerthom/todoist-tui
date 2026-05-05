---
# todoist-tui-g83v
title: bbolt store initialization
status: todo
type: feature
priority: normal
created_at: 2026-05-03T14:58:17Z
updated_at: 2026-05-03T15:11:25Z
parent: todoist-tui-j3br
---

Initialize bbolt embedded KV store (go.etcd.io/bbolt) for local storage. Stores projects, tasks, labels, sections, filters, and sync state (sync_token). Provides the persistence layer that all read operations will use.



## Task Dependency Graph

1. `todoist-tui-umza` Define domain model structs ← F1-dir-structure
2. `todoist-tui-p2ij` bbolt Store Open/Close ← umza
3. `todoist-tui-xbbs` Bucket defs and CRUD helpers ← p2ij
4. `todoist-tui-cn66` SyncToken helpers ← p2ij (parallel with xbbs)
5. `todoist-tui-y2ze` Store tests ← xbbs + cn66

## Cross-feature dependency
- Depends on F1 (project init) completing first
- F4 (sync client) and F5 (TUI) depend on this feature
