package handler

import (
	"bytes"
	"encoding/json"
	"goli/database"
	"goli/middlewares"
	"goli/models"
	"goli/pipeline"
	"goli/queue"
	response_util "goli/utils"
	"io"
	"net/http"
	"strings"
)

// UploadPipelineHandler handles YAML file uploads and creates a pipeline
func UploadPipelineHandler(w http.ResponseWriter, r *http.Request) {
	if !middlewares.VerifyAuth(w, r) {
		return
	}

	// Parse multipart form (max 10MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		response_util.SendBadRequestResponse(w, "Error parsing form: "+err.Error())
		return
	}

	// Get the file
	file, header, err := r.FormFile("yaml_file")
	if err != nil {
		response_util.SendBadRequestResponse(w, "Error retrieving file: "+err.Error())
		return
	}
	defer file.Close()

	// Check file extension
	if !strings.HasSuffix(strings.ToLower(header.Filename), ".yaml") &&
		!strings.HasSuffix(strings.ToLower(header.Filename), ".yml") {
		response_util.SendBadRequestResponse(w, "File must be a YAML file (.yaml or .yml)")
		return
	}

	// Read file content
	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		response_util.SendInternalServerErrorResponse(w, "Error reading file: "+err.Error())
		return
	}

	yamlContent := buf.String()

	// Parse and validate pipeline definition
	pipelineDef, err := pipeline.ParsePipelineDefinition(yamlContent)
	if err != nil {
		response_util.SendBadRequestResponse(w, "Invalid pipeline definition: "+err.Error())
		return
	}

	if err := pipeline.ValidatePipelineDefinition(pipelineDef); err != nil {
		response_util.SendBadRequestResponse(w, "Pipeline validation failed: "+err.Error())
		return
	}

	// Get name and description from form or use defaults
	name := r.FormValue("name")
	if name == "" {
		name = pipelineDef.Name
		if name == "" {
			name = strings.TrimSuffix(header.Filename, ".yaml")
			name = strings.TrimSuffix(name, ".yml")
		}
	}

	description := r.FormValue("description")
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
		response_util.SendInternalServerErrorResponse(w, "Failed to create pipeline: "+err.Error())
		return
	}

	// Optionally run the pipeline immediately if "run" parameter is set
	if r.FormValue("run") == "true" {
		job := &models.Job{
			Name:        name + " - Run",
			PipelineID:  &createdPipeline.ID,
			Status:      models.JobStatusPending,
			TriggeredBy: "UI Upload",
		}

		if err := queue.GetQueue().Enqueue(job); err != nil {
			// Pipeline created but failed to run
			response_util.SendOkResponse(w, "Pipeline created but failed to start: "+err.Error())
			return
		}

		jsonData, _ := json.Marshal(map[string]interface{}{
			"pipeline":    createdPipeline,
			"job_started": true,
			"message":     "Pipeline created and started successfully",
		})
		response_util.SendJsonResponse(w, http.StatusCreated, jsonData)
		return
	}

	jsonData, _ := json.Marshal(createdPipeline)
	response_util.SendJsonResponse(w, http.StatusCreated, jsonData)
}
