---
# todoist-tui-rcpd
title: 'Main panel: task list with compact view'
status: todo
type: feature
priority: normal
created_at: 2026-05-03T14:58:24Z
updated_at: 2026-05-07T19:38:29Z
parent: todoist-tui-96st
blocked_by:
    - todoist-tui-0rtm
---

Read-only main panel showing task list for the selected project or filter. Compact view showing task name, due date, priority indicator, and labels. Data read from bbolt store populated by sync.

## Design decisions

- **Library**: `bubble-table` (evertras) for rendering.
- **Sub-panel architecture**: Task list is its own `tea.Model`. Emits `TaskSelectedMsg` with task ID on selection change.
- **Data access**: Holds `*store.Store` reference. Reads tasks via `ListTasks()` + in-memory filtering by selected project/label.
- **Columns**: `[P] Task Name | Due | Labels` — priority dot first, task name second (widest), due date third, labels fourth.
- **Priority indicator**: Color-coded dot (red=P1, orange=P2, blue=P3, grey=P4).
- **Due date format**: Smart — relative for near dates ("Today", "Tomorrow", "Overdue 2d"), abbreviated weekday+date for this week ("Mon 10"), full date for further out ("May 10").
- **Overdue**: Date text colored with theme error/urgent color only (no prefix/markers).
- **Section grouping**: Insert styled separator rows into the table for section headers (`─── Section Name ───`), non-selectable. Unsectioned tasks get `─── No Section ───`.
- **Completed tasks**: Hidden by default, toggle with `H` key in normal mode.
- **Empty state**: Shows "No tasks in Project Name" with the selected project/filter name.
- **Re-reading**: Re-reads from store on `ProjectSelectedMsg`, `LabelSelectedMsg`, and `SyncCompleteMsg`.
