package database

import (
	"database/sql"
	"goli/models"
	"time"
)

// CreateJob creates a new job in the database
func CreateJob(job *models.Job) (*models.Job, error) {
	query := `INSERT INTO jobs (pipeline_id, name, status, triggered_by, logs) 
			  VALUES (?, ?, ?, ?, ?) RETURNING id, created_at`

	var createdAt time.Time
	err := DB.QueryRow(query, job.PipelineID, job.Name, job.Status, job.TriggeredBy, job.Logs).Scan(&job.ID, &createdAt)
	if err != nil {
		return nil, err
	}
	job.CreatedAt = createdAt
	return job, nil
}

// UpdateJobStatus updates the status of a job
func UpdateJobStatus(id int64, status models.JobStatus, errorMsg string) error {
	now := time.Now()

	var query string
	var args []interface{}

	if status == models.JobStatusRunning {
		query = `UPDATE jobs SET status = ?, started_at = ?, error_message = ? WHERE id = ?`
		args = []interface{}{status, now, errorMsg, id}
	} else if status == models.JobStatusCompleted || status == models.JobStatusFailed || status == models.JobStatusCancelled {
		query = `UPDATE jobs SET status = ?, completed_at = ?, error_message = ? WHERE id = ?`
		args = []interface{}{status, now, errorMsg, id}
	} else {
		query = `UPDATE jobs SET status = ?, error_message = ? WHERE id = ?`
		args = []interface{}{status, errorMsg, id}
	}

	_, err := DB.Exec(query, args...)
	return err
}

// GetRunningJobs retrieves all jobs with status "running"
func GetRunningJobs() ([]*models.Job, error) {
	query := `SELECT id, pipeline_id, name, status, triggered_by, started_at, 
			  completed_at, error_message, logs, created_at 
			  FROM jobs WHERE status = 'running' ORDER BY started_at DESC`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*models.Job
	for rows.Next() {
		job := &models.Job{}
		var startedAt, completedAt sql.NullTime
		err := rows.Scan(
			&job.ID, &job.PipelineID, &job.Name, &job.Status, &job.TriggeredBy,
			&startedAt, &completedAt, &job.ErrorMessage, &job.Logs, &job.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if startedAt.Valid {
			job.StartedAt = &startedAt.Time
		}
		if completedAt.Valid {
			job.CompletedAt = &completedAt.Time
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

// GetJob retrieves a job by ID
func GetJob(id int64) (*models.Job, error) {
	job := &models.Job{}
	query := `SELECT id, pipeline_id, name, status, triggered_by, started_at, 
			  completed_at, error_message, logs, created_at 
			  FROM jobs WHERE id = ?`

	var startedAt, completedAt sql.NullTime
	err := DB.QueryRow(query, id).Scan(
		&job.ID, &job.PipelineID, &job.Name, &job.Status, &job.TriggeredBy,
		&startedAt, &completedAt, &job.ErrorMessage, &job.Logs, &job.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if startedAt.Valid {
		job.StartedAt = &startedAt.Time
	}
	if completedAt.Valid {
		job.CompletedAt = &completedAt.Time
	}

	// Load steps
	steps, err := GetJobSteps(id)
	if err == nil {
		job.Steps = steps
	}

	return job, nil
}

// GetJobSteps retrieves all steps for a job
func GetJobSteps(jobID int64) ([]models.JobStep, error) {
	query := `SELECT id, job_id, step_name, step_order, status, started_at, 
			  completed_at, error_message, logs, created_at 
			  FROM job_steps WHERE job_id = ? ORDER BY step_order`

	rows, err := DB.Query(query, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var steps []models.JobStep
	for rows.Next() {
		step := models.JobStep{}
		var startedAt, completedAt sql.NullTime
		err := rows.Scan(
			&step.ID, &step.JobID, &step.StepName, &step.StepOrder, &step.Status,
			&startedAt, &completedAt, &step.ErrorMessage, &step.Logs, &step.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if startedAt.Valid {
			step.StartedAt = &startedAt.Time
		}
		if completedAt.Valid {
			step.CompletedAt = &completedAt.Time
		}

		steps = append(steps, step)
	}

	return steps, nil
}

// GetRunningStep retrieves the currently running step for a job (if any)
func GetRunningStep(jobID int64) (*models.JobStep, error) {
	query := `SELECT id, job_id, step_name, step_order, status, started_at, 
			  completed_at, error_message, logs, created_at 
			  FROM job_steps WHERE job_id = ? AND status = 'running' 
			  ORDER BY step_order DESC LIMIT 1`

	step := &models.JobStep{}
	var startedAt, completedAt sql.NullTime
	err := DB.QueryRow(query, jobID).Scan(
		&step.ID, &step.JobID, &step.StepName, &step.StepOrder, &step.Status,
		&startedAt, &completedAt, &step.ErrorMessage, &step.Logs, &step.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No running step
		}
		return nil, err
	}

	if startedAt.Valid {
		step.StartedAt = &startedAt.Time
	}
	if completedAt.Valid {
		step.CompletedAt = &completedAt.Time
	}

	return step, nil
}

// GetJobs retrieves all jobs with optional filters (alias for ListJobs for compatibility)
func GetJobs(limit int, offset int, statusFilter string) ([]*models.Job, error) {
	return ListJobs(limit, offset, statusFilter)
}

// ListJobs retrieves all jobs with optional filters
func ListJobs(limit int, offset int, statusFilter string) ([]*models.Job, error) {
	query := `SELECT id, pipeline_id, name, status, triggered_by, started_at, 
			  completed_at, error_message, created_at 
			  FROM jobs`

	var args []interface{}
	if statusFilter != "" {
		query += " WHERE status = ?"
		args = append(args, statusFilter)
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*models.Job
	for rows.Next() {
		job := &models.Job{}
		var startedAt, completedAt sql.NullTime
		err := rows.Scan(
			&job.ID, &job.PipelineID, &job.Name, &job.Status, &job.TriggeredBy,
			&startedAt, &completedAt, &job.ErrorMessage, &job.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if startedAt.Valid {
			job.StartedAt = &startedAt.Time
		}
		if completedAt.Valid {
			job.CompletedAt = &completedAt.Time
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

// UpdateJobLogs appends logs to a job
func UpdateJobLogs(id int64, logs string) error {
	query := `UPDATE jobs SET logs = COALESCE(logs, '') || ? WHERE id = ?`
	_, err := DB.Exec(query, "\n"+logs, id)
	return err
}

// CreateJobStep creates a new step for a job
func CreateJobStep(step *models.JobStep) error {
	query := `INSERT INTO job_steps (job_id, step_name, step_order, status, logs) 
			  VALUES (?, ?, ?, ?, ?) RETURNING id, created_at`

	var createdAt time.Time
	err := DB.QueryRow(query, step.JobID, step.StepName, step.StepOrder, step.Status, step.Logs).Scan(&step.ID, &createdAt)
	if err != nil {
		return err
	}
	step.CreatedAt = createdAt
	return nil
}

// UpdateJobStepStatus updates the status of a job step
func UpdateJobStepStatus(stepID int64, status models.JobStatus, errorMsg string) error {
	now := time.Now()

	var query string
	var args []interface{}

	if status == models.JobStatusRunning {
		query = `UPDATE job_steps SET status = ?, started_at = ?, error_message = ? WHERE id = ?`
		args = []interface{}{status, now, errorMsg, stepID}
	} else if status == models.JobStatusCompleted || status == models.JobStatusFailed || status == models.JobStatusCancelled {
		query = `UPDATE job_steps SET status = ?, completed_at = ?, error_message = ? WHERE id = ?`
		args = []interface{}{status, now, errorMsg, stepID}
	} else {
		query = `UPDATE job_steps SET status = ?, error_message = ? WHERE id = ?`
		args = []interface{}{status, errorMsg, stepID}
	}

	_, err := DB.Exec(query, args...)
	return err
}

// UpdateJobStepLogs appends logs to a job step
func UpdateJobStepLogs(stepID int64, logs string) error {
	query := `UPDATE job_steps SET logs = COALESCE(logs, '') || ? WHERE id = ?`
	_, err := DB.Exec(query, "\n"+logs, stepID)
	return err
}
