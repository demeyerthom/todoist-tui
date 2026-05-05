package config

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// DefaultConfig returns a fully populated Config with sensible defaults
// matching config.toml.example.
func DefaultConfig() *Config {
	return &Config{
		Auth: AuthConfig{
			Token: "",
		},
		Keymap: KeymapConfig{
			Normal: map[string]string{
				"up":            "k",
				"down":          "j",
				"enter":         "Enter",
				"complete":      "x",
				"delete":        "d",
				"quick_add":     "o",
				"edit":          "i",
				"search":        "/",
				"switch_panel":  "Tab",
				"go_top":        "g",
				"go_bottom":     "G",
				"focus_sidebar": "1",
				"focus_main":    "2",
				"focus_detail":  "3",
				"command_mode":  ":",
				"escape":        "Esc",
			},
			Insert: map[string]string{
				"escape": "Esc",
				"next":   "Tab",
			},
			Command: map[string]string{
				"escape": "Esc",
			},
		},
		Theme: ThemeConfig{
			ActiveBorder:         "blue",
			InactiveBorder:       "brightblack",
			SelectedRow:          "blue",
			HoveredRow:           "brightblack",
			Header:               "brightwhite",
			NormalText:           "white",
			MutedText:            "brightblack",
			Error:                "red",
			Success:              "green",
			Warning:              "yellow",
			TaskDueToday:         "yellow",
			TaskOverdue:          "red",
			TaskCompleted:        "brightblack",
			TaskPriorityHigh:    "red",
			TaskPriorityMedium:  "yellow",
			TaskPriorityLow:     "blue",
			SidebarProjectActive: "brightwhite",
			SidebarProjectLabel: "magenta",
			InputBorder:          "blue",
			InputBorderError:     "red",
			InputPlaceholder:     "brightblack",
			CommandBar:           "brightblack",
		},
	}
}

// WriteDefaultConfig writes the default configuration to the given file path
// as TOML. It creates parent directories if needed. It does NOT overwrite
// an existing file.
func WriteDefaultConfig(path string) error {
	if _, err := os.Stat(path); err == nil {
		// File already exists; do not overwrite.
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	data, err := toml.Marshal(DefaultConfig())
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}