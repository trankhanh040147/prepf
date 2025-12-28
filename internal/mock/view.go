package mock

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/trankhanh040147/prepf/internal/config"
	"github.com/trankhanh040147/prepf/internal/ui"
)

// renderConfiguring renders configuration state
func (m *Model) renderConfiguring() string {
	width := m.Width()
	height := m.Height()
	if width == 0 {
		width = 80
	}
	if height == 0 {
		height = 20
	}

	// Render form if available
	formView := ""
	if m.configForm != nil {
		formView = m.configForm.View()
	}

	// Add skip instruction
	skipInstruction := "Press Esc to skip configuration"
	skipStyle := lipgloss.NewStyle().
		Italic(true).
		Padding(1, 0)
	if !m.noColor {
		skipStyle = skipStyle.Foreground(lipgloss.Color("8"))
	}

	content := formView
	if formView != "" {
		content = lipgloss.JoinVertical(lipgloss.Left, formView, "", skipStyle.Render(skipInstruction))
	} else {
		content = skipStyle.Render(skipInstruction)
	}

	// Center the content
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

// renderWaiting renders waiting state
func (m *Model) renderWaiting() string {
	// Ensure viewport has valid size
	width := m.Width()
	height := m.Height()
	if width == 0 {
		width = 80
	}
	if height == 0 {
		height = 20
	}
	m.Viewport().SetSize(width, height)

	// Set loading message and ensure it's visible at the top
	loadingMsg := "Loading interview context..."
	m.SetViewportContent(loadingMsg)
	m.Viewport().GotoTop()

	return m.Viewport().View()
}

// renderAIThinking renders AI thinking state
func (m *Model) renderAIThinking() string {
	statusBar := m.renderStatusBar()
	statusBarHeight := lipgloss.Height(statusBar)

	// Ensure viewport has valid dimensions
	width := m.Width()
	height := m.Height()
	if width == 0 {
		width = 80
	}
	if height == 0 {
		height = 20
	}

	// Adjust viewport height to account for status bar
	viewportHeight := height - statusBarHeight
	if viewportHeight < 1 {
		viewportHeight = 1
	}
	m.Viewport().SetSize(width, viewportHeight)

	// Ensure content is set from buffer (in case of resize after streaming)
	if m.aiResponseBuffer.Len() > 0 {
		m.SetViewportContent(m.aiResponseBuffer.String())
		m.Viewport().GotoBottom()
	}

	viewport := m.Viewport().View()
	return lipgloss.JoinVertical(lipgloss.Left, viewport, statusBar)
}

// renderUserInput renders user input state
func (m *Model) renderUserInput() string {
	statusBar := m.renderStatusBar()
	statusBarHeight := lipgloss.Height(statusBar)

	// Answer input
	inputView := m.answerInput.View()
	inputBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, config.InputBoxPadding).
		Width(m.Width())
	if !m.noColor {
		inputBoxStyle = inputBoxStyle.BorderForeground(lipgloss.Color("62"))
	}
	inputBox := inputBoxStyle.Render(inputView)
	inputBoxHeight := lipgloss.Height(inputBox)

	// Show surrender micro-roast if present
	microRoastHeight := 0
	microRoastText := ""
	if m.surrenderFeedback != "" {
		microRoastText = RenderMicroRoast(m.surrenderFeedback, m.noColor)
		microRoastHeight = lipgloss.Height(microRoastText) + 1 // +1 for spacing
		m.surrenderFeedback = ""                               // Clear after showing
	}

	// Ensure viewport has valid dimensions
	width := m.Width()
	height := m.Height()
	if width == 0 {
		width = 80
	}
	if height == 0 {
		height = 20
	}

	// Adjust viewport height to account for status bar, input box, and micro-roast
	viewportHeight := height - statusBarHeight - inputBoxHeight - microRoastHeight
	if viewportHeight < 1 {
		viewportHeight = 1
	}
	m.Viewport().SetSize(width, viewportHeight)

	// Ensure content is set after resize (viewport might lose content during resize)
	if m.currentQuestion != "" {
		m.SetViewportContent(m.currentQuestion)
		m.Viewport().GotoBottom()
	} else if m.aiResponseBuffer.Len() > 0 {
		// Fallback: if currentQuestion is empty but buffer has content, use buffer
		// This handles edge cases where stream completed but currentQuestion wasn't set yet
		content, _, _ := ParseSignals(m.aiResponseBuffer.String())
		m.SetViewportContent(content)
		m.Viewport().GotoBottom()
		m.currentQuestion = content
	}

	viewport := m.Viewport().View()

	// Build content
	content := viewport
	if microRoastText != "" {
		content = lipgloss.JoinVertical(lipgloss.Left, viewport, "", microRoastText)
	}

	return lipgloss.JoinVertical(lipgloss.Left, content, "", inputBox, statusBar)
}

// renderRoasting renders roasting state
func (m *Model) renderRoasting() string {
	roastView := RenderRoast(m.roastGrade, m.roastPersona, m.roastFeedback, m.Width(), m.noColor)
	buttonsView := RenderRemediationButtons(m.remediationTopics, m.Width(), m.noColor)
	return lipgloss.JoinVertical(lipgloss.Left, roastView, "", buttonsView)
}

// renderStatusBar renders the status bar
func (m *Model) renderStatusBar() string {
	elapsed := time.Since(m.sessionStartTime)
	minutes := int(elapsed.Minutes())
	seconds := int(elapsed.Seconds()) % 60

	status := fmt.Sprintf("Question %d/%d | Time: %02d:%02d", m.questionCount, MaxQuestions, minutes, seconds)

	if m.isSessionLimitReached() {
		alertStyle := lipgloss.NewStyle().Bold(true)
		if !m.noColor {
			alertStyle = alertStyle.Foreground(lipgloss.Color("9")) // Red
		}
		status += " | " + alertStyle.Render("[FINAL QUESTION]")
	}

	statusBarStyle := lipgloss.NewStyle().
		Width(m.Width()).
		BorderTop(true).
		Padding(0, config.StatusBarPadding)

	if !m.noColor {
		statusBarStyle = statusBarStyle.BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("8"))
	}

	return statusBarStyle.Render(status)
}

// renderHelp renders help overlay
func (m *Model) renderHelp() string {
	helpStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2)
	if !m.noColor {
		helpStyle = helpStyle.BorderForeground(lipgloss.Color("62"))
	}

	// Get mock-specific help text
	helpText := MockHelpText

	// Add key bindings
	allBindings := m.keys.FullHelp()
	helpLines := []string{}
	for _, group := range allBindings {
		for _, kb := range group {
			helpLines = append(helpLines, fmt.Sprintf("%-18s %s", kb.Help().Key, kb.Help().Desc))
		}
	}

	content := helpText + "\n\n" + strings.Join(helpLines, "\n")
	return ui.Center(helpStyle.Render(content), m.Width(), m.Height())
}

// toggleHelp toggles help visibility
func (m *Model) toggleHelp() {
	if m.BaseModel.State() == ui.StateHelp {
		m.BaseModel.SetState(ui.StateNormal)
	} else {
		m.BaseModel.SetState(ui.StateHelp)
	}
}

