package teleport

import (
	"errors"
	"testing"
)

func TestErrorKind(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want string
	}{
		{name: "permission", err: errors.New("permission denied"), want: "permission_denied"},
		{name: "unsupported", err: ErrUnsupportedMethod, want: "unsupported_method"},
		{name: "network", err: errors.New("network unreachable"), want: "network_error"},
		{name: "fallback", err: errors.New("x"), want: "internal_error"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ErrorKind(tc.err); got != tc.want {
				t.Fatalf("want %s got %s", tc.want, got)
			}
		})
	}
}
