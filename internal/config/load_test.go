package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfigPath(t *testing.T) {
	tests := []struct {
		name      string
		xdgHome   string
		unsetXDG  bool
		wantSuffix string
	}{
		{
			name:       "uses_xdg_config_home",
			xdgHome:    "/custom/config",
			unsetXDG:   false,
			wantSuffix: filepath.Join("/custom/config", "todoist-tui", "config.toml"),
		},
		{
			name:       "default_when_xdg_unset",
			xdgHome:    "",
			unsetXDG:   true,
			wantSuffix: filepath.Join("todoist-tui", "config.toml"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origXDG, origSet := os.LookupEnv("XDG_CONFIG_HOME")
			defer func() {
				if origSet {
					os.Setenv("XDG_CONFIG_HOME", origXDG)
				} else {
					os.Unsetenv("XDG_CONFIG_HOME")
				}
			}()

			if tt.unsetXDG {
				os.Unsetenv("XDG_CONFIG_HOME")
			} else {
				os.Setenv("XDG_CONFIG_HOME", tt.xdgHome)
			}

			got := ConfigPath()

			if !filepath.IsAbs(got) {
				t.Errorf("ConfigPath() = %q, want absolute path", got)
			}
			if !strings.HasSuffix(got, tt.wantSuffix) {
				t.Errorf("ConfigPath() = %q, want suffix %q", got, tt.wantSuffix)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(dir string) error
		wantErr     error  // sentinel error to check with errors.Is
		errContains string // substring expected in error message
		wantToken   string // expected token on success
		checkFile   bool   // verify config file was created on disk
	}{
		{
			name: "valid_toml_with_token",
			setup: func(dir string) error {
				cfgDir := filepath.Join(dir, "todoist-tui")
				if err := os.MkdirAll(cfgDir, 0o755); err != nil {
					return err
				}
				content := `[auth]
token = "test-token-123"

[keybindings.normal]
up = "k"
down = "j"

[theme]
active_border = "blue"
`
				return os.WriteFile(filepath.Join(cfgDir, "config.toml"), []byte(content), 0o644)
			},
			wantToken: "test-token-123",
		},
		{
			name: "missing_file_creates_default",
			setup: func(dir string) error {
				// No file created; Load should write the default config.
				return nil
			},
			wantErr:   ErrNoToken,
			checkFile: true,
		},
		{
			name: "empty_token_returns_err_no_token",
			setup: func(dir string) error {
				cfgDir := filepath.Join(dir, "todoist-tui")
				if err := os.MkdirAll(cfgDir, 0o755); err != nil {
					return err
				}
				content := `[auth]
token = ""

[theme]
active_border = "blue"
`
				return os.WriteFile(filepath.Join(cfgDir, "config.toml"), []byte(content), 0o644)
			},
			wantErr: ErrNoToken,
		},
		{
			name: "invalid_toml_returns_parse_error",
			setup: func(dir string) error {
				cfgDir := filepath.Join(dir, "todoist-tui")
				if err := os.MkdirAll(cfgDir, 0o755); err != nil {
					return err
				}
				content := `this is not valid toml [[[`
				return os.WriteFile(filepath.Join(cfgDir, "config.toml"), []byte(content), 0o644)
			},
			errContains: "parse config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()

			origXDG, origSet := os.LookupEnv("XDG_CONFIG_HOME")
			os.Setenv("XDG_CONFIG_HOME", dir)
			defer func() {
				if origSet {
					os.Setenv("XDG_CONFIG_HOME", origXDG)
				} else {
					os.Unsetenv("XDG_CONFIG_HOME")
				}
			}()

			if err := tt.setup(dir); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			cfg, err := Load()

			switch {
			case tt.wantErr != nil:
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Load() error = %v, want %v", err, tt.wantErr)
				}
				if cfg != nil {
					t.Errorf("Load() cfg = %v, want nil", cfg)
				}
			case tt.errContains != "":
				if err == nil {
					t.Errorf("Load() expected error containing %q, got nil", tt.errContains)
				} else if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Load() error = %v, want error containing %q", err, tt.errContains)
				}
			default:
				if err != nil {
					t.Errorf("Load() unexpected error = %v", err)
				}
				if cfg == nil {
					t.Fatal("Load() returned nil config")
				}
				if cfg.Auth.Token != tt.wantToken {
					t.Errorf("Load() token = %q, want %q", cfg.Auth.Token, tt.wantToken)
				}
			}

			if tt.checkFile {
				configPath := filepath.Join(dir, "todoist-tui", "config.toml")
				if _, err := os.Stat(configPath); os.IsNotExist(err) {
					t.Errorf("expected config file to be created at %q", configPath)
				}
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr error
	}{
		{
			name:    "empty_token_returns_err_no_token",
			token:   "",
			wantErr: ErrNoToken,
		},
		{
			name:    "valid_token_returns_nil",
			token:   "abc123def456",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Auth: AuthConfig{Token: tt.token},
			}

			err := cfg.Validate()

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Validate() error = %v, want %v", err, tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error = %v", err)
				}
			}
		})
	}
}