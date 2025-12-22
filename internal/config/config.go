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

// Load initializes and loads configuration
func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", ConfigDirName)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("create config directory: %w", err)
	}

	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	// Set defaults
	viper.SetDefault("api_key", "")
	viper.SetDefault("timeout", DefaultTimeout)
	viper.SetDefault("editor", getDefaultEditor())

	// Environment variable overrides
	viper.SetEnvPrefix("PREPF")
	viper.BindEnv("api_key", EnvVarGeminiAPIKey)
	viper.BindEnv("timeout")
	viper.BindEnv("editor", EnvVarEditor)

	// Read config file (ignore if not exists)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("read config file: %w", err)
		}
	}

	// TTY detection
	isTTY := isatty.IsTerminal(os.Stdout.Fd())
	noColor := os.Getenv(EnvVarNoColor) != "" || !isTTY

	cfg := &Config{
		APIKey:      viper.GetString("api_key"),
		Timeout:     viper.GetInt("timeout"),
		Editor:      viper.GetString("editor"),
		NoColor:     noColor,
		IsTTY:       isTTY,
		ConfigDir:   configDir,
		ProfilePath: filepath.Join(configDir, ProfileFileName),
	}

	return cfg, nil
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
