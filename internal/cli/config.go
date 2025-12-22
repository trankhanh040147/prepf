package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/trankhanh040147/prepf/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  "View or edit configuration settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		cmd.Printf("Config Directory: %s\n", cfg.ConfigDir)
		cmd.Printf("Profile Path: %s\n", cfg.ProfilePath)
		cmd.Printf("API Key: %s\n", maskAPIKey(cfg.APIKey))
		cmd.Printf("Timeout: %d seconds\n", cfg.Timeout)
		cmd.Printf("Editor: %s\n", cfg.Editor)
		cmd.Printf("No Color: %v\n", cfg.NoColor)
		cmd.Printf("Is TTY: %v\n", cfg.IsTTY)

		return nil
	},
}

func maskAPIKey(key string) string {
	if key == "" {
		return "(not set)"
	}
	if len(key) <= 8 {
		return "***"
	}
	return key[:4] + "..." + key[len(key)-4:]
}

func init() {
	configCmd.AddCommand(&cobra.Command{
		Use:   "edit",
		Short: "Edit configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			configPath := cfg.ConfigDir + "/config.yaml"
			editor := cfg.Editor

			// Use system editor
			if editor == "" {
				editor = os.Getenv(config.EnvVarEditor)
				if editor == "" {
					return fmt.Errorf("no editor configured and %s not set", config.EnvVarEditor)
				}
			}

			// Open config file in editor
			// This is a placeholder - actual implementation would use exec.Command
			cmd.Printf("Opening %s with %s\n", configPath, editor)
			return nil
		},
	})
}
