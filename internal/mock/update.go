package mock

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/trankhanh040147/prepf/internal/ai"
)

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

	case ConfigSubmittedMsg:
		return m.handleConfigSubmitted(msg)

	case ConfigSkippedMsg:
		return m.handleConfigSkipped()

	case ai.StreamStartedMsg:
		return m.handleStreamStarted(msg)

	case ai.StreamChunkMsg:
		return m.handleStreamChunk(msg)

	case ai.StreamDoneMsg:
		return m.handleStreamDone()

	case ai.StreamErrorMsg:
		return m.handleStreamError(msg)

	case tea.KeyMsg:
		// Handle state-specific keys first to allow form to process all keys
		switch m.state {
		case InterviewConfiguring:
			// Ensure answerInput is not focused - form should receive all keys
			m.answerInput.Blur()

			// Handle skip key (Esc) - intercept this before form
			if key.Matches(msg, m.keys.Skip) {
				return m.handleConfigSkipped()
			}
			// Handle global keys (Quit, Help) - only check these, let form handle everything else
			if key.Matches(msg, m.keys.Quit) {
				m.cancelCtx()
				return m, tea.Quit
			}
			if key.Matches(msg, m.keys.Help) {
				m.toggleHelp()
				return m, nil
			}
			// Let form handle ALL keys (Enter, Tab, Space, arrows, etc.)
			// The form will handle navigation, selection, and submission internally
			if m.configForm != nil {
				var cmd tea.Cmd
				formModel, cmd := m.configForm.Update(msg)
				// Type assert back to *huh.Form
				if form, ok := formModel.(*huh.Form); ok {
					m.configForm = form
					// Check if form completed after any key press
					// Form completes when user presses Enter on the last field
					if m.configForm.State == huh.StateCompleted {
						return m.handleConfigSubmitted(ConfigSubmittedMsg{})
					}
				}
				return m, cmd
			}
			return m, nil

		case InterviewUserInput:
			if key.Matches(msg, m.keys.Surrender) {
				return m.handleSurrender()
			}
			if key.Matches(msg, m.keys.Enter) {
				return m.handleAnswerSubmit()
			}
			// Update text input (handles all other keys including 'q')
			var tiCmd tea.Cmd
			m.answerInput, tiCmd = m.answerInput.Update(msg)
			return m, tiCmd

		case InterviewRoasting:
			// Allow quit only when input not focused
			if !m.answerInput.Focused() && key.Matches(msg, m.keys.Quit) {
				m.cancelCtx()
				return m, tea.Quit
			}

		default:
			// Handle global keys (Quit, Help) for other states
			if !m.answerInput.Focused() {
				if key.Matches(msg, m.keys.Quit) {
					m.cancelCtx()
					return m, tea.Quit
				}
				if key.Matches(msg, m.keys.Help) {
					m.toggleHelp()
					return m, nil
				}
			}
		}
	}

	return m, nil
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
	topicInstructions := BuildTopicInstructions(m.selectedTopics, m.excludedTopics)

	if m.resumeContent != "" {
		return fmt.Sprintf(InitialPromptTemplate, m.resumeContent, topicInstructions)
	}

	// Fallback prompt without resume
	basePrompt := "Conduct a technical interview following these guidelines:\n- Ask questions that a real interviewer would ask for this role/experience level\n- Vary question types (conceptual, practical, problem-solving)\n- Avoid repeating similar questions. Reference conversation history to ensure variety\n- Ask follow-up questions based on the candidate's answers, not generic questions\n- Ask one question at a time\n- When you want to move to the next question, include the signal <NEXT> at the end of your response\n- When you've finished all questions, include the signal <ROAST> at the end of your final response."

	if topicInstructions != "" {
		return topicInstructions + "\n\n" + basePrompt
	}
	return basePrompt
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

	// Add variety instruction to ensure questions don't repeat
	// The AI client maintains conversation history, so this instruction
	// will be applied in context of all previous questions
	varietyInstruction := "\n\n[Based on the conversation so far, ask a different type of question. Avoid repetition. Reference what has been discussed to ensure variety.]"
	answerWithContext := answer + varietyInstruction

	// Send answer to AI
	return m, m.aiClient.StreamStartCmd(m.ctx, answerWithContext)
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

	// Add variety instruction to ensure questions don't repeat
	varietyInstruction := "\n\n[Based on the conversation so far, ask a different type of question. Avoid repetition. Reference what has been discussed to ensure variety.]"
	answerWithContext := msg.Answer + varietyInstruction

	return m, m.aiClient.StreamStartCmd(m.ctx, answerWithContext)
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
	return elapsed >= MaxDurationMinutes*time.Minute
}

// handleConfigSubmitted processes configuration form submission
func (m *Model) handleConfigSubmitted(msg ConfigSubmittedMsg) (*Model, tea.Cmd) {
	// The form has already updated m.selectedTopics and m.excludedTopics via pointers
	// No need to copy from message - form mutates model slices directly
	// Transition to waiting state and load context
	m.state = InterviewWaiting
	return m, LoadContextCmd(m.resumePath)
}

// handleConfigSkipped processes configuration skip
func (m *Model) handleConfigSkipped() (*Model, tea.Cmd) {
	// Set empty topic lists (all topics allowed)
	m.selectedTopics = make([]string, 0)
	m.excludedTopics = make([]string, 0)
	m.skipConfig = true

	// Transition to waiting state and load context
	m.state = InterviewWaiting
	return m, LoadContextCmd(m.resumePath)
}
