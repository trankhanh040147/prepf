package mock

import (
	"github.com/charmbracelet/huh"
)

// buildConfigForm creates a form for topic configuration
func buildConfigForm(selectedTopics *[]string, excludedTopics *[]string) *huh.Form {
	// Create separate options slices for each field to avoid state sharing
	focusOptions := make([]huh.Option[string], len(ValidInterviewTopics))
	excludeOptions := make([]huh.Option[string], len(ValidInterviewTopics))
	for i, topic := range ValidInterviewTopics {
		// Create independent option instances for each field
		focusOptions[i] = huh.NewOption(topic, topic)
		excludeOptions[i] = huh.NewOption(topic, topic)
	}

	// Multi-select for topics to focus on
	focusField := huh.NewMultiSelect[string]().
		Title("Topics to focus on").
		Description("Select topics you want to be asked about (optional, leave empty for all topics). Use ↑/↓ to navigate, Space to toggle selection, Tab to move to next field.").
		Options(focusOptions...).
		Value(selectedTopics)

	// Multi-select for topics to exclude
	excludeField := huh.NewMultiSelect[string]().
		Title("Topics to exclude").
		Description("Select topics you want to avoid (optional). Use ↑/↓ to navigate, Space to toggle selection, Tab to go back, Enter to begin interview.").
		Options(excludeOptions...).
		Value(excludedTopics)

	form := huh.NewForm(
		huh.NewGroup(focusField, excludeField),
	).
		WithTheme(huh.ThemeBase16()).
		WithWidth(80)

	return form
}
