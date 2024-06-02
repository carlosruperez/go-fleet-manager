package common

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/go-fleet-manager/internal/ui/styles"
)

type KeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
	Esc   key.Binding
	Quit  key.Binding
	Help  key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Enter, k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Enter},
		{k.Help, k.Quit},
	}
}

var Keys = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp(styles.KeyBindingStyle.Render("↑/k"), "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp(styles.KeyBindingStyle.Render("↓/j"), "move down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp(styles.KeyBindingStyle.Render("enter"), "select option"),
	),
	Esc: key.NewBinding(
		key.WithKeys("esc", "ctrl+e"),
		key.WithHelp(styles.KeyBindingStyle.Render("esc"), "Go to menu"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp(styles.KeyBindingStyle.Render("q"), "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp(styles.KeyBindingStyle.Render("?"), "toggle help"),
	),
}
