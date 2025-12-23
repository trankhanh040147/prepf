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
)

// Client handles Gemini API interactions
type Client struct {
	apiKey     string
	baseURL    string
	timeout    time.Duration
	httpClient *http.Client
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
	Candidates []Candidate `json:"candidates"`
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

		// Create request with context
		ctx, cancel := context.WithTimeout(ctx, c.timeout)
		defer cancel()

		url := fmt.Sprintf("%s/models/gemini-pro:streamGenerateContent?key=%s", c.baseURL, c.apiKey)

		httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
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
					if len(streamResp.Candidates) > 0 && len(streamResp.Candidates[0].Content.Parts) > 0 {
						ch <- StreamChunk{Text: streamResp.Candidates[0].Content.Parts[0].Text}
					}
				}
			}
		}

		if err := scanner.Err(); err != nil {
			ch <- StreamChunk{Err: fmt.Errorf("read response stream: %w", err)}
			return
		}

		ch <- StreamChunk{Done: true}
	}()

	return ch, nil
}

// GetTokenUsage returns token usage information
func (c *Client) GetTokenUsage() (int, int) {
	// Placeholder - actual implementation would track tokens
	return 0, 0
}

