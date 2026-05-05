// Package ui contains the root Bubbletea application model and shared UI components.
package ui

import (
	_ "github.com/charmbracelet/bubbletea"       // Core TUI framework
	_ "github.com/charmbracelet/lipgloss"        // Layout and styling
	_ "github.com/charmbracelet/bubbles/textinput" // Pre-built components
	_ "github.com/evertras/bubble-table/table"   // Table rendering
)
