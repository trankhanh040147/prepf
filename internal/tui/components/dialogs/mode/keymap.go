package mode

import (
	"charm.land/bubbles/v2/key"
)

type KeyMap struct {
	Select   key.Binding
	Close    key.Binding
	Next     key.Binding
	Previous key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Close: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "close"),
		),
		Next: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("j/↓", "next"),
		),
		Previous: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("k/↑", "previous"),
		),
	}
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Select, k.Close, k.Next, k.Previous}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Select, k.Close},
		{k.Next, k.Previous},
	}
}
