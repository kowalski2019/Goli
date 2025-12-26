package pipeline

import (
	"goli/models"
	"regexp"
	"strconv"
	"strings"

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

// SubstituteVariables substitutes variables in a pipeline definition
// Supports ${VAR_NAME} and {{VAR_NAME}} syntax
func SubstituteVariables(def *models.PipelineDefinition, variables map[string]interface{}) {
	// Substitute in step configs
	for i := range def.Steps {
		substituteInMap(def.Steps[i].Config, variables)
	}
}

// substituteInMap recursively substitutes variables in a map
func substituteInMap(m map[string]interface{}, variables map[string]interface{}) {
	for key, value := range m {
		switch v := value.(type) {
		case string:
			m[key] = substituteString(v, variables)
		case []interface{}:
			for i, item := range v {
				if str, ok := item.(string); ok {
					v[i] = substituteString(str, variables)
				} else if subMap, ok := item.(map[string]interface{}); ok {
					substituteInMap(subMap, variables)
				}
			}
		case map[string]interface{}:
			substituteInMap(v, variables)
		}
	}
}

// substituteString replaces variable placeholders in a string
// Supports ${VAR_NAME} and {{VAR_NAME}} syntax
func substituteString(s string, variables map[string]interface{}) string {
	// Pattern for ${VAR_NAME} or {{VAR_NAME}}
	pattern := regexp.MustCompile(`\$\{([^}]+)\}|\{\{([^}]+)\}\}`)

	return pattern.ReplaceAllStringFunc(s, func(match string) string {
		var varName string
		if strings.HasPrefix(match, "${") {
			varName = strings.TrimSuffix(strings.TrimPrefix(match, "${"), "}")
		} else if strings.HasPrefix(match, "{{") {
			varName = strings.TrimSuffix(strings.TrimPrefix(match, "{{"), "}}")
		} else {
			return match
		}

		// Get variable value
		if val, ok := variables[varName]; ok {
			return toString(val)
		}

		// Variable not found, return original match
		return match
	})
}

// toString converts a value to string
func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	default:
		return ""
	}
}
