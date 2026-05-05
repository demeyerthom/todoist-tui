---
# todoist-tui-nzne
title: Command mode entry, exit, and parsing
status: todo
type: task
created_at: 2026-05-03T15:10:43Z
updated_at: 2026-05-03T15:10:43Z
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
