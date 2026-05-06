---
# todoist-tui-yhti
title: Panel focus tracking and Tab cycling
status: completed
type: task
priority: normal
created_at: 2026-05-03T15:10:18Z
updated_at: 2026-05-06T11:26:03Z
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

## Summary of Changes

Updated internal/ui/app.go with:
- Tab key cycles activePanel forward (sidebar → main → detail → sidebar)
- Shift+Tab cycles backward (detail → main → sidebar → detail)
- Active panel border uses cfg.Theme.ActiveBorder color
- Inactive panel borders use cfg.Theme.InactiveBorder color
- borderColor helper method on Model
- go build and go vet pass
