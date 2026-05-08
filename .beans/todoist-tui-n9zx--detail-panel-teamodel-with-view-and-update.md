---
# todoist-tui-n9zx
title: Detail panel tea.Model with View and Update
status: todo
type: task
created_at: 2026-05-07T19:37:13Z
updated_at: 2026-05-07T19:37:13Z
parent: todoist-tui-ys80
blocked_by:
    - todoist-tui-m8uj
    - todoist-tui-nfjx
---

## Description

Create `internal/ui/detail/model.go` and `internal/ui/detail/view.go` with the `Model` struct, `NewModel()`, `Init()`, `Update()`, and `View()` methods.

## Requirements

### `Model` struct fields
- `store *store.Store` — reads task details, project/section/label names, subtask count
- `cfg *config.Config` — theme access
- `styles *theme.Styles` — lipgloss style definitions
- `width int`, `height int` — panel dimensions
- `taskID string` — currently selected task ID

### `Update(msg tea.Msg) (tea.Model, tea.Cmd)`
- `TaskSelectedMsg{ID}`: set `taskID = id`
- `SyncCompleteMsg`: re-read task from store if `taskID` non-empty
- `tea.WindowSizeMsg`: update width/height

### `View() string`
- No task selected (`taskID == ""`): render "Select a task to view details" centered with `styles.MutedText()`
- Task selected: render key-value pairs vertically:
  - **Content**: task content (bold, `styles.Header()`)
  - **Description**: task description, or "none" with `styles.MutedText()` if empty
  - **Project**: project name (lookup by `task.ProjectID` from store, "Inbox" if not found)
  - **Section**: section name (lookup by `task.SectionID` from store, "none" if empty)
  - **Priority**: "P1"/"P2"/"P3"/"P4" with theme color (`TaskPriorityHigh` etc.)
  - **Due**: formatted via `theme.FormatDueDate(task.Due, now)`
  - **Labels**: comma-separated label names, or "none" if empty
  - **Subtasks**: count of tasks with `ParentID == taskID` from store

### Layout
- Each key-value pair rendered as "Key: Value" on one line
- Key column right-aligned, value left-aligned
- Fits within width/height bounds

## Acceptance Criteria

- `View()` renders key-value pairs correctly for a selected task
- Empty fields show appropriate placeholder ("none")
- No task selected shows placeholder message
- Project/section names looked up from store
- Subtask count computed from store
- `go build ./...` succeeds
