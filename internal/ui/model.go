package ui

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/trankhanh040147/prepf/internal/config"
	"github.com/trankhanh040147/prepf/internal/gemini"
)

// Model represents the main application state for Bubbletea
type Model struct {
	// State machine
	state         State
	previousState State

	// Gemini client
	client       *gemini.Client
	rootCtx      context.Context
	activeCancel context.CancelFunc
	apiKey       string

	// UI components
	spinner  spinner.Model
	viewport viewport.Model
	textarea textarea.Model
	search   *SearchModel
	renderer *Renderer
	keys     KeyMap

	// Content
	content     string
	rawContent  string // Original content without search highlighting
	errorMsg    string
	chatHistory []ChatMessage

	// Dimensions
	width  int
	height int

	// Flags
	ready     bool
	streaming bool
	noColor   bool

	// Streaming channels
	streamChunkChan chan string
	streamErrChan   chan error
	streamDoneChan  chan string

	// Yank state
	yankFeedback string
	lastKeyWasY  bool
}

// NewModel creates a new application model
func NewModel(cfg *config.Config, client *gemini.Client) *Model {
	// Create spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(ColorPrimary)

	// Create textarea for chat input
	ta := textarea.New()
	ta.Placeholder = "Ask a follow-up question..."
	ta.Focus()
	ta.CharLimit = 1000
	ta.SetWidth(80)
	ta.SetHeight(3)
	ta.ShowLineNumbers = false

	// Custom textarea styling
	ta.FocusedStyle.Base = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ColorPrimary)
	ta.BlurredStyle.Base = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ColorSecondary)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.BlurredStyle.CursorLine = lipgloss.NewStyle()

	// Create renderer
	renderer, err := NewRenderer()
	if err != nil {
		// Log warning but continue with empty renderer
		renderer = &Renderer{}
	}

	// Create root context
	rootCtx := context.Background()

	return &Model{
		state:       StateNormal,
		client:      client,
		rootCtx:     rootCtx,
		spinner:     s,
		textarea:    ta,
		search:      NewSearch(),
		renderer:    renderer,
		keys:        DefaultKeyMap(),
		noColor:     cfg.NoColor,
		ready:       false,
		streaming:   false,
		chatHistory: []ChatMessage{},
	}
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return m.spinner.Tick
}

// Run starts the Bubbletea program
func Run(cfg *config.Config, client *gemini.Client) error {
	model := NewModel(cfg, client)
	program := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		return fmt.Errorf("run ui: %w", err)
	}

	return nil
}

// resetStreamState resets streaming state and clears all stream channels
func (m *Model) resetStreamState() {
	m.streaming = false
	m.streamChunkChan = nil
	m.streamErrChan = nil
	m.streamDoneChan = nil
}

// returnToPreviousState returns to the previous state
func (m *Model) returnToPreviousState() {
	m.state = m.previousState
	m.updateViewportHeight()
}

// updateViewportHeight recalculates viewport height based on current state
func (m *Model) updateViewportHeight() {
	headerHeight := 2
	footerHeight := 2

	if m.state == StateChatting {
		footerHeight += 5 // Extra space for textarea
	}

	newHeight := m.height - headerHeight - footerHeight
	if newHeight < MinHeight {
		newHeight = MinHeight
	}

	m.viewport.Height = newHeight
}

// updateViewport updates the viewport content
func (m *Model) updateViewport() {
	rendered := m.renderer.Render(m.content)
	m.viewport.SetContent(rendered)
}

// updateViewportAndScroll updates viewport and scrolls to bottom
func (m *Model) updateViewportAndScroll() {
	m.updateViewport()
	m.viewport.GotoBottom()
}
