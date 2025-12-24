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
	apiKey       string
	baseURL      string
	timeout      time.Duration
	httpClient   *http.Client
	inputTokens  int
	outputTokens int
}

// NewClient creates a new Gemini client
func NewClient(apiKey string, timeout int) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: "https://generativelanguage.googleapis.com/v1beta",
		timeout: time.Duration(timeout) * time.Second,
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
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

	ch := make(chan StreamChunk, 10)

	go func() {
		defer close(ch)

		req := StreamRequest{
			Contents: []Content{
				{
					Role: "user",
					Parts: []Part{
						{Text: prompt},
					},
				},
			},
		}

		reqBody, err := sonic.Marshal(req)
		if err != nil {
			ch <- StreamChunk{Err: fmt.Errorf("marshal request: %w", err)}
			return
		}

		url := fmt.Sprintf("%s/models/gemini-pro:streamGenerateContent?key=%s", c.baseURL, c.apiKey)

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

		// Stream response with a scanner
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
					}

					if len(streamResp.Candidates) > 0 && len(streamResp.Candidates[0].Content.Parts) > 0 {
						ch <- StreamChunk{Text: streamResp.Candidates[0].Content.Parts[0].Text}
					}
				}
			}
		}

		if err = scanner.Err(); err != nil {
			ch <- StreamChunk{Err: fmt.Errorf("read response stream: %w", err)}
			return
		}

		ch <- StreamChunk{Done: true}
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

// GetTokenUsage returns token usage information (input tokens, output tokens)
func (c *Client) GetTokenUsage() (int, int) {
	return c.inputTokens, c.outputTokens
}

// ResetTokenUsage resets token usage counters
func (c *Client) ResetTokenUsage() {
	c.inputTokens = 0
	c.outputTokens = 0
}
