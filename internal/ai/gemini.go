package ai

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	tea "github.com/charmbracelet/bubbletea"
)

// Client handles Gemini API interactions
type Client struct {
	apiKey           string
	baseURL          string
	timeout          time.Duration
	httpClient       *http.Client
	inputTokens      int
	outputTokens     int
	tokenLimit       int
	cumulativeInput  int
	cumulativeOutput int
	history          *History
}

// NewClient creates a new Gemini client
func NewClient(apiKey string, timeout int, tokenLimit int) *Client {
	return &Client{
		apiKey:     apiKey,
		baseURL:    "https://generativelanguage.googleapis.com/v1beta",
		timeout:    time.Duration(timeout) * time.Second,
		tokenLimit: tokenLimit,
		history:    NewHistory(),
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
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

// StreamRequest represents a streaming request
type StreamRequest struct {
	Contents []Content `json:"contents"`
}

// Content represents a message content
type Content struct {
	Role  string `json:"role"`
	Parts []Part `json:"parts"`
}

// Part represents a content part
type Part struct {
	Text string `json:"text"`
}

// StreamResponse represents a streaming response
type StreamResponse struct {
	Candidates    []Candidate    `json:"candidates"`
	UsageMetadata *UsageMetadata `json:"usageMetadata,omitempty"`
}

// UsageMetadata represents token usage information
type UsageMetadata struct {
	PromptTokenCount     int `json:"promptTokenCount"`
	CandidatesTokenCount int `json:"candidatesTokenCount"`
	TotalTokenCount      int `json:"totalTokenCount"`
}

// Candidate represents a response candidate
type Candidate struct {
	Content Content `json:"content"`
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

		// Build contents from history + current prompt
		contents := make([]Content, 0)
		if c.history != nil {
			contents = append(contents, c.history.ToContents()...)
		}
		// Add current user prompt
		contents = append(contents, Content{
			Role: "user",
			Parts: []Part{
				{Text: prompt},
			},
		})

		req := StreamRequest{
			Contents: contents,
		}

		reqBody, err := sonic.Marshal(req)
		if err != nil {
			ch <- StreamChunk{Err: fmt.Errorf("marshal request: %w", err)}
			return
		}

		url := fmt.Sprintf("%s/models/gemini-2.5-flash:streamGenerateContent?key=%s", c.baseURL, c.apiKey)

		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
		if err != nil {
			ch <- StreamChunk{Err: fmt.Errorf("create request: %w", err)}
			return
		}

		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(httpReq)
		if err != nil {
			ch <- StreamChunk{Err: fmt.Errorf("send request: %w", err)}
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			bodyStr := string(body)
			if err != nil {
				// If we can't read the error body, include the read error in the message
				bodyStr = fmt.Sprintf("(failed to read error body: %v)", err)
			}
			ch <- StreamChunk{Err: fmt.Errorf("api error: %s - %s", resp.Status, bodyStr)}
			return
		}

		// Stream response with a scanner and collect for history
		var fullResponse strings.Builder
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				ch <- StreamChunk{Err: fmt.Errorf("context cancelled: %w", ctx.Err())}
				return
			default:
				line := scanner.Text()
				// Parse SSE format: check for "data: " prefix
				if text, found := strings.CutPrefix(line, "data: "); found {
					var streamResp StreamResponse
					if err := sonic.Unmarshal([]byte(text), &streamResp); err != nil {
						ch <- StreamChunk{Err: fmt.Errorf("unmarshal stream chunk: %w", err)}
						continue
					}

					// Update token usage if available
					if streamResp.UsageMetadata != nil {
						c.inputTokens = streamResp.UsageMetadata.PromptTokenCount
						c.outputTokens = streamResp.UsageMetadata.CandidatesTokenCount
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

					if len(streamResp.Candidates) > 0 && len(streamResp.Candidates[0].Content.Parts) > 0 {
						chunkText := streamResp.Candidates[0].Content.Parts[0].Text
						fullResponse.WriteString(chunkText)
						ch <- StreamChunk{Text: chunkText}
					}
				}
			}
		}

		if err = scanner.Err(); err != nil {
			ch <- StreamChunk{Err: fmt.Errorf("read response stream: %w", err)}
			return
		}

		ch <- StreamChunk{Done: true}

		// Add to history after stream completes
		if c.history != nil && fullResponse.Len() > 0 {
			c.history.AddToHistory(prompt, fullResponse.String())
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
