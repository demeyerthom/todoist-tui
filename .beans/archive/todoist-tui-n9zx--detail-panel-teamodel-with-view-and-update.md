---
# todoist-tui-n9zx
title: Detail panel tea.Model with View and Update
status: completed
type: task
priority: normal
created_at: 2026-05-07T19:37:13Z
updated_at: 2026-05-09T14:24:53Z
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

## Summary of Changes

Created the detail panel Bubbletea model (`internal/ui/detail/model.go`) and view (`internal/ui/detail/view.go`), replacing the placeholder `doc.go`.

**model.go** — Model struct with `store`, `cfg`, `styles`, `width`, `height`, `taskID` fields. `NewModel` constructor and `Init()` following the sidebar/tasklist pattern.

**view.go** — `Update()` handles `TaskSelectedMsg`, `SyncCompleteMsg`, `WindowSizeMsg`. `View()` renders task details with key-value pairs: Content (header), Description, Project (store lookup), Section (store lookup), Priority (P1-P4 with theme colors), Due (formatted via `theme.FormatDueDate`), Labels (comma-separated names from store), Subtasks (count from store via list iteration). Empty states use centered "Select a task to view details" / "Task not found" / "none" placeholders.

Build, vet, and tests all pass.

## Summary of Changes

Created internal/ui/detail/model.go with Model struct (store, cfg, styles, width, height, taskID fields), NewModel constructor, and Init() method. Created internal/ui/detail/view.go with Update() handling TaskSelectedMsg, SyncCompleteMsg, and WindowSizeMsg; View() rendering key-value pairs for selected task (Content, Description, Project, Section, Priority, Due, Labels, Subtasks) with store lookups, theme colors, and placeholder text for empty fields. Removed doc.go placeholder. go build and go vet pass.
