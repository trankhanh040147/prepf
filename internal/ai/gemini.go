// Package ai provides Gemini API client functionality using the official SDK.
//
// The client supports streaming responses, token usage tracking, and conversation history management.
package ai

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/sync/errgroup"
	"google.golang.org/genai"
)

// Client handles Gemini API interactions
type Client struct {
	apiKey           string
	timeout          time.Duration
	genaiClient      *genai.Client
	modelID          string
	inputTokens      int
	outputTokens     int
	tokenLimit       int
	cumulativeInput  int
	cumulativeOutput int
	history          *History
}

// NewClient creates a new Gemini client
func NewClient(apiKey string, timeout int, tokenLimit int) *Client {
	ctx := context.Background()
	cfg := &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	}

	genaiClient, err := genai.NewClient(ctx, cfg)
	if err != nil {
		// If client creation fails, we'll return a client with nil genaiClient
		// and handle the error when Stream is called
		genaiClient = nil
	}

	return &Client{
		apiKey:      apiKey,
		timeout:     time.Duration(timeout) * time.Second,
		tokenLimit:  tokenLimit,
		genaiClient: genaiClient,
		modelID:     "gemini-2.5-flash",
		history:     NewHistory(),
	}
}

// SetTokenLimit sets the token limit for the client
func (c *Client) SetTokenLimit(limit int) {
	c.tokenLimit = limit
}

// GetHistory returns the conversation history
func (c *Client) GetHistory() *History {
	return c.history
}

// ClearHistory clears the conversation history
func (c *Client) ClearHistory() {
	if c.history != nil {
		c.history.Clear()
	}
}

// StreamChunk represents a chunk of streaming data
type StreamChunk struct {
	Text string
	Err  error
	Done bool
}

// Stream sends a request and returns a channel of chunks
func (c *Client) Stream(ctx context.Context, prompt string) (<-chan StreamChunk, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("api key not configured")
	}

	if c.genaiClient == nil {
		return nil, fmt.Errorf("genai client not initialized")
	}

	// Check token limit before making request
	if c.tokenLimit > 0 {
		totalUsed := c.cumulativeInput + c.cumulativeOutput
		// Estimate prompt tokens (rough heuristic: 1 token ≈ 4 chars)
		// NOTE: This is an approximation. Actual tokenization may vary.
		// A safety margin is applied to account for estimation inaccuracies.
		estimatedPromptTokens := len(prompt) / 4
		// Apply safety margin (20% buffer by default)
		estimatedWithMargin := estimatedPromptTokens + (estimatedPromptTokens * 20 / 100)
		if totalUsed+estimatedWithMargin > c.tokenLimit {
			return nil, fmt.Errorf("token limit exceeded: %d/%d tokens used, estimated request (with margin) would exceed limit", totalUsed, c.tokenLimit)
		}
	}

	ch := make(chan StreamChunk, 10)

	go func() {
		defer close(ch)

		// Convert history to SDK format
		historyContents := c.historyToGenaiContents()

		// Add current user prompt
		userContent := &genai.Content{
			Role:  "user",
			Parts: []*genai.Part{{Text: prompt}},
		}
		historyContents = append(historyContents, userContent)

		// Create generation config (minimal, using defaults)
		config := c.newGenerationConfig()

		// Start streaming
		iter := c.genaiClient.Models.GenerateContentStream(ctx, c.modelID, historyContents, config)

		// Convert callback pattern to channel pattern
		type iterResult struct {
			resp *genai.GenerateContentResponse
			err  error
		}
		resultCh := make(chan iterResult, 20)

		g, gCtx := errgroup.WithContext(ctx)
		g.Go(func() error {
			defer close(resultCh)
			iter(func(resp *genai.GenerateContentResponse, err error) bool {
				select {
				case resultCh <- iterResult{resp: resp, err: err}:
					return true
				case <-gCtx.Done():
					return false
				}
			})
			return nil
		})

		// Wait for goroutine completion in background
		go func() {
			if err := g.Wait(); err != nil {
				// Error already handled through channel
				select {
				case resultCh <- iterResult{err: err}:
				default:
				}
			}
		}()

		var fullResponse strings.Builder

		for {
			select {
			case <-ctx.Done():
				ch <- StreamChunk{Err: fmt.Errorf("context cancelled: %w", ctx.Err())}
				return
			case result, ok := <-resultCh:
				if !ok {
					// Channel closed, iteration complete
					ch <- StreamChunk{Done: true}

					// Add to history after stream completes
					if c.history != nil && fullResponse.Len() > 0 {
						c.history.AddToHistory(prompt, fullResponse.String())
					}
					return
				}
				if result.err != nil {
					ch <- StreamChunk{Err: fmt.Errorf("stream error: %w", result.err)}
					return
				}

				resp := result.resp
				if resp == nil {
					continue
				}

				// Check for safety block
				if len(resp.Candidates) > 0 && resp.Candidates[0].FinishReason == genai.FinishReasonSafety {
					ch <- StreamChunk{Err: fmt.Errorf("response blocked by safety filters")}
					return
				}

				// Extract text chunk
				chunkText := extractTextFromResponse(resp)
				if chunkText != "" {
					fullResponse.WriteString(chunkText)
					ch <- StreamChunk{Text: chunkText}
				}

				// Update token usage if available
				if resp.UsageMetadata != nil {
					c.inputTokens = int(resp.UsageMetadata.PromptTokenCount)
					c.outputTokens = int(resp.UsageMetadata.CandidatesTokenCount)
					// Update cumulative usage
					c.cumulativeInput += c.inputTokens
					c.cumulativeOutput += c.outputTokens

					// Check if limit exceeded after this request
					if c.tokenLimit > 0 {
						totalUsed := c.cumulativeInput + c.cumulativeOutput
						if totalUsed > c.tokenLimit {
							ch <- StreamChunk{Err: fmt.Errorf("token limit exceeded: %d/%d tokens used", totalUsed, c.tokenLimit)}
							return
						}
					}
				}

			}
		}
	}()

	return ch, nil
}

// StreamStartCmd is a non-blocking command that initiates a stream.
// It returns immediately with a StreamStartedMsg containing the channel,
// allowing the UI to remain responsive while streaming.
func (c *Client) StreamStartCmd(ctx context.Context, prompt string) tea.Cmd {
	return func() tea.Msg {
		ch, err := c.Stream(ctx, prompt)
		if err != nil {
			return StreamErrorMsg{Err: fmt.Errorf("start stream: %w", err)}
		}
		return StreamStartedMsg{Stream: ch}
	}
}

// WaitForStreamChunkCmd waits for the next chunk from the stream channel.
// It reads one chunk and returns it. The model's Update function should handle
// StreamChunkMsg and return WaitForStreamChunkCmd(ch) again to continue the
// streaming loop. The blocking receive is safe here because tea.Cmd functions
// execute asynchronously and don't freeze the UI event loop.
func WaitForStreamChunkCmd(ch <-chan StreamChunk) tea.Cmd {
	return func() tea.Msg {
		chunk, ok := <-ch
		if !ok {
			return StreamDoneMsg{} // Channel closed
		}
		if chunk.Err != nil {
			return StreamErrorMsg{Err: chunk.Err}
		}
		if chunk.Done {
			return StreamDoneMsg{}
		}
		return StreamChunkMsg{Text: chunk.Text}
	}
}

// GetTokenUsage returns token usage information (input tokens, output tokens) for current request
func (c *Client) GetTokenUsage() (int, int) {
	return c.inputTokens, c.outputTokens
}

// GetCumulativeTokenUsage returns cumulative token usage across all requests
func (c *Client) GetCumulativeTokenUsage() (int, int) {
	return c.cumulativeInput, c.cumulativeOutput
}

// GetTotalTokenUsage returns total tokens used (input + output) for current request
func (c *Client) GetTotalTokenUsage() int {
	return c.inputTokens + c.outputTokens
}

// GetTotalCumulativeTokenUsage returns total cumulative tokens used
func (c *Client) GetTotalCumulativeTokenUsage() int {
	return c.cumulativeInput + c.cumulativeOutput
}

// ResetTokenUsage resets token usage counters for current request
func (c *Client) ResetTokenUsage() {
	c.inputTokens = 0
	c.outputTokens = 0
}

// ResetCumulativeTokenUsage resets cumulative token usage counters
func (c *Client) ResetCumulativeTokenUsage() {
	c.cumulativeInput = 0
	c.cumulativeOutput = 0
	c.inputTokens = 0
	c.outputTokens = 0
}

// UsageStats represents token usage statistics
type UsageStats struct {
	InputTokens      int
	OutputTokens     int
	CumulativeInput  int
	CumulativeOutput int
	TokenLimit       int
}

// String returns a formatted string showing token usage
func (u UsageStats) String() string {
	var parts []string

	// Current request usage
	if u.InputTokens > 0 || u.OutputTokens > 0 {
		parts = append(parts, fmt.Sprintf("Request: %d in, %d out", u.InputTokens, u.OutputTokens))
	}

	// Cumulative usage
	if u.CumulativeInput > 0 || u.CumulativeOutput > 0 {
		totalCumulative := u.CumulativeInput + u.CumulativeOutput
		parts = append(parts, fmt.Sprintf("Total: %d in, %d out (%d)", u.CumulativeInput, u.CumulativeOutput, totalCumulative))
	}

	// Token limit if set
	if u.TokenLimit > 0 {
		totalUsed := u.CumulativeInput + u.CumulativeOutput
		percentage := float64(totalUsed) / float64(u.TokenLimit) * 100
		parts = append(parts, fmt.Sprintf("Limit: %d/%d (%.1f%%)", totalUsed, u.TokenLimit, percentage))
	}

	if len(parts) == 0 {
		return "Tokens: 0"
	}

	return "Tokens: " + strings.Join(parts, " | ")
}

// UsageDisplay returns a formatted string showing token usage
func (c *Client) UsageDisplay() string {
	stats := UsageStats{
		InputTokens:      c.inputTokens,
		OutputTokens:     c.outputTokens,
		CumulativeInput:  c.cumulativeInput,
		CumulativeOutput: c.cumulativeOutput,
		TokenLimit:       c.tokenLimit,
	}
	return stats.String()
}

