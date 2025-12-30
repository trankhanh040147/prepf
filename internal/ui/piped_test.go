package ui

import (
	"os"
	"testing"

	"github.com/trankhanh040147/prepf/internal/config"
)

func TestModel_PipedMode(t *testing.T) {
	// Simulate piped mode (not a TTY)
	cfg := &config.Config{
		NoColor: true, // Piped mode should disable colors
		IsTTY:   false,
	}
	model := createTestModel(cfg)

	// Verify colors are disabled
	if !model.noColor {
		t.Error("colors should be disabled in piped mode")
	}

	// View should render without color codes
	view := model.View()
	// Just ensure it doesn't panic
	_ = view
}

func TestConfig_NoColorEnvVar(t *testing.T) {
	// Test that NO_COLOR environment variable is respected
	originalValue := os.Getenv("NO_COLOR")
	defer func() {
		if originalValue == "" {
			os.Unsetenv("NO_COLOR")
		} else {
			os.Setenv("NO_COLOR", originalValue)
		}
	}()

	// Set NO_COLOR
	os.Setenv("NO_COLOR", "1")

	// Config should detect NO_COLOR
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if !cfg.NoColor {
		t.Error("NoColor should be true when NO_COLOR env var is set")
	}
}
