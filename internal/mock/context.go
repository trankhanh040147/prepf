package mock

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadContext loads resume/CV context from a file (.txt or .md)
func LoadContext(path string) (string, error) {
	if path == "" {
		return "", nil
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".txt" && ext != ".md" {
		return "", fmt.Errorf("unsupported file type: %s (expected .txt or .md)", ext)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("resume file not found: %s", path)
		}
		return "", fmt.Errorf("read resume file: %w", err)
	}

	content := strings.TrimSpace(string(data))
	return content, nil
}

