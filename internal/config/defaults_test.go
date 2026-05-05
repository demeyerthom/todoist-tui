package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pelletier/go-toml/v2"
)

func TestDefaultConfigAllFieldsPopulated(t *testing.T) {
	cfg := DefaultConfig()

	// Auth
	if cfg.Auth.Token != "" {
		// Token defaults to empty string, which is fine
	}

	// Keymap — Normal mode
	normalKeys := map[string]string{
		"up": "k", "down": "j", "enter": "Enter",
		"complete": "x", "delete": "d", "quick_add": "o",
		"edit": "i", "search": "/", "switch_panel": "Tab",
		"go_top": "g", "go_bottom": "G",
		"focus_sidebar": "1", "focus_main": "2", "focus_detail": "3",
		"command_mode": ":", "escape": "Esc",
	}
	for k, v := range normalKeys {
		if cfg.Keymap.Normal[k] != v {
			t.Errorf("Keymap.Normal[%q] = %q, want %q", k, cfg.Keymap.Normal[k], v)
		}
	}
	if len(cfg.Keymap.Normal) != len(normalKeys) {
		t.Errorf("Keymap.Normal has %d entries, want %d", len(cfg.Keymap.Normal), len(normalKeys))
	}

	// Keymap — Insert mode
	insertKeys := map[string]string{"escape": "Esc", "next": "Tab"}
	for k, v := range insertKeys {
		if cfg.Keymap.Insert[k] != v {
			t.Errorf("Keymap.Insert[%q] = %q, want %q", k, cfg.Keymap.Insert[k], v)
		}
	}

	// Keymap — Command mode
	if cfg.Keymap.Command["escape"] != "Esc" {
		t.Errorf("Keymap.Command[escape] = %q, want %q", cfg.Keymap.Command["escape"], "Esc")
	}

	// Theme — spot-check a few fields
	if cfg.Theme.ActiveBorder != "blue" {
		t.Errorf("Theme.ActiveBorder = %q, want %q", cfg.Theme.ActiveBorder, "blue")
	}
	if cfg.Theme.CommandBar != "brightblack" {
		t.Errorf("Theme.CommandBar = %q, want %q", cfg.Theme.CommandBar, "brightblack")
	}
}

func TestDefaultConfigRoundTripsTOML(t *testing.T) {
	cfg := DefaultConfig()
	data, err := toml.Marshal(cfg)
	if err != nil {
		t.Fatalf("Marshal default config: %v", err)
	}

	var loaded Config
	if err := toml.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("Unmarshal default config: %v", err)
	}

	// Verify a few fields survive the round trip
	if loaded.Keymap.Normal["up"] != "k" {
		t.Errorf("round-trip: Keymap.Normal[up] = %q, want %q", loaded.Keymap.Normal["up"], "k")
	}
	if loaded.Theme.ActiveBorder != "blue" {
		t.Errorf("round-trip: Theme.ActiveBorder = %q, want %q", loaded.Theme.ActiveBorder, "blue")
	}
}

func TestWriteDefaultConfigCreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "dir", "config.toml")

	if err := WriteDefaultConfig(path); err != nil {
		t.Fatalf("WriteDefaultConfig: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	var loaded Config
	if err := toml.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("Unmarshal written config: %v", err)
	}

	if loaded.Keymap.Normal["down"] != "j" {
		t.Errorf("written config: Keymap.Normal[down] = %q, want %q", loaded.Keymap.Normal["down"], "j")
	}
}

func TestWriteDefaultConfigDoesNotOverwrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	// Write a file with known content
	original := []byte("[auth]\ntoken = 'existing'\n")
	if err := os.WriteFile(path, original, 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	if err := WriteDefaultConfig(path); err != nil {
		t.Fatalf("WriteDefaultConfig should not fail: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	if string(data) != string(original) {
		t.Errorf("WriteDefaultConfig overwrote existing file")
	}
}