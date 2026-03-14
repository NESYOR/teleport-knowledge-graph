package api

// Envelope is the stable API response envelope.
type Envelope struct {
	RequestID string      `json:"request_id"`
	Version   string      `json:"version"`
	Data      interface{} `json:"data,omitempty"`
	Errors    []string    `json:"errors,omitempty"`
}
