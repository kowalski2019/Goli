package handler

import (
	aux "goli/auxiliary"
	response_util "goli/utils"

	"github.com/gin-gonic/gin"
)

// GetConfigHandler returns the current configuration
func GetConfigHandler(c *gin.Context) {
	config := aux.GetAllConfig()

	// Convert setup_complete string to boolean
	setupComplete := config["setup_complete"] == "true"

	response_util.SendJsonResponseGin(c, 200, gin.H{
		"port":            config["port"],
		"auth_key":        config["auth_key"],
		"setup_complete":  setupComplete,
		"gh_username":     config["gh_username"],
		"gh_access_token": config["gh_access_token"],
	})
}

// UpdateConfigHandler updates the configuration
func UpdateConfigHandler(c *gin.Context) {
	var body struct {
		Port          string `json:"port,omitempty"`
		AuthKey       string `json:"auth_key,omitempty"`
		SetupComplete *bool  `json:"setup_complete,omitempty"`
		GHUsername    string `json:"gh_username,omitempty"`
		GHAccessToken string `json:"gh_access_token,omitempty"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		response_util.SendBadRequestResponseGin(c, "Invalid request body: "+err.Error())
		return
	}

	updates := make(map[string]string)

	if body.Port != "" {
		updates["port"] = body.Port
	}

	if body.AuthKey != "" {
		updates["auth_key"] = body.AuthKey
	}

	if body.SetupComplete != nil {
		if *body.SetupComplete {
			updates["setup_complete"] = "true"
			// Invalidate setup password after successful setup
			updates["setup_password"] = ""
		} else {
			updates["setup_complete"] = "false"
		}
	}

	if body.GHUsername != "" {
		updates["gh_username"] = body.GHUsername
	}

	if body.GHAccessToken != "" {
		updates["gh_access_token"] = body.GHAccessToken
	}

	if len(updates) == 0 {
		response_util.SendBadRequestResponseGin(c, "No fields to update")
		return
	}

	if err := aux.UpdateConfig(updates); err != nil {
		response_util.SendInternalServerErrorResponseGin(c, "Failed to update config: "+err.Error())
		return
	}

	// Return updated config
	config := aux.GetAllConfig()
	setupComplete := config["setup_complete"] == "true"

	response_util.SendJsonResponseGin(c, 200, gin.H{
		"port":            config["port"],
		"auth_key":        config["auth_key"],
		"setup_complete":  setupComplete,
		"gh_username":     config["gh_username"],
		"gh_access_token": config["gh_access_token"],
	})
}

