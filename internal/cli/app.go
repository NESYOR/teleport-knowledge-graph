package cli

import "context"

// App is a minimal CLI command runner contract for future command wiring.
type App interface {
	Run(ctx context.Context, args []string) error
}
