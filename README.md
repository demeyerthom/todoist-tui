# todoist-tui

A terminal UI for [Todoist](https://todoist.com), built with Go and the [Charm](https://charm.sh) stack.

## Status

**M0 (Skeleton)** — project initialized with directory structure, dependencies, and minimal entrypoint. Not yet functional.

## Requirements

- Go 1.25+

## Build & Run

```bash
go build ./cmd/todoist-tui
./todoist-tui
# Output: todoist-tui v0.0.1
```

Or via `go run`:

```bash
go run ./cmd/todoist-tui
```

## Configuration

Copy the example config and add your Todoist API token:

```bash
cp config.toml.example config.toml
```

Generate a token at <https://todoist.com/app/settings/integrations/developer>.

See `config.toml.example` for all available options:

- **`[auth]`** — API token
- **`[keybindings.normal]`** — Normal mode keybindings (vim-style navigation)
- **`[keybindings.insert]`** — Insert mode keybindings (form editing)
- **`[keybindings.command]`** — Command mode keybindings (`:` commands)
- **`[theme]`** — Colors and styling (borders, rows, text, task priorities)

## Project Structure

```
cmd/todoist-tui/     # Entrypoint (main.go)
internal/
  config/            # TOML config loading
  sync/              # Todoist Sync API v9 client
  store/             # bbolt storage layer
  model/             # Domain types (Task, Project, Section, Label, Filter)
  queue/             # Offline command queue and replay
  ui/                # Root Bubbletea model and shared UI components
    sidebar/         # Sidebar panel (projects, filters, labels)
    tasklist/        # Task list panel (compact/expanded)
    detail/          # Task detail and edit panel
    quickadd/        # Quick Add modal
    keymap/          # Vim-style keybinding definitions
    theme/           # Lipgloss styling and color scheme
```

## Dependencies

| Package | Purpose |
|---|---|
| `charmbracelet/bubbletea` | Core TUI framework (Elm architecture) |
| `charmbracelet/lipgloss` | Layout and styling |
| `charmbracelet/bubbles` | Pre-built components (text input, etc.) |
| `evertras/bubble-table` | Table rendering for task lists |
| `go.etcd.io/bbolt` | Embedded KV store for local data |
| `pelletier/go-toml/v2` | TOML config parsing |
| `google/uuid` | UUIDs for Sync API idempotency |

Dependencies are pinned via blank imports in `doc.go` files so `go mod tidy` preserves them.

## Architecture

- **3-panel layout**: sidebar | task list | detail — managed by Bubbletea's Elm architecture
- **Vim-style modal keybindings**: normal (navigation), insert (form editing), command (`:` commands)
- **Sync API v9**: full sync on first launch, incremental thereafter (sync token stored in bbolt)
- **Optimistic updates**: local changes use `tmp-` prefixed temp IDs, resolved on sync
- **Offline support**: command queue with replay on reconnect
- **Periodic sync**: every 30s in background

## Key Conventions

- Config format: TOML (`config.toml`)
- Temp IDs for optimistic updates: `tmp-` prefix
- Sync token: stored in bbolt; full sync uses `sync_token=*`
- Keybinding sections: `[keybindings.normal]`, `[keybindings.insert]`, `[keybindings.command]`