---
# todoist-tui-jop7
title: Config loading tests
status: todo
type: task
created_at: 2026-05-03T15:08:53Z
updated_at: 2026-05-03T15:08:53Z
parent: todoist-tui-b52a
blocked_by:
    - todoist-tui-kem6
---

## Description

Table-driven tests for config loading, default generation, and validation.

## Requirements

- Test DefaultConfig() returns all expected values
- Test Load() with a valid TOML file
- Test Load() missing file triggers default creation
- Test Validate() with empty token returns ErrNoToken
- Test Validate() with valid token returns nil
- Use t.TempDir() for test isolation

## Acceptance Criteria

- All tests pass with 'go test ./internal/config/...'
- Coverage of Load, DefaultConfig, Validate, ConfigPath
