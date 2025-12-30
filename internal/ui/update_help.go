package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// updateKeyMsgHelp handles key messages in help state
func (m *Model) updateKeyMsgHelp(msg tea.KeyMsg) (*Model, tea.Cmd) {
	switch {
	// Any key exits help
	case key.Matches(msg, m.keys.Help), key.Matches(msg, m.keys.Back), key.Matches(msg, m.keys.Quit):
		m.returnToPreviousState()
		return m, nil

	default:
		// Any other key also exits help
		m.returnToPreviousState()
		return m, nil
	}
}
