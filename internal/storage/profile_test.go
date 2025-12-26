package storage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStore_Load(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "profile.md")
	store := NewStore(path)

	// Test loading non-existent profile (should return empty profile)
	profile, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if profile == nil {
		t.Fatal("Load() returned nil profile")
	}

	if profile.ID != "" {
		t.Errorf("expected empty ID for new profile, got '%s'", profile.ID)
	}
}

func TestStore_Save(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "profile.md")
	store := NewStore(path)

	profile := &Profile{
		ID:              "test-id",
		CVPath:          "/path/to/cv.md",
		ExperienceLevel: "senior",
	}

	err := store.Save(profile)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("profile file was not created")
	}

	// Reload and verify
	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load() after Save() error = %v", err)
	}

	if loaded.ID != "test-id" {
		t.Errorf("expected ID 'test-id', got '%s'", loaded.ID)
	}

	if loaded.CVPath != "/path/to/cv.md" {
		t.Errorf("expected CVPath '/path/to/cv.md', got '%s'", loaded.CVPath)
	}

	if loaded.ExperienceLevel != "senior" {
		t.Errorf("expected ExperienceLevel 'senior', got '%s'", loaded.ExperienceLevel)
	}
}

func TestStore_SaveCreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	// Create a nested path that doesn't exist
	path := filepath.Join(tmpDir, "nested", "dir", "profile.md")
	store := NewStore(path)

	profile := &Profile{
		CVPath: "/test/cv.md",
	}

	err := store.Save(profile)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify directory was created
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatal("directory was not created")
	}

	// Verify file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("profile file was not created")
	}
}

func TestStore_Load_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "profile.md")
	store := NewStore(path)

	// Create empty file
	if err := os.WriteFile(path, []byte(""), 0644); err != nil {
		t.Fatalf("failed to create empty file: %v", err)
	}

	profile, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if profile == nil {
		t.Fatal("Load() returned nil profile")
	}

	if profile.ID != "" {
		t.Errorf("expected empty ID for empty file, got '%s'", profile.ID)
	}
}

func TestStore_Load_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "profile.md")
	store := NewStore(path)

	// Create file with invalid JSON
	invalidJSON := `{"cv_path": "test", invalid}`
	if err := os.WriteFile(path, []byte(invalidJSON), 0644); err != nil {
		t.Fatalf("failed to create invalid JSON file: %v", err)
	}

	_, err := store.Load()
	if err == nil {
		t.Fatal("Load() should return error for invalid JSON")
	}

	if !strings.Contains(err.Error(), "not valid JSON") {
		t.Errorf("expected error message about invalid JSON, got: %v", err)
	}
}

func TestStore_Load_PlainText(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "profile.md")
	store := NewStore(path)

	// Create file with plain text (resume content)
	plainText := "Danh Son Ha\nPhone: (+84) 357 607"
	if err := os.WriteFile(path, []byte(plainText), 0644); err != nil {
		t.Fatalf("failed to create plain text file: %v", err)
	}

	profile, err := store.Load()
	if err != nil {
		t.Fatalf("Load() should not return error for plain text (markdown/txt): %v", err)
	}

	if profile == nil {
		t.Fatal("Load() returned nil profile")
	}

	// Profile file contains resume content, so CVPath should be set to the file path
	if profile.CVPath != path {
		t.Errorf("expected CVPath to be '%s', got '%s'", path, profile.CVPath)
	}
}
