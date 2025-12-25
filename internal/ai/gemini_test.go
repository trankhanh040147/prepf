package ai

import (
	"context"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-key", 30, 1000000)
	
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}
	
	if client.apiKey != "test-key" {
		t.Errorf("expected apiKey 'test-key', got '%s'", client.apiKey)
	}
	
	if client.tokenLimit != 1000000 {
		t.Errorf("expected tokenLimit 1000000, got %d", client.tokenLimit)
	}
	
	if client.history == nil {
		t.Error("history should be initialized")
	}
}

func TestClient_SetTokenLimit(t *testing.T) {
	client := NewClient("test-key", 30, 1000000)
	
	client.SetTokenLimit(500000)
	if client.tokenLimit != 500000 {
		t.Errorf("expected tokenLimit 500000, got %d", client.tokenLimit)
	}
}

func TestClient_GetTokenUsage(t *testing.T) {
	client := NewClient("test-key", 30, 1000000)
	
	input, output := client.GetTokenUsage()
	if input != 0 || output != 0 {
		t.Errorf("expected initial usage (0, 0), got (%d, %d)", input, output)
	}
}

func TestClient_ResetTokenUsage(t *testing.T) {
	client := NewClient("test-key", 30, 1000000)
	client.inputTokens = 100
	client.outputTokens = 200
	
	client.ResetTokenUsage()
	
	input, output := client.GetTokenUsage()
	if input != 0 || output != 0 {
		t.Errorf("expected usage (0, 0) after reset, got (%d, %d)", input, output)
	}
}

func TestClient_GetCumulativeTokenUsage(t *testing.T) {
	client := NewClient("test-key", 30, 1000000)
	
	input, output := client.GetCumulativeTokenUsage()
	if input != 0 || output != 0 {
		t.Errorf("expected initial cumulative usage (0, 0), got (%d, %d)", input, output)
	}
}

func TestClient_UsageDisplay(t *testing.T) {
	client := NewClient("test-key", 30, 1000000)
	
	// Test empty usage
	display := client.UsageDisplay()
	if display == "" {
		t.Error("UsageDisplay() should return non-empty string")
	}
	
	// Test with usage
	client.inputTokens = 100
	client.outputTokens = 200
	display = client.UsageDisplay()
	if display == "" {
		t.Error("UsageDisplay() should return non-empty string with usage")
	}
}

func TestClient_History(t *testing.T) {
	client := NewClient("test-key", 30, 1000000)
	
	// Test history is initialized
	history := client.GetHistory()
	if history == nil {
		t.Fatal("GetHistory() returned nil")
	}
	
	// Test adding to history
	history.Add("user", "test message")
	if history.Count() != 1 {
		t.Errorf("expected history count 1, got %d", history.Count())
	}
	
	// Test clearing history
	client.ClearHistory()
	if history.Count() != 0 {
		t.Errorf("expected history count 0 after clear, got %d", history.Count())
	}
}

func TestClient_Stream_NoAPIKey(t *testing.T) {
	client := NewClient("", 30, 1000000)
	
	ctx := context.Background()
	_, err := client.Stream(ctx, "test prompt")
	
	if err == nil {
		t.Error("expected error when API key is empty")
	}
}

func TestClient_Stream_TokenLimitExceeded(t *testing.T) {
	client := NewClient("test-key", 30, 100) // Very low limit
	client.cumulativeInput = 50
	client.cumulativeOutput = 50
	
	ctx := context.Background()
	// Large prompt that would exceed limit
	largePrompt := make([]byte, 1000)
	for i := range largePrompt {
		largePrompt[i] = 'a'
	}
	
	_, err := client.Stream(ctx, string(largePrompt))
	
	if err == nil {
		t.Error("expected error when token limit would be exceeded")
	}
}

