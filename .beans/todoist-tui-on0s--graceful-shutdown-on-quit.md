---
# todoist-tui-on0s
title: Graceful shutdown on quit
status: todo
type: task
created_at: 2026-05-03T15:10:56Z
updated_at: 2026-05-03T15:10:56Z
parent: todoist-tui-m0uv
blocked_by:
    - todoist-tui-8czd
    - todoist-tui-p2ij
---

## Description

Ensure bbolt store is cleanly closed and resources released when the app quits.

## Requirements

- On tea.Quit, call store.Close() to close bbolt DB
- Use bubbletea's tea.Quit return in Update()
- Store can be safely closed even if sync is in progress
- Wrap shutdown in model's Update handler

## Acceptance Criteria

- bbolt store is closed on quit
- No goroutine leaks on exit
- App exits cleanly with 'go build' + run
