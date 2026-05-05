---
# todoist-tui-5lz9
title: Ctrl+C quit in normal mode
status: todo
type: task
created_at: 2026-05-03T15:10:47Z
updated_at: 2026-05-03T15:10:47Z
parent: todoist-tui-m0uv
blocked_by:
    - todoist-tui-8czd
---

## Description

Handle Ctrl+C in normal mode to quit the app immediately.

## Requirements

- In ModeNormal: Ctrl+C sends tea.Quit
- In ModeInsert: Ctrl+C also quits (consistent with vim-like behavior)
- In ModeCommand: Ctrl+C cancels command and returns to ModeNormal (not quit)

## Acceptance Criteria

- Ctrl+C in normal mode quits
- Ctrl+C in insert mode quits
- Ctrl+C in command mode cancels back to normal mode
