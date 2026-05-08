---
# todoist-tui-m6qd
title: Sidebar rendering
status: todo
type: task
created_at: 2026-05-07T19:35:42Z
updated_at: 2026-05-07T19:35:42Z
parent: todoist-tui-6ss5
blocked_by:
    - todoist-tui-6hfl
---

## Description

Create `internal/ui/sidebar/view.go` with the `View()` method and helper rendering functions for the sidebar panel.

## Requirements

### `View() string` method on `Model`
- Renders each `SidebarItem` in order with correct styling
- Items beyond visible height are not rendered (scrolling via cursor)
- Active item (at cursor position) styled with `styles.ActiveItem()`
- Inactive items styled with `styles.InactiveItem()`

### Section headers
- Bold text styled with `styles.Header()`
- Show count when collapsed: "Projects (12)"
- Show expand/collapse indicator: "▾ Projects (12)" when expanded, "▸ Projects (12)" when collapsed

### Project items
- Indentation via leading spaces: `Indent * 2` spaces
- Color dot prefix via `theme.TodoistColor(item.Color, fallback)` — colored bullet
- Inbox project: "● Inbox" with pinned indicator
- Expandable projects: "▾ ProjectName" or "▸ ProjectName" indicator
- Nested projects: indented under parent

### Filter and label items
- No color dot for filters
- Labels: colored dot via `theme.TodoistColor(item.Color, fallback)`
- No indentation (flat list)

### Scrolling
- Cursor position determines which item is at the top of the visible area
- Items before cursor are not rendered
- Only `height` items rendered maximum

## Acceptance Criteria

- `View()` renders correctly for various item lists
- Active item visually distinct
- Indentation correct for nested projects
- Color dots render for projects and labels
- `go build ./...` succeeds
