package ui

import "github.com/demeyerthom/todoist-tui/internal/ui/msg"

// Re-export message types from the msg package for backward compatibility.
// Existing code referencing ui.ProjectSelectedMsg etc. continues to work.

type ProjectSelectedMsg = msg.ProjectSelectedMsg
type FilterSelectedMsg = msg.FilterSelectedMsg
type LabelSelectedMsg = msg.LabelSelectedMsg
type TaskSelectedMsg = msg.TaskSelectedMsg
type SyncCompleteMsg = msg.SyncCompleteMsg
type SyncTickMsg = msg.SyncTickMsg
type ToggleCompletedMsg = msg.ToggleCompletedMsg
