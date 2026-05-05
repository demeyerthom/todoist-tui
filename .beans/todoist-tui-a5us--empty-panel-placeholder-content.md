---
# todoist-tui-a5us
title: Empty panel placeholder content
status: todo
type: task
created_at: 2026-05-03T15:10:25Z
updated_at: 2026-05-03T15:10:25Z
parent: todoist-tui-cwq9
blocked_by:
    - todoist-tui-qgud
---

## Description

Render styled placeholder text in each panel to visually distinguish them.

## Requirements

- Sidebar placeholder: centered text 'Projects', 'Filters', 'Labels' (static labels, no data yet)
- Main panel placeholder: centered text 'Select a project or filter' or 'No task selected'
- Detail panel placeholder: centered text 'No task selected'
- Use lipgloss styles for centering (width/height alignment)
- Placeholder content is replaced in M1 with real data

## Acceptance Criteria

- Each panel shows clearly distinguishable placeholder content
- Content is centered within panel bounds
- Layout is readable at 80-column terminal width
