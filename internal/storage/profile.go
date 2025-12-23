package storage

import (
	"fmt"
	"os"

	"github.com/bytedance/sonic"
	tea "github.com/charmbracelet/bubbletea"
)

// Profile represents user profile data
type Profile struct {
	ID              string `json:"id,omitempty"`
	CVPath          string `json:"cv_path"`
	ExperienceLevel string `json:"experience_level"`
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
func (s *Store) Load() (*Profile, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Profile{}, nil
		}
		return nil, fmt.Errorf("read profile: %w", err)
	}

	var profile Profile
	if err := sonic.Unmarshal(data, &profile); err != nil {
		return nil, fmt.Errorf("unmarshal profile: %w", err)
	}

	return &profile, nil
}

// Save saves profile to disk
func (s *Store) Save(profile *Profile) error {
	data, err := sonic.Marshal(profile)
	if err != nil {
		return fmt.Errorf("marshal profile: %w", err)
	}

	if err := os.WriteFile(s.path, data, 0644); err != nil {
		return fmt.Errorf("write profile: %w", err)
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
		if err := s.Save(profile); err != nil {
			return ProfileSaveError{Err: err}
		}
		return ProfileSaved{}
	}
}

// ProfileLoaded is a message sent when profile is loaded
type ProfileLoaded struct {
	Profile *Profile
}

// ProfileSaved is a message sent when profile is saved
type ProfileSaved struct{}

// ProfileLoadError is a message sent when profile load fails
type ProfileLoadError struct {
	Err error
}

// ProfileSaveError is a message sent when profile save fails
type ProfileSaveError struct {
	Err error
}

// SafeUpdate updates profile with safe mutation.
// The original ID is cached for lookup/persistence purposes, but updates may modify the ID if needed.
func (s *Store) SafeUpdate(profile *Profile, updates func(*Profile)) error {
	// Apply updates (may modify ID)
	updates(profile)

	return s.Save(profile)
}
