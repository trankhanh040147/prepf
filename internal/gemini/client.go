package gemini

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

// TokenUsage represents token consumption for a request
type TokenUsage struct {
	PromptTokens     int32
	CompletionTokens int32
	TotalTokens      int32
}

// ModelParams holds generation parameters
type ModelParams struct {
	Temperature float32
	TopP        float32
	TopK        int
}

// Client wraps the genai.Client with conversation state
type Client struct {
	client          *genai.Client
	modelID         string
	lastUsage       *TokenUsage
	modelParams     *ModelParams
	safetyThreshold string
	history         []*genai.Content
	systemPrompt    string
}

// NewClient creates a new Gemini client with the official genai SDK
func NewClient(ctx context.Context, apiKey, modelID string, params *ModelParams) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("api key is required")
	}

	cfg := &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	}

	client, err := genai.NewClient(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create genai client: %w", err)
	}

	// Apply defaults if params not provided
	if params == nil {
		params = &ModelParams{
			Temperature: DefaultModelTemperature,
			TopP:        DefaultModelTopP,
			TopK:        DefaultModelTopK,
		}
	}

	// Use default model if not specified
	if modelID == "" {
		modelID = DefaultModelID
	}

	return &Client{
		client:          client,
		modelID:         modelID,
		modelParams:     params,
		safetyThreshold: DefaultSafetyThreshold,
	}, nil
}

// StartChat initializes a new chat session with a system prompt
func (c *Client) StartChat(systemPrompt string) {
	c.systemPrompt = systemPrompt
	c.history = make([]*genai.Content, 0, DefaultHistoryCapacity)
}

// ClearHistory clears the conversation history
func (c *Client) ClearHistory() {
	c.history = make([]*genai.Content, 0, DefaultHistoryCapacity)
}

// rollbackLastHistoryEntry removes the last entry from history (for error recovery)
func (c *Client) rollbackLastHistoryEntry() {
	if len(c.history) > 0 {
		c.history = c.history[:len(c.history)-1]
	}
}

// GetLastUsage returns token usage from the last request
func (c *Client) GetLastUsage() *TokenUsage {
	return c.lastUsage
}

// GetModelID returns the configured model ID
func (c *Client) GetModelID() string {
	return c.modelID
}

// SetSafetyThreshold sets the safety threshold for content filtering
func (c *Client) SetSafetyThreshold(threshold string) {
	c.safetyThreshold = threshold
}

// GetHistory returns the current conversation history
func (c *Client) GetHistory() []*genai.Content {
	return c.history
}

// HistoryCount returns the number of messages in history
func (c *Client) HistoryCount() int {
	return len(c.history)
}
