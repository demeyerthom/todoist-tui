---
# todoist-tui-0rtm
title: M1 Shared Foundation
status: completed
type: feature
priority: normal
created_at: 2026-05-07T19:34:23Z
updated_at: 2026-05-09T15:03:19Z
parent: todoist-tui-96st
blocked_by:
    - todoist-tui-j3br
---

Shared UI primitives needed by all M1 panels — message types, theme colors/styles, and due date formatting.

## Task Dependency Graph

1. `todoist-tui-m01d` Shared message types ← (none)
2. `todoist-tui-m8uj` Theme package: colors, styles, and due date formatting ← m01d

## Cross-feature dependency
- Sidebar (todoist-tui-6ss5), Task List (todoist-tui-rcpd), and Detail Panel (todoist-tui-ys80) depend on this feature completing first.
- Navigation (todoist-tui-us18) depends on Sidebar, Task List, and Detail Panel completing.
- Periodic Sync (todoist-tui-fykz) depends on Navigation completing.

## Summary of Changes

All 2 shared foundation tasks completed: shared message types (7 message types in internal/ui/msg/), theme package (colors, styles, format due date).
