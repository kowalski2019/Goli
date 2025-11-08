package pipeline

import (
	"fmt"
	"goli/database"
	"goli/models"
	"log"
	"os/exec"
	"strings"
	"time"
)

// logToJob appends a log message to the job's logs in the database
func logToJob(jobID int64, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] %s\n", timestamp, message)
	if err := database.UpdateJobLogs(jobID, logEntry); err != nil {
		log.Printf("Error saving job log: %v", err)
	}
	log.Printf("[Job %d] %s", jobID, message)
}

// logToStep appends a log message to the step's logs in the database
func logToStep(stepID int64, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] %s\n", timestamp, message)
	if err := database.UpdateJobStepLogs(stepID, logEntry); err != nil {
		log.Printf("Error saving step log: %v", err)
	}
	log.Printf("[Step %d] %s", stepID, message)
}

// ExecutePipeline executes a pipeline definition for a job
func ExecutePipeline(job *models.Job, pipelineDef *models.PipelineDefinition) error {
	logToJob(job.ID, fmt.Sprintf("Starting pipeline execution: %s", pipelineDef.Name))
	if pipelineDef.Description != "" {
		logToJob(job.ID, fmt.Sprintf("Description: %s", pipelineDef.Description))
	}
	logToJob(job.ID, fmt.Sprintf("Total steps: %d", len(pipelineDef.Steps)))

	// Create job steps from pipeline definition
	for i, stepDef := range pipelineDef.Steps {
		step := &models.JobStep{
			JobID:     job.ID,
			StepName:  stepDef.Name,
			StepOrder: i + 1,
			Status:    models.JobStatusPending,
		}

		if err := database.CreateJobStep(step); err != nil {
			logToJob(job.ID, fmt.Sprintf("ERROR: Failed to create job step '%s': %v", stepDef.Name, err))
			continue
		}

		logToJob(job.ID, fmt.Sprintf("Created step %d/%d: %s (type: %s, action: %s)", i+1, len(pipelineDef.Steps), stepDef.Name, stepDef.Type, stepDef.Action))

		// Execute the step
		if err := executeStep(step, stepDef, job); err != nil {
			logToJob(job.ID, fmt.Sprintf("ERROR: Step '%s' failed: %v", stepDef.Name, err))
			if stepDef.OnFailure == "stop" {
				logToJob(job.ID, "Pipeline execution stopped due to step failure (on_failure: stop)")
				database.UpdateJobStatus(job.ID, models.JobStatusFailed, err.Error())
				return err
			}
			logToJob(job.ID, "Continuing to next step (on_failure: continue)")
			// Continue to next step if on_failure is "continue"
		} else {
			logToJob(job.ID, fmt.Sprintf("Step '%s' completed successfully", stepDef.Name))
		}
	}

	// All steps completed successfully
	logToJob(job.ID, "All pipeline steps completed successfully")
	if err := database.UpdateJobStatus(job.ID, models.JobStatusCompleted, ""); err != nil {
		logToJob(job.ID, fmt.Sprintf("ERROR: Failed to update job status: %v", err))
		return err
	}
	completedAt := time.Now()
	job.CompletedAt = &completedAt
	job.Status = models.JobStatusCompleted
	return nil
}

// executeStep executes a single pipeline step
func executeStep(step *models.JobStep, stepDef models.PipelineStep, job *models.Job) error {
	logToStep(step.ID, fmt.Sprintf("Starting step execution: %s", stepDef.Name))
	if stepDef.Description != "" {
		logToStep(step.ID, fmt.Sprintf("Description: %s", stepDef.Description))
	}
	logToStep(step.ID, fmt.Sprintf("Type: %s, Action: %s", stepDef.Type, stepDef.Action))

	// Update step status to running
	database.UpdateJobStepStatus(step.ID, models.JobStatusRunning, "")

	var err error
	maxRetries := stepDef.Retry
	if maxRetries == 0 {
		maxRetries = 1
	}
	logToStep(step.ID, fmt.Sprintf("Max retries: %d", maxRetries))

	for attempt := 1; attempt <= maxRetries; attempt++ {
		if attempt > 1 {
			logToStep(step.ID, fmt.Sprintf("Retrying step (attempt %d/%d)", attempt, maxRetries))
			time.Sleep(time.Duration(attempt) * time.Second)
		}

		switch stepDef.Type {
		case "docker":
			err = executeDockerStep(step, stepDef, job)
		case "script":
			err = executeScriptStep(step, stepDef, job)
		case "shell":
			err = executeShellStep(step, stepDef, job)
		default:
			logToStep(step.ID, fmt.Sprintf("WARNING: Unknown step type '%s', defaulting to docker", stepDef.Type))
			err = executeDockerStep(step, stepDef, job) // Default to docker
		}

		if err == nil {
			// Step succeeded
			logToStep(step.ID, "Step completed successfully")
			database.UpdateJobStepStatus(step.ID, models.JobStatusCompleted, "")
			return nil
		}

		logToStep(step.ID, fmt.Sprintf("Step failed (attempt %d/%d): %v", attempt, maxRetries, err))
	}

	// All retries failed
	logToStep(step.ID, fmt.Sprintf("All retry attempts exhausted. Step failed with error: %v", err))
	database.UpdateJobStepStatus(step.ID, models.JobStatusFailed, err.Error())
	return err
}

// executeDockerStep executes a Docker-related step
func executeDockerStep(step *models.JobStep, stepDef models.PipelineStep, job *models.Job) error {
	action := stepDef.Action
	config := stepDef.Config

	logToStep(step.ID, fmt.Sprintf("Executing Docker action: %s", action))
	logToStep(step.ID, fmt.Sprintf("Configuration: %v", config))

	// Map pipeline step to Docker operations
	switch action {
	case "pull":
		image, ok := config["image"].(string)
		if !ok {
			logToStep(step.ID, "ERROR: Missing or invalid 'image' configuration")
			return ErrInvalidConfig
		}
		return executeDockerPull(image, step)
	case "run":
		return executeDockerRun(config, step)
	case "start":
		container, ok := config["container"].(string)
		if !ok {
			logToStep(step.ID, "ERROR: Missing or invalid 'container' configuration")
			return ErrInvalidConfig
		}
		return executeDockerStart(container, step)
	case "stop":
		container, ok := config["container"].(string)
		if !ok {
			logToStep(step.ID, "ERROR: Missing or invalid 'container' configuration")
			return ErrInvalidConfig
		}
		return executeDockerStop(container, step)
	default:
		logToStep(step.ID, fmt.Sprintf("ERROR: Unsupported Docker action: %s", action))
		return ErrUnsupportedAction
	}
}

// executeScriptStep executes a script step
func executeScriptStep(step *models.JobStep, stepDef models.PipelineStep, job *models.Job) error {
	config := stepDef.Config
	script, ok := config["script"].(string)
	if !ok {
		logToStep(step.ID, "ERROR: Missing or invalid 'script' configuration")
		return ErrInvalidConfig
	}

	logToStep(step.ID, "Executing script step")
	logToStep(step.ID, fmt.Sprintf("Script length: %d characters", len(script)))

	// Log script content (truncated if too long)
	if len(script) > 1000 {
		logToStep(step.ID, fmt.Sprintf("Script preview (first 1000 chars):\n%s...", script[:1000]))
	} else {
		logToStep(step.ID, fmt.Sprintf("Script content:\n%s", script))
	}

	// Execute the script
	// For now, we'll execute it as a shell command
	// In production, you might want to handle different script types (bash, python, etc.)
	shell := "sh"
	if shellType, ok := config["shell"].(string); ok {
		shell = shellType
	}

	logToStep(step.ID, fmt.Sprintf("Executing script using: %s", shell))

	cmd := exec.Command(shell, "-c", script)
	output, err := cmd.CombinedOutput()

	if len(output) > 0 {
		logToStep(step.ID, fmt.Sprintf("Script output:\n%s", string(output)))
	}

	if err != nil {
		logToStep(step.ID, fmt.Sprintf("Script execution failed: %v", err))
		return fmt.Errorf("script execution failed: %w", err)
	}

	logToStep(step.ID, "Script executed successfully")
	return nil
}

// executeShellStep executes a shell step
func executeShellStep(step *models.JobStep, stepDef models.PipelineStep, job *models.Job) error {
	config := stepDef.Config
	command, ok := config["command"].(string)
	if !ok {
		logToStep(step.ID, "ERROR: Missing or invalid 'command' configuration")
		return ErrInvalidConfig
	}

	var args []string
	if argsInterface, ok := config["args"]; ok {
		if argsList, ok := argsInterface.([]interface{}); ok {
			for _, arg := range argsList {
				if argStr, ok := arg.(string); ok {
					args = append(args, argStr)
				}
			}
		} else if argsList, ok := argsInterface.([]string); ok {
			args = argsList
		}
	}

	logToStep(step.ID, fmt.Sprintf("Executing shell command: %s", command))
	if len(args) > 0 {
		logToStep(step.ID, fmt.Sprintf("Command arguments: %v", args))
	}

	output, err := executeShellCommand(command, args...)

	if len(output) > 0 {
		logToStep(step.ID, fmt.Sprintf("Command output:\n%s", output))
	}

	if err != nil {
		logToStep(step.ID, fmt.Sprintf("Command execution failed: %v", err))
		return fmt.Errorf("shell command failed: %w", err)
	}

	logToStep(step.ID, "Command executed successfully")
	return nil
}

// Helper functions for Docker operations
func executeDockerPull(image string, step *models.JobStep) error {
	logToStep(step.ID, fmt.Sprintf("Pulling Docker image: %s", image))

	cmd := exec.Command("docker", "pull", image)
	output, err := cmd.CombinedOutput()

	if len(output) > 0 {
		logToStep(step.ID, fmt.Sprintf("Docker pull output:\n%s", string(output)))
	}

	if err != nil {
		logToStep(step.ID, fmt.Sprintf("Docker pull failed: %v", err))
		return fmt.Errorf("docker pull failed: %w", err)
	}

	logToStep(step.ID, "Docker image pulled successfully")
	return nil
}

func executeDockerRun(config map[string]interface{}, step *models.JobStep) error {
	logToStep(step.ID, "Running Docker container")

	image, ok := config["image"].(string)
	if !ok {
		logToStep(step.ID, "ERROR: Missing or invalid 'image' configuration")
		return ErrInvalidConfig
	}

	// Build docker run command
	args := []string{"run", "--rm", "--detach"}

	// Add container name if specified
	if container, ok := config["container"].(string); ok {
		args = append(args, "--name", container)
	}

	// Add environment variables
	if env, ok := config["env"].(map[string]interface{}); ok {
		for key, value := range env {
			args = append(args, "-e", fmt.Sprintf("%s=%v", key, value))
		}
	} else if envList, ok := config["env"].([]interface{}); ok {
		for _, envItem := range envList {
			if envStr, ok := envItem.(string); ok {
				args = append(args, "-e", envStr)
			}
		}
	}

	// Add volumes
	if volumes, ok := config["volumes"].([]interface{}); ok {
		for _, volume := range volumes {
			if volStr, ok := volume.(string); ok {
				args = append(args, "-v", volStr)
			}
		}
	}

	// Add ports
	if ports, ok := config["ports"].([]interface{}); ok {
		for _, port := range ports {
			if portStr, ok := port.(string); ok {
				args = append(args, "-p", portStr)
			}
		}
	}

	// Add command if specified
	if cmd, ok := config["cmd"].(string); ok {
		args = append(args, image, cmd)
	} else if cmdList, ok := config["cmd"].([]interface{}); ok {
		args = append(args, image)
		for _, cmdArg := range cmdList {
			if cmdStr, ok := cmdArg.(string); ok {
				args = append(args, cmdStr)
			}
		}
	} else {
		args = append(args, image)
	}

	logToStep(step.ID, fmt.Sprintf("Executing: docker %s", strings.Join(args, " ")))

	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()

	if len(output) > 0 {
		logToStep(step.ID, fmt.Sprintf("Docker run output:\n%s", string(output)))
	}

	if err != nil {
		logToStep(step.ID, fmt.Sprintf("Docker run failed: %v", err))
		return fmt.Errorf("docker run failed: %w", err)
	}

	logToStep(step.ID, "Docker container ran successfully")
	return nil
}

func executeDockerStart(container string, step *models.JobStep) error {
	logToStep(step.ID, fmt.Sprintf("Starting Docker container: %s", container))

	cmd := exec.Command("docker", "start", container)
	output, err := cmd.CombinedOutput()

	if len(output) > 0 {
		logToStep(step.ID, fmt.Sprintf("Docker start output:\n%s", string(output)))
	}

	if err != nil {
		logToStep(step.ID, fmt.Sprintf("Docker start failed: %v", err))
		return fmt.Errorf("docker start failed: %w", err)
	}

	logToStep(step.ID, "Docker container started successfully")
	return nil
}

func executeDockerStop(container string, step *models.JobStep) error {
	logToStep(step.ID, fmt.Sprintf("Stopping Docker container: %s", container))

	cmd := exec.Command("docker", "stop", container)
	output, err := cmd.CombinedOutput()

	if len(output) > 0 {
		logToStep(step.ID, fmt.Sprintf("Docker stop output:\n%s", string(output)))
	}

	if err != nil {
		logToStep(step.ID, fmt.Sprintf("Docker stop failed: %v", err))
		return fmt.Errorf("docker stop failed: %w", err)
	}

	logToStep(step.ID, "Docker container stopped successfully")
	return nil
}

// executeShellCommand executes a shell command and returns its output
func executeShellCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// Errors
var (
	ErrInvalidConfig     = &PipelineError{Message: "Invalid step configuration"}
	ErrUnsupportedAction = &PipelineError{Message: "Unsupported action"}
)

type PipelineError struct {
	Message string
}

func (e *PipelineError) Error() string {
	return e.Message
}
