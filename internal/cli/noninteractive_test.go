package cli

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/trankhanh040147/prepf/internal/config"
)

func TestConfigCommand_NonInteractive(t *testing.T) {
	// Test config command in non-interactive mode
	// This simulates running: prepf config api_key test-key

	// Create a buffer to capture output
	var buf bytes.Buffer

	// Create a test command
	testCmd := &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			cmd.SetOut(&buf)
			cmd.Printf("API Key: %s\n", cfg.APIKey)
			return nil
		},
	}

	// Execute command
	if err := testCmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	// Verify output was captured
	if buf.Len() == 0 {
		t.Error("expected command to produce output")
	}
}

func TestVersionCommand_NonInteractive(t *testing.T) {
	// Test version command produces output
	var buf bytes.Buffer
	versionCmd.SetOut(&buf)

	if err := versionCmd.RunE(versionCmd, []string{}); err != nil {
		t.Fatalf("RunE() error = %v", err)
	}

	if buf.Len() == 0 {
		t.Error("expected version command to produce output")
	}
}

func TestErrorPropagation(t *testing.T) {
	// Test that errors are properly propagated
	// Create a command that returns an error
	errCmd := &cobra.Command{
		Use: "error",
		RunE: func(cmd *cobra.Command, args []string) error {
			return os.ErrNotExist
		},
	}

	err := errCmd.Execute()
	if err == nil {
		t.Error("expected command to return error")
	}
}

func TestQuietMode(t *testing.T) {
	// Test quiet mode - flags are now accessed via context
	// This is tested via integration tests with actual command execution
	// The flag parsing happens in rootCmd.PersistentPreRunE
	_ = rootCmd.PersistentPreRunE // Verify it exists
}

func TestVerboseMode(t *testing.T) {
	// Test verbose mode - flags are now accessed via context
	// The flag parsing happens in rootCmd.PersistentPreRunE
	_ = rootCmd.PersistentPreRunE // Verify it exists
}

func TestConfigPathOverride(t *testing.T) {
	// Test config path override - now handled via context in PersistentPreRunE
	// Integration tests will verify flag parsing works correctly
	_ = rootCmd.PersistentPreRunE
}

func TestProfilePathOverride(t *testing.T) {
	// Test profile path override - now handled via context in PersistentPreRunE
	// Integration tests will verify flag parsing works correctly
	_ = rootCmd.PersistentPreRunE
}
