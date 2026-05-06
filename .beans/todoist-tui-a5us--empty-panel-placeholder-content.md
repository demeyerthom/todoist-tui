---
# todoist-tui-a5us
title: Empty panel placeholder content
status: completed
type: task
priority: normal
created_at: 2026-05-03T15:10:25Z
updated_at: 2026-05-06T11:37:45Z
parent: todoist-tui-cwq9
blocked_by:
    - todoist-tui-qgud
---

## Description

Render styled placeholder text in each panel to visually distinguish them.

## Todo

- [ ] Add sidebarView helper method with centered headings
- [ ] Add mainView helper method with centered placeholder
- [ ] Add detailView helper method with centered placeholder
- [ ] Update panelsView() to use helpers and compute inner dimensions
- [x] Verify build and vet pass

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

## Summary of Changes

Added styled placeholder content to all three panels in `internal/ui/app.go`:

- **`sidebarView(width, height int)`** — renders centered headings (Projects, Filters, Labels) using `Header` theme color for visual distinction, placed vertically and horizontally centered via `lipgloss.Place`
- **`mainView(width, height int)`** — renders "Select a project or filter" centered in `MutedText` color
- **`detailView(width, height int)`** — renders "No task selected" centered in `MutedText` color
- **`panelsView()`** — updated to compute inner dimensions (minus border), then delegate to the helper methods

Build, vet, and all tests pass.

## Summary of Changes

Updated internal/ui/app.go with:
- sidebarView(width, height) — centered 'Projects', 'Filters', 'Labels' in Header theme color
- mainView(width, height) — centered 'Select a project or filter' in MutedText theme color
- detailView(width, height) — centered 'No task selected' in MutedText theme color
- lipgloss.Place for horizontal and vertical centering within panel bounds
- Inner dimensions (width-2, height-2) account for border characters
- go build and go vet pass
