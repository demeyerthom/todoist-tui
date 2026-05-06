---
# todoist-tui-j3br
title: 'M0: Skeleton'
status: completed
type: milestone
priority: normal
created_at: 2026-05-03T14:58:07Z
updated_at: 2026-05-06T12:08:21Z
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

## Summary of Changes

M0: Skeleton milestone complete. All 6 features implemented:

1. **F1: Project init and directory structure** — Go module, directory layout, dependencies
2. **F2: TOML config loading** — Config structs, XDG path, defaults, validation
3. **F3: bbolt store initialization** — Domain models, CRUD helpers, sync token persistence
4. **F4: Sync API client** — Full/incremental sync, temp ID resolution, 8 tests
5. **F5: Root Bubbletea model** — 3-panel layout, focus cycling, async sync, placeholder content
6. **F6: Quit command** — Command mode (:q/:quit), Ctrl+C, modal keybindings, graceful shutdown

All tests pass (27 total). go build, go vet clean.
