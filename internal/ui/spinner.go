package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SpinnerModel wraps spinner for loading states
type SpinnerModel struct {
	spinner spinner.Model
	active  bool
	message string
}

// NewSpinner creates a new spinner model
func NewSpinner(noColor bool) *SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	if !noColor {
		s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	}

	return &SpinnerModel{
		spinner: s,
		active:  false,
		message: "Loading...",
	}
}

// SetMessage sets the spinner message
func (s *SpinnerModel) SetMessage(msg string) {
	s.message = msg
}

// Message returns the current spinner message
func (s *SpinnerModel) Message() string {
	return s.message
}

// Activate starts the spinner
func (s *SpinnerModel) Activate() {
	s.active = true
}

// Deactivate stops the spinner
func (s *SpinnerModel) Deactivate() {
	s.active = false
}

// IsActive returns whether the spinner is active
func (s *SpinnerModel) IsActive() bool {
	return s.active
}

// Update handles spinner messages
func (s *SpinnerModel) Update(msg tea.Msg) (*SpinnerModel, tea.Cmd) {
	if !s.active {
		return s, nil
	}

	var cmd tea.Cmd
	s.spinner, cmd = s.spinner.Update(msg)
	return s, cmd
}

// View renders the spinner
func (s *SpinnerModel) View() string {
	if !s.active {
		return ""
	}
	return s.spinner.View() + " " + s.message
}

// TickCmd returns a command that ticks the spinner
func (s *SpinnerModel) TickCmd() tea.Cmd {
	if !s.active {
		return nil
	}
	return s.spinner.Tick
}

// SetStyle sets the spinner style
func (s *SpinnerModel) SetStyle(style lipgloss.Style) {
	s.spinner.Style = style
}

