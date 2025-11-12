package handler

import (
	"bytes"
	"context"
	"fmt"
	response_util "goli/utils"
	"log"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthenticateContainerRegistryHandler handles GitHub Container Registry authentication
func AuthenticateContainerRegistryHandler(c *gin.Context) {
	var body struct {
		GHUsername    string `json:"gh_username"`
		GHAccessToken string `json:"gh_access_token"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		response_util.SendBadRequestResponseGin(c, "Invalid request body: "+err.Error())
		return
	}

	if body.GHUsername == "" || body.GHAccessToken == "" {
		response_util.SendBadRequestResponseGin(c, "GitHub username and access token are required")
		return
	}

	// Execute docker login with GitHub Container Registry
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create docker login command with password from stdin for security
	cmd := exec.CommandContext(ctx, "docker", "login", "ghcr.io", "-u", body.GHUsername, "--password-stdin")

	// Write password to stdin
	cmd.Stdin = bytes.NewBufferString(body.GHAccessToken)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	stderrStr := stderr.String()

	if err != nil {
		log.Printf("Docker login failed: %v, stderr: %s", err, stderrStr)
		response_util.SendInternalServerErrorResponseGin(c, fmt.Sprintf("Failed to authenticate with GitHub Container Registry: %s", stderrStr))
		return
	}

	log.Printf("Successfully authenticated to GitHub Container Registry as %s", body.GHUsername)
	response_util.SendOkResponseGin(c, "Successfully authenticated to GitHub Container Registry")
}
