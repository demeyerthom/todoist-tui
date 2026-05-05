---
# todoist-tui-4wyv
title: Project init and directory structure
status: todo
type: feature
priority: normal
created_at: 2026-05-03T14:58:14Z
updated_at: 2026-05-03T15:11:13Z
parent: todoist-tui-j3br
---

Initialize Go module (go mod init), create the standard Go layout directory structure (cmd/todoist-tui/main.go entrypoint, internal/config, internal/sync, internal/store, internal/model, internal/ui/sidebar, internal/ui/tasklist, internal/ui/detail, internal/ui/quickadd, internal/ui/keymap, internal/ui/theme, internal/queue). Go 1.25 minimum. Add all dependencies: bubbletea, lipgloss, bubbles, bubble-table, bbolt, go-toml/v2, google/uuid.



## Task Dependency Graph

1. `todoist-tui-8sbd` Go module init (no deps)
2. `todoist-tui-4yfu` Directory structure ← 8sbd
3. `todoist-tui-kmc2` Add dependencies ← 4yfu
4. `todoist-tui-80nc` Example config.toml ← 4yfu
