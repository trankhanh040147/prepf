package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// updateKeyMsgChatting handles key messages in chatting state
func (m *Model) updateKeyMsgChatting(msg tea.KeyMsg) (*Model, tea.Cmd) {
	// Don't process keys while streaming
	if m.streaming {
		return m, nil
	}

	switch {
	// Exit chat mode
	case key.Matches(msg, m.keys.ExitChat):
		m.textarea.Blur()
		m.returnToPreviousState()
		return m, nil

	// Send message
	case key.Matches(msg, m.keys.SendMessage):
		message := m.textarea.Value()
		if message == "" {
			return m, nil
		}

		// Add user message to history
		m.chatHistory = append(m.chatHistory, ChatMessage{
			Role:    ChatRoleUser,
			Content: message,
		})

		// Clear textarea
		m.textarea.SetValue("")

		// Start streaming
		m.streaming = true
		return m, m.startStreamingChat(message)
	}

	// Pass to textarea
	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

// startStreamingChat initiates a streaming chat request
func (m *Model) startStreamingChat(message string) tea.Cmd {
	return func() tea.Msg {
		// Create channels for streaming
		m.streamChunkChan = make(chan string, 10)
		m.streamErrChan = make(chan error, 1)
		m.streamDoneChan = make(chan string, 1)

		// Start streaming in background
		go func() {
			ctx, cancel := m.rootCtx, func() {}
			if m.activeCancel != nil {
				m.activeCancel()
			}
			m.activeCancel = cancel

			fullResponse, err := m.client.SendMessageStream(ctx, message, func(chunk string) {
				m.streamChunkChan <- chunk
			})

			if err != nil {
				m.streamErrChan <- err
				return
			}
			m.streamDoneChan <- fullResponse
		}()

		return StreamStartMsg{}
	}
}
