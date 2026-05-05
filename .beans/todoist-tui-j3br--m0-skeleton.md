---
# todoist-tui-j3br
title: 'M0: Skeleton'
status: in-progress
type: milestone
priority: normal
created_at: 2026-05-03T14:58:07Z
updated_at: 2026-05-03T14:59:18Z
---

Project scaffolding, config, storage, sync client, and empty TUI shell.



## Key Decisions
- Auth: Personal API token only
- API: Todoist Sync API v9
- Language: Go 1.25 minimum
- Local storage: bbolt (embedded KV store)
- Layout: Sidebar + Main + Detail (3-panel)
- Keybindings: Vim-style modal (normal/insert/command)
- Config format: TOML
- Deps: Full Charm stack (bubbletea, lipgloss, bubbles) + evertras/bubble-table + bbolt + go-toml/v2 + google/uuid
