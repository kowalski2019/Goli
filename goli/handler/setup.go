package handler

import (
	aux "goli/auxiliary"
	response_util "goli/utils"

	"github.com/gin-gonic/gin"
)

// VerifySetupPasswordHandler verifies the setup password (no auth required)
func VerifySetupPasswordHandler(c *gin.Context) {
	var body struct {
		SetupPassword string `json:"setup_password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		response_util.SendBadRequestResponseGin(c, "Invalid request body: "+err.Error())
		return
	}

	config := aux.GetAllConfig()
	setupComplete := config["setup_complete"] == "true"
	storedPassword := config["setup_password"]

	// If setup is already complete, reject any setup password
	if setupComplete {
		response_util.SendUnauthorizedResponseGin(c, "Setup has already been completed")
		return
	}

	// If no setup password is stored, reject
	if storedPassword == "" {
		response_util.SendUnauthorizedResponseGin(c, "Setup password not configured")
		return
	}

	// Verify the password
	if body.SetupPassword != storedPassword {
		response_util.SendUnauthorizedResponseGin(c, "Invalid setup password")
		return
	}

	// Password is valid - return the auth key for setup completion
	authKey := config["auth_key"]
	response_util.SendJsonResponseGin(c, 200, gin.H{
		"message":  "Setup password verified",
		"auth_key": authKey,
	})
}

// GetSetupStatusHandler returns whether setup is complete (no auth required)
func GetSetupStatusHandler(c *gin.Context) {
	config := aux.GetAllConfig()
	setupComplete := config["setup_complete"] == "true"

	response_util.SendJsonResponseGin(c, 200, gin.H{
		"setup_complete": setupComplete,
	})
}

