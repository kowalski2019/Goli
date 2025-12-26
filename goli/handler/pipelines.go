package handler

import (
	"goli/database"
	"goli/models"
	"goli/pipeline"
	"goli/queue"
	response_util "goli/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreatePipelineHandler creates a new pipeline
func CreatePipelineHandler(c *gin.Context) {
	var body struct {
		Name        string                 `json:"name"`
		Description string                 `json:"description"`
		Definition  string                 `json:"definition"`          // YAML content
		Variables   map[string]interface{} `json:"variables,omitempty"` // Map of variable name to {value, is_secret}
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		response_util.SendBadRequestResponseGin(c, "Invalid request body: "+err.Error())
		return
	}

	// Parse and validate pipeline definition
	pipelineDef, err := pipeline.ParsePipelineDefinition(body.Definition)
	if err != nil {
		response_util.SendBadRequestResponseGin(c, "Invalid pipeline definition: "+err.Error())
		return
	}

	if err := pipeline.ValidatePipelineDefinition(pipelineDef); err != nil {
		response_util.SendBadRequestResponseGin(c, "Pipeline validation failed: "+err.Error())
		return
	}

	// Create pipeline in database
	p := &models.Pipeline{
		Name:        body.Name,
		Description: body.Description,
		Definition:  body.Definition,
	}

	createdPipeline, err := database.CreatePipeline(p)
	if err != nil {
		response_util.SendInternalServerErrorResponseGin(c, "Failed to create pipeline: "+err.Error())
		return
	}

	// Add variables if provided
	if body.Variables != nil {
		for name, varData := range body.Variables {
			var value string
			var isSecret bool

			// Handle different formats: string value or object with value/is_secret
			if varStr, ok := varData.(string); ok {
				value = varStr
				isSecret = false
			} else if varMap, ok := varData.(map[string]interface{}); ok {
				if val, ok := varMap["value"].(string); ok {
					value = val
				}
				if secret, ok := varMap["is_secret"].(bool); ok {
					isSecret = secret
				}
			}

			if value != "" {
				if err := database.SetPipelineVariable(createdPipeline.ID, name, value, isSecret); err != nil {
					response_util.SendInternalServerErrorResponseGin(c, "Failed to set variable: "+err.Error())
					return
				}
			}
		}
	}

	// Return created pipeline with variables
	finalPipeline, err := database.GetPipeline(createdPipeline.ID)
	if err != nil {
		response_util.SendInternalServerErrorResponseGin(c, "Failed to retrieve created pipeline: "+err.Error())
		return
	}

	response_util.SendJsonResponseGin(c, 201, finalPipeline)
}

// GetPipelineHandler retrieves a pipeline by ID
func GetPipelineHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response_util.SendBadRequestResponseGin(c, "Invalid pipeline ID")
		return
	}

	p, err := database.GetPipeline(id)
	if err != nil {
		response_util.SendNotFoundResponseGin(c, "Pipeline not found")
		return
	}

	response_util.SendJsonResponseGin(c, 200, p)
}

// ListPipelinesHandler lists all pipelines
func ListPipelinesHandler(c *gin.Context) {
	pipelines, err := database.ListPipelines()
	if err != nil {
		response_util.SendInternalServerErrorResponseGin(c, "Failed to list pipelines: "+err.Error())
		return
	}

	response_util.SendJsonResponseGin(c, 200, pipelines)
}

// RunPipelineHandler creates a job to run a pipeline
func RunPipelineHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response_util.SendBadRequestResponseGin(c, "Invalid pipeline ID")
		return
	}

	// Verify pipeline exists
	_, err = database.GetPipeline(id)
	if err != nil {
		response_util.SendNotFoundResponseGin(c, "Pipeline not found")
		return
	}

	var body struct {
		Name        string `json:"name,omitempty"`
		TriggeredBy string `json:"triggered_by,omitempty"`
	}

	c.ShouldBindJSON(&body)

	jobName := body.Name
	if jobName == "" {
		jobName = "Pipeline Run"
	}

	job := &models.Job{
		Name:        jobName,
		PipelineID:  &id,
		Status:      models.JobStatusPending,
		TriggeredBy: body.TriggeredBy,
	}

	if err := queue.GetQueue().Enqueue(job); err != nil {
		response_util.SendInternalServerErrorResponseGin(c, "Failed to enqueue job: "+err.Error())
		return
	}

	response_util.SendJsonResponseGin(c, 201, job)
}

// UpdatePipelineHandler updates an existing pipeline
func UpdatePipelineHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response_util.SendBadRequestResponseGin(c, "Invalid pipeline ID")
		return
	}

	// Verify pipeline exists
	_, err = database.GetPipeline(id)
	if err != nil {
		response_util.SendNotFoundResponseGin(c, "Pipeline not found")
		return
	}

	var body struct {
		Name        string                 `json:"name"`
		Description string                 `json:"description"`
		Definition  string                 `json:"definition"`          // YAML content
		Variables   map[string]interface{} `json:"variables,omitempty"` // Map of variable name to {value, is_secret}
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		response_util.SendBadRequestResponseGin(c, "Invalid request body: "+err.Error())
		return
	}

	// If definition is provided, parse and validate it
	if body.Definition != "" {
		pipelineDef, err := pipeline.ParsePipelineDefinition(body.Definition)
		if err != nil {
			response_util.SendBadRequestResponseGin(c, "Invalid pipeline definition: "+err.Error())
			return
		}

		if err := pipeline.ValidatePipelineDefinition(pipelineDef); err != nil {
			response_util.SendBadRequestResponseGin(c, "Pipeline validation failed: "+err.Error())
			return
		}
	}

	// Update pipeline
	p := &models.Pipeline{
		ID:          id,
		Name:        body.Name,
		Description: body.Description,
		Definition:  body.Definition,
	}

	// Only update fields that are provided
	if body.Name == "" && body.Description == "" && body.Definition == "" {
		// If only variables are being updated, get the current pipeline
		currentPipeline, err := database.GetPipeline(id)
		if err != nil {
			response_util.SendInternalServerErrorResponseGin(c, "Failed to load pipeline: "+err.Error())
			return
		}
		p.Name = currentPipeline.Name
		p.Description = currentPipeline.Description
		p.Definition = currentPipeline.Definition
	} else {
		// Get current pipeline to preserve fields not being updated
		currentPipeline, err := database.GetPipeline(id)
		if err != nil {
			response_util.SendInternalServerErrorResponseGin(c, "Failed to load pipeline: "+err.Error())
			return
		}
		if body.Name == "" {
			p.Name = currentPipeline.Name
		}
		if body.Description == "" {
			p.Description = currentPipeline.Description
		}
		if body.Definition == "" {
			p.Definition = currentPipeline.Definition
		}
	}

	if err := database.UpdatePipeline(p); err != nil {
		response_util.SendInternalServerErrorResponseGin(c, "Failed to update pipeline: "+err.Error())
		return
	}

	// Update variables if provided
	if body.Variables != nil {
		// Get existing variables to preserve masked secrets
		existingVars, err := database.GetPipelineVariables(id)
		if err != nil {
			response_util.SendInternalServerErrorResponseGin(c, "Failed to load existing variables: "+err.Error())
			return
		}

		// Create a map of existing variable names
		existingVarMap := make(map[string]*database.PipelineVariable)
		for _, v := range existingVars {
			existingVarMap[v.Name] = v
		}

		// Track which variables are in the update
		updatedVarNames := make(map[string]bool)

		// Update or add variables from the request
		for name, varData := range body.Variables {
			updatedVarNames[name] = true
			var value string
			var isSecret bool

			// Handle different formats: string value or object with value/is_secret
			if varStr, ok := varData.(string); ok {
				value = varStr
				isSecret = false
			} else if varMap, ok := varData.(map[string]interface{}); ok {
				if val, ok := varMap["value"].(string); ok {
					value = val
				}
				if secret, ok := varMap["is_secret"].(bool); ok {
					isSecret = secret
				}
			}

			// If value is masked, preserve existing value
			if value == "***MASKED***" {
				if existingVar, exists := existingVarMap[name]; exists {
					// Keep existing value and secret status
					value = existingVar.Value
					isSecret = existingVar.IsSecret
				} else {
					// Variable doesn't exist, skip it
					continue
				}
			}

			// Update or create the variable
			if value != "" {
				if err := database.SetPipelineVariable(id, name, value, isSecret); err != nil {
					response_util.SendInternalServerErrorResponseGin(c, "Failed to set variable: "+err.Error())
					return
				}
			}
		}

		// Delete variables that are not in the update
		for _, existingVar := range existingVars {
			if !updatedVarNames[existingVar.Name] {
				if err := database.DeletePipelineVariable(id, existingVar.Name); err != nil {
					response_util.SendInternalServerErrorResponseGin(c, "Failed to delete variable: "+err.Error())
					return
				}
			}
		}
	}

	// Return updated pipeline
	updatedPipeline, err := database.GetPipeline(id)
	if err != nil {
		response_util.SendInternalServerErrorResponseGin(c, "Failed to retrieve updated pipeline: "+err.Error())
		return
	}

	response_util.SendJsonResponseGin(c, 200, updatedPipeline)
}

// DeletePipelineHandler deletes a pipeline and all related jobs
func DeletePipelineHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response_util.SendBadRequestResponseGin(c, "Invalid pipeline ID")
		return
	}

	// Verify pipeline exists before deletion
	_, err = database.GetPipeline(id)
	if err != nil {
		response_util.SendNotFoundResponseGin(c, "Pipeline not found")
		return
	}

	// Delete pipeline (this will cascade delete related jobs and job steps)
	if err := database.DeletePipeline(id); err != nil {
		response_util.SendInternalServerErrorResponseGin(c, "Failed to delete pipeline: "+err.Error())
		return
	}

	response_util.SendOkResponseGin(c, "Pipeline and all related jobs deleted successfully")
}
