---
# todoist-tui-jub6
title: Responsive layout with stacking threshold
status: completed
type: task
priority: normal
created_at: 2026-05-07T19:37:42Z
updated_at: 2026-05-09T14:52:23Z
parent: todoist-tui-us18
blocked_by:
    - todoist-tui-hwqz
---

## Description

Update `app.go` `panelsView()` to support responsive layout with proportional widths, minimum widths, and vertical stacking below 80 columns.

## Requirements

### `panelsView(panelHeight int) string` method

#### Horizontal layout (width >= 80)
- Proportional widths: sidebar 20%, main 50%, detail 30% of total width
- Minimum widths: sidebar min 20, main min 40, detail min 30
- If total width cannot satisfy minimums, reduce detail first, then main, then sidebar
- Panels rendered side-by-side with `lipgloss.JoinHorizontal(lipgloss.Top, ...)`
- Each panel wrapped in border style (existing `sidebarStyle`, `mainStyle`, `detailStyle`)

#### Vertical stacking (width < 80)
- Panels stacked vertically with `lipgloss.JoinVertical(lipgloss.Top, ...)`
- Focused panel gets `(panelHeight - 2)` lines (full height minus 2 for border)
- Unfocused panels collapse to 1-line header strips:
  - Sidebar: "Projects"
  - Main: "Tasks"
  - Detail: "Details"
- Header strips styled with `styles.Header()` and border

#### Border styles
- Active panel border: `cfg.Theme.ActiveBorder`
- Inactive panel border: `cfg.Theme.InactiveBorder`
- Border style: `lipgloss.RoundedBorder()`

## Acceptance Criteria

- Horizontal layout works correctly at width >= 80
- Vertical stacking activates at width < 80
- Focused panel gets full height in stacked mode
- Unfocused panels show 1-line header strips
- `go build ./...` succeeds

## Summary of Changes

Updated internal/ui/app.go: added stackWidthThreshold constant (80), styles *theme.Styles field to Model struct, panelWidths() helper for proportional width computation with minimum enforcement (sidebar min 20, main min 40, detail min 30), horizontalView() for side-by-side layout at width >= 80, stackedView() for vertical stacking at width < 80 with focused panel full height and unfocused panels as 1-line header strips (Projects/Tasks/Details). Updated WindowSizeMsg handler to branch on stacked vs horizontal mode for sub-panel dimension forwarding. go build, go vet, and go test all pass.
