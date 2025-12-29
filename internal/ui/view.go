package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"
)

// View renders the UI
func (m *Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	switch m.state {
	case StateHelp:
		return m.viewHelp()
	case StateError:
		return m.viewError()
	case StateLoading:
		return m.viewLoading()
	case StateSearch:
		return m.viewSearch()
	case StateChatting:
		return m.viewChatting()
	default:
		return m.viewNormal()
	}
}

// viewNormal renders the normal state
func (m *Model) viewNormal() string {
	var sections []string

	// Main content
	sections = append(sections, m.viewport.View())

	// Status bar
	status := m.renderStatusBar()
	sections = append(sections, status)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// viewSearch renders the search state
func (m *Model) viewSearch() string {
	searchView := m.search.View()
	searchPrompt := "/ " + searchView

	// Style the search prompt
	searchStyle := lipgloss.NewStyle().
		Width(m.width).
		Padding(0, 1).
		BorderBottom(true).
		BorderStyle(lipgloss.RoundedBorder())

	if !m.noColor {
		searchStyle = searchStyle.BorderForeground(ColorPrimary)
	}

	styledSearch := searchStyle.Render(searchPrompt)

	// Combine viewport with search at bottom
	return lipgloss.JoinVertical(lipgloss.Left,
		m.viewport.View(),
		styledSearch,
	)
}

// viewHelp renders the help overlay
func (m *Model) viewHelp() string {
	helpStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2)

	if !m.noColor {
		helpStyle = helpStyle.BorderForeground(ColorPrimary)
	}

	// Flatten nested key binding groups and map to formatted strings
	allBindings := lo.Flatten(m.keys.FullHelp())
	helpLines := lo.Map(allBindings, func(kb key.Binding, _ int) string {
		return fmt.Sprintf("%-18s %s", kb.Help().Key, kb.Help().Desc)
	})

	content := HelpText() + "\n\n" + strings.Join(helpLines, "\n")

	return Center(helpStyle.Render(content), m.width, m.height)
}

// viewError renders the error state
func (m *Model) viewError() string {
	errorStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2)

	if !m.noColor {
		errorStyle = errorStyle.
			BorderForeground(ColorError).
			Foreground(ColorError)
	}

	content := fmt.Sprintf("Error: %s\n\nPress any key to continue...", m.errorMsg)
	return Center(errorStyle.Render(content), m.width, m.height)
}

// viewLoading renders the loading state
func (m *Model) viewLoading() string {
	loadingStyle := lipgloss.NewStyle().Padding(2)

	content := fmt.Sprintf("%s Loading...", m.spinner.View())
	return Center(loadingStyle.Render(content), m.width, m.height)
}

// viewChatting renders the chatting state
func (m *Model) viewChatting() string {
	var sections []string

	// Main content (viewport with chat history)
	sections = append(sections, m.viewport.View())

	// Chat input area
	inputStyle := lipgloss.NewStyle().Padding(0, 1)

	var inputContent string
	if m.streaming {
		inputContent = fmt.Sprintf("%s Generating response...", m.spinner.View())
	} else {
		inputContent = m.textarea.View()
	}

	sections = append(sections, inputStyle.Render(inputContent))

	// Status bar
	status := m.renderChatStatusBar()
	sections = append(sections, status)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderStatusBar renders the status bar
func (m *Model) renderStatusBar() string {
	statusStyle := lipgloss.NewStyle().
		Padding(0, 1)

	if !m.noColor {
		statusStyle = statusStyle.Foreground(ColorSecondary)
	}

	var parts []string

	// Yank feedback
	if m.yankFeedback != "" {
		parts = append(parts, m.yankFeedback)
	}

	// Help hint
	parts = append(parts, "? help • q quit • / search • enter chat")

	return statusStyle.Render(strings.Join(parts, " • "))
}

// renderChatStatusBar renders the chat mode status bar
func (m *Model) renderChatStatusBar() string {
	statusStyle := lipgloss.NewStyle().
		Padding(0, 1)

	if !m.noColor {
		statusStyle = statusStyle.Foreground(ColorSecondary)
	}

	parts := []string{
		"esc exit",
		"alt+enter send",
		"ctrl+x cancel",
	}

	return statusStyle.Render(strings.Join(parts, " • "))
}
