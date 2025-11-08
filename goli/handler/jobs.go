package handler

import (
	"encoding/json"
	"goli/database"
	"goli/middlewares"
	"goli/models"
	"goli/queue"
	response_util "goli/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CreateJobHandler creates a new job
func CreateJobHandler(w http.ResponseWriter, r *http.Request) {
	if !middlewares.VerifyAuth(w, r) {
		return
	}

	var body struct {
		Name        string `json:"name"`
		PipelineID  *int64 `json:"pipeline_id,omitempty"`
		TriggeredBy string `json:"triggered_by,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response_util.SendBadRequestResponse(w, "Invalid request body: "+err.Error())
		return
	}

	job := &models.Job{
		Name:        body.Name,
		PipelineID:  nil,
		Status:      models.JobStatusPending,
		TriggeredBy: body.TriggeredBy,
	}

	if body.PipelineID != nil {
		job.PipelineID = body.PipelineID
	}

	if err := queue.GetQueue().Enqueue(job); err != nil {
		response_util.SendInternalServerErrorResponse(w, "Failed to enqueue job: "+err.Error())
		return
	}

	response_util.SendOkResponse(w, "Job created and enqueued")
}

// GetJobHandler retrieves a job by ID
func GetJobHandler(w http.ResponseWriter, r *http.Request) {
	if !middlewares.VerifyAuth(w, r) {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		response_util.SendBadRequestResponse(w, "Invalid job ID")
		return
	}

	job, err := database.GetJob(id)
	if err != nil {
		response_util.SendNotFoundResponse(w, "Job not found")
		return
	}

	jsonData, _ := json.Marshal(job)
	response_util.SendJsonResponse(w, http.StatusOK, jsonData)
}

// ListJobsHandler lists all jobs
func ListJobsHandler(w http.ResponseWriter, r *http.Request) {
	if !middlewares.VerifyAuth(w, r) {
		return
	}

	limit := 50
	offset := 0
	statusFilter := ""

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}
	if status := r.URL.Query().Get("status"); status != "" {
		statusFilter = status
	}

	jobs, err := database.ListJobs(limit, offset, statusFilter)
	if err != nil {
		response_util.SendInternalServerErrorResponse(w, "Failed to list jobs: "+err.Error())
		return
	}

	jsonData, _ := json.Marshal(jobs)
	response_util.SendJsonResponse(w, http.StatusOK, jsonData)
}
