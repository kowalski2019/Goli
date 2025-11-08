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

// DeletePipeline deletes a pipeline by ID
func DeletePipeline(id int64) error {
	query := `DELETE FROM pipelines WHERE id = ?`
	_, err := DB.Exec(query, id)
	return err
}
