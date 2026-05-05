---
# todoist-tui-kmc2
title: Add all Go dependencies
status: todo
type: task
created_at: 2026-05-03T15:08:07Z
updated_at: 2026-05-03T15:08:07Z
parent: todoist-tui-4wyv
blocked_by:
    - todoist-tui-4yfu
---

## Description

Add all 7 project dependencies via 'go get' and verify with 'go mod tidy'.

## Requirements

- Dependencies: charmbracelet/bubbletea, charmbracelet/lipgloss, charmbracelet/bubbles, evertras/bubble-table, go.etcd.io/bbolt, pelletier/go-toml/v2, google/uuid
- Run 'go mod tidy' after adding
- Verify all imports resolve

## Acceptance Criteria

- All 7 deps appear in go.mod with correct paths
- 'go mod tidy' runs clean
- 'go build ./...' still succeeds
