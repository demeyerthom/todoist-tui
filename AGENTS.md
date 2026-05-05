# todoist-tui

**IMPORTANT**: before you do anything else, run the `beans prime` command and heed its output.

## Project status

Pre-implementation. No Go code exists yet — only `PLAN.md` with architecture and milestone definitions. The project is at M0 (Skeleton).

## Stack & dependencies

- **Go 1.25+** (minimum)
- **Charm stack**: `bubbletea`, `lipgloss`, `bubbles`, `bubble-table` (evertras)
- **bbolt** (`go.etcd.io/bbolt`) for local KV store
- **go-toml/v2** (`pelletier/go-toml/v2`) for config
- **google/uuid** for Sync API idempotency

## Architecture

- Standard Go layout: `cmd/todoist-tui/main.go` entrypoint, `internal/` for all packages
- 3-panel TUI: sidebar, task list, detail — managed by Bubbletea Elm architecture
- Vim-style modal keybindings (normal/insert/command)
- Sync API v9 with optimistic updates using `tmp-` prefixed temp IDs
- Offline: command queue + replay on reconnect

## Task tracking

Uses the **beans CLI** (not TodoWrite/todowrite). See `.beans.yml` and `.opencode/agents/orchestrator.md` for the full workflow.

## Key conventions

- Config format: TOML (`config.toml`)
- Temp IDs for optimistic updates use `tmp-` prefix
- Sync token stored in bbolt; full sync on first launch (`sync_token=*`), incremental thereafter
- Periodic sync every 30s in background

