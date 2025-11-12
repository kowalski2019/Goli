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
		"smtp_host":       config["smtp_host"],
		"smtp_port":       config["smtp_port"],
		"smtp_user":       config["smtp_user"],
		"smtp_pass":       config["smtp_pass"],
		"smtp_from":       config["smtp_from"],
		"smtp_from_name":  config["smtp_from_name"],
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
		SMTPHost      string `json:"smtp_host,omitempty"`
		SMTPPort      string `json:"smtp_port,omitempty"`
		SMTPUser      string `json:"smtp_user,omitempty"`
		SMTPPass      string `json:"smtp_pass,omitempty"`
		SMTPFrom      string `json:"smtp_from,omitempty"`
		SMTPFromName  string `json:"smtp_from_name,omitempty"`
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

	if body.SMTPHost != "" {
		updates["smtp_host"] = body.SMTPHost
	}

	if body.SMTPPort != "" {
		updates["smtp_port"] = body.SMTPPort
	}

	if body.SMTPUser != "" {
		updates["smtp_user"] = body.SMTPUser
	}

	if body.SMTPPass != "" {
		updates["smtp_pass"] = body.SMTPPass
	}

	if body.SMTPFrom != "" {
		updates["smtp_from"] = body.SMTPFrom
	}

	if body.SMTPFromName != "" {
		updates["smtp_from_name"] = body.SMTPFromName
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
		"smtp_host":       config["smtp_host"],
		"smtp_port":       config["smtp_port"],
		"smtp_user":       config["smtp_user"],
		"smtp_pass":       config["smtp_pass"],
		"smtp_from":       config["smtp_from"],
		"smtp_from_name":  config["smtp_from_name"],
	})
}
