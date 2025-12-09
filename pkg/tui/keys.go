package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up      key.Binding
	Down    key.Binding
	Enter   key.Binding
	Back    key.Binding
	Stop    key.Binding
	Start   key.Binding
	Restart key.Binding
	Delete  key.Binding
	Refresh key.Binding
	Quit    key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "details"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Stop: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "stop"),
	),
	Start: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "resume"),
	),
	Restart: key.NewBinding(
		key.WithKeys("R"),
		key.WithHelp("R", "restart"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "refresh"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}