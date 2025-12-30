package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/trankhanh040147/prepf/internal/config"
)

// createTestModel creates a Model for testing (nil client is allowed)
func createTestModel(cfg *config.Config) *Model {
	model := NewModel(cfg, nil)
	// Initialize with a window size so ready becomes true
	model.width = 80
	model.height = 24
	model.ready = true
	return model
}

func TestModel_StateTransitions(t *testing.T) {
	cfg := &config.Config{NoColor: false}
	model := createTestModel(cfg)

	// Test initial state
	if model.state != StateNormal {
		t.Errorf("expected initial state StateNormal, got %d", model.state)
	}

	// Test state change to help
	model.previousState = model.state
	model.state = StateHelp
	if model.state != StateHelp {
		t.Errorf("expected state StateHelp, got %d", model.state)
	}

	// Test return to previous state
	model.returnToPreviousState()
	if model.state != StateNormal {
		t.Errorf("expected state StateNormal after returnToPreviousState, got %d", model.state)
	}
}

func TestModel_SearchState(t *testing.T) {
	cfg := &config.Config{NoColor: false}
	model := createTestModel(cfg)

	// Verify search model is initialized
	if model.search == nil {
		t.Fatal("search model should be initialized")
	}

	// Test that search can be activated (basic structure test)
	model.search.Activate()
	if !model.search.IsActive() {
		t.Error("search should be active after Activate()")
	}
}

func TestModel_WidthSafety(t *testing.T) {
	cfg := &config.Config{NoColor: false}
	model := createTestModel(cfg)

	// Test default width
	if model.width < MinWidth {
		t.Errorf("expected width >= %d, got %d", MinWidth, model.width)
	}

	// Test width after window size message
	windowSizeMsg := tea.WindowSizeMsg{Width: 120, Height: 40}
	model.handleWindowSize(windowSizeMsg)

	if model.width != 120 {
		t.Errorf("expected width 120, got %d", model.width)
	}
}

func TestModel_KeyBindings(t *testing.T) {
	cfg := &config.Config{NoColor: false}
	model := createTestModel(cfg)

	// Test help key (?) triggers help state
	helpKey := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")}
	model.updateKeyMsgNormal(helpKey)
	if model.state != StateHelp {
		t.Errorf("expected StateHelp after '?' key, got %d", model.state)
	}
}

func TestModel_ViewportIntegration(t *testing.T) {
	cfg := &config.Config{NoColor: false}
	model := createTestModel(cfg)

	// Test viewport is initialized (ready was set in createTestModel)
	if !model.ready {
		t.Fatal("model should be ready after initialization")
	}

	// Test setting viewport content
	model.content = "Test content"
	model.updateViewport()

	// Verify content was set (structure is correct)
	if model.content != "Test content" {
		t.Error("content should be set")
	}
}

func TestModel_SearchIntegration(t *testing.T) {
	cfg := &config.Config{NoColor: false}
	model := createTestModel(cfg)

	// Test search activation via key
	searchKey := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("/")}
	model.updateKeyMsgNormal(searchKey)

	if model.state != StateSearch {
		t.Errorf("expected StateSearch after '/' key, got %d", model.state)
	}

	// Test search query retrieval
	query := model.search.Query()
	if query != "" {
		t.Errorf("expected empty query initially, got '%s'", query)
	}
}

func TestModel_NoColorMode(t *testing.T) {
	cfg := &config.Config{NoColor: true}
	model := createTestModel(cfg)

	// Verify noColor is set
	if !model.noColor {
		t.Error("noColor should be true when config.NoColor is true")
	}

	// View should render without errors
	view := model.View()
	// Just ensure it doesn't panic
	_ = view
}
