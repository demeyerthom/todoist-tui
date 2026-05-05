---
# todoist-tui-m0uv
title: Quit command
status: todo
type: feature
priority: normal
created_at: 2026-05-03T14:58:20Z
updated_at: 2026-05-03T15:11:53Z
parent: todoist-tui-j3br
---

Implement quit via :q command mode and Ctrl+C. Establishes the Vim-style modal keybinding foundation (normal/insert/command modes) and command parsing infrastructure.



## Task Dependency Graph

1. `todoist-tui-8czd` Mode system definition ← F5-app-model
2. `todoist-tui-nzne` Command mode entry/exit/parsing ← 8czd
3. `todoist-tui-5lz9` Ctrl+C quit in normal mode ← 8czd
4. `todoist-tui-on0s` Graceful shutdown ← 8czd + F3-store

## Cross-feature dependencies
- Depends on F3 (store) and F5 (TUI model)
- Must be implemented last in M0
