---
# todoist-tui-6ss5
title: 'Sidebar panel: projects, filters, labels navigation'
status: completed
type: feature
priority: normal
created_at: 2026-05-03T14:58:23Z
updated_at: 2026-05-09T15:03:13Z
parent: todoist-tui-96st
blocked_by:
    - todoist-tui-0rtm
---

Read-only sidebar panel listing projects (with color support), filters, and labels. Data read from bbolt store. Projects show hierarchy with sections. Filters use Todoist filter query syntax. Labels show all user labels.

## Design decisions

- **Sub-panel architecture**: Sidebar is its own `tea.Model`. Communicates selections to other panels via `ProjectSelectedMsg`, `FilterSelectedMsg`, `LabelSelectedMsg`.
- **Data access**: Holds `*store.Store` reference, reads directly. Re-reads on `SyncCompleteMsg`.
- **Sections**: Three collapsible sections — Projects (expanded by default), Filters (collapsed), Labels (collapsed). `Enter` toggles collapse on section headers.
- **Project hierarchy**: Load all projects at once via `ListProject()`, build tree in-memory from `ParentID`.
- **Project subtrees**: Collapsible, collapsed by default. `h` collapses, `l` expands, `Enter` toggles.
- **Inbox**: Pinned to top of projects section with separator before remaining alphabetically sorted projects.
- **Project colors**: Map Todoist color names to lipgloss colors via lookup table (~30 entries), render as colored dot prefix.
- **Labels**: Show label name + color dot only, no task count.
- **Selection**: Single-select only. Filters show "not yet supported" in main panel when selected (M1 scope).
- **Navigation**: `j/k` move cursor, `h/l` collapse/expand, `g/G` top/bottom, `Enter` selects leaf item.

## Summary of Changes

All 5 sidebar tasks completed: shared message types, theme package, sidebar tea.Model, sidebar data loading, sidebar rendering, sidebar Update and selection.
