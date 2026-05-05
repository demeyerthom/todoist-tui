package config

import "errors"

// ErrNoToken is returned when the config has an empty auth token.
var ErrNoToken = errors.New("config: auth token is required")

// Config holds the full application configuration loaded from TOML.
type Config struct {
	Auth   AuthConfig   `toml:"auth"`
	Keymap KeymapConfig `toml:"keybindings"`
	Theme  ThemeConfig  `toml:"theme"`
}

// AuthConfig holds authentication settings.
type AuthConfig struct {
	Token string `toml:"token"`
}

// KeymapConfig holds keybinding mappings for each editor mode.
// Each mode maps a logical action name to a key string
// (e.g. "up" → "k", "enter" → "Enter").
type KeymapConfig struct {
	Normal  map[string]string `toml:"normal"`
	Insert  map[string]string `toml:"insert"`
	Command map[string]string `toml:"command"`
}

// ThemeConfig holds visual styling configuration.
// Color values can be standard terminal color names or hex (#RRGGBB).
type ThemeConfig struct {
	ActiveBorder         string `toml:"active_border"`
	InactiveBorder       string `toml:"inactive_border"`
	SelectedRow          string `toml:"selected_row"`
	HoveredRow           string `toml:"hovered_row"`
	Header               string `toml:"header"`
	NormalText           string `toml:"normal_text"`
	MutedText            string `toml:"muted_text"`
	Error                string `toml:"error"`
	Success              string `toml:"success"`
	Warning              string `toml:"warning"`
	TaskDueToday         string `toml:"task_due_today"`
	TaskOverdue          string `toml:"task_overdue"`
	TaskCompleted        string `toml:"task_completed"`
	TaskPriorityHigh     string `toml:"task_priority_high"`
	TaskPriorityMedium   string `toml:"task_priority_medium"`
	TaskPriorityLow      string `toml:"task_priority_low"`
	SidebarProjectActive string `toml:"sidebar_project_active"`
	SidebarProjectLabel  string `toml:"sidebar_project_label"`
	InputBorder          string `toml:"input_border"`
	InputBorderError     string `toml:"input_border_error"`
	InputPlaceholder     string `toml:"input_placeholder"`
	CommandBar           string `toml:"command_bar"`
}

// Validate checks that the configuration is valid.
// It returns ErrNoToken if the auth token is empty.
func (c *Config) Validate() error {
	if c.Auth.Token == "" {
		return ErrNoToken
	}
	return nil
}