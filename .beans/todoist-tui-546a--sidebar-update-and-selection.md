---
# todoist-tui-546a
title: Sidebar Update and selection
status: todo
type: task
created_at: 2026-05-07T19:35:53Z
updated_at: 2026-05-07T19:35:53Z
parent: todoist-tui-6ss5
blocked_by:
    - todoist-tui-m6qd
---

## Description

Create `internal/ui/sidebar/update.go` with the `Update()` method handling key events and selection messages.

## Requirements

### `Update(msg tea.Msg) (tea.Model, tea.Cmd)` method on `Model`

#### Key handling (normal mode, using keymap)
- `j`/`down`: move cursor down one visible item (clamp to last item)
- `k`/`up`: move cursor up one visible item (clamp to 0)
- `Enter`:
  - On section header: toggle `expandedSections` entry, rebuild items
  - On expandable project: toggle `expandedProjects` entry, rebuild items
  - On leaf project: emit `ProjectSelectedMsg{ID: item.ID}`
  - On filter: emit `FilterSelectedMsg{ID: item.ID, Query: filter.Query}`
  - On label: emit `LabelSelectedMsg{Name: item.Name}`
- `l`/`expand`: expand current item if expandable and collapsed
- `h`/`collapse`: collapse current item if expandable and expanded
- `g`/`go_top`: move cursor to first visible item
- `G`/`go_bottom`: move cursor to last visible item

#### Message handling
- `SyncCompleteMsg`: call `loadItems()`, preserve cursor position if valid
- `tea.WindowSizeMsg`: update width/height

### Cursor bounds
- Cursor always stays within `[0, len(items)-1]` range
- Empty items list: cursor stays at 0

## Acceptance Criteria

- Cursor moves correctly with j/k/g/G
- Enter toggles sections/subtrees and selects leaf items
- h/l collapse/expand correctly
- Selection messages emitted with correct data
- `go build ./...` succeeds
