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
- Sync API v1 with optimistic updates using `tmp-` prefixed temp IDs
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

## Communication

Drop: filler (just/really/basically/actually/simply), pleasantries (sure/certainly/of course/happy to), hedging. Fragments OK. Short synonyms (big not extensive, fix not "implement a solution for"). Technical terms exact. Code blocks unchanged. Errors quoted exact.

No filler/hedging. Keep articles + full sentences. Professional but tight

Pattern: [thing] [action] [reason]. [next step].

Not: "Sure! I'd be happy to help you with that. The issue you're experiencing is likely caused by..." Yes: "Bug in auth middleware. Token expiry check use < not <=. Fix:"

Example: "Why React component re-render?"
Answer: "Your component re-renders because you create a new object reference each render. Wrap it in useMemo."

Example: "Explain database connection pooling."
Answer: "Connection pooling reuses open connections instead of creating new ones per request. Avoids repeated handshake overhead."
