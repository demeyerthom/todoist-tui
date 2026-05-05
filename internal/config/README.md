# internal/config — TOML configuration loading

This package handles loading, parsing, validating, and defaulting the todoist-tui TOML configuration.

## API

### `ConfigPath() string`

Returns the XDG-compliant path to the config file:

```
$XDG_CONFIG_HOME/todoist-tui/config.toml
```

Falls back to `$HOME/.config/todoist-tui/config.toml` if `os.UserConfigDir()` fails.

### `Load() (*Config, error)`

Primary entry point. Reads and parses the config file from `ConfigPath()`:

1. If the file does not exist, calls `WriteDefaultConfig` to create it.
2. Reads and unmarshals the TOML file into a `Config` struct.
3. Calls `Config.Validate()` — returns `ErrNoToken` if the auth token is empty.

### `DefaultConfig() *Config`

Returns a fully populated `Config` with sensible defaults matching `config.toml.example`. The auth token defaults to an empty string (user must provide one).

### `WriteDefaultConfig(path string) error`

Writes `DefaultConfig()` as TOML to the given path. Creates parent directories as needed. **Does not overwrite** an existing file — returns `nil` silently if the file already exists.

### `Config.Validate() error`

Validates the config. Currently checks that `Auth.Token` is non-empty, returning `ErrNoToken` if it is.

## Types

### `Config`

Top-level struct with TOML tags:

```go
type Config struct {
    Auth   AuthConfig   `toml:"auth"`
    Keymap KeymapConfig `toml:"keybindings"`  // note: TOML key is "keybindings"
    Theme  ThemeConfig  `toml:"theme"`
}
```

### `AuthConfig`

```go
type AuthConfig struct {
    Token string `toml:"token"`
}
```

### `KeymapConfig`

Each mode is a `map[string]string` mapping logical action names to key strings:

```go
type KeymapConfig struct {
    Normal  map[string]string `toml:"normal"`
    Insert  map[string]string `toml:"insert"`
    Command map[string]string `toml:"command"`
}
```

Using `map[string]string` (rather than a fixed struct) allows users to add or remove bindings without code changes.

### `ThemeConfig`

22 color fields matching `config.toml.example`. Values are terminal color names (`blue`, `brightblack`, etc.) or hex (`#RRGGBB`).

## Design decisions

- **TOML key `keybindings`**: The TOML section is `[keybindings]`, not `[keymap]`. The Go struct field is named `Keymap` for brevity, but the TOML tag is `toml:"keybindings"`.
- **`map[string]string` for keybindings**: Allows flexible key-to-action mappings without requiring struct changes for every new binding.
- **No overwrite on defaults**: `WriteDefaultConfig` preserves user edits. On first launch, `Load()` creates the file; on subsequent launches, it reads the existing one.
- **`ErrNoToken` sentinel**: A dedicated error variable so callers can distinguish "missing token" from other config errors.
- **XDG compliance**: `ConfigPath()` uses `os.UserConfigDir()` for platform-appropriate config directories, with a `$HOME/.config` fallback.

## Tests

10 tests in `defaults_test.go` and `load_test.go` covering:

- `DefaultConfig` field population and TOML round-trip
- `WriteDefaultConfig` file creation and non-overwrite behavior
- `ConfigPath` XDG resolution and fallback
- `Load` with valid config, missing file, empty token, and invalid TOML
- `Validate` with empty and valid tokens