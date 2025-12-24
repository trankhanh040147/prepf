package cli

import (
	"testing"
)

func TestRootCommand(t *testing.T) {
	// Test that root command exists and can be executed
	if rootCmd == nil {
		t.Fatal("root command is nil")
	}

	if rootCmd.Use != "prepf" {
		t.Errorf("expected root command Use to be 'prepf', got '%s'", rootCmd.Use)
	}

	// Test that RunE is set
	if rootCmd.RunE == nil {
		t.Fatal("root command RunE is nil")
	}
}

