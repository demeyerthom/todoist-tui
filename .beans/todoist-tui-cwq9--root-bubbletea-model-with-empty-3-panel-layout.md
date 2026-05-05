---
# todoist-tui-cwq9
title: Root Bubbletea model with empty 3-panel layout
status: todo
type: feature
priority: normal
created_at: 2026-05-03T14:58:19Z
updated_at: 2026-05-03T15:11:39Z
parent: todoist-tui-j3br
---

Create the root Bubbletea model (internal/ui/app.go) implementing the Elm architecture. Empty 3-panel layout: sidebar, main task list, detail panel. Uses charmbracelet/bubbletea for the framework, lipgloss for layout/styling, bubbles for base components.



## Task Dependency Graph

1. `todoist-tui-7h9s` App model struct and constructor ← F1-dir-structure
2. `todoist-tui-qgud` 3-panel layout rendering ← 7h9s
3. `todoist-tui-yhti` Panel focus and Tab cycling ← 7h9s
4. `todoist-tui-szfw` Init command triggers full sync ← 7h9s + F4-full-sync
5. `todoist-tui-a5us` Placeholder panel content ← qgud

## Cross-feature dependencies
- Depends on F1 (project init), F2 (config), F3 (store), F4 (sync)
- F6 (quit) depends on this feature
