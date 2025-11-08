package handler

import (
	"encoding/json"
	"goli/database"
	"goli/middlewares"
	"goli/models"
	"goli/pipeline"
	"goli/queue"
	response_util "goli/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CreatePipelineHandler creates a new pipeline
func CreatePipelineHandler(w http.ResponseWriter, r *http.Request) {
	if !middlewares.VerifyAuth(w, r) {
		return
	}

	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Definition  string `json:"definition"` // YAML content
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response_util.SendBadRequestResponse(w, "Invalid request body: "+err.Error())
		return
	}

	// Parse and validate pipeline definition
	pipelineDef, err := pipeline.ParsePipelineDefinition(body.Definition)
	if err != nil {
		response_util.SendBadRequestResponse(w, "Invalid pipeline definition: "+err.Error())
		return
	}

	if err := pipeline.ValidatePipelineDefinition(pipelineDef); err != nil {
		response_util.SendBadRequestResponse(w, "Pipeline validation failed: "+err.Error())
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
		response_util.SendInternalServerErrorResponse(w, "Failed to create pipeline: "+err.Error())
		return
	}

	jsonData, _ := json.Marshal(createdPipeline)
	response_util.SendJsonResponse(w, http.StatusCreated, jsonData)
}

// GetPipelineHandler retrieves a pipeline by ID
func GetPipelineHandler(w http.ResponseWriter, r *http.Request) {
	if !middlewares.VerifyAuth(w, r) {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		response_util.SendBadRequestResponse(w, "Invalid pipeline ID")
		return
	}

	p, err := database.GetPipeline(id)
	if err != nil {
		response_util.SendNotFoundResponse(w, "Pipeline not found")
		return
	}

	jsonData, _ := json.Marshal(p)
	response_util.SendJsonResponse(w, http.StatusOK, jsonData)
}

// ListPipelinesHandler lists all pipelines
func ListPipelinesHandler(w http.ResponseWriter, r *http.Request) {
	if !middlewares.VerifyAuth(w, r) {
		return
	}

	pipelines, err := database.ListPipelines()
	if err != nil {
		response_util.SendInternalServerErrorResponse(w, "Failed to list pipelines: "+err.Error())
		return
	}

	jsonData, _ := json.Marshal(pipelines)
	response_util.SendJsonResponse(w, http.StatusOK, jsonData)
}

// RunPipelineHandler creates a job to run a pipeline
func RunPipelineHandler(w http.ResponseWriter, r *http.Request) {
	if !middlewares.VerifyAuth(w, r) {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		response_util.SendBadRequestResponse(w, "Invalid pipeline ID")
		return
	}

	// Verify pipeline exists
	_, err = database.GetPipeline(id)
	if err != nil {
		response_util.SendNotFoundResponse(w, "Pipeline not found")
		return
	}

	var body struct {
		Name        string `json:"name,omitempty"`
		TriggeredBy string `json:"triggered_by,omitempty"`
	}

	json.NewDecoder(r.Body).Decode(&body)

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
		response_util.SendInternalServerErrorResponse(w, "Failed to enqueue job: "+err.Error())
		return
	}

	jsonData, _ := json.Marshal(job)
	response_util.SendJsonResponse(w, http.StatusCreated, jsonData)
}
