---
# todoist-tui-r5ys
title: Task list data loading and filtering
status: completed
type: task
priority: normal
created_at: 2026-05-07T19:36:19Z
updated_at: 2026-05-08T20:51:23Z
parent: todoist-tui-rcpd
blocked_by:
    - todoist-tui-nfjx
---

## Description

Create `internal/ui/tasklist/data.go` with `loadTasks()` that reads tasks from the store and applies the current filter.

## Requirements

### `loadTasks()` method on `Model`

#### Filtering
- Calls `store.ListTasks()`, `store.ListSections()`
- If `currentFilter.ProjectID` non-empty: include only tasks with matching `ProjectID`
- If `currentFilter.LabelName` non-empty: include only tasks where `Labels` contains the name
- If `currentFilter.IsFilter`: return empty groups (unsupported in M1)
- Exclude tasks with `IsDeleted == true`
- If `showCompleted == false`: exclude tasks with `Checked == true`

#### Grouping
- Group tasks by `SectionID`
- Build section lookup map from `store.ListSections()`
- Unsectioned tasks (empty `SectionID`) go into a "No Section" group

#### Sorting
- Sort sections by `SectionOrder` ascending
- Within each section group, sort tasks by:
  - Priority ascending (P1 first, P4 last)
  - Then by due date ascending (earliest first, no due date last)
- Unsectioned group always last

#### Output
- Sets `m.groups` to the loaded `[]SectionGroup`
- Rebuilds bubble-table rows from groups

## Acceptance Criteria

- `loadTasks()` returns correct groups for given filter
- Completed tasks excluded when `showCompleted == false`
- Deleted tasks always excluded
- Sections sorted by `SectionOrder`, unsectioned last
- Tasks sorted by priority then due date within each group
- `go build ./...` succeeds

## Summary of Changes

Created internal/ui/tasklist/data.go with loadTasks() method (pointer receiver). Loads tasks and sections from store, filters by ProjectID/LabelName/IsFilter/IsDeleted/Checked, groups tasks by SectionID with section lookup map, sorts sections by SectionOrder (unsectioned last), sorts tasks within each group by priority ascending then due date ascending (nil due dates last). Includes dueDateLess helper for comparing DueDate pointers. go build and go vet pass.
