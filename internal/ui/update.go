package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Update handles messages and updates the model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle quit keys globally
		if key.Matches(msg, m.keys.ForceQuit) {
			if m.activeCancel != nil {
				m.activeCancel()
				m.activeCancel = nil
			}
			return m, tea.Quit
		}

		// Handle cancel request globally
		if key.Matches(msg, m.keys.CancelRequest) {
			if m.activeCancel != nil {
				m.activeCancel()
				m.activeCancel = nil
			}
			if m.state == StateLoading || m.streaming {
				m.resetStreamState()
				m.state = StateError
				m.errorMsg = "Request cancelled"
				return m, nil
			}
		}

		// Route to state-specific handlers
		switch m.state {
		case StateNormal:
			return m.updateKeyMsgNormal(msg)
		case StateSearch:
			return m.updateKeyMsgSearch(msg)
		case StateHelp:
			return m.updateKeyMsgHelp(msg)
		case StateChatting:
			return m.updateKeyMsgChatting(msg)
		case StateError:
			return m.updateKeyMsgError(msg)
		default:
			return m.updateNonKeyMsg(msg)
		}
	default:
		return m.updateNonKeyMsg(msg)
	}
}

// updateNonKeyMsg handles non-key messages
func (m *Model) updateNonKeyMsg(msg tea.Msg) (*Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.handleWindowSize(msg)

	case spinner.TickMsg:
		if cmd := m.handleSpinnerTick(msg); cmd != nil {
			cmds = append(cmds, cmd)
		}

	case StreamChunkMsg:
		return m.handleStreamChunk(msg)

	case StreamDoneMsg:
		return m.handleStreamDone(msg)

	case StreamErrorMsg:
		return m.handleStreamError(msg)

	case YankFeedbackMsg:
		m.yankFeedback = ""
	}

	// Update viewport
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	// Update textarea in chat mode
	if m.state == StateChatting && !m.streaming {
		var taCmd tea.Cmd
		m.textarea, taCmd = m.textarea.Update(msg)
		cmds = append(cmds, taCmd)
	}

	return m, tea.Batch(cmds...)
}

// handleWindowSize handles window resize messages
func (m *Model) handleWindowSize(msg tea.WindowSizeMsg) {
	m.width = msg.Width
	m.height = msg.Height

	m.updateViewportHeight()
	if !m.ready {
		m.viewport = viewport.New(msg.Width, m.viewport.Height)
		m.viewport.Style = lipgloss.NewStyle().Padding(0, 2)
		m.ready = true
	} else {
		m.viewport.Width = msg.Width
	}
	m.textarea.SetWidth(msg.Width - 4)
	m.search.SetWidth(msg.Width - 4)
}

// handleSpinnerTick handles spinner animation
func (m *Model) handleSpinnerTick(msg spinner.TickMsg) tea.Cmd {
	if m.state == StateLoading || m.streaming {
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return cmd
	}
	return nil
}
