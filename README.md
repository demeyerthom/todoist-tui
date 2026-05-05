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

On first launch, `todoist-tui` automatically creates a default config file at the XDG-compliant path:

```
$XDG_CONFIG_HOME/todoist-tui/config.toml    # (or ~/.config/todoist-tui/config.toml)
```

Edit the file and set your Todoist API token under `[auth]`:

```toml
[auth]
token = 'your-api-token-here'
```

Generate a token at <https://todoist.com/app/settings/integrations/developer>.

The app will not start without a valid token — `config.Validate()` returns `ErrNoToken` if the token is empty.

### Config sections

| Section | Description |
|---|---|
| `[auth]` | API token (required) |
| `[keybindings.normal]` | Normal mode — vim-style navigation (`j`/`k`, `Enter`, `x`, etc.) |
| `[keybindings.insert]` | Insert mode — form editing (`Esc` to exit, `Tab` for next field) |
| `[keybindings.command]` | Command mode — `:` commands (`Esc` to exit) |
| `[theme]` | Colors and styling (22 fields: borders, rows, text, task priorities, sidebar, inputs) |

### Behavior

- **Auto-creation**: If `config.toml` doesn't exist, `Load()` writes the defaults first.
- **No overwrite**: `WriteDefaultConfig` will not overwrite an existing config file — manual edits are preserved.
- **Fallback path**: If `os.UserConfigDir()` fails, `ConfigPath()` falls back to `$HOME/.config`.
- **Reference**: See `config.toml.example` in the repo root for the full set of options with comments.

## Project Structure

```
cmd/todoist-tui/     # Entrypoint (main.go)
internal/
  config/            # TOML config loading (Config, Load, DefaultConfig, WriteDefaultConfig, Validate)
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

- Config format: TOML (`config.toml`), loaded via `internal/config` package
- Config path: XDG-compliant (`$XDG_CONFIG_HOME/todoist-tui/config.toml`)
- Config auto-created on first launch; existing files are never overwritten
- TOML key for keybindings: `[keybindings]` (not `[keymap]`)
- KeymapConfig uses `map[string]string` per mode for flexible key-to-action mappings
- ThemeConfig has 22 color fields (terminal color names or `#RRGGBB` hex)
- Temp IDs for optimistic updates: `tmp-` prefix
- Sync token: stored in bbolt; full sync uses `sync_token=*`
- Keybinding sections: `[keybindings.normal]`, `[keybindings.insert]`, `[keybindings.command]`