---
# todoist-tui-szfw
title: Init command triggers full sync
status: completed
type: task
priority: normal
created_at: 2026-05-03T15:10:21Z
updated_at: 2026-05-06T11:26:05Z
parent: todoist-tui-cwq9
blocked_by:
    - todoist-tui-7h9s
    - todoist-tui-hgg5
---

## Description

Implement Init() to return a tea.Cmd that triggers FullSync on startup.

## Requirements

- Init() returns a tea.Cmd that calls syncClient.FullSync(ctx, store)
- Use context.Background() for the initial sync
- Wrap sync result in a custom tea.Msg (SyncDoneMsg or SyncErrMsg)
- Update() handles SyncDoneMsg and SyncErrMsg
- On SyncDoneMsg, store is populated with data (ready for M1)
- On SyncErrMsg, set model.err and display error in View()

## Acceptance Criteria

- Init() kicks off full sync
- Success populates store data
- Error is displayed in the UI
- App doesn't hang during sync (async)

## Summary of Changes

Updated internal/ui/app.go with:
- SyncDoneMsg and SyncErrMsg custom tea.Msg types
- Init() returns async tea.Cmd calling syncClient.FullSync(context.Background(), store)
- Update() handles SyncDoneMsg (no-op) and SyncErrMsg (sets m.err)
- View() displays error banner with cfg.Theme.Error color when m.err is set
- panelsView() helper extracted for clean View() method
- go build and go vet pass
