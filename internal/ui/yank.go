package ui

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
)

// YankContent yanks content and chat history to clipboard
func YankContent(content string, chatHistory []ChatMessage) tea.Cmd {
	return func() tea.Msg {
		var sb strings.Builder

		if content != "" {
			sb.WriteString(content)
		}

		// Add chat history if present
		if len(chatHistory) > 0 {
			sb.WriteString("\n\n---\n\n## Follow-up Chat\n\n")
			for _, msg := range chatHistory {
				if msg.Role == ChatRoleUser {
					sb.WriteString("**You:** ")
					sb.WriteString(msg.Content)
				} else {
					sb.WriteString("**Assistant:**\n")
					sb.WriteString(msg.Content)
				}
				sb.WriteString("\n\n")
			}
		}

		result := sb.String()
		if result == "" {
			return nil
		}

		if err := clipboard.WriteAll(result); err != nil {
			return StreamErrorMsg{Err: fmt.Errorf("clipboard: %w", err)}
		}

		return YankMsg{Type: YankTypeContent}
	}
}

// YankLastResponse yanks only the last assistant response to clipboard
func YankLastResponse(content string, chatHistory []ChatMessage) tea.Cmd {
	return func() tea.Msg {
		var lastResponse string

		// Check chat history for last assistant message
		for i := len(chatHistory) - 1; i >= 0; i-- {
			if chatHistory[i].Role == ChatRoleAssistant {
				lastResponse = chatHistory[i].Content
				break
			}
		}

		// Fall back to main content if no chat history
		if lastResponse == "" {
			lastResponse = content
		}

		if lastResponse == "" {
			return nil
		}

		if err := clipboard.WriteAll(lastResponse); err != nil {
			return StreamErrorMsg{Err: fmt.Errorf("clipboard: %w", err)}
		}

		return YankMsg{Type: YankTypeLastResponse}
	}
}
