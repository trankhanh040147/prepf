package ui

import (
	"github.com/charmbracelet/bubbles/key"
)

// KeyMap defines all keybindings for the TUI
type KeyMap struct {
	// Navigation
	Up     key.Binding
	Down   key.Binding
	Top    key.Binding
	Bottom key.Binding

	// Search
	Search      key.Binding
	NextMatch   key.Binding
	PrevMatch   key.Binding
	SearchEnter key.Binding
	SearchEsc   key.Binding

	// General
	Help      key.Binding
	Quit      key.Binding
	ForceQuit key.Binding
	Enter     key.Binding
	Tab       key.Binding
	Back      key.Binding

	// Chat
	EnterChat     key.Binding
	ExitChat      key.Binding
	SendMessage   key.Binding
	CancelRequest key.Binding

	// Yank
	Yank     key.Binding
	YankLast key.Binding
}

// DefaultKeyMap returns default key bindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		// Navigation
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("k/↑", "scroll up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("j/↓", "scroll down"),
		),
		Top: key.NewBinding(
			key.WithKeys("g", "home"),
			key.WithHelp("g/Home", "go to top"),
		),
		Bottom: key.NewBinding(
			key.WithKeys("G", "end"),
			key.WithHelp("G/End", "go to bottom"),
		),

		// Search
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		NextMatch: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "next match"),
		),
		PrevMatch: key.NewBinding(
			key.WithKeys("N"),
			key.WithHelp("N", "prev match"),
		),
		SearchEnter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm"),
		),
		SearchEsc: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),

		// General
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
		ForceQuit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "force quit"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),

		// Chat
		EnterChat: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "enter chat"),
		),
		ExitChat: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "exit chat"),
		),
		SendMessage: key.NewBinding(
			key.WithKeys("alt+enter"),
			key.WithHelp("alt+enter", "send"),
		),
		CancelRequest: key.NewBinding(
			key.WithKeys("ctrl+x"),
			key.WithHelp("ctrl+x", "cancel"),
		),

		// Yank
		Yank: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("y", "yank"),
		),
		YankLast: key.NewBinding(
			key.WithKeys("Y"),
			key.WithHelp("Y", "yank last"),
		),
	}
}

// ShortHelp returns short help text
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns full help text
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Top, k.Bottom},
		{k.Search, k.NextMatch, k.PrevMatch},
		{k.Help, k.Quit, k.Back},
		{k.EnterChat, k.SendMessage, k.CancelRequest},
		{k.Yank, k.YankLast},
	}
}
