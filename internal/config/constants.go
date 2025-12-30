package config

const (
	// EnvPrefix is the environment variable prefix for prepf
	EnvPrefix = "PREPF"
	// EnvVarGeminiAPIKey is the environment variable name for Gemini API key
	EnvVarGeminiAPIKey = "GEMINI_API_KEY"
	// EnvVarNoColor disables color output
	EnvVarNoColor = "NO_COLOR"
	// EnvVarEditor is the environment variable for editor path
	EnvVarEditor = "EDITOR"
	// EnvVarTimeout is the environment variable for timeout (PREPF_TIMEOUT)
	EnvVarTimeout = "PREPF_TIMEOUT"

	// Config key names (used in Viper and YAML)
	KeyAPIKey     = "api_key"
	KeyTimeout    = "timeout"
	KeyEditor     = "editor"
	KeyTokenLimit = "token_limit"

	// Read-only config keys (display-only, not settable)
	KeyNoColor     = "no_color"
	KeyIsTTY       = "is_tty"
	KeyConfigDir   = "config_dir"
	KeyProfilePath = "profile_path"

	// ConfigDirName is the config directory name
	ConfigDirName = "prepf"
	// ConfigFileName is the config file name
	ConfigFileName = "config.yaml"
	// ProfileFileName is the profile file name
	ProfileFileName = "profile.json"

	// DefaultTimeout is the default network timeout in seconds
	DefaultTimeout = 30
	// DefaultTokenLimit is the default token limit for gemini-pro (1M tokens)
	DefaultTokenLimit = 1_000_000
	// TokenEstimationMargin is the safety margin (as percentage) for token estimation
	// A 20% margin helps account for inaccuracies in the rough len(prompt)/4 heuristic
	TokenEstimationMargin = 20
	// DefaultEditorLinux is the default editor for Linux
	DefaultEditorLinux = "vi"
	// DefaultEditorDarwin is the default editor for macOS
	DefaultEditorDarwin = "open"
	// DefaultEditorWindows is the default editor for Windows
	DefaultEditorWindows = "notepad"
	// DefaultMinWidth is the minimum terminal width
	DefaultMinWidth = 80
	// ConfigFileHeader is the header comment for config files
	ConfigFileHeader = `# prepf configuration
# Example:
# api_key: "your-api-key-here"
# timeout: 30
# editor: "nvim"
`
)
