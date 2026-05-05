---
# todoist-tui-kem6
title: Load and parse config from XDG path
status: completed
type: task
priority: normal
created_at: 2026-05-03T15:08:49Z
updated_at: 2026-05-05T17:58:45Z
parent: todoist-tui-b52a
blocked_by:
    - todoist-tui-unpx
---

## Description

Implement Load() that finds, reads, and parses config.toml from the XDG config directory, with validation.

## Requirements

- ConfigPath() string returns ~/.config/todoist-tui/config.toml (using os.UserConfigDir)
- Load() (*Config, error) reads and parses the file
- If file missing, calls WriteDefaultConfig then loads
- Validate() error on Config checks Token is non-empty (returns ErrNoToken)
- Uses pelletier/go-toml/v2 for parsing

## Examples

```go
func Load() (*Config, error) {
    path := ConfigPath()
    if _, err := os.Stat(path); os.IsNotExist(err) {
        if err := WriteDefaultConfig(path); err != nil {
            return nil, fmt.Errorf("write default config: %w", err)
        }
    }
    // parse with go-toml/v2
}
```

## Acceptance Criteria

- ConfigPath() returns XDG-compliant path
- Load() parses valid TOML into Config struct
- Missing config file triggers default creation then loading
- Validation rejects empty token with ErrNoToken

## Summary of Changes\n\nCreated internal/config/load.go with ConfigPath() returning XDG-compliant path using os.UserConfigDir (with fallback to /home/thomas/.config), and Load() that reads/parses config.toml, creates defaults if missing, and validates the token.
