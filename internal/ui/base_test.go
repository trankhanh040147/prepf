package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/trankhanh040147/prepf/internal/config"
)

func TestBaseModel_StateTransitions(t *testing.T) {
	cfg := &config.Config{NoColor: false}
	model := NewBaseModel(cfg)

	// Test initial state
	if model.State() != StateNormal {
		t.Errorf("expected initial state StateNormal, got %d", model.State())
	}

	// Test help toggle
	model.toggleHelp()
	if model.State() != StateHelp {
		t.Errorf("expected state StateHelp after toggleHelp, got %d", model.State())
	}

	model.toggleHelp()
	if model.State() != StateNormal {
		t.Errorf("expected state StateNormal after second toggleHelp, got %d", model.State())
	}
}

func TestBaseModel_SearchState(t *testing.T) {
	cfg := &config.Config{NoColor: false}
	model := NewBaseModel(cfg)

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

func TestBaseModel_WidthSafety(t *testing.T) {
	cfg := &config.Config{NoColor: false}
	model := NewBaseModel(cfg)

	// Test default width
	width := model.Width()
	if width < config.DefaultMinWidth {
		t.Errorf("expected width >= %d, got %d", config.DefaultMinWidth, width)
	}

	// Test width after window size message
	windowSizeMsg := tea.WindowSizeMsg{Width: 120, Height: 40}
	model.Update(windowSizeMsg)

	if model.Width() != 120 {
		t.Errorf("expected width 120, got %d", model.Width())
	}
}

