package model

import "time"

// Entity is a typed, normalized graph entity.
type Entity struct {
	ID         string              `json:"id"`
	Type       EntityType          `json:"type"`
	Name       string              `json:"name"`
	Cluster    string              `json:"cluster"`
	Namespace  string              `json:"namespace,omitempty"`
	Labels     map[string]string   `json:"labels,omitempty"`
	Traits     map[string][]string `json:"traits,omitempty"`
	Attributes map[string]string   `json:"attributes,omitempty"`
}

// Relationship is a typed link between entities with evidence.
type Relationship struct {
	ID         string           `json:"id"`
	Type       RelationshipType `json:"type"`
	FromID     string           `json:"from_id"`
	ToID       string           `json:"to_id"`
	Evidence   string           `json:"evidence,omitempty"`
	Confidence float64          `json:"confidence"`
}

// CollectorError captures an isolated collector failure.
type CollectorError struct {
	Collector string `json:"collector"`
	Kind      string `json:"kind"`
	Message   string `json:"message"`
}

// CollectorStatus captures collector execution metadata.
type CollectorStatus struct {
	Name       string         `json:"name"`
	DurationMS int64          `json:"duration_ms"`
	StartedAt  time.Time      `json:"started_at"`
	FinishedAt time.Time      `json:"finished_at"`
	Counts     map[string]int `json:"counts,omitempty"`
}
