---
# todoist-tui-zpjy
title: Task list View
status: todo
type: task
created_at: 2026-05-07T19:36:45Z
updated_at: 2026-05-07T19:36:45Z
parent: todoist-tui-rcpd
blocked_by:
    - todoist-tui-r5ys
    - todoist-tui-079h
---

## Description

Create `internal/ui/tasklist/view.go` with the `View()` method using `bubble-table` to render the task list.

## Requirements

### `View() string` method on `Model`

#### Table setup
- Uses `github.com/evertras/bubble-table` to render
- Columns: `[P]` (width 2), `Task` (flex), `Due` (width 12), `Labels` (width 20)
- Column widths calculated from panel width
- Table fits within `width` and `height` bounds

#### Row rendering
- Each `SectionGroup` rendered as:
  - Section separator row (non-selectable, styled header)
  - Task rows from `group.Tasks`
- Rows built from `taskRow` wrapping each task
- Selected row highlighted with `cfg.Theme.SelectedRow`

#### Empty state
- When `groups` is empty or all groups have zero tasks:
  - Render "No tasks in {project/filter/label name}" centered with `styles.MutedText()`
  - Name derived from `currentFilter` context

#### Layout
- Table renders within available width/height
- Scrollbar shown if content overflows

## Acceptance Criteria

- `View()` renders bubble-table with correct columns
- Section separators visible between groups
- Selected row highlighted
- Empty state shows contextual message
- `go build ./...` succeeds
