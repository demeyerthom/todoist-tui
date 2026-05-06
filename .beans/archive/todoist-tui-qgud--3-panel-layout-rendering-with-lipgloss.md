---
# todoist-tui-qgud
title: 3-panel layout rendering with lipgloss
status: completed
type: task
priority: normal
created_at: 2026-05-03T15:10:14Z
updated_at: 2026-05-06T11:26:02Z
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

## Summary of Changes

- Added `lipgloss` import to `internal/ui/app.go`
- Updated `Update()` to handle `tea.WindowSizeMsg`: stores terminal width/height on the model
- Rewrote `View()` to render a 3-panel horizontal layout using lipgloss:
  - Sidebar: ~20% width with rounded border and right border
  - Main: ~50% width with rounded border and right border
  - Detail: ~30% width with rounded border (remainder from rounding)
  - Panels joined via `lipgloss.JoinHorizontal(lipgloss.Top, ...)`
  - Each panel renders placeholder text ("Sidebar", "Tasks", "Detail")
  - Returns "Initializing..." when width/height are zero (before first WindowSizeMsg)
- `go build ./...`, `go vet ./...`, and `go test ./...` all pass

## Summary of Changes

Updated internal/ui/app.go with:
- 3-panel horizontal layout using lipgloss (sidebar 20%, main 50%, detail 30%)
- Rounded borders with right border on sidebar and main panels
- lipgloss.JoinHorizontal for panel layout
- WindowSizeMsg handling for responsive resize
- Width/height guard with 'Initializing...' fallback
- go build and go vet pass
