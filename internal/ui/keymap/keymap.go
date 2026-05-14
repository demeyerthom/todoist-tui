package keymap

// KeyMap holds the action-to-key mappings for each editor mode.
type KeyMap struct {
	Normal  map[string]string
	Insert  map[string]string
	Command map[string]string
}

// DefaultKeyMap returns a KeyMap populated with the default vim-style bindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Normal: map[string]string{
			"up":                "k",
			"down":              "j",
			"enter":             "Enter",
			"complete":          "x",
			"delete":            "d",
			"quick_add":         "o",
			"edit":              "i",
			"search":            "/",
			"switch_panel":      "Tab",
			"go_top":            "g",
			"go_bottom":         "G",
			"focus_sidebar":     "1",
			"focus_main":        "2",
			"focus_detail":      "3",
			"collapse":          "h",
			"expand":            "l",
			"toggle_completed":  "H",
			"command_mode":      ":",
			"escape":            "Esc",
		},
		Insert: map[string]string{
			"escape": "Esc",
			"next":   "Tab",
		},
		Command: map[string]string{
			"escape": "Esc",
		},
	}
}

// KeyFor returns the key bound to the given action in the specified mode.
// It returns ("", false) if no binding is found.
func (km KeyMap) KeyFor(mode Mode, action string) (string, bool) {
	var m map[string]string
	switch mode {
	case ModeNormal:
		m = km.Normal
	case ModeInsert:
		m = km.Insert
	case ModeCommand:
		m = km.Command
	default:
		return "", false
	}

	key, ok := m[action]
	return key, ok
}
