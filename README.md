# todoist-tui

A terminal UI for [Todoist](https://todoist.com), built with Go and the [Charm](https://charm.sh) stack.

## Features

- **3-panel TUI** — sidebar (projects, filters, labels), main task list, and detail panel with lipgloss rounded borders and responsive layout
- **Focus cycling** — Tab / Shift+Tab to rotate between panels; active border highlight from theme config
- **Vim-style modal keybindings** — Normal, Insert, and Command modes via `internal/ui/keymap`
- **Async full sync on startup** — `Init()` launches a background full sync against the Todoist Sync API v9
- **Error banner** — sync failures rendered at the top of the view with the theme's error color
- **Optimistic updates** — local changes use `tmp-` prefixed temp IDs, resolved on sync
- **Offline support** — command queue with replay on reconnect
- **Bearer token auth** — configurable API token with sentinel errors (`ErrAuthFailed`, `ErrSyncFailed`)
- **Context-aware requests** — cancellable HTTP calls with configurable timeout (default 30s)
- **Deleted entity handling** — `IsDeleted` entities are removed from the local store on incremental sync

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
  sync/              # Todoist Sync API v9 client (Client, FullSync, IncrementalSync, resolveTempIDs)
  store/             # bbolt storage layer
  model/             # Domain types (Task, Project, Section, Label, Filter)
  queue/             # Offline command queue and replay
  ui/                # Root Bubbletea model and shared UI components
    app.go           # Root Model, Panel enum, NewModel, Init, Update, View, 3-panel rendering
    sidebar/         # Sidebar panel (projects, filters, labels)
    tasklist/        # Task list panel (compact/expanded)
    detail/          # Task detail and edit panel
    quickadd/        # Quick Add modal
    keymap/          # Vim-style keybinding definitions
      mode.go        # Mode enum (Normal, Insert, Command) with String()
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

### TUI Model (`internal/ui`)

The root Bubbletea model (`Model`) manages the full application lifecycle and 3-panel layout:

| Type | File | Description |
|---|---|---|
| `Model` | `app.go` | Root Bubbletea model holding config, store, sync client, active panel, mode, and dimensions |
| `Panel` | `app.go` | Enum (`PanelSidebar`, `PanelMain`, `PanelDetail`) identifying the focused panel |
| `Mode` | `keymap/mode.go` | Enum (`ModeNormal`, `ModeInsert`, `ModeCommand`) representing the current editor mode |
| `NewModel(cfg, store, client)` | `app.go` | Constructor; starts in Sidebar panel, Normal mode, zero dimensions |
| `Init()` | `app.go` | Kicks off an async `FullSync` via background command; returns `SyncDoneMsg` or `SyncErrMsg` |
| `Update(msg)` | `app.go` | Handles `WindowSizeMsg` (resize), `KeyMsg` (Tab/Shift+Tab focus cycling), and sync result messages |
| `View()` | `app.go` | Renders the 3-panel layout plus an error banner at the top when `m.err` is set |

**Panel layout**:

The `View()` method renders three panels joined horizontally:

- **Sidebar** (~20% width) — placeholder headings for Projects, Filters, Labels (styled with `theme.Header` color)
- **Main** (~50% width) — placeholder "Select a project or filter" (styled with `theme.MutedText`)
- **Detail** (~30% width) — placeholder "No task selected" (styled with `theme.MutedText`)

Ripple-wrap borders are drawn with `lipgloss.RoundedBorder()`. The active panel's border uses `theme.ActiveBorder`; inactive panels use `theme.InactiveBorder`.

**Focus cycling**:

`Tab` advances focus to the next panel (Sidebar → Main → Detail → Sidebar). `Shift+Tab` moves backward. Focus cycling wraps around using modulo arithmetic on the `activePanel` field.

**Vim-style modes**:

The `Mode` type in `keymap/mode.go` defines three editor modes with a `String()` method:

| Constant | String | Purpose |
|---|---|---|
| `ModeNormal` | `"NORMAL"` | Default navigation and command mode |
| `ModeInsert` | `"INSERT"` | Text input mode (form fields) |
| `ModeCommand` | `"COMMAND"` | Command-line mode for ex-style (colon) commands |

**Async full sync**:

`Init()` returns a `tea.Cmd` that runs `syncClient.FullSync(context.Background(), store)` in a background goroutine. On success, `SyncDoneMsg` is dispatched — the store is already populated by the sync client. On failure, `SyncErrMsg` carries the error, which is displayed as a red banner at the top of the view.

**Responsive layout**:

The model tracks terminal dimensions via `tea.WindowSizeMsg`. Panel widths are recalculated on every render using percentage-based splits (20/50/30). If dimensions are zero or negative (e.g., before the first `WindowSizeMsg`), `View()` returns `"Initializing..."` (or the error message if present).

### Sync Client (`internal/sync`)

The sync package provides a fully tested HTTP client for the Todoist Sync API v9:

| Type | File | Description |
|---|---|---|
| `ClientConfig` | `client.go` | Configuration: API token, timeout, endpoint override |
| `Client` | `client.go` | HTTP client with Bearer auth and context-aware requests |
| `NewClient(cfg)` | `client.go` | Constructor; defaults timeout to 30s, endpoint to `SyncEndpoint` |
| `DoSync(ctx, req)` | `client.go` | Low-level sync request; returns `ErrAuthFailed` on 401, `ErrSyncFailed` on other errors |
| `FullSync(ctx, store)` | `sync.go` | Fetches all data (`sync_token=*`), writes every entity to bbolt, persists sync token |
| `IncrementalSync(ctx, store)` | `sync.go` | Delta sync using stored token; falls back to `FullSync` if no token; removes `IsDeleted` entities |
| `resolveTempIDs(store, mapping)` | `tempid.go` | Replaces `tmp-` prefixed IDs with server-assigned IDs; updates cross-entity references; removes orphans |
| `SyncRequest` / `SyncResponse` | `types.go` | Request/response types for the Sync API |
| `Command` | `types.go` | Single command payload for optimistic updates |
| `ErrAuthFailed` / `ErrSyncFailed` | `types.go` | Sentinel errors for auth failure (401) and general sync failures |

**Sync flow**:

1. On first launch, `FullSync` sends `sync_token=*` to fetch all items, projects, sections, labels, and filters.
2. On subsequent syncs, `IncrementalSync` sends the stored sync token for delta updates.
3. Entities with `IsDeleted=true` are removed from the local store.
4. The `temp_id_mapping` from the response is processed by `resolveTempIDs`, which replaces `tmp-` placeholder IDs with real server IDs and updates all cross-entity references (e.g., `Task.ProjectID`, `Task.SectionID`, `Task.ParentID`, `Section.ProjectID`).
5. Any orphaned entities still carrying a `tmp-` prefix after resolution are cleaned up.

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