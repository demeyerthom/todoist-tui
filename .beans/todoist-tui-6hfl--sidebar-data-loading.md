---
# todoist-tui-6hfl
title: Sidebar data loading
status: todo
type: task
created_at: 2026-05-07T19:35:31Z
updated_at: 2026-05-07T19:35:31Z
parent: todoist-tui-6ss5
blocked_by:
    - todoist-tui-4tsx
---

## Description

Create `internal/ui/sidebar/data.go` with the `loadItems()` method that reads from the store and builds the flat `[]SidebarItem` list with correct ordering, hierarchy, and collapse state.

## Requirements

### `loadItems()` method on `Model`
- Calls `store.ListProjects()`, `store.ListFilters()`, `store.ListLabels()`
- Excludes soft-deleted items (`IsDeleted == true`)
- Filters excluded from `IsArchived` if applicable

### Projects section
- Section header row: "Projects" with count of visible projects
- Inbox project pinned at top (sorted by `IsInbox` descending)
- Remaining projects sorted alphabetically by name
- Build tree from `ParentID` — top-level projects have `Indent: 0`, children `Indent: 1`, etc.
- Respects `expandedProjects` map — collapsed subtrees only show the parent project row
- Project items have `Kind: "project"`, `Expandable: true` if they have children

### Filters section
- Section header row: "Filters" with count
- Sorted by name
- Each filter item has `Kind: "filter"`

### Labels section
- Section header row: "Labels" with count
- Sorted by name
- Each label item has `Kind: "label"`

### Collapse state
- Respects `expandedSections` — collapsed section only shows the header row
- Respects `expandedProjects` — collapsed subtree only shows the parent row
- Items list is rebuilt each call (not incrementally updated)

## Acceptance Criteria

- `loadItems()` returns correct `[]SidebarItem` for given store data
- Inbox always appears first in projects section
- Nested projects appear with correct indent
- Collapsed sections/subtrees hide children
- `go build ./...` succeeds
