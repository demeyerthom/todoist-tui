---
# todoist-tui-079h
title: Task list bubble-table rows
status: completed
type: task
priority: normal
created_at: 2026-05-07T19:36:33Z
updated_at: 2026-05-09T13:21:55Z
parent: todoist-tui-rcpd
blocked_by:
    - todoist-tui-m8uj
    - todoist-tui-nfjx
---

## Description

Create `internal/ui/tasklist/row.go` with the `taskRow` type that implements `bubble-table`'s `Row` interface and renders priority dot, task name, due date, and labels columns.

## Requirements

### `taskRow` struct
- Wraps `model.Task`
- Implements `bubble-table` Row interface: `Field(key string) string`
- Fields: "priority", "content", "due", "labels"

### Priority column
- Colored dot via `theme.TodoistColor()`:
  - P1 → `cfg.Theme.TaskPriorityHigh` (red)
  - P2 → `cfg.Theme.TaskPriorityMedium` (yellow)
  - P3 → `cfg.Theme.TaskPriorityLow` (blue)
  - P4 → `cfg.Theme.NormalText` (grey/no priority)
- Renders as "●" with appropriate color

### Task name column
- Plain text from `task.Content`
- Truncated to column width by bubble-table

### Due date column
- Uses `theme.FormatDueDate(task.Due, time.Now())`
- Overdue dates: styled with `cfg.Theme.TaskOverdue`
- Today's dates: styled with `cfg.Theme.TaskDueToday`
- No due date: empty string

### Labels column
- Comma-separated: `strings.Join(task.Labels, ", ")`
- Truncated by bubble-table if too wide

### Section separator rows
- Non-selectable rows with styled text: "─── Section Name ───"
- Unsectioned group: "─── No Section ───"
- Implemented as a special row type that bubble-table renders but doesn't allow selection on

## Acceptance Criteria

- `taskRow` implements bubble-table Row interface correctly
- Priority dots render with correct colors
- Due dates formatted correctly with smart format
- Section separator rows are non-selectable
- `go build ./...` succeeds

## Summary of Changes

Created internal/ui/tasklist/row.go with four functions: taskToRow (converts Task to bubble-table Row with priority dot, content, due date, labels, and hidden task_id/is_section metadata), sectionRow (creates section separator row with is_section:true metadata and SectionSep styling), buildRows (iterates SectionGroups to build flat row list), and defaultColumns (returns column definitions: priority(3), content(flex), due(12), labels(20)). Includes taskPriorityColor helper for P1-P4 priority color mapping. go build and go vet pass.
