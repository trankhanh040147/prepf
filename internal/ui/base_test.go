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

	// Test zero width handling (should default to min width)
	zeroWidthMsg := tea.WindowSizeMsg{Width: 0, Height: 20}
	model.Update(zeroWidthMsg)
	if model.Width() < config.DefaultMinWidth {
		t.Errorf("expected width >= %d for zero width, got %d", config.DefaultMinWidth, model.Width())
	}
}

func TestBaseModel_KeyBindings(t *testing.T) {
	cfg := &config.Config{NoColor: false}
	model := NewBaseModel(cfg)

	// Test help key (?) triggers help state
	helpKey := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")}
	_, _ = model.Update(helpKey)
	if model.State() != StateHelp {
		t.Errorf("expected StateHelp after '?' key, got %d", model.State())
	}

	// Test quit key (q) returns quit command
	// Reset to normal state first
	model.SetState(StateNormal)
	quitKey := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}
	_, cmd := model.Update(quitKey)
	// Check if cmd is tea.Quit (it's a function that returns tea.Quit())
	if cmd == nil {
		t.Error("quit key should return tea.Quit command")
	}
}

func TestBaseModel_ViewportIntegration(t *testing.T) {
	cfg := &config.Config{NoColor: false}
	model := NewBaseModel(cfg)

	// Test viewport is initialized
	if model.Viewport() == nil {
		t.Fatal("viewport should be initialized")
	}

	// Test setting viewport content
	content := "Test content"
	model.SetViewportContent(content)
	
	// Viewport content is set (we can't easily verify without rendering, but structure is correct)
	if model.Viewport() == nil {
		t.Error("viewport should still exist after setting content")
	}
}

func TestBaseModel_SearchIntegration(t *testing.T) {
	cfg := &config.Config{NoColor: false}
	model := NewBaseModel(cfg)

	// Test search activation via key
	searchKey := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("/")}
	_, _ = model.Update(searchKey)
	
	if model.State() != StateSearch {
		t.Errorf("expected StateSearch after '/' key, got %d", model.State())
	}

	// Test search query retrieval
	query := model.SearchQuery()
	if query != "" {
		t.Errorf("expected empty query initially, got '%s'", query)
	}
}

func TestBaseModel_NoColorMode(t *testing.T) {
	cfg := &config.Config{NoColor: true}
	model := NewBaseModel(cfg)

	// Verify noColor is set
	if !model.noColor {
		t.Error("noColor should be true when config.NoColor is true")
	}

	// Set some content to viewport so View() returns something
	model.SetViewportContent("Test content")
	
	// View should still render without errors
	view := model.View()
	// Viewport might be empty initially, but View() should not panic
	_ = view // Just ensure it doesn't panic
}

