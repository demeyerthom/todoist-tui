---
# todoist-tui-m0uv
title: Quit command
status: completed
type: feature
priority: normal
created_at: 2026-05-03T14:58:20Z
updated_at: 2026-05-06T12:05:52Z
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

## Summary of Changes

Implemented quit via :q command mode and Ctrl+C in internal/ui/:

1. **keymap/keymap.go** — KeyMap struct with Normal/Insert/Command sub-maps, DefaultKeyMap() with vim-style bindings, KeyFor() lookup method
2. **app.go** — Command mode entry (: enters ModeCommand), keystroke accumulation, Enter executes (:q/:quit → tea.Quit), Esc cancels, Backspace deletes char or cancels
3. **app.go** — Ctrl+C handling: ModeNormal/ModeInsert → tea.Quit, ModeCommand → cancel back to normal
4. **app.go** — Command bar rendered at bottom of View() in ModeCommand using Theme.CommandBar color
5. **app.go** — Cleanup() method for graceful store.Close()
6. **main.go** — Full entrypoint: config loading, store init, sync client, Bubbletea program with alt screen, Cleanup() on exit

All 4 tasks completed. go build, go vet, and go test pass.
