---
# todoist-tui-8czd
title: Mode system definition
status: todo
type: task
created_at: 2026-05-03T15:10:37Z
updated_at: 2026-05-03T15:10:37Z
parent: todoist-tui-m0uv
blocked_by:
    - todoist-tui-7h9s
---

## Description

Define the Mode enum and KeyMap type for vim-style modal keybindings in internal/ui/keymap/.

## Requirements

- Mode type: ModeNormal, ModeInsert, ModeCommand (int enum with String() method)
- KeyMap struct with Normal, Insert, Command sub-maps (map[string]string for now, extensible)
- DefaultKeyMap() returns the default vim-style bindings from PLAN.md
- Mode field added to app.Model (if not already present from F5)

## Acceptance Criteria

- Mode type defined as int with String() method
- KeyMap struct defined with sub-maps for each mode
- DefaultKeyMap() returns planned bindings
- Model.mode field exists
