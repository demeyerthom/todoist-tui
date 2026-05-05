---
# todoist-tui-4yfu
title: Create directory structure and minimal main.go
status: todo
type: task
created_at: 2026-05-03T15:08:03Z
updated_at: 2026-05-03T15:08:03Z
parent: todoist-tui-4wyv
blocked_by:
    - todoist-tui-8sbd
---

## Description

Create the full directory structure per PLAN.md and a minimal cmd/todoist-tui/main.go entrypoint that compiles.

## Requirements

- Directories: cmd/todoist-tui/, internal/config/, internal/sync/, internal/store/, internal/model/, internal/ui/, internal/ui/sidebar/, internal/ui/tasklist/, internal/ui/detail/, internal/ui/quickadd/, internal/ui/keymap/, internal/ui/theme/, internal/queue/
- cmd/todoist-tui/main.go: minimal main() that compiles (e.g., prints version or just returns)
- Each internal/ package gets a placeholder doc.go or a minimal .go file so 'go build' works

## Acceptance Criteria

- All directories exist
- 'go build ./...' succeeds
- 'go vet ./...' passes
