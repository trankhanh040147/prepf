package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"
	"github.com/trankhanh040147/prepf/internal/config"
)

// State represents UI state
type State int

const (
	StateNormal State = iota
	StateHelp
	StateLoading
	StateSearch
)

// BaseModel provides base functionality for TUI models
type BaseModel struct {
	keys          KeyMap
	state         State
	previousState State
	helpVisible   bool
	viewport      *ViewportModel
	search        *SearchModel
	width         int
	height        int
	noColor       bool
}

// NewBaseModel creates a new base model
func NewBaseModel(cfg *config.Config) *BaseModel {
	return &BaseModel{
		keys:          DefaultKeyMap(),
		state:         StateNormal,
		previousState: StateNormal,
		helpVisible:   false,
		viewport:      NewViewport(),
		search:        NewSearch(),
		width:         config.DefaultMinWidth,
		height:        20,
		noColor:       cfg.NoColor,
	}
}

// Update handles base messages
func (m *BaseModel) Update(msg tea.Msg) (*BaseModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.SetSize(msg.Width, msg.Height)
		if m.search != nil {
			m.search.SetWidth(msg.Width)
		}
		return m, nil

	case tea.KeyMsg:
		if key.Matches(msg, m.keys.Help) && m.state != StateSearch {
			m.toggleHelp()
			return m, nil
		}
		if key.Matches(msg, m.keys.Search) && m.state == StateNormal {
			m.previousState = m.state
			m.state = StateSearch
			m.search.Activate()
			var cmd tea.Cmd
			m.search, cmd = m.search.Update(msg)
			return m, cmd
		}
		if key.Matches(msg, m.keys.Quit) && m.state == StateNormal {
			return m, tea.Quit
		}
	}

	// Handle search state
	if m.state == StateSearch && m.search != nil {
		var cmd tea.Cmd
		m.search, cmd = m.search.Update(msg)
		if !m.search.IsActive() {
			// Search completed or cancelled
			m.state = m.previousState
			// Query is available via m.search.Query() for parent models to use
		}
		return m, cmd
	}

	return m, nil
}

// toggleHelp toggles help visibility
func (m *BaseModel) toggleHelp() {
	if m.helpVisible {
		m.helpVisible = false
		m.state = m.previousState
	} else {
		m.previousState = m.state
		m.helpVisible = true
		m.state = StateHelp
	}
}

// returnToPreviousState returns to previous state
func (m *BaseModel) returnToPreviousState() {
	m.state = m.previousState
	m.helpVisible = false
}

// View renders base UI
func (m *BaseModel) View() string {
	if m.helpVisible {
		return m.renderHelp()
	}
	if m.state == StateSearch && m.search != nil {
		return m.renderSearch()
	}
	return m.viewport.View()
}

// renderSearch renders search input
func (m *BaseModel) renderSearch() string {
	searchView := m.search.View()
	searchPrompt := "/ " + searchView

	// Style the search prompt
	searchStyle := lipgloss.NewStyle().Width(m.Width()).Padding(0, 1).BorderBottom(true)
	if m.noColor {
		searchStyle = searchStyle.BorderStyle(lipgloss.NormalBorder())
	} else {
		searchStyle = searchStyle.
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))
	}

	styledSearch := searchStyle.Render(searchPrompt)

	// Combine viewport with search at bottom
	viewportHeight := m.height - lipgloss.Height(styledSearch)
	if viewportHeight < 1 {
		viewportHeight = 1
	}
	m.viewport.SetSize(m.width, viewportHeight)

	return lipgloss.JoinVertical(lipgloss.Left,
		m.viewport.View(),
		styledSearch,
	)
}

// SearchQuery returns the current search query (empty if not searching)
func (m *BaseModel) SearchQuery() string {
	if m.search == nil {
		return ""
	}
	return m.search.Query()
}

// renderHelp renders help overlay
func (m *BaseModel) renderHelp() string {
	helpStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2)
	if !m.noColor {
		helpStyle = helpStyle.BorderForeground(lipgloss.Color("62"))
	}

	// Flatten nested key binding groups and map to formatted strings
	allBindings := lo.Flatten(m.keys.FullHelp())
	helpLines := lo.Map(allBindings, func(kb key.Binding, _ int) string {
		return fmt.Sprintf("%-18s %s", kb.Help().Key, kb.Help().Desc)
	})

	content := HelpText() + "\n\n" + strings.Join(helpLines, "\n")

	return Center(helpStyle.Render(content), m.width, m.height)
}

// SetState sets the current state
func (m *BaseModel) SetState(state State) {
	m.previousState = m.state
	m.state = state
}

// State returns current state
func (m *BaseModel) State() State {
	return m.state
}

// Width returns terminal width
func (m *BaseModel) Width() int {
	if m.width == 0 {
		return config.DefaultMinWidth
	}
	return m.width
}

// Height returns terminal height
func (m *BaseModel) Height() int {
	return m.height
}
