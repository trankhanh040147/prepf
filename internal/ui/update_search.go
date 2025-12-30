package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// updateKeyMsgSearch handles key messages in search state
func (m *Model) updateKeyMsgSearch(msg tea.KeyMsg) (*Model, tea.Cmd) {
	// Handle search-specific keys first
	switch {
	case key.Matches(msg, m.keys.SearchEsc):
		m.search.Deactivate()
		m.returnToPreviousState()
		return m, nil

	case key.Matches(msg, m.keys.SearchEnter):
		// Search confirmed, apply query
		query := m.search.Query()
		if query != "" {
			m.highlightSearchResults(query)
		}
		m.search.Deactivate()
		m.returnToPreviousState()
		return m, nil
	}

	// Pass to search model
	var cmd tea.Cmd
	m.search, cmd = m.search.Update(msg)

	// Check if search was completed/cancelled
	if !m.search.IsActive() {
		m.returnToPreviousState()
	}

	return m, cmd
}

// highlightSearchResults highlights search matches in content
func (m *Model) highlightSearchResults(query string) {
	// Store raw content if not already stored
	if m.rawContent == "" {
		m.rawContent = m.content
	}

	// TODO: Implement search highlighting
	// For now, just update the viewport
	m.updateViewport()
}

// clearSearchHighlight removes search highlighting
func (m *Model) clearSearchHighlight() {
	if m.rawContent != "" {
		m.content = m.rawContent
		m.updateViewport()
	}
}
