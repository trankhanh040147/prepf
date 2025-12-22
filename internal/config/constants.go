package config

const (
	// EnvVarGeminiAPIKey is the environment variable name for Gemini API key
	EnvVarGeminiAPIKey = "GEMINI_API_KEY"
	// EnvVarNoColor disables color output
	EnvVarNoColor = "NO_COLOR"
	// EnvVarEditor is the environment variable for editor path
	EnvVarEditor = "EDITOR"

	// ConfigDirName is the config directory name
	ConfigDirName = "prepf"
	// ConfigFileName is the config file name
	ConfigFileName = "config.yaml"
	// ProfileFileName is the profile file name
	ProfileFileName = "profile.json"

	// DefaultTimeout is the default network timeout in seconds
	DefaultTimeout = 30
	// DefaultEditorLinux is the default editor for Linux
	DefaultEditorLinux = "vi"
	// DefaultEditorDarwin is the default editor for macOS
	DefaultEditorDarwin = "open"
	// DefaultEditorWindows is the default editor for Windows
	DefaultEditorWindows = "notepad"
	// DefaultMinWidth is the minimum terminal width
	DefaultMinWidth = 80
)
