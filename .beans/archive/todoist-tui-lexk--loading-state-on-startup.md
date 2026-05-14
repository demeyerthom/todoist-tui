---
# todoist-tui-lexk
title: Loading state on startup
status: completed
type: task
priority: normal
created_at: 2026-05-07T19:37:54Z
updated_at: 2026-05-09T14:56:20Z
parent: todoist-tui-us18
blocked_by:
    - todoist-tui-jub6
---

## Description

Update `app.go` to show "Loading..." in all panels while the initial sync is in progress.

## Requirements

### `Model` struct changes
- Add `synced bool` field (initially `false`)

### `View()` changes
- When `synced == false`:
  - Render "Loading..." centered in each panel area with `styles.MutedText()`
  - Still render borders around each panel
  - Error banner still shown at top if `m.err != nil`
  - Command bar still shown at bottom if in command mode

### `Update()` changes
- On `SyncDoneMsg`:
  - Set `synced = true`
  - Return init commands for each sub-panel (they can now render with data)
- On `SyncErrMsg`:
  - Set `m.err = msg.Err`
  - Do NOT set `synced = true` — panels stay in loading state
  - User can still see the error and quit

### Startup flow
1. `Init()` kicks off full sync
2. UI shows "Loading..." in all panels
3. On success (`SyncDoneMsg`): `synced = true`, sub-panels initialize with data
4. On failure (`SyncErrMsg`): error shown, panels stay loading

## Acceptance Criteria

- "Loading..." shown in all panels during initial sync
- On sync success: panels render with data
- On sync failure: error shown, panels stay in loading state
- `go build ./...` succeeds

## Summary of Changes

Updated internal/ui/app.go: added !m.synced early-return blocks in horizontalView() and stackedView(). When synced is false, all three panels render centered Loading... text using m.styles.MutedText() with lipgloss.Place, wrapped in the same border styles as normal content. Error banner and command bar still render regardless of sync state. go build, go vet, and go test all pass.
