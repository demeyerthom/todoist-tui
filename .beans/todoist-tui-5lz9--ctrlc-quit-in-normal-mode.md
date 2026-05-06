---
# todoist-tui-5lz9
title: Ctrl+C quit in normal mode
status: completed
type: task
priority: normal
created_at: 2026-05-03T15:10:47Z
updated_at: 2026-05-06T12:00:48Z
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

## Summary of Changes

Implemented in internal/ui/app.go Update() function:

- ModeNormal: Ctrl+C returns tea.Quit
- ModeInsert: Ctrl+C returns tea.Quit
- ModeCommand: Ctrl+C cancels command, clears commandBuf, returns to ModeNormal (does NOT quit)

go build, go vet, and go test all pass.

## Summary of Changes

Updated internal/ui/app.go with:
- Ctrl+C in ModeNormal → tea.Quit (immediate exit)
- Ctrl+C in ModeInsert → tea.Quit (consistent vim-like behavior)
- Ctrl+C in ModeCommand → cancel command mode, return to ModeNormal (not quit)
- go build and go vet pass
