package ai

// HistoryEntry represents a single message in conversation history
type HistoryEntry struct {
	Role    string
	Content string
}

// History manages conversation context for multi-turn interactions
type History struct {
	entries []HistoryEntry
}

// NewHistory creates a new conversation history
func NewHistory() *History {
	return &History{
		entries: make([]HistoryEntry, 0),
	}
}

// Add adds a message to the history
func (h *History) Add(role, content string) {
	h.entries = append(h.entries, HistoryEntry{
		Role:    role,
		Content: content,
	})
}

// Clear clears all history entries
func (h *History) Clear() {
	h.entries = make([]HistoryEntry, 0)
}

// Entries returns all history entries
func (h *History) Entries() []HistoryEntry {
	// Return a copy to prevent external mutation
	result := make([]HistoryEntry, len(h.entries))
	copy(result, h.entries)
	return result
}

// Count returns the number of entries in history
func (h *History) Count() int {
	return len(h.entries)
}

// ToContents converts history to Content slice for API requests
func (h *History) ToContents() []Content {
	contents := make([]Content, 0, len(h.entries))
	for _, entry := range h.entries {
		contents = append(contents, Content{
			Role: entry.Role,
			Parts: []Part{
				{Text: entry.Content},
			},
		})
	}
	return contents
}

// AddToHistory adds a user message and assistant response to history
func (h *History) AddToHistory(userMsg, assistantMsg string) {
	h.Add("user", userMsg)
	if assistantMsg != "" {
		h.Add("assistant", assistantMsg)
	}
}

