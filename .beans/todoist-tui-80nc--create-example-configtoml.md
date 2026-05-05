---
# todoist-tui-80nc
title: Create example config.toml
status: todo
type: task
created_at: 2026-05-03T15:08:11Z
updated_at: 2026-05-03T15:08:11Z
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
