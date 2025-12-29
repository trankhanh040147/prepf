package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// updateKeyMsgError handles key messages in error state
func (m *Model) updateKeyMsgError(msg tea.KeyMsg) (*Model, tea.Cmd) {
	switch {
	// Quit
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit

	// Back to previous state
	case key.Matches(msg, m.keys.Back), key.Matches(msg, m.keys.Enter):
		m.errorMsg = ""
		m.returnToPreviousState()
		return m, nil

	default:
		// Any other key clears error
		m.errorMsg = ""
		m.returnToPreviousState()
		return m, nil
	}
}
