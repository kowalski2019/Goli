package utils

import (
	"bytes"
	"context"
	aux "goli/auxiliary"
	"log"
	"os/exec"
	"strings"
	"time"
)

// AuthenticateGitHubContainerRegistry authenticates Docker with GitHub Container Registry
// using credentials from the config file. Returns error if authentication fails.
func AuthenticateGitHubContainerRegistry() error {
	config := aux.GetAllConfig()
	username := strings.TrimSpace(config["gh_username"])
	token := strings.TrimSpace(config["gh_access_token"])

	// Default/placeholder values that should not be used for authentication
	const defaultUsername = "dummy_gh_user"
	const defaultToken = "ghp_xxxxxxxxxxxxxxxxxxxxxxx"

	// If credentials are not set or are default/placeholder values, skip authentication
	if username == "" || token == "" {
		log.Println("GitHub credentials not configured, skipping GitHub Container Registry authentication")
		return nil
	}

	// Check if username is the default placeholder
	if username == defaultUsername {
		log.Println("GitHub username is set to default placeholder, skipping GitHub Container Registry authentication")
		return nil
	}

	// Check if token is the default placeholder
	if token == defaultToken || strings.HasPrefix(token, "ghp_xxxxxxxx") {
		log.Println("GitHub access token is set to default placeholder, skipping GitHub Container Registry authentication")
		return nil
	}

	log.Printf("Authenticating Docker with GitHub Container Registry as %s...", username)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create docker login command with password from stdin for security
	cmd := exec.CommandContext(ctx, "docker", "login", "ghcr.io", "-u", username, "--password-stdin")

	// Write password to stdin
	cmd.Stdin = bytes.NewBufferString(token)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	stderrStr := stderr.String()

	if err != nil {
		log.Printf("Docker login to GitHub Container Registry failed: %v, stderr: %s", err, stderrStr)
		return err
	}

	log.Printf("Successfully authenticated to GitHub Container Registry as %s", username)
	return nil
}

// EnsureGitHubAuthForImage checks if an image is from GitHub Container Registry (ghcr.io)
// and ensures authentication if needed. Returns true if authentication was attempted.
func EnsureGitHubAuthForImage(image string) bool {
	// Check if image is from GitHub Container Registry
	if strings.HasPrefix(image, "ghcr.io/") {
		err := AuthenticateGitHubContainerRegistry()
		if err != nil {
			log.Printf("Warning: Failed to authenticate with GitHub Container Registry for image %s: %v", image, err)
			return false
		}
		return true
	}
	return false
}
