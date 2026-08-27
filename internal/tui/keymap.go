package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines the keyboard shortcuts for the TUI.
type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Tab    key.Binding
	Enter  key.Binding
	Hint   key.Binding
	Reset  key.Binding
	Search key.Binding
	Quit   key.Binding
	Esc    key.Binding
	Help   key.Binding
}

// DefaultKeyMap returns the default keyboard bindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "switch pane"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "run/select"),
		),
		Hint: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "hint"),
		),
		Reset: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "reset"),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Esc: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back/cancel"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
	}
}
