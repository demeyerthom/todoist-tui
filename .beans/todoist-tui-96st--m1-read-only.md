---
# todoist-tui-96st
title: 'M1: Read-only'
status: todo
type: milestone
priority: normal
created_at: 2026-05-03T14:58:08Z
updated_at: 2026-05-03T14:59:15Z
blocked_by:
    - todoist-tui-j3br
---

3-panel read-only TUI: sidebar navigation, task list, detail view with periodic sync.

## Architecture

- **Sub-panels**: Each is its own `tea.Model`. Root model delegates `Update()` to focused panel. Cross-panel communication via typed `tea.Msg` (`ProjectSelectedMsg`, `TaskSelectedMsg`, `SyncCompleteMsg`, etc.).
- **Data access**: Sub-panels hold `*store.Store` reference, read directly. Messages carry IDs only, panels look up from store.
- **Layout**: Proportional (20/50/30) with minimum widths. Below 80 columns, stack vertically — focused panel full height, others as header strips.

## Implementation order

1. **todoist-tui-6ss5** Sidebar — foundation, drives other panels
2. **todoist-tui-rcpd** Task list — reacts to sidebar selection
3. **todoist-tui-ys80** Detail panel — reacts to task selection
4. **todoist-tui-us18** Navigation — panels exist to navigate between
5. **todoist-tui-fykz** Periodic sync — independent, slots in anytime
