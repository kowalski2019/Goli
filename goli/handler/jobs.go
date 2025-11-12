package handler

import (
	"goli/database"
	"goli/models"
	"goli/queue"
	response_util "goli/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateJobHandler creates a new job
func CreateJobHandler(c *gin.Context) {
	var body struct {
		Name        string `json:"name"`
		PipelineID  *int64 `json:"pipeline_id,omitempty"`
		TriggeredBy string `json:"triggered_by,omitempty"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		response_util.SendBadRequestResponseGin(c, "Invalid request body: "+err.Error())
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
		response_util.SendInternalServerErrorResponseGin(c, "Failed to enqueue job: "+err.Error())
		return
	}

	response_util.SendOkResponseGin(c, "Job created and enqueued")
}

// GetJobHandler retrieves a job by ID
func GetJobHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response_util.SendBadRequestResponseGin(c, "Invalid job ID")
		return
	}

	job, err := database.GetJob(id)
	if err != nil {
		response_util.SendNotFoundResponseGin(c, "Job not found")
		return
	}

	response_util.SendJsonResponseGin(c, 200, job)
}

// ListJobsHandler lists all jobs
func ListJobsHandler(c *gin.Context) {
	limit := 50
	offset := 0
	statusFilter := ""

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}
	if status := c.Query("status"); status != "" {
		statusFilter = status
	}

	jobs, err := database.GetJobs(limit, offset, statusFilter)
	if err != nil {
		response_util.SendInternalServerErrorResponseGin(c, "Failed to list jobs: "+err.Error())
		return
	}

	response_util.SendJsonResponseGin(c, 200, jobs)
}

// CancelJobHandler cancels a running or pending job
func CancelJobHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response_util.SendBadRequestResponseGin(c, "Invalid job ID")
		return
	}

	if err := queue.GetQueue().CancelJob(id); err != nil {
		response_util.SendBadRequestResponseGin(c, "Failed to cancel job: "+err.Error())
		return
	}

	response_util.SendOkResponseGin(c, "Job cancelled successfully")
}
