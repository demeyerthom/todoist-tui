---
# todoist-tui-umza
title: Define domain model structs
status: completed
type: task
priority: normal
created_at: 2026-05-03T15:09:02Z
updated_at: 2026-05-05T18:34:28Z
parent: todoist-tui-g83v
blocked_by:
    - todoist-tui-4yfu
---

## Description

Define Task, Project, Section, Label, Filter, and SyncMeta structs in internal/model/ with JSON tags matching the Todoist Sync API v9 response format.

## Requirements

- Task struct: id, content, description, project_id, section_id, parent_id, labels []string, priority int, due *DueDate, completed bool, checked bool, is_deleted bool, added_at, completed_at string, user_id
- Project struct: id, name, color, is_favorite, is_inbox, view_style, is_deleted, is_archived, user_id
- Section struct: id, name, project_id, section_order, is_deleted, is_archived
- Label struct: id, name, color, is_favorite, is_deleted, item_order
- Filter struct: id, name, color, query, is_deleted
- SyncMeta struct: key (string, e.g. 'sync_token'), value (string)
- DueDate struct: date, is_recurring, datetime, timezone
- All fields have 'json' tags matching Sync API v9 snake_case names (e.g. `json:"project_id"`)
- Each type has a Bucket() string method returning the bbolt bucket name (e.g. Task.Bucket() returns "tasks")

## Examples

```go
type Task struct {
    ID          string   `json:"id"`
    Content     string   `json:"content"`
    ProjectID   string   `json:"project_id"`
    // ...
}

func (Task) Bucket() string { return "tasks" }
```

## Acceptance Criteria

- All structs defined in internal/model/
- JSON tags match Sync API v9 field names
- Bucket() method on each type
- 'go build ./...' succeeds

## Summary of Changes\n\nImplemented all domain model structs in internal/model/:\n- Task (with DueDate nested struct) in task.go\n- Project in project.go\n- Section in section.go\n- Label in label.go\n- Filter in filter.go\n- SyncMeta in sync_meta.go\n\nAll structs have JSON tags matching Sync API v9 snake_case field names and Bucket() methods returning bbolt bucket names. DueDate.Bucket() returns "tasks" since it's always stored as part of a Task. Existing doc.go preserved. Build and vet pass cleanly.
