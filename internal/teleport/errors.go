package teleport

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrPermissionDenied indicates Teleport API denied access.
	ErrPermissionDenied = errors.New("permission denied")
	// ErrUnsupportedMethod indicates the API/version does not support the operation.
	ErrUnsupportedMethod = errors.New("unsupported method")
	// ErrAuthUnavailable indicates auth material was unavailable.
	ErrAuthUnavailable = errors.New("auth unavailable")
)

// ErrorKind maps errors to stable collector error kinds.
func ErrorKind(err error) string {
	if err == nil {
		return ""
	}
	msg := strings.ToLower(err.Error())
	switch {
	case errors.Is(err, ErrPermissionDenied), strings.Contains(msg, "access denied"), strings.Contains(msg, "permission"):
		return "permission_denied"
	case errors.Is(err, ErrUnsupportedMethod), strings.Contains(msg, "unsupported"):
		return "unsupported_method"
	case errors.Is(err, ErrAuthUnavailable), strings.Contains(msg, "identity"), strings.Contains(msg, "profile"):
		return "auth_error"
	case strings.Contains(msg, "timeout"):
		return "timeout"
	case strings.Contains(msg, "connection"), strings.Contains(msg, "network"):
		return "network_error"
	default:
		return "internal_error"
	}
}

// Wrap annotates an integration error with operation context.
func Wrap(operation string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("teleport %s: %w", operation, err)
}
