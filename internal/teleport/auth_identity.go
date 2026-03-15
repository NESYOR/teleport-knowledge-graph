package teleport

import (
	"fmt"
	"os"
	"path/filepath"
)

// LoadIdentityAuth resolves identity-file based auth material.
func LoadIdentityAuth(identityFile string) (AuthSource, error) {
	if identityFile == "" {
		return AuthSource{}, Wrap("identity auth", fmt.Errorf("%w: identity file path empty", ErrAuthUnavailable))
	}
	if _, err := os.Stat(identityFile); err != nil {
		return AuthSource{}, Wrap("identity auth", fmt.Errorf("%w: %s", ErrAuthUnavailable, err))
	}
	return AuthSource{Mode: "identity", Detail: filepath.Clean(identityFile)}, nil
}
