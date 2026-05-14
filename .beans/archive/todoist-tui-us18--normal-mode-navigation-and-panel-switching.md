---
# todoist-tui-us18
title: Normal mode navigation and panel switching
status: completed
type: feature
priority: normal
created_at: 2026-05-03T14:58:26Z
updated_at: 2026-05-09T15:03:15Z
parent: todoist-tui-96st
blocked_by:
    - todoist-tui-6ss5
    - todoist-tui-rcpd
    - todoist-tui-ys80
---

Vim-style normal mode navigation: j/k for up/down movement, Tab to switch between panels, g/G for top/bottom, 1/2/3 to jump directly to sidebar/main/detail panel. Part of the modal keybinding system (internal/ui/keymap).

## Design decisions

- **Keybindings**: `j/k` up/down within panel, `Tab/Shift+Tab` cycle panel focus, `g/G` top/bottom, `1/2/3` jump to sidebar/main/detail, `h/l` collapse/expand in sidebar, `Enter` toggle sections/subtrees or select items, `H` toggle completed task visibility in task list.
- **Focus indication**: Border color change only (existing behavior in `app.go`).
- **Layout**: Proportional widths (20/50/30) with minimum widths per panel. Below 80 columns terminal width, panels stack vertically — focused panel gets full height, unfocused panels collapse to header strips.
- **Root model delegation**: Root `Update()` delegates to focused panel. Cross-panel messages (`ProjectSelectedMsg`, `TaskSelectedMsg`, etc.) routed by root model to affected panels.

## Summary of Changes

All 3 navigation tasks completed: root model integration and keymap, responsive layout with stacking threshold, loading state on startup.
