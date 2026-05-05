---
# todoist-tui-8sbd
title: Go module init and Go version
status: todo
type: task
created_at: 2026-05-03T15:07:58Z
updated_at: 2026-05-03T15:07:58Z
parent: todoist-tui-4wyv
---

## Description

Initialize the Go module with 'go mod init github.com/demeyerthom/todoist-tui' and set minimum Go version to 1.25 in go.mod.

## Requirements

- Module path: github.com/demeyerthom/todoist-tui
- Go version: >= 1.25 in go.mod
- No dependencies added yet (separate task)

## Acceptance Criteria

- 'go mod init' completed
- go.mod contains 'go 1.25'
- 'go build ./...' succeeds with no errors
