package models

import "time"

// Pipeline represents a deployment pipeline definition
type Pipeline struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Definition  string    `json:"definition"` // YAML or JSON string
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PipelineDefinition represents the parsed pipeline structure
type PipelineDefinition struct {
	Name        string                 `yaml:"name" json:"name"`
	Description string                 `yaml:"description" json:"description,omitempty"`
	Steps       []PipelineStep         `yaml:"steps" json:"steps"`
	Variables   map[string]interface{} `yaml:"variables" json:"variables,omitempty"`
}

// PipelineStep represents a single step in a pipeline
type PipelineStep struct {
	Name        string                 `yaml:"name" json:"name"`
	Description string                 `yaml:"description" json:"description,omitempty"`
	Type        string                 `yaml:"type" json:"type"`     // docker, script, etc.
	Action      string                 `yaml:"action" json:"action"` // run, pull, start, stop, etc.
	Config      map[string]interface{} `yaml:"config" json:"config"`
	OnFailure   string                 `yaml:"on_failure" json:"on_failure,omitempty"` // continue, stop, rollback
	Retry       int                    `yaml:"retry" json:"retry,omitempty"`
}
