package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bytedance/sonic"
	tea "github.com/charmbracelet/bubbletea"
)

// ValidExperienceLevels contains the list of valid experience levels
var ValidExperienceLevels = []string{"junior", "mid", "senior", "principal", "architect"}

// Profile represents user profile data
type Profile struct {
	ID              string `json:"id,omitempty"`
	CVPath          string `json:"cv_path"`
	ExperienceLevel string `json:"experience_level"`
}

// Validate validates profile fields
// NOTE: File existence checks (e.g., CVPath) are deferred to actual usage to avoid I/O in validation.
func (p *Profile) Validate() error {
	if p.ExperienceLevel != "" {
		found := false
		for _, level := range ValidExperienceLevels {
			if strings.EqualFold(p.ExperienceLevel, level) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("invalid experience level '%s', must be one of: %v", p.ExperienceLevel, ValidExperienceLevels)
		}
	}

	// CVPath validation: only check format, not file existence
	// File existence should be checked when the file is actually read
	if p.CVPath != "" {
		// Basic path validation (non-empty, no other checks to avoid I/O)
		if strings.TrimSpace(p.CVPath) == "" {
			return fmt.Errorf("CV path cannot be empty or whitespace")
		}
	}

	return nil
}

// Store handles profile persistence
type Store struct {
	path string
}

// NewStore creates a new profile store
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Load loads profile from disk
// If the file is JSON, it unmarshals as Profile struct.
// If the file is markdown/txt (resume content), it returns a Profile with CVPath set to the file path.
func (s *Store) Load() (*Profile, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Profile{}, nil
		}
		return nil, fmt.Errorf("read profile: %w", err)
	}

	// Check if file is empty or whitespace-only
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" {
		return &Profile{}, nil
	}

	// Try to unmarshal as JSON first
	var profile Profile
	if err := sonic.Unmarshal(data, &profile); err != nil {
		// If it's not JSON (doesn't start with { or [), treat it as markdown/txt resume content
		if !strings.HasPrefix(trimmed, "{") && !strings.HasPrefix(trimmed, "[") {
			// Profile file contains resume content, so CVPath is the file itself
			return &Profile{
				CVPath: s.path,
			}, nil
		}
		// Invalid JSON format
		return nil, fmt.Errorf("profile file is not valid JSON: %w", err)
	}

	return &profile, nil
}

// Save saves profile to disk atomically
func (s *Store) Save(profile *Profile) error {
	// Validate before saving
	if err := profile.Validate(); err != nil {
		return fmt.Errorf("validate profile: %w", err)
	}

	data, err := sonic.Marshal(profile)
	if err != nil {
		return fmt.Errorf("marshal profile: %w", err)
	}

	// Ensure directory exists before writing
	if err := os.MkdirAll(filepath.Dir(s.path), 0755); err != nil {
		return fmt.Errorf("create profile directory: %w", err)
	}

	// Atomic write: write to temp file then rename
	tmpPath := s.path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("write temp profile: %w", err)
	}

	// Rename is atomic on most filesystems
	if err := os.Rename(tmpPath, s.path); err != nil {
		// Clean up temp file on error
		_ = os.Remove(tmpPath)
		return fmt.Errorf("rename profile: %w", err)
	}

	return nil
}

// LoadCmd returns a tea.Cmd for loading profile
func (s *Store) LoadCmd() tea.Cmd {
	return func() tea.Msg {
		profile, err := s.Load()
		if err != nil {
			return ProfileLoadError{Err: err}
		}
		return ProfileLoaded{Profile: profile}
	}
}

// SaveCmd returns a tea.Cmd for saving profile
func (s *Store) SaveCmd(profile *Profile) tea.Cmd {
	return func() tea.Msg {
		originalID := profile.ID
		if err := s.Save(profile); err != nil {
			return ProfileSaveError{Err: err}
		}
		return ProfileSaved{OriginalID: originalID, NewID: profile.ID}
	}
}

// ProfileLoaded is a message sent when profile is loaded
type ProfileLoaded struct {
	Profile *Profile
}

// ProfileSaved is a message sent when profile is saved
type ProfileSaved struct {
	OriginalID string // ID before update (for tracking/deletion)
	NewID      string // ID after update
}

// ProfileLoadError is a message sent when profile load fails
type ProfileLoadError struct {
	Err error
}

// ProfileSaveError is a message sent when profile save fails
type ProfileSaveError struct {
	Err error
}

// safeUpdateResult holds the result of a safe update operation
type safeUpdateResult struct {
	originalID string
	newID      string
	err        error
}

// safeUpdate applies updates to a profile and saves it
// Returns the original ID, new ID, and any error
func (s *Store) safeUpdate(profile *Profile, updates func(*Profile)) (string, string, error) {
	// Cache original ID before any mutations
	originalID := profile.ID

	// Apply updates (may modify ID)
	updates(profile)

	// Validate after updates
	if err := profile.Validate(); err != nil {
		return originalID, profile.ID, fmt.Errorf("validation failed: %w", err)
	}

	if err := s.Save(profile); err != nil {
		return originalID, profile.ID, err
	}

	return originalID, profile.ID, nil
}

// SafeUpdateCmd applies updates to a profile and saves it
// Caches the original ID before applying updates to ensure safe mutation
// (important for future deletion operations that need the original ID)
func (s *Store) SafeUpdateCmd(profile *Profile, updates func(*Profile)) tea.Cmd {
	return func() tea.Msg {
		originalID, newID, err := s.safeUpdate(profile, updates)
		if err != nil {
			return ProfileSaveError{Err: err}
		}
		return ProfileSaved{OriginalID: originalID, NewID: newID}
	}
}

// SafeUpdate applies updates to a profile and saves it (non-TUI version)
// Caches the original ID before applying updates to ensure safe mutation
func (s *Store) SafeUpdate(profile *Profile, updates func(*Profile)) error {
	_, _, err := s.safeUpdate(profile, updates)
	return err
}
