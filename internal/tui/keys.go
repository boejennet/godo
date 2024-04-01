package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Help            key.Binding
	Quit            key.Binding
	Up              key.Binding
	Down            key.Binding
	Refresh         key.Binding
	NewTodo         key.Binding
	ToggleCompleted key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Refresh, k.NewTodo, k.ToggleCompleted, k.Help, k.Quit}
}

var keys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("up", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("down", "move down"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	NewTodo: key.NewBinding(
		key.WithKeys("ctrl+n"),
		key.WithHelp("ctrl+n", "new todo"),
	),
	ToggleCompleted: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "toggle completed"),
	),
}
