package keys

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines all key bindings for the application
type KeyMap struct {
	// Navigation
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
	Back  key.Binding

	// Global
	Quit key.Binding
	Help key.Binding

	// View controls
	TogglePreview     key.Binding
	ToggleSortDirection key.Binding

	// Scene switching
	BrowserScene  key.Binding
	MountsScene   key.Binding
	SettingsScene key.Binding
}

// DefaultKeyMap returns the default key bindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		// Navigation
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter", "l"),
			key.WithHelp("enter/l", "select/open"),
		),
		Back: key.NewBinding(
			key.WithKeys("backspace", "h"),
			key.WithHelp("←/h", "go back"),
		),

		// Global
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),

		// View controls
		TogglePreview: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "toggle preview"),
		),
		ToggleSortDirection: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "toggle sort direction"),
		),

		// Scene switching
		BrowserScene: key.NewBinding(
			key.WithKeys("f1"),
			key.WithHelp("F1", "browser"),
		),
		MountsScene: key.NewBinding(
			key.WithKeys("f2"),
			key.WithHelp("F2", "mounts"),
		),
		SettingsScene: key.NewBinding(
			key.WithKeys("f3"),
			key.WithHelp("F3", "settings"),
		),
	}
}
