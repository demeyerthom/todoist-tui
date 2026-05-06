package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/demeyerthom/todoist-tui/internal/config"
	"github.com/demeyerthom/todoist-tui/internal/store"
	"github.com/demeyerthom/todoist-tui/internal/sync"
	"github.com/demeyerthom/todoist-tui/internal/ui"
)

const version = "0.0.1"

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: loading config: %v\n", err)
		os.Exit(1)
	}

	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	dbPath := store.DBPath()
	s, err := store.New(store.StoreConfig{Path: dbPath})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: opening store: %v\n", err)
		os.Exit(1)
	}

	syncClient := sync.NewClient(sync.ClientConfig{
		Token:  cfg.Auth.Token,
		Timeout: 30 * time.Second,
	})

	model := ui.NewModel(cfg, s, syncClient)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		_ = model.Cleanup()
		os.Exit(1)
	}

	if err := model.Cleanup(); err != nil {
		fmt.Fprintf(os.Stderr, "error: closing store: %v\n", err)
		os.Exit(1)
	}
}
