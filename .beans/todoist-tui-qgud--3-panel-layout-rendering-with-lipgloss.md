---
# todoist-tui-qgud
title: 3-panel layout rendering with lipgloss
status: todo
type: task
created_at: 2026-05-03T15:10:14Z
updated_at: 2026-05-03T15:10:14Z
parent: todoist-tui-cwq9
blocked_by:
    - todoist-tui-7h9s
---

## Description

Implement View() to render a 3-panel horizontal layout using lipgloss: sidebar (20%), main (50%), detail (30%).

## Requirements

- Use lipgloss for styling, borders, and layout
- Sidebar: ~20% of terminal width, left-aligned
- Main panel: ~50% of terminal width
- Detail panel: ~30% of terminal width
- Borders between panels (right border on sidebar and main)
- Responsive: recalculate widths on tea.WindowSizeMsg
- Each panel renders placeholder text: 'Sidebar', 'Tasks', 'Detail'
- Panels are joined horizontally with lipgloss.JoinHorizontal

## Examples

```
┌─────────┬──────────────────────┬──────────────┐
│ Sidebar │        Tasks         │   Detail     │
│         │                      │              │
│         │                      │              │
└─────────┴──────────────────────┴──────────────┘
```

## Acceptance Criteria

- View() renders 3 side-by-side panels with borders
- Layout responds to terminal resize
- Panel widths approximate 20/50/30 split
