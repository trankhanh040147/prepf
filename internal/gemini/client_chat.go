package gemini

import (
	"context"
	"fmt"
	"log"
	"strings"

	"golang.org/x/sync/errgroup"
	"google.golang.org/genai"
)

// StreamCallback is called for each chunk during streaming
type StreamCallback func(string)

// SendMessage sends a message and returns the full response
func (c *Client) SendMessage(ctx context.Context, message string) (string, error) {
	_, config := c.appendUserMessageAndPrepareTurn(message)

	resp, err := c.client.Models.GenerateContent(ctx, c.modelID, c.history, config)
	if err != nil {
		c.rollbackLastHistoryEntry()
		return "", fmt.Errorf("generate content: %w", err)
	}

	// Check for safety block
	if len(resp.Candidates) > 0 && resp.Candidates[0].FinishReason == genai.FinishReasonSafety {
		c.rollbackLastHistoryEntry()
		return "", fmt.Errorf("response blocked by safety filters")
	}

	// Add assistant response to history
	if len(resp.Candidates) > 0 {
		c.history = append(c.history, resp.Candidates[0].Content)
	}

	c.lastUsage = extractUsage(resp)
	return extractText(resp), nil
}

// SendMessageStream streams responses while maintaining history
func (c *Client) SendMessageStream(ctx context.Context, message string, callback StreamCallback) (string, error) {
	_, config := c.appendUserMessageAndPrepareTurn(message)

	iter := c.client.Models.GenerateContentStream(ctx, c.modelID, c.history, config)

	var fullResponseText string
	var finalContent *genai.Content

	// Convert iter.Seq2 to channel-based consumption pattern
	type iterResult struct {
		resp *genai.GenerateContentResponse
		err  error
	}
	ch := make(chan iterResult, StreamChannelBufferSize)

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		defer close(ch)
		iter(func(resp *genai.GenerateContentResponse, err error) bool {
			select {
			case ch <- iterResult{resp: resp, err: err}:
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
			log.Printf("Warning: errgroup returned error: %v", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			c.rollbackLastHistoryEntry()
			return fullResponseText, fmt.Errorf("stream cancelled: %w", ctx.Err())
		case result, ok := <-ch:
			if !ok {
				goto done
			}
			if result.err != nil {
				c.rollbackLastHistoryEntry()
				return fullResponseText, fmt.Errorf("stream error: %w", result.err)
			}

			resp := result.resp
			// Check for safety block
			if len(resp.Candidates) > 0 && resp.Candidates[0].FinishReason == genai.FinishReasonSafety {
				c.rollbackLastHistoryEntry()
				return fullResponseText, fmt.Errorf("response blocked by safety filters")
			}

			chunk := extractText(resp)
			fullResponseText += chunk
			if callback != nil {
				callback(chunk)
			}

			if len(resp.Candidates) > 0 {
				finalContent = resp.Candidates[0].Content
				c.lastUsage = extractUsage(resp)
			}
		}
	}

done:
	// Add final response to history
	if finalContent != nil {
		c.history = append(c.history, finalContent)
	}

	return fullResponseText, nil
}

// GenerateContent performs a one-off generation without maintaining history
func (c *Client) GenerateContent(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	config := c.newGenerationConfig(systemPrompt)

	resp, err := c.client.Models.GenerateContent(ctx, c.modelID, genai.Text(userPrompt), config)
	if err != nil {
		return "", fmt.Errorf("generate content: %w", err)
	}

	c.lastUsage = extractUsage(resp)
	return extractText(resp), nil
}

// extractText extracts text from response candidates
func extractText(resp *genai.GenerateContentResponse) string {
	if resp == nil || len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return ""
	}
	if len(resp.Candidates) > 1 {
		log.Printf("Warning: Multiple candidates (%d) in response, using first only", len(resp.Candidates))
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

// extractUsage extracts token usage from response
func extractUsage(resp *genai.GenerateContentResponse) *TokenUsage {
	if resp == nil || resp.UsageMetadata == nil {
		return nil
	}
	return &TokenUsage{
		PromptTokens:     resp.UsageMetadata.PromptTokenCount,
		CompletionTokens: resp.UsageMetadata.CandidatesTokenCount,
		TotalTokens:      resp.UsageMetadata.TotalTokenCount,
	}
}
