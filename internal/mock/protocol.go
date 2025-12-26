package mock

import (
	"regexp"
	"strings"
)

var (
	nextSignalRegex = regexp.MustCompile(`(?i)<NEXT>`)
	roastSignalRegex = regexp.MustCompile(`(?i)<ROAST>`)
)

// ParseSignals extracts <NEXT> and <ROAST> signals from text
// Returns cleaned content and boolean flags for state transitions
func ParseSignals(text string) (content string, hasNext bool, hasRoast bool) {
	content = text
	hasNext = nextSignalRegex.MatchString(content)
	hasRoast = roastSignalRegex.MatchString(content)

	// Remove signals from content
	content = nextSignalRegex.ReplaceAllString(content, "")
	content = roastSignalRegex.ReplaceAllString(content, "")
	
	// Clean up whitespace
	content = strings.TrimSpace(content)
	
	return content, hasNext, hasRoast
}

