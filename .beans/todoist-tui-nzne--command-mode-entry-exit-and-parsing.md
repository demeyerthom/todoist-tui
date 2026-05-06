---
# todoist-tui-nzne
title: Command mode entry, exit, and parsing
status: completed
type: task
priority: normal
created_at: 2026-05-03T15:10:43Z
updated_at: 2026-05-06T12:00:45Z
parent: todoist-tui-m0uv
blocked_by:
    - todoist-tui-8czd
---

## Description

Implement command mode: ':' enters command mode, Esc cancels, Enter executes. Parse :q and :quit.

## Requirements

- ':' key in ModeNormal → set mode to ModeCommand, show ':' prompt
- In ModeCommand: accumulate keystrokes into command string
- Enter → parse command string, execute, return to ModeNormal
- Esc → clear command string, return to ModeNormal
- Backspace → delete last char (or Esc if empty)
- Command string displayed at bottom of View() like Vim's command line
- Supported commands: ':q' and ':quit' → tea.Quit
- Unknown commands → show error message briefly, return to ModeNormal

## Acceptance Criteria

- ':' switches to command mode and shows ':' prompt
- Enter executes parsed command
- ':q' and ':quit' quit the app
- Esc cancels command mode
- Unknown commands show error feedback

## Summary of Changes

Implemented in internal/ui/app.go:

- Added `commandBuf string` field to Model struct
- Added `fmt` import for error formatting
- NewModel() initializes commandBuf to ""
- Mode-aware key handling in Update():
  - ModeNormal: ':' enters command mode, Ctrl+C quits, Tab/Shift+Tab cycles panels
  - ModeInsert: Ctrl+C quits
  - ModeCommand: Enter parses/executes commands, Esc cancels, Backspace deletes char (or cancels if empty), Ctrl+C cancels to normal, printable chars accumulate
- Supported commands: :q and :quit quit via tea.Quit
- Unknown commands set m.err with formatted error message
- Command bar displayed at bottom in View() using CommandBar theme color, with adjusted panel height
- panelsView() now accepts panelHeight parameter

go build, go vet, and go test all pass.

## Summary of Changes

Updated internal/ui/app.go with:
- commandBuf string field on Model for accumulating command input
- ':' key in ModeNormal enters ModeCommand and clears commandBuf
- In ModeCommand: Enter parses commandBuf (q/quit → tea.Quit, else error)
- Esc cancels command mode, Backspace deletes last char or cancels if empty
- Printable runes appended to commandBuf in ModeCommand
- Command bar rendered at bottom of View() using Theme.CommandBar color
- Unknown commands set m.err with error message and return to ModeNormal
- go build and go vet pass
