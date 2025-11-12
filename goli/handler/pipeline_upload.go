package handler

import (
	"bytes"
	"goli/database"
	"goli/models"
	"goli/pipeline"
	"goli/queue"
	response_util "goli/utils"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

// UploadPipelineHandler handles YAML file uploads and creates a pipeline
func UploadPipelineHandler(c *gin.Context) {
	// Parse multipart form (max 10MB)
	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		response_util.SendBadRequestResponseGin(c, "Error parsing form: "+err.Error())
		return
	}

	// Get the file
	file, header, err := c.Request.FormFile("yaml_file")
	if err != nil {
		response_util.SendBadRequestResponseGin(c, "Error retrieving file: "+err.Error())
		return
	}
	defer file.Close()

	// Check file extension
	if !strings.HasSuffix(strings.ToLower(header.Filename), ".yaml") &&
		!strings.HasSuffix(strings.ToLower(header.Filename), ".yml") {
		response_util.SendBadRequestResponseGin(c, "File must be a YAML file (.yaml or .yml)")
		return
	}

	// Read file content
	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		response_util.SendInternalServerErrorResponseGin(c, "Error reading file: "+err.Error())
		return
	}

	yamlContent := buf.String()

	// Parse and validate pipeline definition
	pipelineDef, err := pipeline.ParsePipelineDefinition(yamlContent)
	if err != nil {
		response_util.SendBadRequestResponseGin(c, "Invalid pipeline definition: "+err.Error())
		return
	}

	if err := pipeline.ValidatePipelineDefinition(pipelineDef); err != nil {
		response_util.SendBadRequestResponseGin(c, "Pipeline validation failed: "+err.Error())
		return
	}

	// Get name and description from form or use defaults
	name := c.PostForm("name")
	if name == "" {
		name = pipelineDef.Name
		if name == "" {
			name = strings.TrimSuffix(header.Filename, ".yaml")
			name = strings.TrimSuffix(name, ".yml")
		}
	}

	description := c.PostForm("description")
	if description == "" {
		description = pipelineDef.Description
	}

	// Create pipeline in database
	p := &models.Pipeline{
		Name:        name,
		Description: description,
		Definition:  yamlContent,
	}

	createdPipeline, err := database.CreatePipeline(p)
	if err != nil {
		response_util.SendInternalServerErrorResponseGin(c, "Failed to create pipeline: "+err.Error())
		return
	}

	// Optionally run the pipeline immediately if "run" parameter is set
	if c.PostForm("run") == "true" {
		job := &models.Job{
			Name:        name + " - Run",
			PipelineID:  &createdPipeline.ID,
			Status:      models.JobStatusPending,
			TriggeredBy: "UI Upload",
		}

		if err := queue.GetQueue().Enqueue(job); err != nil {
			// Pipeline created but failed to run
			response_util.SendOkResponseGin(c, "Pipeline created but failed to start: "+err.Error())
			return
		}

		response_util.SendJsonResponseGin(c, 201, gin.H{
			"pipeline":    createdPipeline,
			"job_started": true,
			"message":     "Pipeline created and started successfully",
		})
		return
	}

	response_util.SendJsonResponseGin(c, 201, createdPipeline)
}
