---
# todoist-tui-80nc
title: Create example config.toml
status: completed
type: task
priority: normal
created_at: 2026-05-03T15:08:11Z
updated_at: 2026-05-05T15:42:52Z
parent: todoist-tui-4wyv
blocked_by:
    - todoist-tui-4yfu
---

## Description

Create a root-level config.toml.example with documented defaults for all config fields.

## Requirements

- Sections: [auth], [keybindings], [theme]
- [auth]: token = '' (placeholder, documented as required)
- [keybindings]: documented defaults matching the vim-style modal keybinding plan from PLAN.md
- [theme]: documented color/style defaults
- Comments explaining each field

## Examples

```toml
[auth]
token = ''  # Todoist personal API token (required)

[keybindings.normal]
up = 'k'
down = 'j'
# ...

[theme]
active_border = 'blue'
# ...
```

## Acceptance Criteria

- config.toml.example exists at project root
- All planned config fields are present with comments
- File is valid TOML

## Summary of Changes\n\n- Created config.toml.example at project root\n- Sections: [auth], [keybindings.normal], [keybindings.insert], [keybindings.command], [theme]\n- All 16 normal-mode keybindings documented with defaults from PLAN.md\n- 18 theme color/style fields with inline comments\n- Valid TOML format

## Review Note\n\nReviewer found a naming inconsistency: config.toml.example uses  but downstream bean  used . Resolved by updating bean  to use  to match the example config.
