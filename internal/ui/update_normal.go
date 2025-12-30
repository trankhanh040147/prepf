package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// updateKeyMsgNormal handles key messages in normal state
func (m *Model) updateKeyMsgNormal(msg tea.KeyMsg) (*Model, tea.Cmd) {
	switch {
	// Quit
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit

	// Help toggle
	case key.Matches(msg, m.keys.Help):
		m.previousState = m.state
		m.state = StateHelp
		return m, nil

	// Search
	case key.Matches(msg, m.keys.Search):
		m.previousState = m.state
		m.state = StateSearch
		m.search.Activate()
		return m, nil

	// Enter chat mode
	case key.Matches(msg, m.keys.EnterChat):
		m.previousState = m.state
		m.state = StateChatting
		m.textarea.Focus()
		m.updateViewportHeight()
		return m, nil

	// Yank content
	case key.Matches(msg, m.keys.Yank):
		if m.lastKeyWasY {
			// Double y - yank full content
			m.lastKeyWasY = false
			return m, YankContent(m.rawContent, m.chatHistory)
		}
		m.lastKeyWasY = true
		return m, nil

	// Yank last response
	case key.Matches(msg, m.keys.YankLast):
		return m, YankLastResponse(m.rawContent, m.chatHistory)

	// Navigation - let viewport handle it
	default:
		m.lastKeyWasY = false
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}
}
