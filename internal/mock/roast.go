package mock

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// CalculateGrade calculates grade based on surrender count and total questions
func CalculateGrade(surrenders int, questions int) string {
	switch {
	case surrenders <= GradeASurrenders:
		return "A"
	case surrenders <= GradeBSurrenders:
		return "B"
	case surrenders <= GradeCSurrenders:
		return "C"
	case surrenders <= GradeDSurrenders:
		return "D"
	default:
		return "F"
	}
}

// GetPersonaLabel returns the persona label for a grade
func GetPersonaLabel(grade string) string {
	if label, ok := PersonaLabels[grade]; ok {
		return label
	}
	return "UNKNOWN"
}

// RenderRoast renders the final roast/grade screen
func RenderRoast(grade, persona, feedback string, width int, noColor bool) string {
	// Grade box styling
	gradeStyle := lipgloss.NewStyle().
		Width(50).
		Padding(2, 4).
		Border(lipgloss.RoundedBorder()).
		Align(lipgloss.Center)

	if !noColor {
		gradeStyle = gradeStyle.BorderForeground(lipgloss.Color("62"))
	}

	// Persona label styling
	personaStyle := lipgloss.NewStyle().
		Bold(true).
		Padding(0, 2)

	if !noColor {
		personaColor := getPersonaColor(grade)
		personaStyle = personaStyle.Foreground(lipgloss.Color(personaColor))
	}

	// Grade text (large, bold)
	gradeText := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Center)

	if !noColor {
		gradeColor := getPersonaColor(grade)
		gradeText = gradeText.Foreground(lipgloss.Color(gradeColor))
	}

	// Build content
	content := []string{
		gradeText.Render(fmt.Sprintf("[%s]", grade)),
		personaStyle.Render(persona),
	}

	if feedback != "" {
		feedbackStyle := lipgloss.NewStyle().
			Padding(1, 0).
			Width(50)
		content = append(content, feedbackStyle.Render(feedback))
	}

	gradeBox := gradeStyle.Render(strings.Join(content, "\n\n"))

	// Center the box
	return lipgloss.Place(width, 20, lipgloss.Center, lipgloss.Center, gradeBox)
}

// RenderRemediationButtons renders remediation topic buttons (UI only, non-functional in v0.1.1)
func RenderRemediationButtons(topics []string, width int, noColor bool) string {
	if len(topics) == 0 {
		return ""
	}

	// Limit to 3 buttons
	if len(topics) > 3 {
		topics = topics[:3]
	}

	buttonStyle := lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		Align(lipgloss.Center)

	if !noColor {
		buttonStyle = buttonStyle.BorderForeground(lipgloss.Color("62"))
	}

	buttons := make([]string, len(topics))
	for i, topic := range topics {
		buttons[i] = buttonStyle.Width(20).Render(topic)
	}

	// Join buttons horizontally
	buttonsRow := lipgloss.JoinHorizontal(lipgloss.Center, buttons...)
	return lipgloss.Place(width, 10, lipgloss.Center, lipgloss.Bottom, buttonsRow)
}

// RenderMicroRoast renders an inline micro-roast (for surrenders)
func RenderMicroRoast(feedback string, noColor bool) string {
	style := lipgloss.NewStyle().Bold(true)
	if !noColor {
		style = style.Foreground(lipgloss.Color("9")) // Red
	}
	return style.Render(feedback)
}

// getPersonaColor returns color code for grade
func getPersonaColor(grade string) string {
	switch grade {
	case "A":
		return "2" // Green
	case "B":
		return "10" // Bright green
	case "C":
		return "3" // Yellow
	case "D":
		return "11" // Bright yellow
	case "F":
		return "9" // Red
	default:
		return "7" // Gray
	}
}

