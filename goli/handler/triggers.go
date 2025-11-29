package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"goli/types"
	response_util "goli/utils"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Helper function to decode request body and validate
func decodeAndValidateBodyGin(c *gin.Context) (*types.GoliRequestBody, bool) {
	var body types.GoliRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		response_util.SendBadRequestResponseGin(c, fmt.Sprintf("Invalid request body: %v", err))
		return nil, false
	}
	return &body, true
}

// Helper function to handle container action handlers
func handleContainerAction(c *gin.Context, action string) {
	body, ok := decodeAndValidateBodyGin(c)
	if !ok {
		return
	}

	if strings.TrimSpace(body.Name) == "" {
		response_util.SendBadRequestResponseGin(c, "Container name is required")
		return
	}

	res, err := DoDockerContainerAction(body.Name, action)
	if err != nil {
		response_util.SendInternalServerErrorResponseGin(c, err.Error())
		return
	}
	response_util.SendOkResponseGin(c, res)
}

// Helper function to handle image action handlers
func handleImageAction(c *gin.Context, action string) {
	body, ok := decodeAndValidateBodyGin(c)
	if !ok {
		return
	}

	if strings.TrimSpace(body.Image) == "" {
		response_util.SendBadRequestResponseGin(c, "Image name is required")
		return
	}

	res, err := DoDockerImageAction(body.Image, action)
	if err != nil {
		response_util.SendInternalServerErrorResponseGin(c, err.Error())
		return
	}
	response_util.SendOkResponseGin(c, res)
}

// Helper function to execute docker commands with timeout and proper error handling
func executeDockerCommand(ctx context.Context, args ...string) (string, string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	if err != nil {
		// Combine stderr with error message for better debugging
		if stderrStr != "" {
			return stdoutStr, stderrStr, fmt.Errorf("%v: %s", err, stderrStr)
		}
		return stdoutStr, stderrStr, err
	}

	return stdoutStr, stderrStr, nil
}

func StartADockerOrchestra(c *gin.Context) {
	response_util.SendOkResponseGin(c, "Is fine we can start the docker compose")
}

func StopADockerOrchestra(c *gin.Context) {
	response_util.SendOkResponseGin(c, "Is fine we can stop the docker compose")
}

func StartADocker(c *gin.Context) {
	handleContainerAction(c, "start")
}

func StopADocker(c *gin.Context) {
	handleContainerAction(c, "stop")
}

func RemoveADocker(c *gin.Context) {
	handleContainerAction(c, "rm")
}

func PauseADocker(c *gin.Context) {
	handleContainerAction(c, "pause")
}

func UnPauseADocker(c *gin.Context) {
	handleContainerAction(c, "unpause")
}

func InspectADocker(c *gin.Context) {
	handleContainerAction(c, "inspect")
}

func GetADockerLogs(c *gin.Context) {
	handleContainerAction(c, "logs")
}

func GetDockerPS(c *gin.Context) {
	stdout, stderr, err := executeDockerCommand(c.Request.Context(), "ps", "-a")
	if err != nil {
		errorMsg := stderr
		if errorMsg == "" {
			errorMsg = err.Error()
		}
		response_util.SendInternalServerErrorResponseGin(c, errorMsg)
		return
	}
	response_util.SendOkResponseGin(c, stdout)
}

func GetDockerImages(c *gin.Context) {
	stdout, stderr, err := executeDockerCommand(c.Request.Context(), "images")
	if err != nil {
		errorMsg := stderr
		if errorMsg == "" {
			errorMsg = err.Error()
		}
		response_util.SendInternalServerErrorResponseGin(c, errorMsg)
		return
	}
	response_util.SendOkResponseGin(c, stdout)
}

func RemoveAnDockerImage(c *gin.Context) {
	handleImageAction(c, "rm")
}

func PullAnDockerImage(c *gin.Context) {
	handleImageAction(c, "pull")
}

func RunDockerContainer(c *gin.Context) {
	body, ok := decodeAndValidateBodyGin(c)
	if !ok {
		return
	}

	// Validate required fields
	if strings.TrimSpace(body.Name) == "" {
		response_util.SendBadRequestResponseGin(c, "Container name is required")
		return
	}
	if strings.TrimSpace(body.Image) == "" {
		response_util.SendBadRequestResponseGin(c, "Image name is required")
		return
	}

	containerName := body.Name
	containerImage := body.Image
	containerExists := checkDockerExistence(c.Request.Context(), containerName)

	if !containerExists {
		log.Printf("Container %s does not exist, creating new one", containerName)
		res, err := createContainer(c.Request.Context(), containerImage, containerName, body.Network,
			body.Port_Ex, body.Port_In, body.Volume_Ex, body.Volume_In, body.V_Map, body.Opts)
		if err != nil {
			response_util.SendInternalServerErrorResponseGin(c, err.Error())
			return
		}
		response_util.SendOkResponseGin(c, res)
		return
	}

	// Container exists - stop, remove, pull image, and recreate
	log.Printf("Container %s already exists, recreating", containerName)
	var results []string

	// Stop container
	if res, err := DoDockerContainerAction(containerName, "stop"); err != nil {
		response_util.SendInternalServerErrorResponseGin(c, fmt.Sprintf("Failed to stop container: %v", err))
		return
	} else {
		results = append(results, res)
	}

	// Remove container
	if res, err := DoDockerContainerAction(containerName, "rm"); err != nil {
		response_util.SendInternalServerErrorResponseGin(c, fmt.Sprintf("Failed to remove container: %v", err))
		return
	} else {
		results = append(results, res)
	}

	// Remove old image (ignore errors - image might not exist)
	if res, err := DoDockerImageAction(containerImage, "rm"); err == nil {
		results = append(results, res)
	}

	// Pull latest image
	if res, err := DoDockerImageAction(containerImage, "pull"); err != nil {
		response_util.SendInternalServerErrorResponseGin(c, fmt.Sprintf("Failed to pull image: %v", err))
		return
	} else {
		results = append(results, res)
	}

	// Create new container
	res, err := createContainer(c.Request.Context(), containerImage, containerName, body.Network,
		body.Port_Ex, body.Port_In, body.Volume_Ex, body.Volume_In, body.V_Map, body.Opts)
	if err != nil {
		response_util.SendInternalServerErrorResponseGin(c, err.Error())
		return
	}
	results = append(results, res)

	response_util.SendOkResponseGin(c, strings.Join(results, "\n"))
}

func checkDockerExistence(ctx context.Context, name string) bool {
	// Use docker ps with filter to check if container exists (more efficient than logs)
	stdout, _, err := executeDockerCommand(ctx, "ps", "-a", "--filter", fmt.Sprintf("name=^%s$", name), "--format", "{{.Names}}")
	if err != nil {
		return false
	}
	return strings.TrimSpace(stdout) == name
}

func DoDockerContainerAction(container string, action string) (string, error) {
	// Validate container name to prevent command injection
	if strings.TrimSpace(container) == "" {
		return "", errors.New("container name cannot be empty")
	}

	var args []string
	switch action {
	case "start":
		args = []string{"start", container}
	case "stop":
		args = []string{"stop", container}
	case "rm":
		args = []string{"rm", "-f", container}
	case "pause":
		args = []string{"pause", container}
	case "unpause":
		args = []string{"unpause", container}
	case "inspect":
		args = []string{"inspect", container}
	case "logs":
		args = []string{"logs", container}
	default:
		return "", fmt.Errorf("unknown action: %s", action)
	}

	log.Printf("Executing: docker %s", strings.Join(args, " "))
	ctx := context.Background()
	stdout, stderr, err := executeDockerCommand(ctx, args...)
	if err != nil {
		if stderr != "" {
			return "", fmt.Errorf("%v: %s", err, stderr)
		}
		return "", err
	}
	return stdout, nil
}

func DoDockerImageAction(image string, action string) (string, error) {
	// Validate image name to prevent command injection
	if strings.TrimSpace(image) == "" {
		return "", errors.New("image name cannot be empty")
	}

	// Ensure GitHub Container Registry authentication if pulling from ghcr.io
	if action == "pull" {
		response_util.EnsureGitHubAuthForImage(image)
	}

	var args []string
	switch action {
	case "rm":
		args = []string{"rmi", "-f", image}
	case "pull":
		args = []string{"pull", image}
	default:
		return "", fmt.Errorf("unknown action: %s", action)
	}

	log.Printf("Executing: docker %s", strings.Join(args, " "))
	ctx := context.Background()
	stdout, stderr, err := executeDockerCommand(ctx, args...)
	if err != nil {
		if stderr != "" {
			return "", fmt.Errorf("%v: %s", err, stderr)
		}
		return "", err
	}
	return stdout, nil
}

func createContainer(ctx context.Context, image string, name string, network string, portEx string, portIn string, volumeEx string, volumeIn string, vMap bool, opts string) (string, error) {
	// Validate required fields
	if strings.TrimSpace(image) == "" {
		return "", errors.New("image name cannot be empty")
	}
	if strings.TrimSpace(name) == "" {
		return "", errors.New("container name cannot be empty")
	}

	// Build docker run command arguments
	args := []string{"run", "--detach", "--name", name}

	// Add network configuration
	if network == "host" {
		args = append(args, "--network", "host")
	} else if strings.TrimSpace(portEx) != "" && strings.TrimSpace(portIn) != "" {
		// Only add port mapping if not using host network
		portMapping := fmt.Sprintf("%s:%s", portEx, portIn)
		args = append(args, "-p", portMapping)
	}

	// Add volume mapping if enabled
	if vMap && strings.TrimSpace(volumeEx) != "" && strings.TrimSpace(volumeIn) != "" {
		volumeMapping := fmt.Sprintf("%s:%s", volumeEx, volumeIn)
		args = append(args, "-v", volumeMapping)
	}

	// Parse and add additional options
	if strings.TrimSpace(opts) != "" {
		// Split opts by space and append (basic parsing - could be improved)
		optParts := strings.Fields(opts)
		args = append(args, optParts...)
	}

	// Add image name at the end
	args = append(args, image)

	log.Printf("Executing: docker %s", strings.Join(args, " "))

	stdout, stderr, err := executeDockerCommand(ctx, args...)
	if err != nil {
		errorMsg := stderr
		if errorMsg == "" {
			errorMsg = err.Error()
		}
		return "", fmt.Errorf("failed to create container: %s", errorMsg)
	}

	containerID := strings.TrimSpace(stdout)
	return fmt.Sprintf("Docker container %s (%s) successfully started!", name, containerID), nil
}
