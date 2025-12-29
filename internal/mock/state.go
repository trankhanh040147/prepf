package mock

import tea "github.com/charmbracelet/bubbletea"

// InterviewState represents the current state of the interview
type InterviewState int

const (
	// InterviewConfiguring is the state for pre-mock topic configuration
	InterviewConfiguring InterviewState = iota
	// InterviewWaiting is the initial state, waiting for context to load
	InterviewWaiting
	// InterviewAIThinking indicates AI is generating a question/response
	InterviewAIThinking
	// InterviewUserInput indicates waiting for user to submit an answer
	InterviewUserInput
	// InterviewRoasting indicates showing the final roast/grade
	InterviewRoasting
)

// ContextLoadedMsg is sent when resume context has been loaded
type ContextLoadedMsg struct {
	Content string
	Err     error
}

// QuestionReceivedMsg is sent when a question has been fully received from AI
type QuestionReceivedMsg struct {
	Question string
	HasNext  bool
	HasRoast bool
}

// AnswerSubmittedMsg is sent when user submits an answer
type AnswerSubmittedMsg struct {
	Answer string
}

// SurrenderTriggeredMsg is sent when user presses Tab to surrender
type SurrenderTriggeredMsg struct{}

// RoastTriggeredMsg is sent when it's time to show the roast
type RoastTriggeredMsg struct{}

// SessionExpiredMsg is sent when session time or question limit is reached
type SessionExpiredMsg struct{}

// TimeTickMsg is sent periodically to update the session timer
type TimeTickMsg struct{}

// ConfigSubmittedMsg is sent when configuration form is submitted
// Form updates model slices directly via pointers, so no payload needed
type ConfigSubmittedMsg struct{}

// ConfigSkippedMsg is sent when configuration is skipped
type ConfigSkippedMsg struct{}

// LoadContextCmd returns a tea.Cmd that loads context from a file
func LoadContextCmd(path string) tea.Cmd {
	return func() tea.Msg {
		content, err := LoadContext(path)
		return ContextLoadedMsg{Content: content, Err: err}
	}
}
