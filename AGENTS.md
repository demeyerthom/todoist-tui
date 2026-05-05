# todoist-tui

**IMPORTANT**: before you do anything else, run the `beans prime` command and heed its output.

## Stack & dependencies

- **Go 1.25+** (minimum)
- **Charm stack**: `bubbletea`, `lipgloss`, `bubbles`, `bubble-table` (evertras)
- **bbolt** (`go.etcd.io/bbolt`) for local KV store
- **go-toml/v2** (`pelletier/go-toml/v2`) for config
- **google/uuid** for Sync API idempotency

Dependencies are pinned via blank imports in `doc.go` files so `go mod tidy` preserves them.

## Architecture

- Standard Go layout: `cmd/todoist-tui/main.go` entrypoint, `internal/` for all packages
- 3-panel TUI: sidebar, task list, detail — managed by Bubbletea Elm architecture
- Vim-style modal keybindings (normal/insert/command)
- Sync API v9 with optimistic updates using `tmp-` prefixed temp IDs
- Offline: command queue + replay on reconnect

## Task tracking

Uses the **beans CLI** (not TodoWrite/todowrite). See `.beans.yml` and `.opencode/agents/orchestrator.md` for the full workflow.

## Key conventions

- Config format: TOML (`config.toml`), loaded via `internal/config` package
- Config path: XDG-compliant (`$XDG_CONFIG_HOME/todoist-tui/config.toml`), fallback to `$HOME/.config`
- `Load()` auto-creates default config if missing; `WriteDefaultConfig` never overwrites existing files
- TOML key for keybindings is `[keybindings]` (not `[keymap]`)
- KeymapConfig sub-structs use `map[string]string` for flexible key-to-action mappings
- ThemeConfig has 22 color fields matching `config.toml.example`
- `ErrNoToken` sentinel error returned when auth token is empty
- Temp IDs for optimistic updates use `tmp-` prefix
- Sync token stored in bbolt; full sync on first launch (`sync_token=*`), incremental thereafter
- Periodic sync every 30s in background
- Documenter subagent must only update `README.md` — no additional doc files

## Verification

```bash
go test ./...          # Run tests (only internal/config has test files currently)
go vet ./...           # Static analysis
go build ./cmd/todoist-tui  # Build check
```