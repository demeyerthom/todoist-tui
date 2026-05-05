---
# todoist-tui-unpx
title: Default config generation
status: todo
type: task
created_at: 2026-05-03T15:08:37Z
updated_at: 2026-05-03T15:08:37Z
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
