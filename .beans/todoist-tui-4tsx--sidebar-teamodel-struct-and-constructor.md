---
# todoist-tui-4tsx
title: Sidebar tea.Model struct and constructor
status: todo
type: task
created_at: 2026-05-07T19:35:18Z
updated_at: 2026-05-07T19:35:18Z
parent: todoist-tui-6ss5
blocked_by:
    - todoist-tui-m01d
---

## Description

Create `internal/ui/sidebar/model.go` defining the sidebar `Model` struct as a `tea.Model` and the `SidebarItem` type representing a renderable row.

## Requirements

### `Model` struct fields
- `store *store.Store` — reads projects, filters, labels
- `cfg *config.Config` — theme and keymap access
- `styles *theme.Styles` — lipgloss style definitions
- `width int`, `height int` — panel dimensions
- `cursor int` — current cursor position in visible items
- `expandedSections map[string]bool` — "projects", "filters", "labels" → expanded state
- `expandedProjects map[string]bool` — project ID → subtree expanded state
- `items []SidebarItem` — flat renderable item list (rebuilt on data load)

### `SidebarItem` struct fields
- `Kind` string: "project", "filter", "label", "section-header", "subtree-header"
- `ID` string — project/filter/label ID (empty for headers)
- `Name` string — display name
- `Color` string — Todoist color name (for projects/labels)
- `Indent` int — indentation level (0 for top-level, 1+ for nested)
- `IsInbox` bool — true for Inbox project
- `Expandable` bool — true for section headers and subtree headers
- `Expanded` bool — current expansion state

### Constructor and Init
- `NewModel(cfg *config.Config, store *store.Store, styles *theme.Styles) Model`
- `Init() tea.Cmd` returns nil

## Acceptance Criteria

- `internal/ui/sidebar/model.go` exists and compiles
- `NewModel()` returns properly initialized model
- `go build ./...` succeeds
