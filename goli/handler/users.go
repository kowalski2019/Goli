package handler

import (
	"goli/database"
	"goli/models"
	response_util "goli/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListUsersHandler lists all users
func ListUsersHandler(c *gin.Context) {
	users, err := database.ListUsers()
	if err != nil {
		response_util.SendInternalServerErrorResponseGin(c, "Failed to list users: "+err.Error())
		return
	}

	response_util.SendJsonResponseGin(c, 200, users)
}

// CreateUserHandler creates a new user
func CreateUserHandler(c *gin.Context) {
	var body models.UserCreateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		response_util.SendBadRequestResponseGin(c, "Invalid request body: "+err.Error())
		return
	}

	if body.Username == "" || body.Password == "" {
		response_util.SendBadRequestResponseGin(c, "Username and password are required")
		return
	}

	if len(body.Password) < 6 {
		response_util.SendBadRequestResponseGin(c, "Password must be at least 6 characters")
		return
	}

	user := &models.User{
		Username:          body.Username,
		Email:             body.Email,
		Phone:             body.Phone,
		Password:          body.Password,
		Role:              body.Role,
		TwoFAEmailEnabled: body.TwoFAEmailEnabled,
		TwoFASmsEnabled:   body.TwoFASmsEnabled,
	}

	if user.Role == "" {
		user.Role = "user"
	}

	createdUser, err := database.CreateUser(user)
	if err != nil {
		response_util.SendInternalServerErrorResponseGin(c, "Failed to create user: "+err.Error())
		return
	}

	// Don't return password
	createdUser.Password = ""
	createdUser.PasswordHash = ""

	response_util.SendJsonResponseGin(c, 201, createdUser)
}

// UpdateUserHandler updates an existing user
func UpdateUserHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response_util.SendBadRequestResponseGin(c, "Invalid user ID")
		return
	}

	var body models.UserUpdateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		response_util.SendBadRequestResponseGin(c, "Invalid request body: "+err.Error())
		return
	}

	// Get existing user
	user, err := database.GetUser(id)
	if err != nil {
		response_util.SendNotFoundResponseGin(c, "User not found")
		return
	}

	// Update fields
	if body.Email != "" {
		user.Email = body.Email
	}
	if body.Phone != "" {
		user.Phone = body.Phone
	}
	if body.Role != "" {
		user.Role = body.Role
	}
	if body.TwoFAEmailEnabled != nil {
		user.TwoFAEmailEnabled = *body.TwoFAEmailEnabled
	}
	if body.TwoFASmsEnabled != nil {
		user.TwoFASmsEnabled = *body.TwoFASmsEnabled
	}

	updatePassword := body.Password != ""
	if updatePassword {
		if len(body.Password) < 6 {
			response_util.SendBadRequestResponseGin(c, "Password must be at least 6 characters")
			return
		}
		user.Password = body.Password
	}

	if err := database.UpdateUser(user, updatePassword); err != nil {
		response_util.SendInternalServerErrorResponseGin(c, "Failed to update user: "+err.Error())
		return
	}

	// Don't return password
	user.Password = ""
	user.PasswordHash = ""

	response_util.SendJsonResponseGin(c, 200, user)
}

// DeleteUserHandler deletes a user
func DeleteUserHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response_util.SendBadRequestResponseGin(c, "Invalid user ID")
		return
	}

	if err := database.DeleteUser(id); err != nil {
		response_util.SendInternalServerErrorResponseGin(c, "Failed to delete user: "+err.Error())
		return
	}

	response_util.SendOkResponseGin(c, "User deleted successfully")
}

