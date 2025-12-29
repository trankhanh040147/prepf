package ai

import (
	"strings"

	"google.golang.org/genai"
)

// historyToGenaiContents converts History to []*genai.Content
// Maps "assistant" role to "model" as required by the SDK
// Valid roles for SDK are: "user" and "model"
func (c *Client) historyToGenaiContents() []*genai.Content {
	if c.history == nil {
		return nil
	}

	entries := c.history.Entries()
	contents := make([]*genai.Content, 0, len(entries))

	for _, entry := range entries {
		role := entry.Role
		// SDK expects "model" instead of "assistant"
		// Valid roles are: "user" and "model"
		if role == "assistant" {
			role = "model"
		}
		// Ensure we only send valid roles (user or model)
		// If role is invalid, default to "user" for safety
		if role != "user" && role != "model" {
			role = "user"
		}
		contents = append(contents, &genai.Content{
			Role:  role,
			Parts: []*genai.Part{{Text: entry.Content}},
		})
	}

	return contents
}

// newGenerationConfig creates a minimal GenerateContentConfig
func (c *Client) newGenerationConfig() *genai.GenerateContentConfig {
	return &genai.GenerateContentConfig{
		// Use defaults for temperature, topP, topK
		// Safety settings will use SDK defaults
	}
}

// extractTextFromResponse extracts text from a GenerateContentResponse
func extractTextFromResponse(resp *genai.GenerateContentResponse) string {
	if resp == nil || len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return ""
	}

	var parts []string
	for _, part := range resp.Candidates[0].Content.Parts {
		if part == nil {
			continue
		}
		if txt := part.Text; txt != "" {
			parts = append(parts, txt)
		}
	}
	return strings.Join(parts, "")
}
