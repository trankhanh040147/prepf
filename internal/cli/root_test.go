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

func TestRootCommand_Flags(t *testing.T) {
	// Test that flags are defined
	configFlag := rootCmd.PersistentFlags().Lookup(FlagConfig)
	if configFlag == nil {
		t.Error("--config flag should be defined")
	}

	profileFlag := rootCmd.PersistentFlags().Lookup(FlagProfile)
	if profileFlag == nil {
		t.Error("--profile flag should be defined")
	}

	verboseFlag := rootCmd.PersistentFlags().Lookup(FlagVerbose)
	if verboseFlag == nil {
		t.Error("--verbose flag should be defined")
	}

	quietFlag := rootCmd.PersistentFlags().Lookup(FlagQuiet)
	if quietFlag == nil {
		t.Error("--quiet flag should be defined")
	}
}

func TestRootCommand_ContextInjection(t *testing.T) {
	// Test that PersistentPreRunE is set
	if rootCmd.PersistentPreRunE == nil {
		t.Fatal("root command PersistentPreRunE is nil")
	}

	// Context helpers (GetConfig, GetFlags) are tested via integration tests
	// with actual command execution
}
