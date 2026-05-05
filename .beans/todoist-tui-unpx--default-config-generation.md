---
# todoist-tui-unpx
title: Default config generation
status: completed
type: task
priority: normal
created_at: 2026-05-03T15:08:37Z
updated_at: 2026-05-05T17:57:11Z
parent: todoist-tui-b52a
blocked_by:
    - todoist-tui-a0b4
---

## Description

Implement DefaultConfig() returning a fully populated Config with sensible defaults, and a function to write it to a file path.

## Requirements

- DefaultConfig() *Config returns hardcoded defaults matching config.toml.example
- WriteDefaultConfig(path string) error writes defaults to file (creates parent dirs if needed)
- Defaults match the vim-style keybinding plan from PLAN.md
- Does NOT overwrite existing config

## Acceptance Criteria

- DefaultConfig() returns all fields with values
- WriteDefaultConfig creates dirs and writes valid TOML
- Existing files are not overwritten

## Summary of Changes\n\nCreated internal/config/defaults.go with DefaultConfig() returning fully populated Config matching config.toml.example, and WriteDefaultConfig(path) that creates parent dirs, writes valid TOML, and does NOT overwrite existing files. Created internal/config/defaults_test.go with 4 tests covering field population, TOML round-trip, file creation with dirs, and no-overwrite guard.
