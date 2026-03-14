package model

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
)

// DeterministicID creates a deterministic identifier for scoped entities.
func DeterministicID(prefix string, parts ...string) string {
	norm := strings.Join(parts, ":")
	h := sha1.Sum([]byte(norm))
	return fmt.Sprintf("%s:%s", prefix, hex.EncodeToString(h[:8]))
}
