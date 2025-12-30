package ui

import tea "github.com/charmbracelet/bubbletea"

// handleStreamChunk handles incoming stream chunks
func (m *Model) handleStreamChunk(msg StreamChunkMsg) (*Model, tea.Cmd) {
	// Append chunk to content
	m.content += msg.Text
	m.updateViewportAndScroll()

	// Continue listening for more chunks
	return m, m.waitForStreamChunk()
}

// handleStreamDone handles stream completion
func (m *Model) handleStreamDone(msg StreamDoneMsg) (*Model, tea.Cmd) {
	m.streaming = false

	// Add assistant response to chat history
	m.chatHistory = append(m.chatHistory, ChatMessage{
		Role:    ChatRoleAssistant,
		Content: msg.FullResponse,
	})

	// Refocus textarea
	m.textarea.Focus()
	m.updateViewportAndScroll()

	return m, nil
}

// handleStreamError handles stream errors
func (m *Model) handleStreamError(msg StreamErrorMsg) (*Model, tea.Cmd) {
	m.resetStreamState()
	m.state = StateError
	m.errorMsg = msg.Err.Error()
	return m, nil
}

// waitForStreamChunk creates a command to wait for the next stream message
func (m *Model) waitForStreamChunk() tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-m.streamMsgChan
		if !ok {
			return nil
		}
		switch msg.Type {
		case StreamMsgChunk:
			return StreamChunkMsg{Text: msg.Chunk}
		case StreamMsgDone:
			return StreamDoneMsg{FullResponse: msg.FullResponse}
		case StreamMsgError:
			return StreamErrorMsg{Err: msg.Err}
		}
		return nil
	}
}
