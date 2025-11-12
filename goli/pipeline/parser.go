package pipeline

import (
	"goli/models"
	"strconv"

	"gopkg.in/yaml.v3"
)

// ParsePipelineDefinition parses a YAML pipeline definition
func ParsePipelineDefinition(yamlContent string) (*models.PipelineDefinition, error) {
	var def models.PipelineDefinition

	if err := yaml.Unmarshal([]byte(yamlContent), &def); err != nil {
		return nil, err
	}

	return &def, nil
}

// ValidatePipelineDefinition validates a pipeline definition
func ValidatePipelineDefinition(def *models.PipelineDefinition) error {
	if def.Name == "" {
		return &PipelineError{Message: "Pipeline name is required"}
	}

	if len(def.Steps) == 0 {
		return &PipelineError{Message: "Pipeline must have at least one step"}
	}

	for i, step := range def.Steps {
		if step.Name == "" {
			return &PipelineError{Message: "Step name is required for step " + strconv.Itoa(i+1)}
		}
		if step.Type == "" {
			return &PipelineError{Message: "Step type is required for step " + step.Name}
		}
		if step.Action == "" {
			return &PipelineError{Message: "Step action is required for step " + step.Name}
		}
	}

	return nil
}
