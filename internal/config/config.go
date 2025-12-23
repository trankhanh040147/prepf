package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mattn/go-isatty"
	"github.com/spf13/viper"
)

// Config holds application configuration
type Config struct {
	APIKey      string
	Timeout     int
	Editor      string
	NoColor     bool
	IsTTY       bool
	ConfigDir   string
	ProfilePath string
}

// setupViper initializes Viper configuration
func setupViper() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", ConfigDirName)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("create config directory: %w", err)
	}

	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	// Set defaults
	viper.SetDefault(KeyAPIKey, "")
	viper.SetDefault(KeyTimeout, DefaultTimeout)
	viper.SetDefault(KeyEditor, getDefaultEditor())

	// Environment variable overrides
	viper.SetEnvPrefix(EnvPrefix)
	viper.BindEnv(KeyAPIKey, EnvVarGeminiAPIKey)
	viper.BindEnv(KeyTimeout, EnvVarTimeout)
	viper.BindEnv(KeyEditor, EnvVarEditor)

	return configDir, nil
}

// Load initializes and loads configuration
func Load() (*Config, error) {
	configDir, err := setupViper()
	if err != nil {
		return nil, err
	}

	// Read config file (ignore if not exists, warn if malformed)
	if err := viper.ReadInConfig(); err != nil {
		// If file not found, continue with defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// File doesn't exist, use defaults
		} else {
			// File exists but is malformed - warn user but continue with defaults
			// This allows 'config edit' to work even with corrupted YAML
			fmt.Fprintf(os.Stderr, "warning: config file is malformed: %v\n", err)
			fmt.Fprintf(os.Stderr, "warning: using default values. Run 'prepf config edit' to fix.\n")
		}
	}

	// TTY detection
	isTTY := isatty.IsTerminal(os.Stdout.Fd())
	noColor := os.Getenv(EnvVarNoColor) != "" || !isTTY

	cfg := &Config{
		APIKey:      viper.GetString(KeyAPIKey),
		Timeout:     viper.GetInt(KeyTimeout),
		Editor:      viper.GetString(KeyEditor),
		NoColor:     noColor,
		IsTTY:       isTTY,
		ConfigDir:   configDir,
		ProfilePath: filepath.Join(configDir, ProfileFileName),
	}

	return cfg, nil
}

// Save saves configuration to file
func Save(cfg *Config) error {
	// Set only writable fields on the existing viper instance
	viper.Set(KeyAPIKey, cfg.APIKey)
	viper.Set(KeyTimeout, cfg.Timeout)
	viper.Set(KeyEditor, cfg.Editor)

	// Get config path from viper if it was read successfully, otherwise construct it
	configPath := viper.ConfigFileUsed()
	if configPath == "" {
		// Fallback if config file didn't exist before
		configPath = filepath.Join(cfg.ConfigDir, ConfigFileName)
	}

	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("write config file: %w", err)
	}

	return nil
}

// getDefaultEditor returns platform-specific default editor
func getDefaultEditor() string {
	switch runtime.GOOS {
	case "linux":
		return DefaultEditorLinux
	case "darwin":
		return DefaultEditorDarwin
	case "windows":
		return DefaultEditorWindows
	default:
		return DefaultEditorLinux
	}
}

// InitialConfigContent returns the initial YAML content for a new config file
func InitialConfigContent() string {
	return fmt.Sprintf(`%s%s: ""
%s: %d
%s: ""
`, ConfigFileHeader, KeyAPIKey, KeyTimeout, DefaultTimeout, KeyEditor)
}
