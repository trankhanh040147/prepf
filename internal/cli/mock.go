package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/trankhanh040147/prepf/internal/ai"
	"github.com/trankhanh040147/prepf/internal/mock"
)

var (
	mockResumePath string
)

var mockCmd = &cobra.Command{
	Use:   "mock",
	Short: "Start a mock technical interview",
	Long:  "Start a mock technical interview (The Gauntlet). Answer questions one at a time.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := GetConfig(cmd)
		if cfg == nil {
			return fmt.Errorf("config not loaded")
		}

		// Validate API key
		if cfg.APIKey == "" {
			return fmt.Errorf("API key not configured. Set it with 'prepf config api_key <key>' or GEMINI_API_KEY env var")
		}

		// Determine resume path
		resumePath := mockResumePath
		if resumePath == "" {
			// Use profile file path directly (profile file contains resume content)
			resumePath = cfg.ProfilePath
		}

		// Create AI client
		aiClient := ai.NewClient(cfg.APIKey, cfg.Timeout, cfg.TokenLimit)

		// Create model (cfg is *config.Config from GetConfig)
		model := mock.NewModel(cfg, aiClient, resumePath)

		// Launch TUI
		program := tea.NewProgram(model, tea.WithAltScreen())
		if _, err := program.Run(); err != nil {
			return fmt.Errorf("run program: %w", err)
		}

		return nil
	},
}

func init() {
	mockCmd.Flags().StringVarP(&mockResumePath, "resume", "r", "", "Path to resume/CV file (.txt or .md)")
	rootCmd.AddCommand(mockCmd)
}
