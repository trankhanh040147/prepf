package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Stream messages for AI interactions

// StreamStartMsg signals that streaming has started
type StreamStartMsg struct{}

// StreamChunkMsg contains a chunk of streamed response
type StreamChunkMsg struct {
	Text string
}

// StreamDoneMsg signals that streaming is complete
type StreamDoneMsg struct {
	FullResponse string
}

// StreamErrorMsg contains an error from streaming
type StreamErrorMsg struct {
	Err error
}

// streamChunkCmd creates a command to listen for chunks from a channel
func streamChunkCmd(chunkChan chan string) tea.Cmd {
	return func() tea.Msg {
		chunk, ok := <-chunkChan
		if !ok {
			return nil
		}
		return StreamChunkMsg{Text: chunk}
	}
}

// streamDoneCmd creates a command to listen for completion from a channel
func streamDoneCmd(doneChan chan string) tea.Cmd {
	return func() tea.Msg {
		fullResponse, ok := <-doneChan
		if !ok {
			return nil
		}
		return StreamDoneMsg{FullResponse: fullResponse}
	}
}

// streamErrorCmd creates a command to listen for errors from a channel
func streamErrorCmd(errChan chan error) tea.Cmd {
	return func() tea.Msg {
		err, ok := <-errChan
		if !ok {
			return nil
		}
		return StreamErrorMsg{Err: err}
	}
}

// Chat messages

// ChatResponseMsg contains a response to a follow-up question
type ChatResponseMsg struct {
	Response string
}

// ChatErrorMsg contains an error from a chat interaction
type ChatErrorMsg struct {
	Err error
}

// ChatMessage represents a single message in chat history
type ChatMessage struct {
	Role    string
	Content string
}

// Chat roles
const (
	ChatRoleUser      = "user"
	ChatRoleAssistant = "assistant"
)

// UI feedback messages

// YankType represents the type of yanked content
type YankType int

const (
	YankTypeContent YankType = iota
	YankTypeLastResponse
)

// String returns the string representation of YankType
func (t YankType) String() string {
	switch t {
	case YankTypeContent:
		return "content"
	case YankTypeLastResponse:
		return "last response"
	default:
		return "unknown"
	}
}

// YankMsg signals that content was yanked to clipboard
type YankMsg struct {
	Type YankType
}

// YankFeedbackMsg signals that yank feedback should be cleared
type YankFeedbackMsg struct{}

// ClearYankFeedbackCmd creates a command to clear yank feedback after a delay
func ClearYankFeedbackCmd(duration time.Duration) tea.Cmd {
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		return YankFeedbackMsg{}
	})
}

// YankFeedbackDuration is the duration to show yank feedback
const YankFeedbackDuration = 2 * time.Second

// Navigation messages

// TickMsg is used for spinner animation
type TickMsg struct{}

// QuitMsg signals the program should quit
type QuitMsg struct{}

// ErrorMsg contains a general error
type ErrorMsg struct {
	Err error
}
