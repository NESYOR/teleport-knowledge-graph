package ai

import (
	"context"
	"fmt"
)

// StaticAdapter is a deterministic adapter for tests/local runs.
type StaticAdapter struct {
	Model string
}

// Generate returns deterministic content based on request prompt.
func (a StaticAdapter) Generate(ctx context.Context, req GenerateRequest) (GenerateResponse, error) {
	select {
	case <-ctx.Done():
		return GenerateResponse{}, ctx.Err()
	default:
	}
	if req.Prompt == "" {
		return GenerateResponse{}, fmt.Errorf("prompt is required")
	}
	model := a.Model
	if model == "" {
		model = "static-v1"
	}
	return GenerateResponse{Content: "deterministic-response", Model: model, Metadata: map[string]string{"length": fmt.Sprintf("%d", len(req.Prompt))}}, nil
}
