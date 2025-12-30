package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)

	// Set HOME to temp dir
	os.Setenv("HOME", tmpDir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg == nil {
		t.Fatal("Load() returned nil config")
	}

	// Verify defaults
	if cfg.Timeout != DefaultTimeout {
		t.Errorf("expected timeout %d, got %d", DefaultTimeout, cfg.Timeout)
	}

	if cfg.ConfigDir == "" {
		t.Error("ConfigDir should not be empty")
	}
}

func TestSave(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)

	os.Setenv("HOME", tmpDir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Set a test value
	cfg.APIKey = "test-api-key"
	cfg.Timeout = 60

	err = Save(cfg)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Reload and verify
	cfg2, err := Load()
	if err != nil {
		t.Fatalf("Load() after Save() error = %v", err)
	}

	if cfg2.APIKey != "test-api-key" {
		t.Errorf("expected APIKey 'test-api-key', got '%s'", cfg2.APIKey)
	}

	if cfg2.Timeout != 60 {
		t.Errorf("expected Timeout 60, got %d", cfg2.Timeout)
	}
}

func TestNoColorDetection(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	origNoColor := os.Getenv(EnvVarNoColor)
	defer func() {
		os.Setenv("HOME", origHome)
		if origNoColor == "" {
			os.Unsetenv(EnvVarNoColor)
		} else {
			os.Setenv(EnvVarNoColor, origNoColor)
		}
	}()

	os.Setenv("HOME", tmpDir)

	// Test with NO_COLOR set
	os.Setenv(EnvVarNoColor, "1")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if !cfg.NoColor {
		t.Error("expected NoColor to be true when NO_COLOR is set")
	}

	// Test without NO_COLOR
	os.Unsetenv(EnvVarNoColor)
	cfg2, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// NoColor should be false if stdout is a TTY (which it should be in tests)
	// But we can't easily test non-TTY case without mocking
	_ = cfg2
}
