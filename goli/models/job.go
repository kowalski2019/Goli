package models

import "time"

// JobStatus represents the status of a job
type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
	JobStatusCancelled JobStatus = "cancelled"
)

// Job represents a deployment job
type Job struct {
	ID           int64      `json:"id"`
	PipelineID   *int64     `json:"pipeline_id,omitempty"`
	Name         string     `json:"name"`
	Status       JobStatus  `json:"status"`
	TriggeredBy  string     `json:"triggered_by,omitempty"`
	StartedAt    *time.Time `json:"started_at,omitempty"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	ErrorMessage string     `json:"error_message,omitempty"`
	Logs         string     `json:"logs,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	Steps        []JobStep  `json:"steps,omitempty"`
}

// JobStep represents a single step in a job
type JobStep struct {
	ID           int64      `json:"id"`
	JobID        int64      `json:"job_id"`
	StepName     string     `json:"step_name"`
	StepOrder    int        `json:"step_order"`
	Status       JobStatus  `json:"status"`
	StartedAt    *time.Time `json:"started_at,omitempty"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	ErrorMessage string     `json:"error_message,omitempty"`
	Logs         string     `json:"logs,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}
