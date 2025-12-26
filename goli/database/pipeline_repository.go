package database

import (
	"database/sql"
	"goli/models"
)

// CreatePipeline creates a new pipeline in the database
func CreatePipeline(pipeline *models.Pipeline) (*models.Pipeline, error) {
	query := `INSERT INTO pipelines (name, description, definition) 
			  VALUES (?, ?, ?) RETURNING id, created_at, updated_at`

	err := DB.QueryRow(query, pipeline.Name, pipeline.Description, pipeline.Definition).Scan(
		&pipeline.ID, &pipeline.CreatedAt, &pipeline.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// GetPipeline retrieves a pipeline by ID
func GetPipeline(id int64) (*models.Pipeline, error) {
	pipeline := &models.Pipeline{}
	query := `SELECT id, name, description, definition, created_at, updated_at 
			  FROM pipelines WHERE id = ?`

	err := DB.QueryRow(query, id).Scan(
		&pipeline.ID, &pipeline.Name, &pipeline.Description, &pipeline.Definition,
		&pipeline.CreatedAt, &pipeline.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Load variables (secrets are masked)
	variables, err := GetPipelineVariables(id)
	if err == nil && len(variables) > 0 {
		pipeline.Variables = make(map[string]interface{})
		for _, v := range variables {
			if v.IsSecret {
				pipeline.Variables[v.Name] = "***MASKED***"
			} else {
				pipeline.Variables[v.Name] = v.Value
			}
		}
	}

	return pipeline, nil
}

// GetPipelineWithSecrets retrieves a pipeline by ID including secret values (for execution)
func GetPipelineWithSecrets(id int64) (*models.Pipeline, error) {
	pipeline := &models.Pipeline{}
	query := `SELECT id, name, description, definition, created_at, updated_at 
			  FROM pipelines WHERE id = ?`

	err := DB.QueryRow(query, id).Scan(
		&pipeline.ID, &pipeline.Name, &pipeline.Description, &pipeline.Definition,
		&pipeline.CreatedAt, &pipeline.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Load variables with actual secret values (for execution)
	variables, err := GetPipelineVariables(id)
	if err == nil && len(variables) > 0 {
		pipeline.Variables = make(map[string]interface{})
		for _, v := range variables {
			pipeline.Variables[v.Name] = v.Value
		}
	}

	return pipeline, nil
}

// GetPipelineByName retrieves a pipeline by name
func GetPipelineByName(name string) (*models.Pipeline, error) {
	pipeline := &models.Pipeline{}
	query := `SELECT id, name, description, definition, created_at, updated_at 
			  FROM pipelines WHERE name = ?`

	err := DB.QueryRow(query, name).Scan(
		&pipeline.ID, &pipeline.Name, &pipeline.Description, &pipeline.Definition,
		&pipeline.CreatedAt, &pipeline.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// ListPipelines retrieves all pipelines
func ListPipelines() ([]*models.Pipeline, error) {
	query := `SELECT id, name, description, definition, created_at, updated_at 
			  FROM pipelines ORDER BY created_at DESC`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pipelines []*models.Pipeline
	for rows.Next() {
		pipeline := &models.Pipeline{}
		err := rows.Scan(
			&pipeline.ID, &pipeline.Name, &pipeline.Description, &pipeline.Definition,
			&pipeline.CreatedAt, &pipeline.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Load variables (secrets are masked)
		variables, err := GetPipelineVariables(pipeline.ID)
		if err == nil && len(variables) > 0 {
			pipeline.Variables = make(map[string]interface{})
			for _, v := range variables {
				if v.IsSecret {
					pipeline.Variables[v.Name] = "***MASKED***"
				} else {
					pipeline.Variables[v.Name] = v.Value
				}
			}
		}

		pipelines = append(pipelines, pipeline)
	}

	return pipelines, nil
}

// UpdatePipeline updates an existing pipeline
func UpdatePipeline(pipeline *models.Pipeline) error {
	query := `UPDATE pipelines SET name = ?, description = ?, definition = ?, updated_at = CURRENT_TIMESTAMP 
			  WHERE id = ?`

	result, err := DB.Exec(query, pipeline.Name, pipeline.Description, pipeline.Definition, pipeline.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// DeletePipeline deletes a pipeline by ID and all related jobs and job steps (cascade delete)
func DeletePipeline(id int64) error {
	// Start a transaction to ensure atomicity
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// First, get all job IDs for this pipeline
	jobRows, err := tx.Query(`SELECT id FROM jobs WHERE pipeline_id = ?`, id)
	if err != nil {
		return err
	}
	defer jobRows.Close()

	var jobIDs []int64
	for jobRows.Next() {
		var jobID int64
		if err := jobRows.Scan(&jobID); err != nil {
			return err
		}
		jobIDs = append(jobIDs, jobID)
	}
	jobRows.Close()

	// Delete all job steps for these jobs
	if len(jobIDs) > 0 {
		// Build placeholders for IN clause
		placeholders := ""
		args := make([]interface{}, len(jobIDs))
		for i, jobID := range jobIDs {
			if i > 0 {
				placeholders += ","
			}
			placeholders += "?"
			args[i] = jobID
		}

		_, err = tx.Exec(`DELETE FROM job_steps WHERE job_id IN (`+placeholders+`)`, args...)
		if err != nil {
			return err
		}
	}

	// Delete all jobs for this pipeline
	_, err = tx.Exec(`DELETE FROM jobs WHERE pipeline_id = ?`, id)
	if err != nil {
		return err
	}

	// Finally, delete the pipeline itself
	_, err = tx.Exec(`DELETE FROM pipelines WHERE id = ?`, id)
	if err != nil {
		return err
	}

	// Commit the transaction
	return tx.Commit()
}

// PipelineVariable represents a pipeline variable or secret
type PipelineVariable struct {
	ID         int64  `json:"id"`
	PipelineID int64  `json:"pipeline_id"`
	Name       string `json:"name"`
	Value      string `json:"value"`
	IsSecret   bool   `json:"is_secret"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// GetPipelineVariables retrieves all variables for a pipeline
func GetPipelineVariables(pipelineID int64) ([]*PipelineVariable, error) {
	query := `SELECT id, pipeline_id, name, value, is_secret, created_at, updated_at 
			  FROM pipeline_variables WHERE pipeline_id = ? ORDER BY name`

	rows, err := DB.Query(query, pipelineID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var variables []*PipelineVariable
	for rows.Next() {
		var v PipelineVariable
		var isSecret int
		err := rows.Scan(
			&v.ID, &v.PipelineID, &v.Name, &v.Value, &isSecret,
			&v.CreatedAt, &v.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		v.IsSecret = isSecret == 1
		variables = append(variables, &v)
	}

	return variables, nil
}

// SetPipelineVariable sets or updates a pipeline variable
func SetPipelineVariable(pipelineID int64, name string, value string, isSecret bool) error {
	isSecretInt := 0
	if isSecret {
		isSecretInt = 1
	}

	query := `INSERT INTO pipeline_variables (pipeline_id, name, value, is_secret, updated_at)
			  VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
			  ON CONFLICT(pipeline_id, name) DO UPDATE SET
			  value = excluded.value,
			  is_secret = excluded.is_secret,
			  updated_at = CURRENT_TIMESTAMP`

	_, err := DB.Exec(query, pipelineID, name, value, isSecretInt)
	return err
}

// DeletePipelineVariable deletes a pipeline variable
func DeletePipelineVariable(pipelineID int64, name string) error {
	query := `DELETE FROM pipeline_variables WHERE pipeline_id = ? AND name = ?`
	_, err := DB.Exec(query, pipelineID, name)
	return err
}

// DeletePipelineVariables deletes all variables for a pipeline
func DeletePipelineVariables(pipelineID int64) error {
	query := `DELETE FROM pipeline_variables WHERE pipeline_id = ?`
	_, err := DB.Exec(query, pipelineID)
	return err
}
