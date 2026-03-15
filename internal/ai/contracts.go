package ai

import "context"

// Adapter abstracts provider-specific LLM calls without binding to a single provider.
type Adapter interface {
	Generate(ctx context.Context, req GenerateRequest) (GenerateResponse, error)
}

// GenerateRequest is a provider-neutral generation request.
type GenerateRequest struct {
	Prompt      string            `json:"prompt"`
	Temperature float64           `json:"temperature,omitempty"`
	MaxTokens   int               `json:"max_tokens,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// GenerateResponse is a provider-neutral generation response.
type GenerateResponse struct {
	Content  string            `json:"content"`
	Model    string            `json:"model,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// ContextBuilder builds compact graph context for AI prompts.
type ContextBuilder interface {
	Build(question string, maxNodes int) (ContextBundle, error)
}

// Summarizer emits deterministic summaries for entities and analysis outputs.
type Summarizer interface {
	SummarizeUser(userID string, maxResources int) (UserSummary, error)
	SummarizeRole(roleID string) (RoleSummary, error)
	SummarizeResource(resourceID string) (ResourceSummary, error)
}

// RuleReasoner provides deterministic, local reasoning for AI-readable explanations.
type RuleReasoner interface {
	WhyUserCanAccess(userID, resourceID string) (ReasoningResult, error)
}

// ContextBundle is a bounded graph/context payload ready for prompt generation.
type ContextBundle struct {
	Question      string         `json:"question"`
	SchemaVersion string         `json:"schema_version"`
	Nodes         []ContextNode  `json:"nodes"`
	Edges         []ContextEdge  `json:"edges"`
	Facts         []string       `json:"facts,omitempty"`
	Metadata      map[string]any `json:"metadata,omitempty"`
}

// ContextNode represents a compact node in AI context output.
type ContextNode struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

// ContextEdge represents a compact typed edge in AI context output.
type ContextEdge struct {
	Type string `json:"type"`
	From string `json:"from"`
	To   string `json:"to"`
}

// UserSummary is deterministic summary output for users.
type UserSummary struct {
	UserID           string   `json:"user_id"`
	Roles            []string `json:"roles"`
	ReachableTargets []string `json:"reachable_targets"`
}

// RoleSummary is deterministic summary output for roles.
type RoleSummary struct {
	RoleID        string   `json:"role_id"`
	AllowedLogins []string `json:"allowed_logins"`
	Capabilities  []string `json:"capabilities"`
}

// ResourceSummary is deterministic summary output for resources.
type ResourceSummary struct {
	ResourceID string   `json:"resource_id"`
	ExposedTo  []string `json:"exposed_to"`
	Labels     []string `json:"labels"`
}

// ReasoningResult is rule-based explanation output.
type ReasoningResult struct {
	Question   string   `json:"question"`
	Conclusion string   `json:"conclusion"`
	Path       []string `json:"path"`
	Evidence   []string `json:"evidence"`
}
