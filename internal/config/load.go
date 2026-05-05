package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// ConfigPath returns the XDG-compliant path to the configuration file.
// It uses os.UserConfigDir to find the user's configuration directory
// (typically ~/.config on Linux) and appends "todoist-tui/config.toml".
func ConfigPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		// Fallback to a sensible default if UserConfigDir fails.
		dir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return filepath.Join(dir, "todoist-tui", "config.toml")
}

// Load reads and parses the configuration file from the XDG config path.
// If the file does not exist, it writes the default config first, then
// loads it. The returned Config is validated before being returned.
func Load() (*Config, error) {
	path := ConfigPath()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := WriteDefaultConfig(path); err != nil {
			return nil, fmt.Errorf("write default config: %w", err)
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}