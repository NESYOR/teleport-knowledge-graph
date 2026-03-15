package model

import "time"

// SnapshotSource describes how a snapshot was collected.
type SnapshotSource struct {
	ClusterName string `json:"cluster_name"`
	AuthSource  string `json:"auth_source"`
	Caller      string `json:"caller"`
}

// Snapshot stores a normalized digital twin capture.
type Snapshot struct {
	SnapshotID      string            `json:"snapshot_id"`
	CollectedAt     time.Time         `json:"collected_at"`
	Source          SnapshotSource    `json:"source"`
	CollectorStatus []CollectorStatus `json:"collector_status"`
	CollectorErrors []CollectorError  `json:"collector_errors,omitempty"`
	Entities        []Entity          `json:"entities"`
	Relationships   []Relationship    `json:"relationships"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}
