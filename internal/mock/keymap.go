package mock

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/trankhanh040147/prepf/internal/ui"
)

// MockKeyMap extends the base keymap with mock-specific bindings
type MockKeyMap struct {
	ui.KeyMap
	Surrender key.Binding
}

// DefaultMockKeyMap returns the default keymap for mock interviews
func DefaultMockKeyMap() MockKeyMap {
	base := ui.DefaultKeyMap()
	return MockKeyMap{
		KeyMap: base,
		Surrender: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "surrender"),
		),
	}
}

// ShortHelp returns short help text for mock interviews
func (k MockKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns full help text for mock interviews
func (k MockKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Top, k.Bottom},
		{k.Enter, k.Surrender},
		{k.Help, k.Quit},
	}
}

