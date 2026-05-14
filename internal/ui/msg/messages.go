package msg

// ProjectSelectedMsg is emitted when a project is selected in the sidebar.
type ProjectSelectedMsg struct {
	ID string
}

// FilterSelectedMsg is emitted when a filter is selected in the sidebar.
type FilterSelectedMsg struct {
	ID    string
	Query string
}

// LabelSelectedMsg is emitted when a label is selected in the sidebar.
type LabelSelectedMsg struct {
	Name string
}

// TaskSelectedMsg is emitted when a task row is selected in the task list.
type TaskSelectedMsg struct {
	ID string
}

// SyncCompleteMsg is emitted when a sync (full or incremental) completes
// successfully. The root model sets m.synced on the first SyncCompleteMsg
// to dismiss the "Loading..." placeholder.
type SyncCompleteMsg struct{}

// SyncTickMsg is emitted every 30 seconds by the periodic sync ticker to trigger
// an incremental sync against the Todoist Sync API.
type SyncTickMsg struct{}

// ToggleCompletedMsg is emitted when the user toggles completed task visibility
// with the configured keybinding (default: H).
type ToggleCompletedMsg struct{}
