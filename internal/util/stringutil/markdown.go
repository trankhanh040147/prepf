package stringutil

import (
	"github.com/charmbracelet/glamour"
)

// RenderMarkdown renders markdown content to styled ANSI text for terminal display.
// Returns the original text if rendering fails (graceful degradation).
func RenderMarkdown(content string, width int) string {
	if content == "" {
		return content
	}

	// Create glamour renderer with custom width
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		// Graceful degradation: return original content if renderer creation fails
		return content
	}

	rendered, err := renderer.Render(content)
	if err != nil {
		// Graceful degradation: return original content if rendering fails
		return content
	}

	return rendered
}
