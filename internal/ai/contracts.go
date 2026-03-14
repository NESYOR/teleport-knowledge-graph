package ai

import "context"

// Adapter abstracts provider-specific AI model calls.
type Adapter interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

// ContextBuilder builds compact graph context for AI prompts.
type ContextBuilder interface {
	Build(question string, maxNodes int) (string, error)
}
