---
# todoist-tui-ys80
title: 'Detail panel: task detail view'
status: todo
type: feature
priority: normal
created_at: 2026-05-03T14:58:25Z
updated_at: 2026-05-07T19:38:29Z
parent: todoist-tui-96st
blocked_by:
    - todoist-tui-0rtm
---

Read-only detail panel showing full information for the selected task: content, description, due date, priority, labels, project, section, subtask count. Updates when selection changes in the task list.

## Design decisions

- **Sub-panel architecture**: Detail panel is its own `tea.Model`. Receives `TaskSelectedMsg` with task ID, looks up full task from store.
- **Data access**: Holds `*store.Store` reference. Looks up task by ID on selection change. Re-reads on `SyncCompleteMsg`.
- **Layout**: Key-value pairs stacked vertically — `Project: Inbox`, `Due: 2026-05-10`, `Priority: P1`, `Labels: work, urgent`, etc.
- **No selection state**: Shows "Select a task to view details" as muted text.
