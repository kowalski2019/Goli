package middlewares

import (
	aux "goli/auxiliary"
	"goli/database"
	response_util "goli/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var auth_key = aux.GetFromConfig("constants.auth_key")

// AuthMiddleware returns a Gin middleware that verifies authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response_util.SendUnauthorizedResponseGin(c, "Missing Authorization header")
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 {
			response_util.SendUnauthorizedResponseGin(c, "Invalid Authorization header")
			c.Abort()
			return
		}
		scheme := parts[0]
		cred := parts[1]

		// New: Bearer session token
		if strings.EqualFold(scheme, "Bearer") {
			session, err := database.GetSessionByToken(cred)
			if err != nil || session == nil {
				response_util.SendUnauthorizedResponseGin(c, "Invalid session")
				c.Abort()
				return
			}
			if time.Now().After(session.ExpiresAt) {
				_ = database.DeleteSession(cred)
				response_util.SendUnauthorizedResponseGin(c, "Session expired")
				c.Abort()
				return
			}
			// Store session info in context for handlers to use
			c.Set("session_token", cred)
			c.Set("user_id", session.UserID)
			c.Next()
			return
		}

		// Legacy support
		if strings.EqualFold(scheme, "Goli-Auth-Key") && cred == auth_key {
			c.Next()
			return
		}

		response_util.SendUnauthorizedResponseGin(c, "Unauthorized")
		c.Abort()
	}
}
