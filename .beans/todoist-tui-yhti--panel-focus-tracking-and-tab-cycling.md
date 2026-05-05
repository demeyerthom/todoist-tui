---
# todoist-tui-yhti
title: Panel focus tracking and Tab cycling
status: todo
type: task
created_at: 2026-05-03T15:10:18Z
updated_at: 2026-05-03T15:10:18Z
parent: todoist-tui-cwq9
blocked_by:
    - todoist-tui-7h9s
---

## Description

Track which panel is active (focused) and cycle focus with Tab. Highlight the active panel border.

## Requirements

- activePanel field on Model (PanelSidebar, PanelMain, PanelDetail)
- Tab key cycles: sidebar → main → detail → sidebar
- Shift+Tab cycles in reverse
- Active panel gets highlighted border style (e.g., bold or blue)
- Inactive panels get dimmed border style
- Update() handles tea.KeyMsg for Tab and Shift+Tab

## Acceptance Criteria

- Tab cycles active panel forward
- Shift+Tab cycles backward
- Active panel has visually distinct border
- Inactive panels have dimmed border
