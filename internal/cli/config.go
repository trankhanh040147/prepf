package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/trankhanh040147/prepf/internal/config"
	"github.com/trankhanh040147/prepf/internal/stringext"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  "View or edit configuration settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := GetConfig(cmd)
		if cfg == nil {
			return fmt.Errorf("config not loaded")
		}

		// If no args, show all config
		if len(args) == 0 {
			return showAllConfig(cmd, cfg)
		}

		key := args[0]

		// If 2 args, set value
		if len(args) == 2 {
			return setConfigValue(cmd, cfg, key, args[1])
		}

		// If 1 arg, show value
		if displayFunc, ok := configKeyDisplayMap[key]; ok {
			displayFunc(cmd, cfg)
			return nil
		}

		// Unknown key - show error with fuzzy match suggestion
		keys := lo.Keys(configKeyDisplayMap)
		slices.Sort(keys)
		return newInvalidKeyError(key, keys, "displayed")
	},
}

// newInvalidKeyError creates a standardized error for unknown/unsupported keys with fuzzy match suggestions
func newInvalidKeyError(providedKey string, validKeys []string, action string) error {
	availableKeysStr := strings.Join(validKeys, ", ")
	msg := fmt.Sprintf("key '%s' cannot be %s", providedKey, action)

	// Add fuzzy match suggestion if available
	if closest := stringext.FuzzyMatch(providedKey, validKeys); closest != "" {
		msg += fmt.Sprintf(", did you mean '%s'?", closest)
	}

	return fmt.Errorf("%s. Available keys: %s", msg, availableKeysStr)
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
			cfg := cmd.Context().Value(configKey{}).(*config.Config)

			configPath := filepath.Join(cfg.ConfigDir, config.ConfigFileName)
			editor := cfg.Editor

			// Fallback to environment variable if editor not in config file
			// This allows 'config edit' to work even when config file doesn't exist yet
			// or when editor hasn't been configured via 'prepf config editor <path>'
			if editor == "" {
				editor = os.Getenv(config.EnvVarEditor)
				if editor == "" {
					return fmt.Errorf("no editor configured and %s not set", config.EnvVarEditor)
				}
			}

			return openEditor(editor, configPath)
		},
	})
}

// openEditor opens the config file in the specified editor
func openEditor(editor, filePath string) error {
	// Ensure config file exists (create with template if not)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		initialContent := config.InitialConfigContent()
		if err := os.WriteFile(filePath, []byte(initialContent), 0644); err != nil {
			return fmt.Errorf("create config file: %w", err)
		}
	}

	// Parse editor command (handle cases like "nvim", "vi", "open -a TextEdit")
	editorParts := strings.Fields(editor)
	if len(editorParts) == 0 {
		return fmt.Errorf("invalid editor command: contains only whitespace")
	}
	editorCmd := editorParts[0]
	editorArgs := editorParts[1:]

	// Standard editor: append file path as argument
	editorArgs = append(editorArgs, filePath)

	cmd := exec.Command(editorCmd, editorArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("execute editor: %w", err)
	}

	return nil
}

// showAllConfig displays all configuration values
func showAllConfig(cmd *cobra.Command, cfg *config.Config) error {
	// Iterate over configKeyDisplayMap to ensure single source of truth
	keys := lo.Keys(configKeyDisplayMap)
	slices.Sort(keys)
	lo.ForEach(keys, func(key string, _ int) {
		configKeyDisplayMap[key](cmd, cfg)
	})
	return nil
}

// configKeyDisplayMap maps config keys to their display functions
// Display functions show only cfg.* values (single source of truth after Viper processing)
// Read-only keys (no_color, is_tty, config_dir, profile_path) are display-only and not in configKeySetterMap
var configKeyDisplayMap = map[string]func(*cobra.Command, *config.Config){
	config.KeyAPIKey: func(cmd *cobra.Command, cfg *config.Config) {
		cmd.Printf("API Key: %s\n", maskAPIKey(cfg.APIKey))
	},
	config.KeyTimeout: func(cmd *cobra.Command, cfg *config.Config) {
		cmd.Printf("Timeout: %d seconds\n", cfg.Timeout)
	},
	config.KeyEditor: func(cmd *cobra.Command, cfg *config.Config) {
		cmd.Printf("Editor: %s\n", cfg.Editor)
	},
	config.KeyNoColor: func(cmd *cobra.Command, cfg *config.Config) {
		cmd.Printf("No Color: %v\n", cfg.NoColor)
	},
	config.KeyIsTTY: func(cmd *cobra.Command, cfg *config.Config) {
		cmd.Printf("Is TTY: %v\n", cfg.IsTTY)
	},
	config.KeyConfigDir: func(cmd *cobra.Command, cfg *config.Config) {
		cmd.Printf("Config Directory: %s\n", cfg.ConfigDir)
	},
	config.KeyProfilePath: func(cmd *cobra.Command, cfg *config.Config) {
		cmd.Printf("Profile Path: %s\n", cfg.ProfilePath)
	},
}

// configKeySetterMap maps config keys to their setter functions (only writable keys)
var configKeySetterMap = map[string]func(*config.Config, string) error{
	config.KeyAPIKey: func(cfg *config.Config, value string) error {
		if value == "" {
			return fmt.Errorf("api_key cannot be empty")
		}
		cfg.APIKey = value
		return nil
	},
	config.KeyTimeout: func(cfg *config.Config, value string) error {
		timeout, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("timeout must be a number: %w", err)
		}
		if timeout <= 0 {
			return fmt.Errorf("timeout must be greater than 0")
		}
		cfg.Timeout = timeout
		return nil
	},
	config.KeyEditor: func(cfg *config.Config, value string) error {
		if value == "" {
			return fmt.Errorf("editor cannot be empty")
		}
		cfg.Editor = value
		return nil
	},
}

// setConfigValue sets a config value and saves it
func setConfigValue(cmd *cobra.Command, cfg *config.Config, key, value string) error {
	setterFunc, ok := configKeySetterMap[key]
	if !ok {
		keys := lo.Keys(configKeySetterMap)
		slices.Sort(keys)
		return newInvalidKeyError(key, keys, "set")
	}

	// Set the value
	if err := setterFunc(cfg, value); err != nil {
		return fmt.Errorf("invalid value for %s: %w", key, err)
	}

	// Save config
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	// Show updated value
	if displayFunc, ok := configKeyDisplayMap[key]; ok {
		displayFunc(cmd, cfg)
	}

	return nil
}
