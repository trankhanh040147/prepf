package mock

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/trankhanh040147/prepf/internal/ai"
	"github.com/trankhanh040147/prepf/internal/config"
	"github.com/trankhanh040147/prepf/internal/ui"
)

// Model represents the mock interview TUI model
type Model struct {
	*ui.BaseModel
	keys      MockKeyMap
	state     InterviewState
	aiClient  *ai.Client
	ctx       context.Context
	cancelCtx context.CancelFunc

	// Session metadata
	sessionStartTime time.Time
	questionCount    int
	surrenderCount   int

	// Context
	resumeContent string
	resumePath    string // Path to resume file (used for loading)
	contextLoaded bool

	// AI streaming
	aiResponseBuffer strings.Builder
	stream           <-chan ai.StreamChunk
	currentQuestion  string

	// User input
	answerInput   textinput.Model
	currentAnswer string

	// Roast data
	roastGrade        string
	roastPersona      string
	roastFeedback     string
	remediationTopics []string

	// Surrender micro-roast
	surrenderFeedback string
	isSurrenderMode   bool // Track if we're waiting for surrender micro-roast

	// Config
	noColor bool
}

// NewModel creates a new mock interview model
func NewModel(cfg *config.Config, aiClient *ai.Client, resumePath string) *Model {
	ctx, cancel := context.WithCancel(context.Background())

	ti := textinput.New()
	ti.Placeholder = "Type your answer here... (Enter to submit, Tab to surrender)"
	ti.CharLimit = 2000
	ti.Width = 80
	ti.Focus()

	base := ui.NewBaseModel(cfg)
	m := &Model{
		BaseModel:        base,
		keys:             DefaultMockKeyMap(),
		state:            InterviewWaiting,
		aiClient:         aiClient,
		ctx:              ctx,
		cancelCtx:        cancel,
		sessionStartTime: time.Now(),
		answerInput:      ti,
		noColor:          cfg.NoColor,
		resumePath:       resumePath,
	}

	return m
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		LoadContextCmd(m.resumePath),
		tea.Tick(time.Second, func(time.Time) tea.Msg {
			return TimeTickMsg{}
		}),
	)
}

// Update handles messages
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle window resize
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Update base model for window size
		m.BaseModel.Update(msg)
		// Update input width (account for border and padding)
		inputWidth := m.Width()
		if inputWidth > 4 {
			inputWidth -= 4
		}
		m.answerInput.Width = inputWidth
		return m, nil

	case ContextLoadedMsg:
		return m.handleContextLoaded(msg)

	case QuestionReceivedMsg:
		return m.handleQuestionReceived(msg)

	case AnswerSubmittedMsg:
		return m.handleAnswerSubmitted(msg)

	case SurrenderTriggeredMsg:
		return m.handleSurrender()

	case RoastTriggeredMsg:
		return m.handleRoast()

	case SessionExpiredMsg:
		return m.handleSessionExpired()

	case TimeTickMsg:
		return m.handleTimeTick()

	case ai.StreamStartedMsg:
		return m.handleStreamStarted(msg)

	case ai.StreamChunkMsg:
		return m.handleStreamChunk(msg)

	case ai.StreamDoneMsg:
		return m.handleStreamDone()

	case ai.StreamErrorMsg:
		return m.handleStreamError(msg)

	case tea.KeyMsg:
		if key.Matches(msg, m.keys.Quit) {
			m.cancelCtx()
			return m, tea.Quit
		}

		// Handle help
		if key.Matches(msg, m.keys.Help) {
			m.toggleHelp()
			return m, nil
		}

		// Handle state-specific keys
		switch m.state {
		case InterviewUserInput:
			if key.Matches(msg, m.keys.Surrender) {
				return m.handleSurrender()
			}
			if key.Matches(msg, m.keys.Enter) {
				return m.handleAnswerSubmit()
			}
			// Update text input
			var tiCmd tea.Cmd
			m.answerInput, tiCmd = m.answerInput.Update(msg)
			return m, tiCmd

		case InterviewRoasting:
			// Allow quit only
			if key.Matches(msg, m.keys.Quit) {
				m.cancelCtx()
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

// View renders the UI
func (m *Model) View() string {
	if m.BaseModel.State() == ui.StateHelp {
		return m.renderHelp()
	}

	switch m.state {
	case InterviewWaiting:
		return m.renderWaiting()
	case InterviewAIThinking:
		return m.renderAIThinking()
	case InterviewUserInput:
		return m.renderUserInput()
	case InterviewRoasting:
		return m.renderRoasting()
	default:
		return ""
	}
}

// handleContextLoaded processes loaded resume context
func (m *Model) handleContextLoaded(msg ContextLoadedMsg) (*Model, tea.Cmd) {
	if msg.Err != nil {
		// Continue without context
		m.contextLoaded = false
		m.resumeContent = ""
	} else {
		m.contextLoaded = true
		m.resumeContent = msg.Content
	}

	// Start the interview
	prompt := m.buildInitialPrompt()
	m.state = InterviewAIThinking
	m.aiResponseBuffer.Reset()
	return m, m.aiClient.StreamStartCmd(m.ctx, prompt)
}

// buildInitialPrompt builds the initial interview prompt
func (m *Model) buildInitialPrompt() string {
	if m.resumeContent != "" {
		return fmt.Sprintf(InitialPromptTemplate, m.resumeContent)
	}
	return "Conduct a technical interview. Ask one question at a time. When you want to move to the next question, include the signal <NEXT> at the end of your response. When you've finished all questions, include the signal <ROAST> at the end of your final response."
}

// handleStreamStarted processes stream start
func (m *Model) handleStreamStarted(msg ai.StreamStartedMsg) (*Model, tea.Cmd) {
	m.stream = msg.Stream
	return m, ai.WaitForStreamChunkCmd(m.stream)
}

// handleStreamChunk processes a stream chunk
func (m *Model) handleStreamChunk(msg ai.StreamChunkMsg) (*Model, tea.Cmd) {
	m.aiResponseBuffer.WriteString(msg.Text)
	m.SetViewportContent(m.aiResponseBuffer.String())
	m.Viewport().GotoBottom()
	return m, ai.WaitForStreamChunkCmd(m.stream)
}

// handleStreamDone processes stream completion
func (m *Model) handleStreamDone() (*Model, tea.Cmd) {
	fullResponse := m.aiResponseBuffer.String()

	// Check if this was a surrender (micro-roast) response
	if m.isSurrenderMode {
		return m.handleSurrenderStreamDone()
	}

	content, _, hasRoast := ParseSignals(fullResponse)

	// Ensure viewport has valid size before setting content
	// Use basic dimensions - render functions will set proper size
	width := m.Width()
	height := m.Height()
	if width == 0 {
		width = 80
	}
	if height == 0 {
		height = 20
	}
	// Set a basic size to ensure viewport is initialized (render will adjust)
	m.Viewport().SetSize(width, height)

	// Update viewport with cleaned content
	m.SetViewportContent(content)
	m.Viewport().GotoBottom()
	m.currentQuestion = content

	// Check for session limits
	if m.isSessionLimitReached() {
		return m, tea.Batch(
			func() tea.Msg { return SessionExpiredMsg{} },
		)
	}

	// Transition based on signals
	if hasRoast {
		return m, func() tea.Msg { return RoastTriggeredMsg{} }
	}

	// Move to user input state
	m.state = InterviewUserInput
	m.answerInput.Focus()
	m.answerInput.SetValue("")
	m.currentAnswer = ""

	return m, nil
}

// handleStreamError processes stream errors
func (m *Model) handleStreamError(msg ai.StreamErrorMsg) (*Model, tea.Cmd) {
	errorText := fmt.Sprintf("Error: %v\n\nPress q to quit.", msg.Err)
	m.SetViewportContent(errorText)
	m.Viewport().GotoBottom()
	return m, nil
}

// handleAnswerSubmit submits the user's answer
func (m *Model) handleAnswerSubmit() (*Model, tea.Cmd) {
	answer := strings.TrimSpace(m.answerInput.Value())
	if answer == "" {
		return m, nil
	}

	m.currentAnswer = answer
	m.questionCount++
	m.state = InterviewAIThinking
	m.aiResponseBuffer.Reset()
	m.answerInput.Blur()

	// Send answer to AI
	return m, m.aiClient.StreamStartCmd(m.ctx, answer)
}

// handleAnswerSubmitted processes submitted answer (alternative entry point)
func (m *Model) handleAnswerSubmitted(msg AnswerSubmittedMsg) (*Model, tea.Cmd) {
	if msg.Answer == "" {
		return m, nil
	}

	m.currentAnswer = msg.Answer
	m.questionCount++
	m.state = InterviewAIThinking
	m.aiResponseBuffer.Reset()
	m.answerInput.Blur()

	return m, m.aiClient.StreamStartCmd(m.ctx, msg.Answer)
}

// handleQuestionReceived processes a received question
func (m *Model) handleQuestionReceived(msg QuestionReceivedMsg) (*Model, tea.Cmd) {
	m.currentQuestion = msg.Question
	m.SetViewportContent(msg.Question)

	if msg.HasRoast {
		return m, func() tea.Msg { return RoastTriggeredMsg{} }
	}

	m.state = InterviewUserInput
	m.answerInput.Focus()
	return m, nil
}

// handleSurrender processes surrender action
func (m *Model) handleSurrender() (*Model, tea.Cmd) {
	m.surrenderCount++
	m.questionCount++

	// Generate micro-roast
	m.state = InterviewAIThinking
	m.isSurrenderMode = true
	m.aiResponseBuffer.Reset()
	surrenderPrompt := ShadowPromptSurrender

	return m, m.aiClient.StreamStartCmd(m.ctx, surrenderPrompt)
}

// handleSurrenderStreamDone handles completion of surrender stream (micro-roast)
func (m *Model) handleSurrenderStreamDone() (*Model, tea.Cmd) {
	microRoast := m.aiResponseBuffer.String()
	m.surrenderFeedback = strings.TrimSpace(microRoast)
	m.isSurrenderMode = false

	// Move to next question
	prompt := "Continue with the next question."
	m.aiResponseBuffer.Reset()
	return m, m.aiClient.StreamStartCmd(m.ctx, prompt)
}

// handleRoast processes roast trigger
func (m *Model) handleRoast() (*Model, tea.Cmd) {
	m.state = InterviewRoasting
	m.roastGrade = CalculateGrade(m.surrenderCount, m.questionCount)
	m.roastPersona = GetPersonaLabel(m.roastGrade)
	m.roastFeedback = "Great job completing the interview!"
	m.remediationTopics = []string{"Data Structures", "Algorithms", "System Design"}
	m.answerInput.Blur()
	return m, nil
}

// handleSessionExpired processes session expiry
func (m *Model) handleSessionExpired() (*Model, tea.Cmd) {
	return m.handleRoast()
}

// handleTimeTick processes time tick for session timer
func (m *Model) handleTimeTick() (*Model, tea.Cmd) {
	// Check if session expired
	if m.isSessionLimitReached() && m.state != InterviewRoasting {
		return m, func() tea.Msg { return SessionExpiredMsg{} }
	}

	// Continue ticking
	return m, tea.Tick(time.Second, func(time.Time) tea.Msg {
		return TimeTickMsg{}
	})
}

// isSessionLimitReached checks if session limits are reached
func (m *Model) isSessionLimitReached() bool {
	if m.questionCount >= MaxQuestions {
		return true
	}
	elapsed := time.Since(m.sessionStartTime)
	if elapsed >= MaxDurationMinutes*time.Minute {
		return true
	}
	return false
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

	// Set loading message if viewport is empty
	loadingMsg := "Loading interview context..."
	m.SetViewportContent(loadingMsg)

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
		Padding(0, 1).
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
		Padding(0, 1)

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
