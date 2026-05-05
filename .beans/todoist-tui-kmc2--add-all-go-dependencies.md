---
# todoist-tui-kmc2
title: Add all Go dependencies
status: completed
type: task
priority: normal
created_at: 2026-05-03T15:08:07Z
updated_at: 2026-05-05T15:42:45Z
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

## Summary of Changes\n\n- Added all 7 dependencies via go get: bubbletea, lipgloss, bubbles, bubble-table, bbolt, go-toml/v2, google/uuid\n- go mod tidy ran clean (deps recorded in go.sum; will appear in go.mod require block once code imports them)\n- go build ./... succeeds

## Summary of Changes\n\n- Added all 7 dependencies as direct (not indirect) imports via blank imports in appropriate packages\n- Fixed bubble-table import path to github.com/evertras/bubble-table/table\n- Fixed bubbles import to github.com/charmbracelet/bubbles/textinput\n- go mod tidy runs clean, go build ./... succeeds, go vet ./... passes
