package ai

import (
	"encoding/json"
	"fmt"
)

// PromptTemplate renders prompts from deterministic context bundles.
type PromptTemplate interface {
	Render(context ContextBundle) (string, error)
}

// AccessWhyTemplate renders a prompt for access explanation.
type AccessWhyTemplate struct{}

// Render returns deterministic prompt text.
func (AccessWhyTemplate) Render(context ContextBundle) (string, error) {
	payload, err := json.Marshal(context)
	if err != nil {
		return "", fmt.Errorf("marshal context: %w", err)
	}
	return "You are a security reasoning assistant. Use only the JSON context to answer.\nCONTEXT:" + string(payload), nil
}
