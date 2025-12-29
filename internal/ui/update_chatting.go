package ui

import (
	"context"

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
		// Create unified streaming channel
		m.streamMsgChan = make(chan StreamMsg, 10)

		// Start streaming in background
		go func() {
			ctx, cancel := context.WithCancel(m.rootCtx)
			if m.activeCancel != nil {
				m.activeCancel()
			}
			m.activeCancel = cancel

			fullResponse, err := m.client.SendMessageStream(ctx, message, func(chunk string) {
				m.streamMsgChan <- StreamMsg{Type: StreamMsgChunk, Chunk: chunk}
			})

			if err != nil {
				m.streamMsgChan <- StreamMsg{Type: StreamMsgError, Err: err}
				return
			}
			m.streamMsgChan <- StreamMsg{Type: StreamMsgDone, FullResponse: fullResponse}
		}()

		return StreamStartMsg{}
	}
}
