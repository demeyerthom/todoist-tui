package keymap

// Mode represents the current editor mode (vim-style modal keybindings).
type Mode int

const (
	ModeNormal  Mode = iota // default navigation/command mode
	ModeInsert              // text input mode
	ModeCommand             // command-line mode for ex-style commands
)

// String returns a human-readable representation of the mode.
func (m Mode) String() string {
	switch m {
	case ModeNormal:
		return "NORMAL"
	case ModeInsert:
		return "INSERT"
	case ModeCommand:
		return "COMMAND"
	default:
		return "UNKNOWN"
	}
}
