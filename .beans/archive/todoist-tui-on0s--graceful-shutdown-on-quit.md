---
# todoist-tui-on0s
title: Graceful shutdown on quit
status: completed
type: task
priority: normal
created_at: 2026-05-03T15:10:56Z
updated_at: 2026-05-06T12:05:46Z
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

## Summary of Changes

Updated internal/ui/app.go with Cleanup() method that calls store.Close().
Rewrote cmd/todoist-tui/main.go as full entrypoint:
- Loads config and validates auth token
- Opens bbolt store at XDG-compliant path
- Creates sync client with 30s timeout
- Creates ui.Model and runs Bubbletea program with alt screen
- Calls model.Cleanup() after p.Run() returns to close bbolt
- Error handling on both success and failure paths
go build, go vet, and go test pass.
