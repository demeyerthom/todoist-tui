---
# todoist-tui-a0b4
title: Define Config structs with TOML tags
status: todo
type: task
priority: normal
created_at: 2026-05-03T15:08:29Z
updated_at: 2026-05-05T15:39:56Z
parent: todoist-tui-b52a
blocked_by:
    - todoist-tui-4yfu
---

## Description

Define the Go structs for configuration in internal/config/config.go with TOML tags for serialization.

## Requirements

- Config struct containing: Auth (AuthConfig), Keymap (KeymapConfig), Theme (ThemeConfig)
- AuthConfig: Token string
- KeymapConfig: Normal, Insert, Command sub-structs with key-to-action mappings (map[string]string)
- ThemeConfig: ActiveBorder string, InactiveBorder string, etc.
- All fields have 'toml' tags
- Sentinel error: ErrNoToken for validation
- Config.Validate() error method

## Examples

```go
type Config struct {
    Auth   AuthConfig   `toml:"auth"`
    Keymap KeymapConfig `toml:"keybindings"`
    Theme  ThemeConfig  `toml:"theme"`
}

type AuthConfig struct {
    Token string `toml:"token"`
}
```

## Acceptance Criteria

- Structs defined in internal/config/config.go
- All fields have TOML tags
- ErrNoToken sentinel error defined
- Config.Validate() method checks token is non-empty
