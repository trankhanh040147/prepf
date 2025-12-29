package gemini

import (
	"github.com/samber/lo"
	"google.golang.org/genai"
)

// buildSafetySettings constructs safety settings from threshold string
func (c *Client) buildSafetySettings() []*genai.SafetySetting {
	threshold := mapSafetyThreshold(c.safetyThreshold)
	categories := getAllHarmCategories()
	return lo.Map(categories, func(category genai.HarmCategory, _ int) *genai.SafetySetting {
		return &genai.SafetySetting{
			Category:  category,
			Threshold: threshold,
		}
	})
}

// newGenerationConfig creates a GenerateContentConfig with the given system prompt
func (c *Client) newGenerationConfig(systemPrompt string) *genai.GenerateContentConfig {
	temp := c.modelParams.Temperature
	topP := c.modelParams.TopP
	topK := float32(c.modelParams.TopK)

	return &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: systemPrompt}},
		},
		Temperature:    &temp,
		TopP:           &topP,
		TopK:           &topK,
		SafetySettings: c.buildSafetySettings(),
	}
}

// appendUserMessageAndPrepareTurn creates user message, appends it to history,
// and prepares the generation config
func (c *Client) appendUserMessageAndPrepareTurn(message string) (*genai.Content, *genai.GenerateContentConfig) {
	userMsg := &genai.Content{
		Role:  "user",
		Parts: []*genai.Part{{Text: message}},
	}

	config := c.newGenerationConfig(c.systemPrompt)
	c.history = append(c.history, userMsg)

	return userMsg, config
}
