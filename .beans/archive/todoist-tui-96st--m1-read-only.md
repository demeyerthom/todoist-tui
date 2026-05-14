---
# todoist-tui-96st
title: 'M1: Read-only'
status: completed
type: milestone
priority: normal
created_at: 2026-05-03T14:58:08Z
updated_at: 2026-05-09T15:03:22Z
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

## Summary of Changes

M1: Read-only milestone complete. All 5 features implemented:

1. **M1 Shared Foundation** — 7 message types, theme package with colors/styles/format
2. **Sidebar panel** — projects/filters/labels navigation with expand/collapse, color dots, selection
3. **Task list panel** — bubble-table with section grouping, priority colors, due dates, filtering
4. **Detail panel** — task detail view with store lookups for project/section/label names
5. **Navigation** — root model integration, cross-panel message routing, responsive layout (stacks below 80 cols), loading state
6. **Periodic sync** — 30-second background incremental sync, resolveTempIDs wiring, concurrent sync prevention

The app now has a fully functional read-only 3-panel TUI with sidebar navigation, task list, detail view, and periodic sync.
