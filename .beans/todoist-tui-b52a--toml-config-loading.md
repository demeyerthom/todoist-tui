---
# todoist-tui-b52a
title: TOML config loading
status: todo
type: feature
priority: normal
created_at: 2026-05-03T14:58:16Z
updated_at: 2026-05-03T15:11:20Z
parent: todoist-tui-j3br
---

Load config.toml with API token, keybindings, and theme settings. Uses pelletier/go-toml/v2 for parsing. Config stores auth token (personal API token only), customizable keybindings, and theme preferences.



## Task Dependency Graph

1. `todoist-tui-a0b4` Define Config structs ← F1-dir-structure
2. `todoist-tui-unpx` Default config generation ← a0b4
3. `todoist-tui-kem6` Load and parse config from XDG ← unpx
4. `todoist-tui-jop7` Config loading tests ← kem6

## Cross-feature dependency
- Depends on F1 (project init) completing first
