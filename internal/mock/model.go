package mock

import (
	"context"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
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
	surrenderFeedback     string
	isSurrenderMode       bool // Track if we're waiting for surrender micro-roast
	showSurrenderFeedback bool // Track if micro-roast should be displayed

	// Topic configuration
	selectedTopics []string
	excludedTopics []string
	configForm     *huh.Form
	skipConfig     bool

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
	// Don't focus text input initially - will be focused when entering InterviewUserInput state

	base := ui.NewBaseModel(cfg)
	m := &Model{
		BaseModel:        base,
		keys:             DefaultMockKeyMap(),
		state:            InterviewConfiguring,
		aiClient:         aiClient,
		ctx:              ctx,
		cancelCtx:        cancel,
		sessionStartTime: time.Now(),
		answerInput:      ti,
		noColor:          cfg.NoColor,
		resumePath:       resumePath,
		selectedTopics:   make([]string, 0),
		excludedTopics:   make([]string, 0),
		skipConfig:       false,
	}

	return m
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	// Initialize configuration form if in configuring state
	if m.state == InterviewConfiguring {
		// Ensure answerInput is blurred so form can receive keys
		m.answerInput.Blur()
		m.configForm = buildConfigForm(&m.selectedTopics, &m.excludedTopics)
		return tea.Batch(
			m.configForm.Init(),
			tea.Tick(time.Second, func(time.Time) tea.Msg {
				return TimeTickMsg{}
			}),
		)
	}

	// Otherwise, start loading context (for InterviewWaiting state)
	return tea.Batch(
		LoadContextCmd(m.resumePath),
		tea.Tick(time.Second, func(time.Time) tea.Msg {
			return TimeTickMsg{}
		}),
	)
}

// View renders the UI
func (m *Model) View() string {
	if m.BaseModel.State() == ui.StateHelp {
		return m.renderHelp()
	}

	switch m.state {
	case InterviewConfiguring:
		return m.renderConfiguring()
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
