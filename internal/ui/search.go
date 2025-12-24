package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// SearchModel wraps textinput for search functionality
type SearchModel struct {
	textinput textinput.Model
	query     string
	active    bool
}

// NewSearch creates a new search model
func NewSearch() *SearchModel {
	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.CharLimit = 200
	ti.Width = 80

	return &SearchModel{
		textinput: ti,
		query:     "",
		active:    false,
	}
}

// Update handles search messages
func (s *SearchModel) Update(msg tea.Msg) (*SearchModel, tea.Cmd) {
	if !s.active {
		return s, nil
	}

	var cmd tea.Cmd
	s.textinput, cmd = s.textinput.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			s.query = s.textinput.Value()
			s.active = false
			s.textinput.Blur()
			return s, cmd
		case tea.KeyEsc:
			s.query = ""
			s.textinput.SetValue("")
			s.active = false
			s.textinput.Blur()
			return s, cmd
		}
	}

	return s, cmd
}

// View renders the search input
func (s *SearchModel) View() string {
	if !s.active {
		return ""
	}
	return s.textinput.View()
}

// Activate activates search mode
func (s *SearchModel) Activate() {
	s.active = true
	s.textinput.Focus()
	s.textinput.SetValue("")
}

// Deactivate deactivates search mode
func (s *SearchModel) Deactivate() {
	s.active = false
	s.textinput.Blur()
}

// IsActive returns whether search is active
func (s *SearchModel) IsActive() bool {
	return s.active
}

// Query returns the current search query
func (s *SearchModel) Query() string {
	return s.query
}

// SetWidth sets the search input width
func (s *SearchModel) SetWidth(width int) {
	s.textinput.Width = width
}
