package teleport

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AuthSource describes the auth method used for Teleport API calls.
type AuthSource struct {
	Mode   string `json:"mode"`
	Detail string `json:"detail"`
}

// LoadProfileAuth resolves profile-based auth material.
func LoadProfileAuth(profilePath string) (AuthSource, error) {
	if profilePath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return AuthSource{}, Wrap("profile auth", fmt.Errorf("resolve home dir: %w", err))
		}
		profilePath = filepath.Join(home, ".tsh")
	}
	if _, err := os.Stat(profilePath); err != nil {
		return AuthSource{}, Wrap("profile auth", fmt.Errorf("%w: %s", ErrAuthUnavailable, err))
	}
	return AuthSource{Mode: "profile", Detail: profilePath}, nil
}

// ParseProfileName extracts a profile name hint from a path.
func ParseProfileName(profilePath string) string {
	base := filepath.Base(profilePath)
	if base == ".tsh" {
		return "default"
	}
	return strings.TrimSuffix(base, filepath.Ext(base))
}
