package ui

import "github.com/charmbracelet/lipgloss"

// UI State constants
type State int

const (
	StateNormal State = iota
	StateHelp
	StateLoading
	StateSearch
	StateChatting
	StateError
	StateQuitting
)

// Colors - using adaptive colors for light/dark theme support
var (
	ColorPrimary   = lipgloss.AdaptiveColor{Light: "#7C3AED", Dark: "#A78BFA"}
	ColorSecondary = lipgloss.AdaptiveColor{Light: "#4B5563", Dark: "#9CA3AF"}
	ColorError     = lipgloss.AdaptiveColor{Light: "#DC2626", Dark: "#F87171"}
	ColorSuccess   = lipgloss.AdaptiveColor{Light: "#059669", Dark: "#34D399"}
	ColorBorder    = lipgloss.AdaptiveColor{Light: "#E5E7EB", Dark: "#374151"}
)

// Styles - global style definitions
var (
	StyleBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder)

	StyleFocused = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary)

	StyleError = lipgloss.NewStyle().
			Foreground(ColorError)

	StyleSuccess = lipgloss.NewStyle().
			Foreground(ColorSuccess)

	StyleDim = lipgloss.NewStyle().
			Foreground(ColorSecondary)
)

// Layout constants
const (
	DefaultPadding = 1
	MinWidth       = 40
	MinHeight      = 10
)
