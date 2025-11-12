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
		Name        string `json:"name"`
		Description string `json:"description"`
		Definition  string `json:"definition"` // YAML content
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

	response_util.SendJsonResponseGin(c, 201, createdPipeline)
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

