package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// Styles for the UI
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			MarginBottom(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(ColorError).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess)

	helpStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			MarginTop(1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(1, 2)
)

// Renderer handles markdown rendering with glamour
type Renderer struct {
	glamour *glamour.TermRenderer
}

// NewRenderer creates a new markdown renderer
func NewRenderer() (*Renderer, error) {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(100),
	)
	if err != nil {
		return nil, fmt.Errorf("create glamour renderer: %w", err)
	}

	return &Renderer{glamour: r}, nil
}

// Render renders markdown content for the terminal
func (r *Renderer) Render(content string) string {
	if r == nil || r.glamour == nil {
		return content
	}
	rendered, err := r.glamour.Render(content)
	if err != nil {
		return content
	}
	return rendered
}

// RenderMarkdown is an alias for Render for clarity
func (r *Renderer) RenderMarkdown(content string) (string, error) {
	if r == nil || r.glamour == nil {
		return content, nil
	}
	return r.glamour.Render(content)
}

// RenderTitle renders a title
func RenderTitle(text string) string {
	return titleStyle.Render(text)
}

// RenderSubtitle renders a subtitle
func RenderSubtitle(text string) string {
	return subtitleStyle.Render(text)
}

// RenderError renders an error message
func RenderError(text string) string {
	return errorStyle.Render("✗ " + text)
}

// RenderSuccess renders a success message
func RenderSuccess(text string) string {
	return successStyle.Render("✓ " + text)
}

// RenderHelp renders help text
func RenderHelp(text string) string {
	return helpStyle.Render(text)
}

// RenderBox renders content in a box
func RenderBox(content string) string {
	return boxStyle.Render(content)
}

// RenderDivider renders a divider line
func RenderDivider(width int) string {
	safeWidth := width
	if safeWidth < 0 {
		safeWidth = 0
	}
	return StyleDim.Render(strings.Repeat("─", safeWidth))
}

// RenderTokenUsage renders token usage information
func RenderTokenUsage(prompt, completion, total int32) string {
	return subtitleStyle.Render(fmt.Sprintf(
		"📊 Token Usage: %d prompt + %d completion = %d total",
		prompt, completion, total,
	))
}
